package pool

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v8/sdk/access"
)

type Client struct {
	Inference string `json:"inference"`
}

type ClientRegistry struct {
	path     string
	mu       sync.Mutex
	modified time.Time
	size     int64
	clients  map[string]Client
}

func Fingerprint(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func ReadClients(path string) (map[string]Client, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, 1024*1024+1))
	if err != nil || len(raw) > 1024*1024 {
		return nil, fmt.Errorf("cannot read client registry")
	}
	var clients map[string]Client
	if err = json.Unmarshal(raw, &clients); err != nil {
		return nil, fmt.Errorf("invalid client registry")
	}
	seen := make(map[string]bool)
	for name, client := range clients {
		if strings.TrimSpace(name) == "" || len(name) > 80 {
			return nil, fmt.Errorf("invalid client name")
		}
		if len(client.Inference) < 32 || seen[client.Inference] {
			return nil, fmt.Errorf("invalid or repeated client key")
		}
		seen[client.Inference] = true
	}
	return clients, nil
}

// Snapshot notices atomic registry replacements without caching revoked keys on errors.
func (r *ClientRegistry) Snapshot() (map[string]Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	info, err := os.Stat(r.path)
	if err != nil {
		return nil, err
	}
	if r.clients == nil || !info.ModTime().Equal(r.modified) || info.Size() != r.size {
		clients, errRead := ReadClients(r.path)
		if errRead != nil {
			return nil, errRead
		}
		r.clients, r.modified, r.size = clients, info.ModTime(), info.Size()
	}
	out := make(map[string]Client, len(r.clients))
	for name, client := range r.clients {
		out[name] = client
	}
	return out, nil
}

func bearer(request *http.Request, inference bool) string {
	value := request.Header.Get("Authorization")
	if len(value) >= 7 && strings.EqualFold(value[:7], "Bearer ") {
		value = strings.TrimSpace(value[7:])
	}
	if inference && request.Header.Get("X-Api-Key") != "" {
		value = request.Header.Get("X-Api-Key")
	}
	return value
}

func (r *ClientRegistry) match(key string) (string, error) {
	clients, err := r.Snapshot()
	if err != nil {
		return "", err
	}
	wanted := Fingerprint(key)
	for name, client := range clients {
		if Fingerprint(client.Inference) == wanted {
			return name, nil
		}
	}
	return "", nil
}

func (r *ClientRegistry) Identifier() string { return "pool-clients" }

func (r *ClientRegistry) Authenticate(_ context.Context, request *http.Request) (*access.Result, *access.AuthError) {
	key := bearer(request, true)
	if key == "" {
		return nil, access.NewNoCredentialsError()
	}
	name, err := r.match(key)
	if err != nil {
		return nil, access.NewInvalidCredentialError()
	}
	if name == "" {
		return nil, access.NewNotHandledError()
	}
	return &access.Result{Provider: r.Identifier(), Principal: key, Metadata: map[string]string{"name": name}}, nil
}

func NewClient() (Client, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return Client{}, err
	}
	return Client{Inference: base64.RawURLEncoding.EncodeToString(data)}, nil
}

func ValidName(name string) bool { return strings.TrimSpace(name) != "" && len(name) <= 80 }

// UpdateClients edits the registry under an exclusive lock and writes it when change succeeds.
func UpdateClients(path string, change func(map[string]Client) error) error {
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	clients, err := ReadClients(path)
	if os.IsNotExist(err) {
		clients = map[string]Client{}
	} else if err != nil {
		return err
	}
	if err = change(clients); err != nil {
		return err
	}
	raw, err := json.Marshal(clients)
	if err != nil {
		return err
	}
	return WritePrivateFile(path, raw)
}

func AddClient(clients map[string]Client, name string) (Client, error) {
	if !ValidName(name) {
		return Client{}, fmt.Errorf("choose a name of 1–80 characters")
	}
	if _, exists := clients[name]; exists {
		return Client{}, fmt.Errorf("client already exists")
	}
	client, err := NewClient()
	if err == nil {
		clients[name] = client
	}
	return client, err
}

func RenameClient(clients map[string]Client, name, newName string) error {
	client, exists := clients[name]
	if !exists {
		return fmt.Errorf("unknown client")
	}
	if !ValidName(newName) {
		return fmt.Errorf("choose a name of 1–80 characters")
	}
	if _, exists = clients[newName]; exists {
		return fmt.Errorf("new name already exists")
	}
	delete(clients, name)
	clients[newName] = client
	return nil
}

func RevokeClient(clients map[string]Client, name string) error {
	if _, exists := clients[name]; !exists {
		return fmt.Errorf("unknown client")
	}
	delete(clients, name)
	return nil
}

func WritePrivateFile(path string, data []byte) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".pool-client-*")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if existing, errStat := os.Stat(path); errStat == nil {
		// Only root editing the service's file needs this. Sandboxed services may not chown.
		if owner, ok := existing.Sys().(*syscall.Stat_t); ok && (int(owner.Uid) != os.Geteuid() || int(owner.Gid) != os.Getegid()) {
			if err = file.Chown(int(owner.Uid), int(owner.Gid)); err != nil {
				file.Close()
				return err
			}
		}
	}
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	if errClose := file.Close(); err == nil {
		err = errClose
	}
	if err != nil {
		return err
	}
	if err = os.Rename(name, path); err != nil {
		return err
	}
	directory, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}
