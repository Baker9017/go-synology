package filestation

// CheckPermissionRequest 用于检查用户对目录的写权限（上传权限）
type CheckPermissionRequest struct {
	Path      string `form:"path" url:"path"`           // 目录路径，如 /Temp/12
	Filename  string `form:"filename" url:"filename"`   // 要上传的文件名
	Size      int64  `form:"size" url:"size"`           // 文件大小（字节）
	Overwrite bool   `form:"overwrite" url:"overwrite"` // 是否覆盖已存在文件
}

// CheckPermissionResponse SYNO.FileStation.CheckPermission 接口返回
// blSkip: false 表示有写权限，true 表示无写权限需跳过操作
type CheckPermissionResponse struct {
	BlSkip bool `json:"blSkip"`
}
