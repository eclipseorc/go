package lib

import (
	"io/ioutil"
	"os"
)

func ReadConfig(name string) ([]byte, error) {
	jsonFile, err := os.Open(name)
	if err != nil {
		panic("打开文件错误，请查看：" + name)
	}
	defer jsonFile.Close()
	return ioutil.ReadAll(jsonFile)
}
