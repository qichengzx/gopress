package xlib

import (
	"io/ioutil"
	"os"
)

func makeFile(c []byte, file string) {
	err := ioutil.WriteFile(file, c, 0644)
	if err != nil {
		panic(err)
	}
}

func clearDir(path string) {
	os.RemoveAll(path)
	os.Mkdir(path, 0777)
}
