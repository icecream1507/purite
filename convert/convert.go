package convert

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type articleLink struct {
	title string
	link  string
}

var (
	articleLinks []articleLink
	config       interface{}
)

// Convert 从md文件生成完整的html页面
func Convert() error {
	// 读取配置文件
	if err := getConfig(); err != nil {
		return err
	}
	// 读取博客目录下的文件列表
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	// 遍历文件列表并进行处理
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			// 将此目录下的md文件转为html并记入索引
			os.Chdir(fileName)
			pages, err := ioutil.ReadDir(".")
			if err != nil {
				return err
			}
			for _, page := range pages {
				pageName := page.Name()
				if !page.IsDir() && path.Ext(pageName) == ".md" {
					pageTitle, err := getTitle(pageName)
					if err != nil {
						return err
					}
					pageLink, err := getLink(pageName)
					if err != nil {
						return err
					}
					articleLinks = append(articleLinks, articleLink{pageTitle, pageLink})
					if err = toHTML(pageName); err != nil {
						return err
					}
				}
			}
			os.Chdir("..")
		}
	}

	if err := generateIndex(); err != nil {
		return err
	}
	if err = toHTML("index.md"); err != nil {
		return err
	}

	return nil
}

// getConfig 读取配置文件
func getConfig() error {
	buf, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, &config); err != nil {
		return err
	}

	return nil
}

// toHTML 将md文件转为html
func toHTML(fileName string) error {
	// 获取页面标题
	title, err := getTitle(fileName)
	if err != nil {
		return err
	}
	// 转为html
	unprocessedPage := strings.TrimSuffix(fileName, ".md") + "_old.html"
	cmd := exec.Command("pandoc", fileName, "-f", "markdown", "-t", "html", "-s", "-o", unprocessedPage)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	// 添加样式
	processedPage := strings.TrimSuffix(fileName, ".md") + ".html"
	buf, err := ioutil.ReadFile(unprocessedPage)
	if err != nil {
		return err
	}
	newStr := strings.Replace(string(buf), "<head>", config.(map[string]interface{})["HTMLHead"].(string), 1)
	newStr = strings.Replace(newStr, "<title>"+strings.TrimSuffix(fileName, ".md")+"</title>", "<title>"+title+"</title>", 1)
	err = ioutil.WriteFile(processedPage, []byte(newStr), 0644)
	if err != nil {
		return err
	}

	// 删除旧页面
	if err = os.Remove(unprocessedPage); err != nil {
		return err
	}

	return nil
}

// getTitle 获取文章标题
func getTitle(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	firstline, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(firstline, "# "), nil
}

// getLink 获取文章相对链接
func getLink(fileName string) (string, error) {
	absPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absPath = strings.ReplaceAll(absPath, "\\", "/")
	absPath = strings.ReplaceAll(absPath, "\\\\", "/")
	slice := strings.Split(absPath, "/")
	currentFolder := slice[len(slice)-1]
	return "./" + currentFolder + "/" + strings.TrimSuffix(fileName, ".md") + ".html", nil
}

// generateIndex 生成索引
func generateIndex() error {
	// 生成索引列表字符串
	var indexStr string
	for _, item := range articleLinks {
		indexStr += "- [" + item.title + "](" + item.link + ")\n"
	}

	// 对模版进行填充
	buf, err := ioutil.ReadFile("index-template.md")
	if err != nil {
		return err
	}
	newStr := strings.ReplaceAll(string(buf), "%purite:article_list%", indexStr)
	err = ioutil.WriteFile("index.md", []byte(newStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
