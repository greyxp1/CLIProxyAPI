package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/pool"
)

var errUsage = errors.New("usage: cliproxy-client [--clients-file PATH] list|add NAME|rename NAME NEW_NAME|revoke NAME|export NAME DIRECTORY")

func run() error {
	flags := flag.NewFlagSet("cliproxy-client", flag.ContinueOnError)
	path := flags.String("clients-file", "/var/lib/cliproxyapi/clients.json", "Named client registry")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}
	args := flags.Args()
	if len(args) == 0 {
		return errUsage
	}
	switch {
	case args[0] == "list" && len(args) == 1:
		clients, err := pool.ReadClients(*path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		names := make([]string, 0, len(clients))
		for name := range clients {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Println(name)
		}
		return nil
	case args[0] == "export" && len(args) == 3:
		clients, err := pool.ReadClients(*path)
		if err != nil {
			return err
		}
		client, exists := clients[args[1]]
		if !exists {
			return fmt.Errorf("unknown client")
		}
		if err = os.MkdirAll(args[2], 0700); err != nil {
			return err
		}
		if err = pool.WritePrivateFile(filepath.Join(args[2], "inference-key"), []byte(client.Inference+"\n")); err != nil {
			return err
		}
		fmt.Println("Key saved in " + args[2])
		return nil
	}
	err := pool.UpdateClients(*path, func(clients map[string]pool.Client) error {
		switch {
		case args[0] == "add" && len(args) == 2:
			_, err := pool.AddClient(clients, args[1])
			return err
		case args[0] == "rename" && len(args) == 3:
			return pool.RenameClient(clients, args[1], args[2])
		case args[0] == "revoke" && len(args) == 2:
			return pool.RevokeClient(clients, args[1])
		}
		return errUsage
	})
	if err != nil {
		return err
	}
	fmt.Println("Saved. The running proxy reloads named keys automatically.")
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
