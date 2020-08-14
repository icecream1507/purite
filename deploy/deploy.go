package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

var config interface{}

// Deploy 部署至Git仓库
func Deploy() error {
	var message string
	fmt.Print("请输入本次提交的备注：")
	if _, err := fmt.Scanln(&message); err != nil {
		return err
	}

	cmd1 := exec.Command("git", "add", ".")
	cmd2 := exec.Command("git", "commit", "-m", message)
	cmd3 := exec.Command("git", "push", "origin", "master")

	if err := cmd1.Run(); err != nil {
		return err
	}
	if err := cmd2.Run(); err != nil {
		return err
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd3.Stdout = &out
	cmd3.Stderr = &stderr
	if err := cmd3.Run(); err != nil {
		return errors.New(stderr.String())
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
