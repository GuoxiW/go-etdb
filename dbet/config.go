// 定义设置结构体。
// 具体的设置存贮在 ```conf.json``` 中。

package main

import (
	"os"
	"encoding/json"
)

var (
	config Configuration
)

type Configuration struct {
	DatabaseConfiguration `json:"databaseConfiguration"`
	FloConfiguration      `json:"floConfiguration"`
}

type FloConfiguration struct {
	FloAddress string  `json:"floAddress"`
	RpcAddress string  `json:"rpcAddress"`
	RpcUser    string  `json:"rpcUser"`
	RpcPass    string  `json:"rpcPass"`
	TxFeePerKb float64 `json:"txFeePerKb"`
}

type DatabaseConfiguration struct {
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
