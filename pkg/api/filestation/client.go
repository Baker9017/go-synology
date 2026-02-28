package filestation

import (
	"context"
	"fmt"
	pathpkg "path"
	"slices"
	"time"

	"github.com/synology-community/go-synology/pkg/api"
	"github.com/synology-community/go-synology/pkg/api/filestation/methods"
	"github.com/synology-community/go-synology/pkg/models"
	"github.com/synology-community/go-synology/pkg/util/form"
)

type Client struct {
	client api.Api
}

func New(client api.Api) Api {
	return &Client{client: client}
}

type FileNotFoundError struct {
	Path string
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("File not found: %s", e.Path)
}

// List implements FileStationApi.
func (f *Client) List(ctx context.Context, folderPath string) (*models.FileList, error) {
	return api.Get[models.FileList](f.client, ctx, &models.FileListRequest{
		FolderPath: folderPath,
		Additional: []string{
			"real_path",
			"size",
			"owner",
			"time",
			"perm",
			"mount_point_type",
			"type",
			"fileid",
		},
		FileType: "all",
	}, methods.List)
}

func (f *Client) Get(ctx context.Context, path string) (*models.File, error) {
	folder := pathpkg.Dir(path)
	resp, err := f.List(ctx, folder)
	if err != nil {
		return nil, fmt.Errorf("unable to get file, got error: %s", err)
	}
	if resp.Files == nil {
		return nil, fmt.Errorf("files is nil")
	}
	if len(resp.Files) == 0 {
		return nil, fmt.Errorf("result is empty")
	}
	i := slices.IndexFunc(resp.Files, func(f models.File) bool {
		return f.Path == path
	})
	if i == -1 {
		return nil, FileNotFoundError{Path: path}
	}
	return &resp.Files[i], nil
}

func (f *Client) Delete(
	ctx context.Context,
	paths []string,
	accurateProgress bool,
) (*DeleteStatusResponse, error) {
	// Start Delete the file
	rdel, err := f.DeleteStart(ctx, paths, true)
	if err != nil {
		return nil, fmt.Errorf("unable to delete file, got error: %s", err)
	}
	return f.DeleteStatus(ctx, rdel.TaskID)
}

func (f *Client) DeleteStart(
	ctx context.Context,
	paths []string,
	accurateProgress bool,
) (*DeleteStartResponse, error) {
	method := methods.DeleteStart
	return api.Get[DeleteStartResponse](f.client, ctx, &DeleteStartRequest{
		Paths:            paths,
		AccurateProgress: accurateProgress,
	}, method)
}

func (f *Client) DeleteStatus(ctx context.Context, taskID string) (*DeleteStatusResponse, error) {
	return api.Get[DeleteStatusResponse](f.client, ctx, &DeleteStatusRequest{
		TaskID: taskID,
	}, methods.DeleteStatus)
}

func (f *Client) MD5(ctx context.Context, path string) (*MD5StatusResponse, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(60 * time.Second)
	}

	rmd5, err := f.MD5Start(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to get file md5, got error: %s", err)
	}

	delay := 1 * time.Second
	for {
		task, err := f.MD5Status(ctx, rmd5.TaskID)
		if err != nil && task == nil {
			return nil, err
		}
		if task.Finished {
			return task, nil
		}
		if time.Now().After(deadline.Add(delay)) {
			return nil, fmt.Errorf("timeout waiting for task to complete")
		}
		time.Sleep(delay)
	}
}

func (f *Client) MD5Start(ctx context.Context, path string) (*MD5StartResponse, error) {
	return api.Get[MD5StartResponse](f.client, ctx, &MD5StartRequest{
		Path: path,
	}, methods.MD5Start)
}

func (f *Client) MD5Status(ctx context.Context, taskID string) (*MD5StatusResponse, error) {
	return api.Get[MD5StatusResponse](f.client, ctx, &MD5StatusRequest{
		TaskID: taskID,
	}, methods.MD5Status)
}

// Download implements FileStationApi.
func (f *Client) Download(ctx context.Context, path string, mode string) (*form.File, error) {
	return api.Get[form.File](f.client, ctx, &DownloadRequest{
		Path: path,
		Mode: mode,
	}, methods.Download)
}

// DownloadFolder 打包下载文件夹，将指定路径的文件夹压缩为 zip 后返回
// dlName 可选，不传则取 paths[0] 最后一层文件夹名加 .zip
func (f *Client) DownloadFolder(ctx context.Context, paths []string, dlName ...string) (*form.File, error) {
	name := ""
	if len(dlName) > 0 && dlName[0] != "" {
		name = dlName[0]
	} else if len(paths) > 0 {
		name = pathpkg.Base(paths[0]) + ".zip"
	}
	file, err := api.Get[form.File](f.client, ctx, &DownloadFolderRequest{
		Path:    paths,
		DlName:  name,
		Mode:    "download",
		Stdhtml: "false",
	}, methods.Download)
	if err != nil {
		return nil, err
	}
	if file != nil {
		file.Name = name
	}
	return file, nil
}

// Rename implements FileStationApi.
// path 为父目录路径，name 为当前文件名，newName 为新的文件名。
func (f *Client) Rename(
	ctx context.Context,
	path string,
	name string,
	newName string,
) (*models.FileList, error) {
	fullPath := pathpkg.Join(path, name)
	return api.Get[models.FileList](f.client, ctx, &RenameRequest{
		Path: []string{fullPath},
		Name: []string{newName},
	}, methods.Rename)
}

// CreateFolder implements FileStationApi.
func (f *Client) CreateFolder(
	ctx context.Context,
	paths []string,
	names []string,
	forceParent bool,
) (*models.FolderList, error) {
	return api.Get[models.FolderList](f.client, ctx, &CreateFolderRequest{
		Paths:       paths,
		Names:       names,
		ForceParent: forceParent,
	}, methods.CreateFolder)
}

// ListShares implements FileStationApi.
func (f *Client) ListShares(ctx context.Context) (*models.ShareList, error) {
	return api.Get[models.ShareList](f.client, ctx, &ListShareRequest{}, methods.ListShares)
}

// CheckWritePermission 检查用户对指定目录是否有写权限（上传权限）
// path 为目录路径，如 /Temp/12；filename 为要上传的文件名；size 为文件大小（字节）
// 返回 true 表示有写权限，false 表示无写权限
func (f *Client) CheckWritePermission(ctx context.Context, path, filename string, size int64, overwrite bool) (bool, error) {
	resp, err := api.Get[CheckPermissionResponse](f.client, ctx, &CheckPermissionRequest{
		Path:      path,
		Filename:  filename,
		Size:      size,
		Overwrite: overwrite,
	}, methods.CheckPermission)
	if err != nil {
		return false, err
	}
	// blSkip: false 表示有权限，true 表示无权限需跳过
	return !resp.BlSkip, nil
}

// Upload implements FileStationApi.
// mtime 可选，指定上传后文件的修改时间（毫秒时间戳），传 nil 则不设置
func (f *Client) Upload(
	ctx context.Context,
	path string,
	file form.File,
	createParents bool,
	overwrite bool,
	mtime *time.Time,
) (*UploadResponse, error) {
	req := &UploadRequest{
		Path:          path,
		File:          file,
		CreateParents: createParents,
		Overwrite:     overwrite,
	}
	if mtime != nil {
		ms := mtime.UnixMilli()
		req.Mtime = &ms
	}
	return api.PostFile[UploadResponse](f.client, ctx, req, methods.Upload)
}
