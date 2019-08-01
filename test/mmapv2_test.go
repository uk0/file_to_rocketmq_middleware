package main

import (
	"fmt"
	"golang.org/x/exp/mmap"
	"testing"
	"time"
)

func Test_mmap2(test *testing.T) {
	var start,end time.Time;

	at, err := mmap.Open("/Users/zhangjianxin/home/GO_LIB/src/github.com/uk0/readCSV/test.csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	buff := make([]byte, at.Len())
	//读入的长度为slice预设的长度，0是offset。预设长度过长将会用0填充。
	at.ReadAt(buff, 0)
	//fmt.Println(string(buff))
	at.Close()
	fmt.Println("time ",start.Sub(end))
}