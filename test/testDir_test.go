package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test(test *testing.T)  {
	folder:="/Users/zhangjianxin/home/GO_LIB/src/github.com/uk0/readCSV/test"
	ListDir(folder)
	FloatMath()
}

func ListDir(folder string)  {
	files, _ := ioutil.ReadDir(folder) //specify the current dir
	for _,file := range files{
		if file.IsDir(){
			ListDir(folder + "/" + file.Name())
		}else{
			fmt.Println(folder + "/" + file.Name())
		}
	}
}

func FloatMath(){
	fmt.Println(fmt.Sprintf("%d", 1000))
}