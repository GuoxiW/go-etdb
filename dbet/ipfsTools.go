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
	//fmt.Println(ipfsHashes) //map[testseries:{QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9 QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx }]
	//fmt.Println(b) //[123 10 32 34 116 101 115 116 115 101 114 105 101 115 34 58 32 123 10 32 32 34 100 34 58 32 34 81 109 102 68 97 116 65 65 104 49 75 122 80 117 70 86 115 65 97 72 111 49 86 83 55 98 52 87 80 101 72 111 52 67 112 119 65 88 83 51 85 50 120 115 87 69 34 44 10 32 32 34 107 34 58 32 34 81 109 81 55 98 77 119 67 89 90 116 118 111 114 67 103 97 81 115 89 97 49 104 82 82 116 67 57 120 118 115 74 56 49 106 71 71 112 71 66 71 51 107 90 71 57 34 44 10 32 32 34 99 34 58 32 34 81 109 102 74 120 119 69 66 67 98 102 101 53 83 82 81 112 80 49 84 49 106 97 74 114 67 76 77 117 66 83 119 80 56 70 103 112 101 86 53 52 114 115 80 76 120 34 10 32 125 10 125]
	return ioutil.WriteFile("ipfsHashes.json", b, 0644)
}

func ipfsPinPath(path string, name string) (string, error) {  // ipfs的pin命令在本地仓库中固定(或解除固定)ipfs对象。
	fmt.Println("Pinning " + name)
	//fmt.Println(path) // /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries
	bin := "ipfs"
	args := []string{"add", "-r", "-p=false", "--nocopy", path}  // -r 递归添加目录内容  -p 流式输出过程数据  --nocopy 使用filestore添加文件，实验特性
	//fmt.Println(args)  // [add -r -p=false --nocopy /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries]

	ial := exec.Command(bin, args...)
	//fmt.Println(ial) // &{/snap/bin/ipfs [ipfs add -r -p=false --nocopy /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries] []  <nil> <nil> <nil> [] <nil> <nil> <nil> <nil> <nil> false [] [] [] [] <nil> <nil>}
	//ial.Env = append(ial.Env, "IPFS_PATH=/home/guoxi/snap/ipfs")
	//ial.Env = append(ial.Env, "/home/guoxi/snap/ipfs/common/.ipfs")
	out, err := ial.CombinedOutput()

	fmt.Println(string(out)) 
	//added bafkreig65botpnfaoyaqw6y4fum42mtpjv7uwpn5jzmug52mxuhwnnr26m testseries/file_123/testfname.png
	//added bafkreihgaa3gsbzplqpa5ejg4x27be23dw2oclb4ovpxjebuuox2lh7r6m testseries/keyimg_testseries.jpg
	//added bafkreihgaa3gsbzplqpa5ejg4x27be23dw2oclb4ovpxjebuuox2lh7r6m testseries/keyimg_testseries_s.jpg
	//added QmNXJpmkQz7E15bN2xaE5xLZ5MQKxLcgoY9G1AeGq6AxoX testseries/keymov_testseries.flv
	//added bafkreicwdkuz6ttg57xvthypclt4yaxv5zo5idkcvjc637hgyvqv3a3txi testseries/rawdata/testfname.mp4
	//added QmPt7pF3tW5ED86Gp6cVsgfZZfMSFmsmLWUhpVbpNzz2fY testseries/file_123
	//added QmV8KjBom7dTA8Mrv4yscgmJfxWzLq8nieWtvHkbrFML6W testseries/rawdata
	//added QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE testseries
	
	//fmt.Println(err) //<nil>
	if err != nil {
		return string(out), err
	}

	lines := strings.Split(string(out), "\n")
	//fmt.Println(lines) //[added bafkreig65botpnfaoyaqw6y4fum42mtpjv7uwpn5jzmug52mxuhwnnr26m testseries/file_123/testfname.png added bafkreihgaa3gsbzplqpa5ejg4x27be23dw2oclb4ovpxjebuuox2lh7r6m testseries/keyimg_testseries.jpg added bafkreihgaa3gsbzplqpa5ejg4x27be23dw2oclb4ovpxjebuuox2lh7r6m testseries/keyimg_testseries_s.jpg added QmNXJpmkQz7E15bN2xaE5xLZ5MQKxLcgoY9G1AeGq6AxoX testseries/keymov_testseries.flv added bafkreicwdkuz6ttg57xvthypclt4yaxv5zo5idkcvjc637hgyvqv3a3txi testseries/rawdata/testfname.mp4 added QmPt7pF3tW5ED86Gp6cVsgfZZfMSFmsmLWUhpVbpNzz2fY testseries/file_123 added QmV8KjBom7dTA8Mrv4yscgmJfxWzLq8nieWtvHkbrFML6W testseries/rawdata added QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE testseries ]
	last := lines[len(lines)-2]
	//fmt.Println(last) //added QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE testseries
	words := strings.Split(last, " ")
	//fmt.Println(words)  //[added QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE testseries]

	if words[0] == "added" && words[2] == name {
		fmt.Println("Pinned. " + words[1])
		return words[1], nil
	} else {
		fmt.Println(words)
		return string(out), errors.New("ipfs hash not found")
	}
}

func ipfsAddLink(dirHash string, name string, link string) (string, error) {  // 为指定对象加入一个新的链接,为指定的一个数据链接到keymov。
	bin := "ipfs"
	args := []string{"object", "patch", "add-link", dirHash, name, link}
	//fmt.Println(args) //[object patch add-link QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE keymov_testseries.mp4 QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9]

	ial := exec.Command(bin, args...)
	//ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	//fmt.Println(string(out)) // QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx
	//fmt.Println(strings.TrimSpace(string(out))) // QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func ipfsNewUnixFsDir() (string, error) {  //新建类unix文件系统
	bin := "ipfs"
	args := []string{"object", "new", "unixfs-dir"}

	ial := exec.Command(bin, args...)
	//ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	//fmt.Println(out) //[81 109 85 78 76 76 115 80 65 67 67 122 49 118 76 120 81 86 107 88 113 113 76 88 53 82 49 88 51 52 53 113 113 102 72 98 115 102 54 55 104 118 65 51 78 110 10]
	//fmt.Println(string(out)) //QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn
	//fmt.Println(strings.TrimSpace(string(out))) //QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn
	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func containsEmptyFolder(cid string) (bool, error) {
	bin := "ipfs"
	args := []string{"object", "links", cid}
	fmt.Println(args) //[object links QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE]

	ial := exec.Command(bin, args...)
	fmt.Println(ial) //&{/snap/bin/ipfs [ipfs object links QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE] []  <nil> <nil> <nil> [] <nil> <nil> <nil> <nil> <nil> false [] [] [] [] <nil> <nil>}
	//ial.Env = append(ial.Env, "/home/guoxi/blockchain/tomography/.ipfs")
	out, err := ial.CombinedOutput()
	fmt.Println(out) //数字
	if err != nil {
		return false, err
	}

	return strings.Contains(string(out), "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"), nil
}