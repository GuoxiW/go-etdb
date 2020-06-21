// 定义了将数据发送给区块链和分块发送的函数。

package main

import (
	"strings"
	"strconv"
)

const maxDataSize = 1040
const maxPrefixNoRef = 200  // 有前缀的最大值
const maxPrefixRef = 250 // 没有前缀的最大值
const dataChunk1 = maxDataSize - maxPrefixNoRef // 有前缀时的组块大小
const dataChunkX = maxDataSize - maxPrefixRef // 没有前缀时的组块大小

func sendToBlockchain(data string) ([]string, error) {
	l := len(data)

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

	p1 := strconv.FormatInt(part, 10)
	p2 := strconv.FormatInt(of, 10)

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
