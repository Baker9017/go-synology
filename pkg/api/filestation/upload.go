package filestation

import (
	"net/url"

	"github.com/synology-community/go-synology/pkg/api"
	"github.com/synology-community/go-synology/pkg/util/form"
)

type UploadRequest struct {
	Path          string    `form:"path"           url:"path"`
	CreateParents bool      `form:"create_parents" url:"create_parents"`
	Overwrite     bool      `form:"overwrite"      url:"overwrite"`
	Mtime         *int64    `form:"mtime"          url:"mtime"` // 修改时间，毫秒时间戳，可选
	File          form.File `form:"file"                                kind:"file"`
}

func (l UploadRequest) EncodeValues(_ string, _ *url.Values) error {
	return nil
}

type UploadResponse struct {
	BlSkip   bool   `json:"blSkip"`
	File     string `json:"file"`
	Pid      int    `json:"pid"`
	Progress int    `json:"progress"`
}

var _ api.Request = (*UploadRequest)(nil)
