package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestReadMMAP(test *testing.T) {
	var err error
	var fn string
	var fd int
	var fi os.FileInfo
	var data []byte
	//打开文件
	fn = "/Users/zhangjianxin/home/GO_LIB/src/github.com/uk0/readCSV/test.csv"
	fd, err = syscall.Open(fn, syscall.O_RDWR, 0)
	if nil != err {
		fmt.Println("open file fail!!!")
	}
	//获取文件大小
	fi, err = os.Stat(fn)
	if nil != err {
		fmt.Println("get file size fail!!!")
	}
	//映射到内存
	data, err = syscall.Mmap(fd, 0, int(fi.Size()), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_SHARED)

	//
	if nil != err {
		fmt.Println("mmap fail!!!")
		return
	}
	allData := string(data)
	Lines :=strings.Split(allData,"\n")

	fmt.Println(len(Lines))
	addr := &data[0]
	fmt.Println("mmap success,addr=", addr, "size=", len(data))
	//取消映射
	syscall.Munmap(data)
}