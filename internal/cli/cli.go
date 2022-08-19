package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go.nhat.io/vanityrender/internal/config"
	"go.nhat.io/vanityrender/internal/git"
	"go.nhat.io/vanityrender/internal/renderer"
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
	cfg, err := config.HydrateFile(configFile, initConfigHydrators()...)
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

	r, err := renderer.NewHandlebarsRenderder(homepageSrc, templates.EmbeddedRepository(), outputPath)
	if err != nil {
		return err
	}

	return r.Render(cfg)
}

func initConfigHydrators() []config.Hydrator {
	return []config.Hydrator{
		git.NewHydrator(),
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
