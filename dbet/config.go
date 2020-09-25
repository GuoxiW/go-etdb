// 将 conf.json 中的数据读出来。
// 具体的设置存贮在 ```conf.json``` 中。

package main

import (
	"encoding/json"
	"os"
)

var (
	config Configuration
)

type Configuration struct { // 指向本目录下的 config.go
	DatabaseConfiguration `json:"databaseConfiguration"`
	FloConfiguration      `json:"floConfiguration"`
}

type FloConfiguration struct { //读取 config.go 中的 floConfiguration
	FloAddress string  `json:"floAddress"`
	RpcAddress string  `json:"rpcAddress"`
	RpcUser    string  `json:"rpcUser"`
	RpcPass    string  `json:"rpcPass"`
	TxFeePerKb float64 `json:"txFeePerKb"`
}

type DatabaseConfiguration struct {  //读取 config.go 中的 databaseConfiguration
	User     string `json:"user"`
	Password string `json:"password"`
	Net      string `json:"net"`
	Address  string `json:"address"`
	Name     string `json:"name"`
}

func init() {
	file, err := os.Open("./conf.json")
	if err != nil { // nil 为零值
		panic(err) // 输出错误
	}
	defer file.Close() // 延迟，当函数执行到最后时，这些defer语句会按照逆序执行，最后该函数返回。
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
}
