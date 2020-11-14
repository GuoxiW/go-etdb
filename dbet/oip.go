// 定义了将数据发送给区块链和分块发送的函数。

package main

import (
	"strings"
	"strconv"
	"fmt"
)

const maxDataSize = 1040
const maxPrefixNoRef = 200  // 没有编号时前缀的最大值
const maxPrefixRef = 250 // 有编号时前缀的最大值
const dataChunk1 = maxDataSize - maxPrefixNoRef // 没有编号时的组块大小
const dataChunkX = maxDataSize - maxPrefixRef // 有编号时的组块大小

func sendToBlockchain(data string) ([]string, error) {  //设置交易费用, 将数据发送给flo地址。

	l := len(data)
	//fmt.Println(data)
	//json:{"oip042":{"publish":{"artifact":{"floAddress":"ofMvqGLqxjdJr784cVGRquV3edJA5jEykd","timestamp":1605167830,"type":"research","subtype":"tomogram","info":{"title":"testtitle","tags":"etdb,jensen.lab,tomogram,electron.tomography","description":"Auto imported from etdb"},"details":{"date":1577836800,"NCBItaxID":1,"artNotes":"Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n","scopeName":"testscope","speciesName":"testsname","strain":"teststrain","tiltSingleDual":1,"defocus":0.1,"dosage":0.3,"tiltMin":0.4,"tiltMax":2,"tiltStep":0.1,"magnification":0.2,"emdb":"testemdb","microscopist":"testuname","institution":"Caltech","lab":"Jensen Lab","sid":"testseries"},"storage":{"network":"ipfs","location":"QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx","files":[{"dname":"testfname.png","fname":"file_123/testfname.png","fsize":14,"type":"tomogram","subtype":"snapshot","fNotes":"testnotes"},{"software":"testacquisition","dname":"testfname.mp4","fname":"rawdata/testfname.mp4","fsize":14,"type":"tomogram","subtype":"tiltSeries","fNotes":"testfname.mp4"},{"fname":"keyimg_testseries_s.jpg","fsize":24,"type":"image","subtype":"thumbnail","cType":"image/jpeg"},{"fname":"keyimg_testseries.jpg","fsize":24,"type":"tomogram","subtype":"keyimg","cType":"image/jpeg"},{"fname":"keymov_testseries.mp4","fsize":1990544,"type":"tomogram","subtype":"keymov","cType":"video/mp4"},{"fname":"keymov_testseries.flv","fsize":6389760,"type":"tomogram","subtype":"keymov","cType":"video/x-flv"}]},"signature":"Hx4Lf+StYF01XvPvULBk8jgqBerF51Bi6aqZKL3pCht+UcLKkDqHMgCVzfcBDls/1iGYnVZy/NPa0G6VFTF+JlQ="}}}}
	fmt.Println(l) //1621

	err := setTxFee(config.TxFeePerKb)
	if err != nil {
		return []string{}, nil
	}

	// send as a single part
	if l <= maxDataSize {
		txid, err := sendToAddress(config.FloAddress, 0.1, data)
		if err != nil {
			return []string{}, err
		}
		return []string{txid}, nil
	}

	var ret []string

	var i int64 = 0
	var chunkCount = int64((l-dataChunk1)/dataChunkX + 1) // 分块数量
	remainder := data

	// send first master chunk
	chunk := remainder[:dataChunk1] // 分块1数据
	remainder = remainder[dataChunk1:] // 剩余部分数据
	ref, err := sendPart(i, chunkCount, "", chunk)
	if err != nil {
		return ret, err
	}
	ret = append(ret, ref)

	for i++; i <= chunkCount; i++ {
		// if the last chunk don't out-of-bounds
		c := dataChunkX
		if c > len(remainder) {
			c = len(remainder)
		}
		// slice off a chunk to send
		chunk = remainder[:c]
		remainder = remainder[c:]

		txid, err := sendPart(i, chunkCount, ref, chunk)
		if err != nil {
			return ret, err
		}

		ret = append(ret, txid)
	}

	return ret, nil
}

func sendPart(part int64, of int64, reference string, data string) (string, error) {
	prefix := "oip-mp("
	suffix := "):"

	p1 := strconv.FormatInt(part, 10)  // 类型转换, 将part转换为10进制
	p2 := strconv.FormatInt(of, 10)  // 类型转换, 将of转换为10进制

	pi := []string{p1, p2, config.FloAddress, reference, data}
	preImage := strings.Join(pi, "-")

	sig, err := signMessage(config.FloAddress, preImage)
	if err != nil {
		return "", err
	}

	meta := []string{p1, p2, config.FloAddress, reference, sig}
	floData := prefix + strings.Join(meta, ",") + suffix + data

	txid, err := sendToAddress(config.FloAddress, 0.1, floData)
	if err != nil {
		return "", err
	}

	return txid, nil
}
