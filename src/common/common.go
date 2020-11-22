package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type StringSlice []string

func (ss StringSlice) ToString(name string) string {
	s := name + ": ["
	for _, v := range ss {
		s += "'" + v + "'"
	}
	s += "]\n"
	return s
}

type StringMap map[string]string

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		panic("Exit: Failure")
	}
}

func LoadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetRootPath() string {
	rootPath, exists := os.LookupEnv("ROOT_PATH")
	if !exists {
		fmt.Fprintf(os.Stderr, "ROOT_PATH environment variable is not defined")
		os.Exit(1)
	}
	return rootPath
}

func JoinRootPathIfNotAbsolute(in string) string {
	if path.IsAbs(in) {
		return in
	}
	return path.Join(GetRootPath(), in)
}
