package methods

import "github.com/synology-community/go-synology/pkg/api"

var (
	Search = api.Method{
		API:     "SYNO.Finder.FileIndexing.Search",
		Method:  "search",
		Version: 1,
	}
)
