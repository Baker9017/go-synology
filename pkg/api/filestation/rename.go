package filestation

type RenameRequest struct {
	Path []string `url:"path,json"` // 文件完整路径（含文件名），如 ["/Temp/12/old.pdf"]
	Name []string `url:"name,json"` // 新文件名（不含路径），如 ["new.pdf"]
}
