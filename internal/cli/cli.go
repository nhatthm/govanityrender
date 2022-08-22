package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go.nhat.io/vanityrender/internal/config"
	"go.nhat.io/vanityrender/internal/git"
	"go.nhat.io/vanityrender/internal/github"
	"go.nhat.io/vanityrender/internal/site"
	"go.nhat.io/vanityrender/templates"
)

// Execute is the entrypoint for the cli.
func Execute() int {
	var (
		configFile  string
		homepageTpl string
		outputPath  string
	)

	flag.StringVar(&configFile, "config", "config.json", "config file")
	flag.StringVar(&homepageTpl, "homepage-tpl", "", "template file")
	flag.StringVar(&outputPath, "out", "build", "output path")

	flag.Parse()

	err := runRender(configFile, homepageTpl, outputPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)

		return 1
	}

	return 0
}

func runRender(configFile string, homepageTpl string, outputPath string) error {
	cfg, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	homepageSrc, err := initHomepageSrc(homepageTpl)
	if err != nil {
		return err
	}

	outputPath, err = initOutputDir(outputPath)
	if err != nil {
		return err
	}

	siteCfg, err := initSiteConfig(cfg)
	if err != nil {
		return err
	}

	r, err := site.NewHandlebarsRenderder(homepageSrc, templates.EmbeddedRepository(), outputPath)
	if err != nil {
		return err
	}

	return r.Render(*siteCfg)
}

func initConfigHydrators() []site.Hydrator {
	return []site.Hydrator{
		github.NewHydrator(git.NewModuleFinder()),
	}
}

func initHomepageSrc(homepageTpl string) (string, error) {
	if len(homepageTpl) > 0 {
		data, err := os.ReadFile(filepath.Clean(homepageTpl))
		if err != nil {
			return "", fmt.Errorf("could not read homepage template: %w", err)
		}

		return string(data), nil
	}

	return templates.EmbeddedHomepage(), nil
}

func initOutputDir(outputPath string) (string, error) {
	fi, err := os.Stat(filepath.Clean(outputPath))
	if err == nil {
		if !fi.IsDir() {
			return "", fmt.Errorf("output path %q is not a directory", outputPath) // nolint: goerr113
		}

		return outputPath, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("could not stat output path: %w", err)
	}

	if err := os.MkdirAll(filepath.Clean(outputPath), 0o755); err != nil { // nolint: gosec
		return "", fmt.Errorf("could not create output directory %q: %w", outputPath, err)
	}

	return outputPath, nil
}

func initSiteConfig(cfg config.Config) (*site.Site, error) {
	s := site.Site{
		PageTitle:       cfg.PageTitle,
		PageDescription: cfg.PageDescription,
		Hostname:        cfg.Host,
		Repositories:    make([]site.Repository, len(cfg.Repositories)),
	}

	for i, r := range cfg.Repositories {
		s.Repositories[i] = site.Repository{
			Name:          r.Name,
			Path:          r.Path,
			Deprecated:    r.Deprecated,
			RepositoryURL: r.Repository,
			Ref:           r.Ref,
		}
	}

	err := site.Hydrate(&s, initConfigHydrators()...)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
