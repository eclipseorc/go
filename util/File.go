package util

/******************************************************************************
Copyright:cloud
Author:cloudapex@126.com
Version:1.0
Date:2014-10-18
Description:utd文件接口
https://studygolang.com/articles/3154
******************************************************************************/
import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FileTimeMod(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.ModTime().Unix(), nil
}

func FileSizefi(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

func FileRename(file string, to string) error {
	return os.Rename(file, to)
}

func FileRemove(file string) error {
	return os.Remove(file)
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func FileCutTo(srcFile, desFile string) error {
	if err := FileCopyTo(srcFile, desFile); err != nil {
		return err
	}
	return os.Remove(srcFile)
}

func FileCopyTo(srcFile, desFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	des, err := os.Create(desFile)
	if err != nil {
		return err
	}
	defer des.Close()

	_, err = io.Copy(des, src)
	return err
}

func DirList(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	pthSep := string(os.PathSeparator) //分隔符
	suffix = strings.ToUpper(suffix)   //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if suffix == "*" || suffix == "" {
			files = append(files, dirPth+pthSep+fi.Name())
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+pthSep+fi.Name())
		}
	}
	return files, nil
}

func DirWalk(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if suffix == "*" || suffix == "" {
			files = append(files, filename)
		} else if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
