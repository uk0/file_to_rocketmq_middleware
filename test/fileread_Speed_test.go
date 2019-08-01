package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

func read1(path string,blocksize int){
	fi,err := os.Open(path)
	if err != nil{
		panic(err)
	}
	defer fi.Close()
	block := make([]byte,blocksize)
	for{
		n,err := fi.Read(block)
		if err != nil && err != io.EOF{panic(err)}
		if 0 ==n {break}
	}
}

func read2(path string,blocksize int){
	fi,err := os.Open(path)
	if err != nil{panic(err)}
	defer fi.Close()
	r := bufio.NewReader(fi)
	block := make([]byte,blocksize)
	for{
		n,err := r.Read(block)
		if err != nil && err != io.EOF{panic(err)}

		if 0 ==n {break}
	}
}

func read3(path string){
	fi,err := os.Open(path)
	if err != nil{panic(err)}
	defer fi.Close()
	_,err = ioutil.ReadAll(fi)
}

func read4(path string){
	_,err := ioutil.ReadFile(path)
	if err != nil{panic(err)}
}
func read5(path string){
	fi, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		fmt.Println(string(a))
	}
}

func TestRead( test *testing.T){
	flag.Parse()
	file1 := "./test.csv"
	file2 := "./test.csv"
	file3 := "./test.csv"
	file4 := "./test.csv"
	blocksize,_ :=strconv.Atoi(flag.Arg(0))
	var start,end time.Time
	start = time.Now()
	read1(file1,blocksize)
	end = time.Now()
	fmt.Printf("file/Read() cost time %v\n",end.Sub(start))
	start = time.Now()
	read2(file2,blocksize)
	end = time.Now()
	fmt.Printf("bufio/Read() cost time %v\n",end.Sub(start))
	start = time.Now()
	read3(file3)
	end = time.Now()
	fmt.Printf("ioutil.ReadAll() cost time %v\n",end.Sub(start))
	start = time.Now()
	read4(file4)
	end = time.Now()
	fmt.Printf("ioutil.ReadFile() cost time %v\n",end.Sub(start))
	read5(file4)
	end = time.Now()
	fmt.Printf("ioutil.ReadFile(ReadLine) cost time %v\n",end.Sub(start))

}