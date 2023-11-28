package config

type BrandingConfig struct {
	Name string `mapstructure:"name"`
}

type StoredConnection struct {
	Tag            string `mapstructure:"tag"`
	Host           string `mapstructure:"host"`
	User           string `mapstructure:"user"`
	Port           int    `mapstructure:"por"`
	SshKeyPath     string `mapstructure:"ssh_key_path"`
	Provider       string `mapstructure:"provider"`
	DefaltDatabase string `mapstructure:"default_database"`
}

type SeraphimConfig struct {
	Version           string             `mapstructure:"version"`
	BrandingConfig    BrandingConfig     `mapstructure:"branding"`
	StoredConnections []StoredConnection `mapstructure:"stored_connections"`
}
