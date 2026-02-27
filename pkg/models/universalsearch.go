package models

type SearchRequest struct {
	Keyword          string `url:"keyword"`
	SearchWeightList string `url:"search_weight_list"`
	QueryID          string `url:"query_id"`
	Agent            string `url:"agent"`
	Fields           string `url:"fields"`
	From             int    `url:"from"`
	Size             int    `url:"size"`
}

type SearchData struct {
	HasError bool        `json:"has_error"`
	Hits     []SearchHit `json:"hits"`
	Total    int         `json:"total"`
	Unavail  []string    `json:"unavail"`
}

type SearchHit struct {
	SYNOMDExtension               string              `json:"SYNOMDExtension,omitempty"`
	SYNOMDContentModificationDate string              `json:"SYNOMDContentModificationDate,omitempty"`
	SYNOMDFSCreationDate          string              `json:"SYNOMDFSCreationDate,omitempty"`
	SYNOMDLastUsedDate            string              `json:"SYNOMDLastUsedDate,omitempty"`
	SYNOMDFSName                  string              `json:"SYNOMDFSName,omitempty"`
	SYNOMDFSSize                  string              `json:"SYNOMDFSSize,omitempty"`
	SYNOMDIsDir                   string              `json:"SYNOMDIsDir,omitempty"`
	SYNOMDOwnerUserID             interface{}         `json:"SYNOMDOwnerUserID,omitempty"`
	SYNOMDPath                    string              `json:"SYNOMDPath,omitempty"`
	SYNOMDSharePath               string              `json:"SYNOMDSharePath,omitempty"`
	Additional                    SearchHitAdditional `json:"additional,omitempty"`
	Readable                      bool                `json:"readable,omitempty"`
}

type SearchHitAdditional struct {
	DocID        int     `json:"doc_id"`
	IsMountPoint bool    `json:"is_mount_point"`
	Keyword      string  `json:"keyword"`
	Score        float64 `json:"score"`
}
