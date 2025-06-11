package docs

import (
	"embed"
)

//go:embed swagger.json
var swagger embed.FS

// GetSwaggerJSON 返回swagger.json文件内容
func GetSwaggerJSON() ([]byte, error) {
	return swagger.ReadFile("swagger.json")
}
