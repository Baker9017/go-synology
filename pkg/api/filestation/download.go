package filestation

type DownloadRequest struct {
	Path string `form:"path" url:"path"`
	Mode string `form:"mode" url:"mode"`
}

// DownloadFolderRequest 用于打包下载文件夹，将文件夹压缩为 zip 后下载
type DownloadFolderRequest struct {
	Path    []string `url:"path,json"` // 要下载的文件夹路径列表，如 ["/Temp/测试/12345"]
	DlName  string   `url:"dlname"`    // 下载的 zip 文件名，如 "12345.zip"
	Mode    string   `url:"mode"`      // 固定为 "download"
	Stdhtml string   `url:"stdhtml"`   // 固定为 "false"
}
