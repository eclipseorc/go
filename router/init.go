package router

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type RunParam struct {
	Port string `json:"run_port"`
	File string `json:"file_path"`
}

var Run RunParam

func InitRunParam() {
	path := "./config/run.json"
	jsonFile, err := os.Open(path)
	if err != nil {
		panic("初始化端口文件错误，请查看：" + path)
	}
	defer jsonFile.Close()
	data, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(data, &Run)
}
