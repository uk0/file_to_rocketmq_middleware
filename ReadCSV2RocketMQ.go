package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

type Message struct {
	line         string
	topic        string
	FileName     string
	TotalCount   int
	ProcessCount int
}

//////////////////base/////////////////////
var data = make(chan Message, 100000)
var waitgroup sync.WaitGroup

const STATUS_FILE_SUCC = "_load_mem_success_"

//////////////////base/////////////////////

// 内存映射文件
func ReadFileUesMMAP(conf FtpSender, path string, dataChan chan Message) {
	var err error
	var fn string
	var fd int
	var fi os.FileInfo
	var data []byte
	//打开文件
	fn = path
	fd, err = syscall.Open(fn, syscall.O_RDWR, 0)
	if nil != err {
		logs.Debug("open file fail!!!", err)
	}
	//获取文件大小
	fi, err = os.Stat(fn)
	if nil != err {
		logs.Debug("get file size fail!!!", err)
	}
	//映射到内存
	data, err = syscall.Mmap(fd, 0, int(fi.Size()), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_SHARED)

	if nil != err {
		logs.Debug("mmap fail!!!", err)
		return
	}

	allData := string(data)
	Lines := strings.Split(allData, "\n")
	LineCount := len(Lines)

	var count = 0;
	r := bytes.NewReader(data)
	br := bufio.NewReader(r)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		//内存文件映射 提升文件读取速度 顺便处理数据
		var dataLine = string([]rune(string(a)));
		SuffSymbolLine := DeleteSymbol(dataLine)
		if count%1000 == 0 {
			logs.Debug(fmt.Sprintf("DeleteSymbolDebug [1] Data %s", string(a)))
			logs.Debug(fmt.Sprintf("DeleteSymbolDebug [2] Data %s", SuffSymbolLine))
		}
		count++;
		// 包含数据 所属Topic 处理过后的行数据以及 FileName 总行数 已经处理行数
		dataChan <- Message{topic: conf.SendTopic, line: SuffSymbolLine, FileName: path, TotalCount: LineCount, ProcessCount: count}
	}

	addr := &data[0]

	logs.Debug(fmt.Sprintf("MMAP success,addr= %b Size=%d", addr, len(data)))
	// 获取文件名修改文件名
	FileReadSuccessModifyNames(path, conf)
	// 加载完成删除数据
	//FileDelete(path)
	//取消映射
	syscall.Munmap(data)
	waitgroup.Done()
}

func GetFileName(path string) string {
	arrayStr := strings.Split(path, "/")
	return arrayStr[len(arrayStr)-1]
}

func FileReadSuccessModifyNames(path string, conf FtpSender) bool {
	// 随机10个字符串
	err := os.Rename(path, path+STATUS_FILE_SUCC+RandStringBytesMaskImprSrc(10)) //重命名
	if err != nil {
		logs.Debug("file rename Error!")
		fmt.Printf("%s", err)
		return false
	} else {
		logs.Debug("file rename OK!")
		return true
	}
	return false
}

func FileDelete(file string) bool {
	if file != "/" {
		return false
	}
	if strings.Contains(file, "root") {
		return false
	}

	err := os.Remove(file) //删除文件test.txt
	if err != nil {
		//如果删除失败则输出 file remove Error!
		logs.Debug("file remove Error!")
		//输出错误详细信息
		logs.Debug(err)
		return false
	} else {
		//如果删除成功则输出 file remove OK!
		logs.Debug("file remove OK!")
		return true
	}
	return false
}

func WatchFiles(conf FtpSender) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					logs.Debug("Create file result: %s \n", event.Name)
					// 检查到新文件 创建协程
					if strings.HasSuffix(event.Name, ".csv") {
						waitgroup.Add(1) //检查到一个文件启动一个携程
						//go ReadFiles(conf,event.Name, data)
						go ReadFileUesMMAP(conf, event.Name, data)
					}
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	var Paths = strings.Split(conf.WatchDir, ",")
	for index := range Paths {
		err = watcher.Add(Paths[index])
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
}

func getLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "trace":
		return logs.LevelTrace
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "error":
		return logs.LevelError
	default:
		return logs.LevelDebug
	}
}

func InitLog() (err error) {
	//初始化日志库
	config := make(map[string]interface{})
	config["filename"] = "./running.log"
	config["level"] = getLevel("debug")
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println(" json.Marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func ListDir(folder string, conf FtpSender) {
	files, _ := ioutil.ReadDir(folder) //specify the current dir
	for _, file := range files {
		if file.IsDir() {
			ListDir(folder+"/"+file.Name(), conf)
		} else {
			filePath := folder + "/" + file.Name()
			// 检查到新文件 创建协程
			if strings.HasSuffix(filePath, ".csv") {
				logs.Debug("第一次初始化已经启动 MMAP映射 ：", filePath)
				waitgroup.Add(1) //检查到一个文件启动一个携程
				//go ReadFiles(conf,event.Name, data)
				go ReadFileUesMMAP(conf, filePath, data)
			}
		}
	}
	waitgroup.Done()
}

func InitLoadData(conf FtpSender) {
	logs.Debug("初始化目录")
	var Paths = strings.Split(conf.WatchDir, ",")
	for index := range Paths {
		logs.Debug("InitLoadData :", Paths[index])
		ListDir(Paths[index], conf)
	}

}

func main() {
	// 设置使用所有CPU
	runtime.GOMAXPROCS(runtime.NumCPU())

	_ = InitLog()

	waitgroup.Add(3)
	TConfig := GetConfig()

	InitLoadData(TConfig)
	go func() {
		WatchFiles(TConfig)
		waitgroup.Done()
	}()

	_, _ = NewMQSender(TConfig.SendGroup, TConfig.RocketMQNameserver)

	waitgroup.Wait() //.Wait()这里会发生阻塞，直到队列中所有的任务结束就会解除阻塞

}
