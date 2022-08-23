package site

// Site is the site configuration.
type Site struct {
	PageTitle       string       `json:"page_title"`
	PageDescription string       `json:"page_description"`
	Hostname        string       `json:"hostname"`
	Repositories    []Repository `json:"repositories"`
}

// Repository is a repository configuration.
type Repository struct {
	Name           string   `json:"name"`
	Path           string   `json:"path"`
	Deprecated     string   `json:"deprecated"`
	RepositoryURL  string   `json:"repository_url"`
	RepositoryName string   `json:"repository_name"`
	Ref            string   `json:"ref"`
	LatestVersion  string   `json:"latest_version"`
	Modules        []Module `json:"modules"`
}

// Module is a module configuration.
type Module struct {
	Path          string `json:"path"`
	ImportPrefix  string `json:"import_prefix"`
	VCS           string `json:"vcs"`
	RepositoryURL string `json:"repository_url"`
	HomeURL       string `json:"home_url"`
	DirectoryURL  string `json:"directory_url"`
	FileURL       string `json:"file_url"`
}
