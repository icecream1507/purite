package build

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Config 站点配置数据结构体
type Config struct {
	GitRepo  string
	HTMLHead string
}

var (
	gitRepo, htmlHead string
)

// Build 初始化
func Build() error {
	if err := getInfo(); err != nil {
		return err
	}
	htmlHead := "<head>\n<link rel=\"stylesheet\" href=\"https://f-y-blog.oss-cn-shenzhen.aliyuncs.com/style.css\" />"
	config := Config{gitRepo, htmlHead}

	if err := generateJSON(config); err != nil {
		return err
	}

	if err := gitInit(); err != nil {
		return err
	}

	return nil
}

// getInfo 读取用户输入
func getInfo() error {
	fmt.Printf("请输入您的站点Git仓库地址并确保其可以访问：")
	if _, err := fmt.Scanln(&gitRepo); err != nil {
		return err
	}

	return nil
}

// generateJSON 生成JSON配置文件
func generateJSON(config Config) error {
	file, err := os.OpenFile("config.json", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}
	return nil
}

func gitInit() error {
	cmd1 := exec.Command("git", "init")
	cmd2 := exec.Command("git", "add", ".")
	cmd3 := exec.Command("git", "commit", "-m", "init")
	cmd4 := exec.Command("git", "remote", "add", "origin", gitRepo)
	cmd5 := exec.Command("git", "push", "-u", "origin", "master")

	if err := cmd1.Run(); err != nil {
		return err
	}
	if err := cmd2.Run(); err != nil {
		return err
	}
	if err := cmd3.Run(); err != nil {
		return err
	}
	if err := cmd4.Run(); err != nil {
		return err
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd5.Stdout = &out
	cmd5.Stderr = &stderr
	if err := cmd5.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
