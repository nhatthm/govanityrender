package config

// Hydrator hydrates configuration.
type Hydrator interface {
	Hydrate(cfg *Config) error
}
