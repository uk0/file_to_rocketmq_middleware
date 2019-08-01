# FileToRocketMQ Middleware



### 具体功能
  * 计算发送的TPS以及当前处理的文件以及进度。
  * 监控配置文件内的`watchDir` 里面如果有新文件被丢入，程序会将文件映射到内存通过`mmap`在进行启动多个协程发送到MQ