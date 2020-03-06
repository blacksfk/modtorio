package api

// returned from mods.factorio.com/api/mods
type ModListResponse struct {
	Pagination Pagination
	Results    []*Result
}

// pagination data (sub-struct) of a ModListResponse
type Pagination struct {
	Count, Page           int
	Page_count, Page_size int
	Links                 struct {
		First, Prev string
		Next, Last  string
	}
}

// mod data
type Result struct {
	Downloads_count              uint
	Name, Owner, Summary         string
	Title, Changelog, Created_at string
	Description, Github_path     string
	Category, Homepage           string
	Latest_release               Release
	Releases                     []Release
	Tag                          []Tag
}

// pretty print a mod's information
func (r *Result) String() string {
	return fmt.Sprintf("%s\n\tName: %s\n\tOwner: %s\n\tCategory: %s\n\tSummary: %s\n", r.Title, r.Name, r.Owner, r.Category, r.Summary)
}

// specific release information of a mod
type Release struct {
	Download_url, File_name    string
	Released_at, Version, sha1 string
	Info_json                  struct {
		Factorio_version string
	}
}

// mod tags (refer to array above)
type Tag struct {
	Id                             int
	Name, Title, Description, Type string
}

// JSON API errors
type apiError struct {
	message string
}
