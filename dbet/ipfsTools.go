// 定义了保存ipfs哈希值等ipfs操作的函数，todo。

package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
	"fmt"
	"errors"
)

var ipfsHashes map[string]ipfsHash

type ipfsHash struct {
	Data     string `json:"d,omitempty"` // 指定如果字段具有空值，定义为false，0，零指针，nil接口值以及任何空数组，切片，映射或字符串，则该字段应从编码中省略
	KeyMov   string `json:"k,omitempty"`
	Combined string `json:"c,omitempty"`
	Caps     string `json:"caps,omitempty"`
}

func init() {
	file, err := os.Open("./ipfsHashes.json")
	// 如果没有json文件就直接return
	if os.IsNotExist(err) {
		ipfsHashes = make(map[string]ipfsHash, 100)  // 100可能是map的容量
		//fmt.Println(ipfsHashes) // map[]
		return
	}
	// 如果有json文件进行decode
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ipfsHashes)
	if err != nil {
		panic(err)
	}
}

func saveIpfsHashes() error {  // 存储ipfs哈希值
	b, _ := json.MarshalIndent(ipfsHashes, "", " ")
	return ioutil.WriteFile("ipfsHashes.json", b, 0644)
}

func ipfsPinPath(path string, name string) (string, error) {  // ipfs的pin命令在本地仓库中固定(或解除固定)ipfs对象。
	fmt.Println("Pinning " + name)
	//fmt.Println(path) // /home/guoxi/blockchain/tomography/data/testseries
	bin := "ipfs"
	args := []string{"add", "-r", "-p=false", "--nocopy", path}  // -r 递归添加目录内容  -p 流式输出过程数据  --nocopy 使用filestore添加文件，实验特性
	//fmt.Println(args)  // [add -r -p=false --nocopy /home/guoxi/blockchain/tomography/data/testseries]

	ial := exec.Command(bin, args...)
	//fmt.Println(ial) // &{/snap/bin/ipfs [ipfs add -r -p=false --nocopy /home/guoxi/blockchain/tomography/data/testseries] []  <nil> <nil> <nil> [] <nil> <nil> <nil> <nil> <nil> false [] [] [] [] <nil> <nil>}
	ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	fmt.Println(string(out)) // internal error, please report: running "ipfs" failed: cannot parse environment entry: "/home/guoxi/blockchain/tomography/.ipfs"
	//fmt.Println(err) // exit status 46
	if err != nil {
		return string(out), err
	}

	lines := strings.Split(string(out), "\n")
	last := lines[len(lines)-2]
	words := strings.Split(last, " ")

	if words[0] == "added" && words[2] == name {
		fmt.Println("Pinned. " + words[1])
		return words[1], nil
	} else {
		fmt.Println(words)
		return string(out), errors.New("ipfs hash not found")
	}
}

func ipfsAddLink(dirHash string, name string, link string) (string, error) {  // 为指定对象加入一个新的链接。
	bin := "ipfs"
	args := []string{"object", "patch", "add-link", dirHash, name, link}

	ial := exec.Command(bin, args...)
	ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func ipfsNewUnixFsDir() (string, error) {  //新建类unix文件系统
	bin := "ipfs"
	args := []string{"object", "new", "unixfs-dir"}

	ial := exec.Command(bin, args...)
	ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func containsEmptyFolder(cid string) (bool, error) {
	bin := "ipfs"
	args := []string{"object", "links", cid}
	//fmt.Println(args) // [object links ]

	ial := exec.Command(bin, args...)
	//fmt.Println(ial) // &{ipfs [ipfs object links ] []  <nil> <nil> <nil> [] <nil> <nil> <nil> <nil> 0xc420202600 false [] [] [] [] <nil> <nil>}
	ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	//fmt.Println(out) // []
	if err != nil {
		return false, err
	}

	return strings.Contains(string(out), "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"), nil
}