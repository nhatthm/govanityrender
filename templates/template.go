// Package templates provides the embedded templates.
package templates

import _ "embed"

//go:embed homepage.html.hbs
var homepageTpl string

//go:embed repository.html.hbs
var repositoryTpl string

// EmbeddedHomepage provides the homepage template.
func EmbeddedHomepage() string {
	return homepageTpl
}

// EmbeddedRepository provides the repository template.
func EmbeddedRepository() string {
	return repositoryTpl
}
