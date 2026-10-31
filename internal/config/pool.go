package config

// PoolConfig enables the shared Claude and Codex subscription pool.
type PoolConfig struct {
	Enabled     bool   `yaml:"enabled" json:"enabled"`
	Address     string `yaml:"address" json:"address"`
	Origin      string `yaml:"origin" json:"origin"`
	ClientsFile string `yaml:"clients-file" json:"clients-file"`
	Database    string `yaml:"database" json:"database"`
	Maintenance bool   `yaml:"maintenance" json:"maintenance"`
}
