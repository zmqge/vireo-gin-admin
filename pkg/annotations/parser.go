package annotations

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

// RouteMeta 路由元数据
type RouteMeta struct {
	Method         string     // HTTP 方法
	Path           string     // 路由路径
	HandlerName    string     // 处理函数名称
	Middlewares    []string   // 中间件列表
	ControllerType string     // 控制器类型名
	ImportPath     string     // 控制器导入路径
	Permission     string     // 权限码
	Group          string     // 路由分组
	GroupMeta      *GroupMeta // 路由分组元数据
	FilePath       string     // 新增字段，记录控制器的文件路径
}

// PermissionAnnotation 权限注解
type PermissionAnnotation struct {
	Name        string   // 权限名称
	Code        string   // 权限码
	Description string   // 权限描述
	Group       string   // 权限分组
	Codes       []string // 多个权限码
}

// GroupMeta 分组元数据
type GroupMeta struct {
	Name        string // 分组名称
	Path        string // 分组路径
	Description string // 分组描述
}

// ParseRouteAndPermission 解析文件中的路由、权限和分组注解
func ParseRouteAndPermission(filePath string) ([]RouteMeta, []PermissionAnnotation, map[string]*GroupMeta, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, nil, err
	}

	var routes []RouteMeta
	var permissions []PermissionAnnotation
	groupMetaMap := make(map[string]*GroupMeta)

	// 获取控制器导入路径
	importPath := getImportPath(filePath)

	// 解析控制器类型声明上的 @Group 注解
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					controllerType := typeSpec.Name.Name
					// 检查 GenDecl.Doc
					if genDecl.Doc != nil {
						for _, comment := range genDecl.Doc.List {
							if strings.HasPrefix(comment.Text, "// @Group") {
								groupMetaMap[controllerType] = parseGroupAnnotation(comment.Text)
								break // 只解析第一个类型声明上的 @Group 注解
							}
						}
					}

					// 检查 TypeSpec.Doc
					if typeSpec.Doc != nil {
						for _, comment := range typeSpec.Doc.List {
							if strings.HasPrefix(comment.Text, "// @Group") {
								groupMetaMap[controllerType] = parseGroupAnnotation(comment.Text)
							}
						}
					}
				}
			}
		}
	}

	// 解析控制器方法上的 @Route 和 @Permission 注解
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Doc != nil {
			var route RouteMeta
			var permission PermissionAnnotation

			for _, comment := range fn.Doc.List {
				// 解析 @Route 注解
				if strings.HasPrefix(comment.Text, "// @Route") {
					route = parseRouteAnnotation(comment.Text)
					route.HandlerName = fn.Name.Name
					route.ControllerType = getControllerType(node)
					route.ImportPath = importPath
					route.GroupMeta = groupMetaMap[route.ControllerType] // 设置路由分组元数据
				}

				// 解析 @Permission 注解
				if strings.HasPrefix(comment.Text, "// @Permission") {
					permission = parsePermissionAnnotation(comment.Text)
					route.Permission = permission.Code
				}
			}

			if route.Method != "" && route.Path != "" {
				routes = append(routes, route)
			}
			if permission.Code != "" {
				permissions = append(permissions, permission)
			}
		}
	}

	return routes, permissions, groupMetaMap, nil
}

// ParseFile 解析文件中的权限注解
func ParseFile(filePath string) ([]PermissionAnnotation, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var permissions []PermissionAnnotation
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Doc != nil {
			for _, comment := range fn.Doc.List {
				if strings.HasPrefix(comment.Text, "// @Permission") {
					permission := parsePermissionAnnotation(comment.Text)
					permissions = append(permissions, permission)
				}
			}
		}
	}
	return permissions, nil
}

// parsePermissionAnnotation 解析权限注解
func parsePermissionAnnotation(comment string) PermissionAnnotation {
	re := regexp.MustCompile(`@Permission\(([^)]+)\)`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) < 2 {
		return PermissionAnnotation{}
	}

	params := strings.Split(matches[1], ",")
	permission := PermissionAnnotation{}
	for _, param := range params {
		kv := strings.SplitN(strings.TrimSpace(param), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(kv[1], `"`)
			switch key {
			case "name":
				permission.Name = value
			case "code":
				permission.Code = value
			case "description":
				permission.Description = value
			case "group":
				permission.Group = value
			}
		} else {
			permission.Codes = append(permission.Codes, strings.Trim(kv[0], `"`))
		}
	}
	return permission
}

// parseRouteAnnotation 解析 @Route 注解
func parseRouteAnnotation(comment string) RouteMeta {
	re := regexp.MustCompile(`@Route\(([^)]*)\)`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) < 2 {
		return RouteMeta{}
	}
	paramsStr := matches[1]
	// 支持逗号分隔但中括号内不分割
	var params []string
	bracket := 0
	start := 0
	for i, ch := range paramsStr {
		switch ch {
		case '[':
			bracket++
		case ']':
			bracket--
		case ',':
			if bracket == 0 {
				params = append(params, strings.TrimSpace(paramsStr[start:i]))
				start = i + 1
			}
		}
	}
	if start < len(paramsStr) {
		params = append(params, strings.TrimSpace(paramsStr[start:]))
	}

	route := RouteMeta{}
	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(kv[1], "\"'")
		switch key {
		case "method":
			route.Method = value
		case "path":
			route.Path = value
		case "middlewares":
			// 兼容各种写法 ["jwt",'dataPerm']、[ 'jwt' , "dataPerm" ]、["jwt" , 'dataPerm'] 等
			value = strings.Trim(value, "[] ")
			if value == "" {
				route.Middlewares = nil
			} else {
				// 用正则提取所有被单引号或双引号包裹的内容
				reMw := regexp.MustCompile(`['"][^'"]+['"]`)
				matches := reMw.FindAllString(value, -1)
				for _, m := range matches {
					mw := strings.Trim(m, "'\"")
					if mw != "" {
						route.Middlewares = append(route.Middlewares, mw)
					}
				}
			}
		case "group":
			route.Group = value
		}
	}
	return route
}

// getControllerType 从 AST 节点中提取控制器类型名
func getControllerType(node *ast.File) string {
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					return typeSpec.Name.Name
				}
			}
		}
	}
	return ""
}

// getImportPath 从文件路径中提取控制器的导入路径
func getImportPath(filePath string) string {
	// 假设项目根目录为 "github.com/zmqge/vireo-gin-admin"
	// 根据文件路径生成导入路径
	projectRoot := "github.com/zmqge/vireo-gin-admin"
	relativePath := strings.TrimPrefix(filePath, "app/")
	importPath := projectRoot + "/" + relativePath
	return strings.TrimSuffix(importPath, ".go")
}

// 解析分组注解
func parseGroupAnnotation(comment string) *GroupMeta {
	groupMeta := &GroupMeta{}
	comment = strings.TrimPrefix(comment, "// @Group(")
	comment = strings.TrimSuffix(comment, ")")
	pairs := strings.Split(comment, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(kv[1], `"`)
		switch key {
		case "name":
			groupMeta.Name = value
		case "path":
			groupMeta.Path = value
		case "description":
			groupMeta.Description = value
		}
	}
	return groupMeta
}
