package main

import (
	"fmt"
	"os"
	"purite/build"
	"purite/convert"
	"purite/deploy"
)

func main() {
	if _, err := os.Stat("config.json"); err != nil && os.IsNotExist(err) {
		if err = build.Build(); err != nil {
			fmt.Println("初始化站点过程中出现了问题。")
			fmt.Println(err)
			return
		}
	}
	if err := convert.Convert(); err != nil {
		fmt.Println("HTML页面生成过程中出现了问题。")
		fmt.Println(err)
		return
	}
	if err := deploy.Deploy(); err != nil {
		fmt.Println("部署到Git仓库过程中出现了问题。")
		fmt.Println(err)
		return
	}
	fmt.Println("您的站点部署/更新完成，可以打开对应地址查看了。")
}
