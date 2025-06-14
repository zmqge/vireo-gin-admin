package annotations

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

type PermissionMeta struct {
	Code        string // 权限码
	Name        string // 权限名称
	Description string // 权限描述
	Module      string // 权限模块
}

// ParsePermissionAnnotations 解析文件中的权限注解，返回所有权限元信息
func ParsePermissionAnnotations(filePath string) ([]PermissionMeta, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var permissions []PermissionMeta
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Doc != nil {
			for _, comment := range fn.Doc.List {
				if strings.HasPrefix(comment.Text, "// @Permission") {
					perm := parsePermissionMeta(comment.Text)
					if perm.Code != "" {
						permissions = append(permissions, perm)
					}
				}
			}
		}
	}
	return permissions, nil
}

// parsePermissionMeta 解析单条权限注解
func parsePermissionMeta(comment string) PermissionMeta {
	re := regexp.MustCompile(`@Permission\(([^)]+)\)`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) < 2 {
		return PermissionMeta{}
	}
	params := strings.Split(matches[1], ",")
	perm := PermissionMeta{}
	for _, param := range params {
		kv := strings.SplitN(strings.TrimSpace(param), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(kv[1], `"`)
			switch key {
			case "code":
				perm.Code = value
			case "name":
				perm.Name = value
			case "desc":
				perm.Description = value
			case "description":
				perm.Description = value
			case "modules":
				perm.Module = value
			}
		}
	}
	return perm
}
