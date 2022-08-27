package cli

import (
	"crypto/sha1" // nolint: gosec
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-colorable"

	"go.nhat.io/vanityrender/internal/config"
	"go.nhat.io/vanityrender/internal/git"
	"go.nhat.io/vanityrender/internal/github"
	"go.nhat.io/vanityrender/internal/service/sitecache"
	"go.nhat.io/vanityrender/internal/service/sitefragment"
	"go.nhat.io/vanityrender/internal/site"
	"go.nhat.io/vanityrender/templates"
)

// Execute is the entrypoint for the cli.
func Execute() int {
	var (
		configFile  string
		homepageTpl string
		outputPath  string
		modulesVal  string
		noColor     bool
	)

	flag.StringVar(&configFile, "config", "config.json", "config file")
	flag.StringVar(&homepageTpl, "homepage-tpl", "", "template file")
	flag.StringVar(&outputPath, "out", "build", "output path")
	flag.StringVar(&modulesVal, "modules", "", "rebuild only the listed modules, comma separated")
	flag.BoolVar(&noColor, "no-color", false, "do not use colors in output")

	flag.Parse()

	modules := split(strings.Trim(modulesVal, "\r\n "), ",")

	out := colorable.NewNonColorable(os.Stdout)
	if !noColor {
		out = colorable.NewColorable(os.Stdout)
	}

	err := runRender(out, configFile, homepageTpl, outputPath, modules)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)

		return 1
	}

	return 0
}

func runRender(out io.Writer, configFile string, homepageTpl string, outputPath string, modules []string) error {
	checksum, err := checksum(configFile)
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

	siteCfg, err := initSiteConfig(out, configFile, checksum, modules)
	if err != nil {
		return err
	}

	r, err := initRenderer(out, homepageSrc, outputPath, checksum)
	if err != nil {
		return err
	}

	return r.Render(*siteCfg)
}

func initConfigHydrators(out io.Writer, checksum string, modules []string) []site.Hydrator {
	return []site.Hydrator{
		sitefragment.NewHydrator(
			sitecache.NewMetadataHydrator(checksum, sitecache.WithOutput(out)),
			github.NewHydrator(git.NewModuleFinder(), github.WithOutput(out)),
			modules,
			sitefragment.WithOutput(out),
		),
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

func initSiteConfig(out io.Writer, configFile, checksum string, modules []string) (*site.Site, error) {
	cfg, err := config.FromFile(configFile)
	if err != nil {
		return nil, err
	}

	s := site.Site{
		PageTitle:       cfg.PageTitle,
		PageDescription: cfg.PageDescription,
		Hostname:        cfg.Host,
		SourceURL:       cfg.SourceURL,
		Repositories:    make([]site.Repository, len(cfg.Repositories)),
	}

	for i, r := range cfg.Repositories {
		s.Repositories[i] = site.Repository{
			Name:          r.Name,
			Path:          r.Path,
			Deprecated:    r.Deprecated,
			Hidden:        r.Hidden,
			RepositoryURL: r.Repository,
			Ref:           r.Ref,
		}
	}

	err = site.Hydrate(&s, initConfigHydrators(out, checksum, modules)...)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func initRenderer(out io.Writer, homepageSrc, outputPath, checksum string) (site.Renderder, error) {
	var r site.Renderder

	r, err := site.NewHandlebarsRenderder(homepageSrc, templates.EmbeddedNotFound(), templates.EmbeddedRepository(), outputPath, site.WithOutput(out))
	if err != nil {
		return nil, err
	}

	r = sitecache.NewRenderder(r, outputPath, checksum, sitecache.WithOutput(out))

	return r, nil
}

func checksum(path string) (string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("could not open config file: %w", err)
	}

	defer func() {
		_ = f.Close() // nolint: errcheck
	}()

	h := sha1.New() // nolint: gosec
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("could not calculate checksum: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func split(s, sep string) []string {
	var r []string

	for _, str := range strings.Split(s, sep) {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}
