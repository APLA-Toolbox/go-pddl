package common

import (
	"context"
	"fmt"
	"io/ioutil"
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

type Domain struct {

}

type Problem struct {
	
}
