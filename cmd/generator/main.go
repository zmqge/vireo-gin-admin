package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// 获取项目根目录绝对路径
func getProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	for dir := cwd; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Abs(dir)
		}
	}

	return "", fmt.Errorf("无法定位项目根目录(未找到go.mod)")
}

// 从模板文件加载模板
func loadTemplate(filePath string) (*template.Template, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取模板文件 %s 失败: %v", filePath, err)
	}
	return template.New(filepath.Base(filePath)).Parse(string(content))
}

// 确保首字母大写
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func main() {
	entity := flag.String("entity", "", "实体名称（如User）")
	module := flag.String("module", "", "模块名称（留空则自动检测）")
	outputPath := flag.String("path", "", "输出路径（如app/admin），相对于项目根目录")
	flag.Parse()

	// 处理实体名称，确保首字母大写
	*entity = capitalizeFirstLetter(*entity)

	if *entity == "" {
		fmt.Println("请指定实体名称")
		return
	}

	// 获取项目根目录
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("获取项目根目录失败: %v\n", err)
		return
	}

	// 自动检测模块路径
	if *module == "" {
		cmd := exec.Command("go", "list", "-m")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("自动检测模块路径失败: %v\n", err)
			fmt.Println("请手动指定 -module 参数")
			return
		}
		*module = strings.TrimSpace(string(output))
		if *module == "" {
			fmt.Println("无法获取模块路径，请确保项目已初始化 (go mod init)")
			return
		}
		fmt.Printf("自动检测模块路径: %s\n", *module)
	}

	// 规范化输出路径（基于项目根目录）
	normalizedPath := filepath.Join(projectRoot, *outputPath)
	if *outputPath == "" {
		normalizedPath = projectRoot // 默认输出到项目根目录
	}

	// 创建基础输出目录
	if err := os.MkdirAll(normalizedPath, 0755); err != nil {
		fmt.Printf("创建输出目录 %s 失败: %v\n", normalizedPath, err)
		return
	}

	// 创建MVC目录结构
	dirs := []string{
		"models", "controllers", "services", "repositories",
	}
	for _, dir := range dirs {
		dirPath := filepath.Join(normalizedPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Printf("创建目录 %s 失败: %v\n", dirPath, err)
			return
		}
	}

	// 计算相对路径（用于导入语句）
	relPath, err := filepath.Rel(projectRoot, normalizedPath)
	if err != nil {
		relPath = normalizedPath
	}
	relPath = strings.ReplaceAll(relPath, "\\", "/")

	// 生成路由路径（复数形式，小写）
	// entityPath := strings.ToLower(pluralize(*entity))
	entityPath := strings.ToLower(*entity)

	// 生成权限代码前缀（小写）
	entityPermission := strings.ToLower(*entity)

	// 生成文件
	data := struct {
		Entity           string
		Module           string
		RelPath          string
		EntityPath       string
		EntityPermission string
	}{
		Entity:           *entity,
		Module:           *module,
		RelPath:          relPath,
		EntityPath:       entityPath,
		EntityPermission: entityPermission,
	}

	// 模型文件名添加Model后缀
	controllerFileName := toCamelCase(*entity) + "Controller.go"
	modelFileName := toCamelCase(*entity) + "Model.go"
	serviceFileName := toCamelCase(*entity) + "Service.go"
	repositoryFileName := toCamelCase(*entity) + "Repository.go"

	// 定义模板文件路径（相对于项目根目录）
	templateFiles := []struct {
		outputRelativePath string // 输出文件相对路径（相对于normalizedPath）
		templatePath       string // 模板文件相对路径（相对于项目根目录）
	}{
		{filepath.Join("models", modelFileName), "cmd/generator/templates/model.tmpl"},
		{filepath.Join("controllers", controllerFileName), "cmd/generator/templates/controller.tmpl"},
		{filepath.Join("services", serviceFileName), "cmd/generator/templates/service.tmpl"},
		{filepath.Join("repositories", repositoryFileName), "cmd/generator/templates/repository.tmpl"},
	}

	for _, tf := range templateFiles {
		// 输出文件完整路径（基于项目根目录）
		outputFile := filepath.Join(normalizedPath, tf.outputRelativePath)
		// 模板文件完整路径（基于项目根目录）
		templateFile := filepath.Join(projectRoot, tf.templatePath)

		// 创建父目录（如果不存在）
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			fmt.Printf("创建目录 %s 失败: %v\n", filepath.Dir(outputFile), err)
			return
		}

		// 创建文件
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("创建文件 %s 失败: %v\n", outputFile, err)
			return
		}
		defer f.Close()

		// 加载并执行模板
		tmpl, err := loadTemplate(templateFile)
		if err != nil {
			fmt.Printf("加载模板 %s 失败: %v\n", templateFile, err)
			return
		}

		if err := tmpl.Execute(f, data); err != nil {
			fmt.Printf("生成文件 %s 失败: %v\n", outputFile, err)
			return
		}

		// 显示相对路径（相对于项目根目录）
		relOutputPath, _ := filepath.Rel(projectRoot, outputFile)
		fmt.Printf("生成文件: %s\n", relOutputPath)
	}

	fmt.Println("MVC文件生成完成!")
}

// pluralize 将单数名词转为复数（简化版）
func pluralize(word string) string {
	word = strings.TrimSpace(word)
	if len(word) == 0 {
		return ""
	}

	lastChar := strings.ToLower(word[len(word)-1:])
	lastTwoChars := strings.ToLower(word[len(word)-2:])

	switch {
	case lastChar == "s" || lastChar == "x" || lastChar == "z" || lastTwoChars == "ch" || lastTwoChars == "sh":
		return word + "es"
	case lastChar == "y" && len(word) > 1 && !isVowel(word[len(word)-2]):
		return word[:len(word)-1] + "ies"
	default:
		return word + "s"
	}
}

// isVowel 判断字符是否为元音字母
func isVowel(c byte) bool {
	c = strings.ToLower(string(c))[0]
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}

// toCamelCase 将字符串转为小驼峰（首字母小写）
func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// 首字母小写
	return strings.ToLower(s[:1]) + s[1:]
}
