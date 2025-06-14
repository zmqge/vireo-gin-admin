package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"github.com/zmqge/vireo-gin-admin/config"
	perm_parser "github.com/zmqge/vireo-gin-admin/pkg/annotations"
	"github.com/zmqge/vireo-gin-admin/pkg/database"
)

func main() {
	// 初始化配置和数据库
	config.Init()
	db := database.InitDB()
	defer database.Close()

	// 获取 controller 路径
	controllerDirs := config.App.ControllerDirs
	if len(controllerDirs) == 0 {
		controllerDirs = []string{"app/admin/controllers"}
	}

	// 收集所有 go 文件
	var files []string
	for _, dir := range controllerDirs {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		})
	}

	// 清空 permissions 表
	db.Exec("TRUNCATE TABLE permissions")

	// 扫描所有 controller 文件，提取权限注解
	codeFilesMap := make(map[string][]string) // code -> []file
	codePermMap := make(map[string]models.Permission)
	var duplicateCodes []string

	for _, file := range files {
		perms, err := perm_parser.ParsePermissionAnnotations(file)
		if err != nil {
			fmt.Printf("Parse %s error: %v\n", file, err)
			continue
		}
		for _, p := range perms {
			if p.Code == "" {
				continue
			}
			codeFilesMap[p.Code] = append(codeFilesMap[p.Code], file)
			// 只保留第一个注解内容（合并重复项时使用第一个项的信息）
			if _, exists := codePermMap[p.Code]; !exists {
				codePermMap[p.Code] = models.Permission{
					Code:        p.Code,
					Name:        p.Name,
					Description: p.Description,
					Module:      p.Module,
				}
			}
		}
	}

	// 检查重复 code 并输出所有位置
	for code, files := range codeFilesMap {
		if len(files) > 1 {
			duplicateCodes = append(duplicateCodes, code)
			fmt.Printf("权限 code 重复: %s\n", code)
			for _, f := range files {
				fmt.Printf("位置: %s\n", f)
			}
		}
	}
	if len(duplicateCodes) > 0 {
		fmt.Println("存在重复 code，已合并，仅保留一条写入数据库。请检查上方提示！")
	}
	// 合并后写入数据库
	var permissions []models.Permission
	for _, perm := range codePermMap {
		permissions = append(permissions, perm)
	}
	// 批量写入数据库，只写入 code、name、description 字段
	if len(permissions) > 0 {
		db.Model(&models.Permission{}).Select("code", "name", "description", "module").Create(&permissions)
		fmt.Printf("写入权限 %d 条\n", len(permissions))
	} else {
		fmt.Println("未发现任何权限注解")
	}
}
