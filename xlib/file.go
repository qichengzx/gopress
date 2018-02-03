package xlib

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func makeFile(c []byte, file string) {
	dir := filepath.Dir(file)
	os.MkdirAll(dir, 0777)

	err := ioutil.WriteFile(file, c, 0644)
	if err != nil {
		panic(err)
	}
}

func clearDir(path string) {
	os.RemoveAll(path)
	os.Mkdir(path, 0777)
}
