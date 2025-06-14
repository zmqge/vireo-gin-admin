package annotations

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Permission 表示权限的结构体
type Permission struct {
	Name        string // 权限名称
	Code        string // 权限代码
	Description string // 权限描述
}

// ScanDirectories 扫描多个目录中的权限注解
func ScanDirectories() ([]Permission, error) {
	// 从配置文件中获取需要扫描的目录路径
	controllerDirs := viper.GetStringSlice("controller_dirs")

	// 扫描目录并返回结果
	var permissions []Permission
	for _, dir := range controllerDirs {
		files, err := scanDirectory(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			// 解析文件中的权限注解
			perms, err := parsePermissions(file)
			if err != nil {
				return nil, err
			}
			permissions = append(permissions, perms...)
		}
	}

	return permissions, nil
}

// scanDirectory 扫描单个目录中的权限注解
func scanDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// parsePermissions 解析文件中的权限注解
func parsePermissions(file string) ([]Permission, error) {
	// 实现解析逻辑，返回权限列表
	return []Permission{}, nil
}
