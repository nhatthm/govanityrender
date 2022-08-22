package site

// Site is the site configuration.
type Site struct {
	PageTitle       string
	PageDescription string
	Hostname        string
	Repositories    []Repository
}

// Repository is a repository configuration.
type Repository struct {
	Name           string
	Path           string
	Deprecated     string
	RepositoryURL  string
	RepositoryName string
	Ref            string
	LatestVersion  string
	Modules        []Module
}

// Module is a module configuration.
type Module struct {
	Path          string
	ImportPrefix  string
	VCS           string
	RepositoryURL string
	HomeURL       string
	DirectoryURL  string
	FileURL       string
}
