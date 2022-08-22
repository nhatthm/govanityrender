package site

// Hydrator hydrates configuration.
type Hydrator interface {
	Hydrate(s *Site) error
}

// Hydrate hydrates the configuration.
func Hydrate(s *Site, hydrators ...Hydrator) error {
	for _, hydrator := range hydrators {
		if err := hydrator.Hydrate(s); err != nil {
			return err
		}
	}

	return nil
}
