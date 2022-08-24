// Package templates provides the embedded templates.
package templates

import _ "embed"

//go:embed homepage.html.hbs
var homepageTpl string

//go:embed 404.html.hbs
var notFoundTpl string

//go:embed repository.html.hbs
var repositoryTpl string

// EmbeddedHomepage provides the homepage template.
func EmbeddedHomepage() string {
	return homepageTpl
}

// EmbeddedNotFound provides the 404 Not Found template.
func EmbeddedNotFound() string {
	return notFoundTpl
}

// EmbeddedRepository provides the repository template.
func EmbeddedRepository() string {
	return repositoryTpl
}
