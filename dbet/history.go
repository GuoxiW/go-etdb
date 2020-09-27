// 在 ```history.json``` 中存着信息，如果没有会自动创建并返回，有的话解析文件。
// 定义了保存记录的函数。

package main

import (
	"encoding/json"
	"io/ioutil" // I/O实用程序函数
	"os"
)

var (
	history map[string][]string // Map 是一种无序的键值对的集合 history 是定义的变量
)

func init() { // 生成一个记录 history 的 map

	file, err := os.Open("./history.json")
	if os.IsNotExist(err) {
		history = make(map[string][]string, 100)
		return
	}

	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	if err != nil {
		panic(err)
	}
}

func saveHistory() error {  // 记录 history
	b, _ := json.MarshalIndent(history, "", " ")
	return ioutil.WriteFile("history.json", b, 0644)
}
