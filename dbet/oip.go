// 定义了将数据发送给区块链和分块发送的函数。

package main

import (
	"strings"
	"strconv"
)

const maxDataSize = 1040
const maxPrefixNoRef = 200  // 没有编号时前缀的最大值
const maxPrefixRef = 250 // 有编号时前缀的最大值
const dataChunk1 = maxDataSize - maxPrefixNoRef // 没有编号时的组块大小
const dataChunkX = maxDataSize - maxPrefixRef // 有编号时的组块大小

func sendToBlockchain(data string) ([]string, error) {  //设置交易费用, 将数据发送给flo地址。

	l := len(data)
	//fmt.Println(data)
	//json:{"oip042":{"publish":{"artifact":{"floAddress":"oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ","timestamp":1605495868,"type":"research","subtype":"tomogram","info":{"title":"testtitle","tags":"etdb,jensen.lab,tomogram,electron.tomography","description":"Auto imported from etdb"},"details":{"date":1577836800,"NCBItaxID":1,"artNotes":"Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n","scopeName":"testscope","speciesName":"testsname","strain":"teststrain","tiltSingleDual":1,"defocus":0.1,"dosage":0.3,"tiltMin":0.4,"tiltMax":2,"tiltStep":0.1,"magnification":0.2,"emdb":"testemdb","microscopist":"testuname","institution":"Caltech","lab":"Jensen Lab","sid":"testseries"},"storage":{"network":"ipfs","location":"QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx","files":[{"dname":"testfname.png","fname":"file_123/testfname.png","fsize":14,"type":"tomogram","subtype":"snapshot","fNotes":"testnotes"},{"software":"testacquisition","dname":"testfname.mp4","fname":"rawdata/testfname.mp4","fsize":14,"type":"tomogram","subtype":"tiltSeries","fNotes":"testfname.mp4"},{"fname":"keyimg_testseries_s.jpg","fsize":24,"type":"image","subtype":"thumbnail","cType":"image/jpeg"},{"fname":"keyimg_testseries.jpg","fsize":24,"type":"tomogram","subtype":"keyimg","cType":"image/jpeg"},{"fname":"keymov_testseries.mp4","fsize":1990544,"type":"tomogram","subtype":"keymov","cType":"video/mp4"},{"fname":"keymov_testseries.flv","fsize":6389760,"type":"tomogram","subtype":"keymov","cType":"video/x-flv"}]},"signature":"H0K7FDE8YCykiwxGYB4x8Br4N13GPB4RDGM7lC0LMEIpAjhdTSgVI4bw8BgERADtejFQI0xOylysyLuWKIgzmPk="}}}}
	//fmt.Println(l) //1621

	//fmt.Println(config.TxFeePerKb) //0.001
	//err := setTxFee(config.TxFeePerKb)
	//if err != nil {
	//	return []string{}, nil
	//}

	// send as a single part
	//相当于如果不用分块的话，信息已经被私钥签署过，就直接发送。
	//如果需要分块，每个分块分别用私钥签署。
	if l <= maxDataSize {
		txid, err := sendToAddress(config.FloAddress, 0.1, data)
		//fmt.Println(txid) //空
		//fmt.Println(err) //空
		if err != nil {
			return []string{}, err
		}
		return []string{txid}, nil
	}

	var ret []string

	var i int64 = 0
	var chunkCount = int64((l-dataChunk1)/dataChunkX + 1) // 分块数量
	//fmt.Println(dataChunk1) //840
	//fmt.Println(dataChunkX) //790
	//fmt.Println(chunkCount) //1 相当于再发一次就结束
	remainder := data

	// send first master chunk
	chunk := remainder[:dataChunk1] // 分块1数据
	remainder = remainder[dataChunk1:] // 剩余部分数据
	ref, err := sendPart(i, chunkCount, "", chunk)
	//fmt.Println(ref) 返回 txid
	//1d135c385c790b5aca36ee7c148b2c14bb702e2547d1114247db3f3879a7c2b8

	if err != nil {
		return ret, err
	}
	ret = append(ret, ref)
	//fmt.Println(ret)
	//[1d135c385c790b5aca36ee7c148b2c14bb702e2547d1114247db3f3879a7c2b8]

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
	
	//fmt.Println(preImage)
	//0-1-oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ--json:{"oip042":{"publish":{"artifact":{"floAddress":"oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ","timestamp":1605498085,"type":"research","subtype":"tomogram","info":{"title":"testtitle","tags":"etdb,jensen.lab,tomogram,electron.tomography","description":"Auto imported from etdb"},"details":{"date":1577836800,"NCBItaxID":1,"artNotes":"Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n","scopeName":"testscope","speciesName":"testsname","strain":"teststrain","tiltSingleDual":1,"defocus":0.1,"dosage":0.3,"tiltMin":0.4,"tiltMax":2,"tiltStep":0.1,"magnification":0.2,"emdb":"testemdb","microscopist":"testuname","institution":"Caltech","lab":"Jensen Lab","sid":"testseries"},"storage":{"network":"ipfs","location":"QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx","files":[{"dname":"testfname.png","fname":"file_123/te
	//1-1-oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ-9ee65d7e4f4eb301fd5447b0fafb1c2092748937a7fd24dda4c9ad69c6c931ae-stfname.png","fsize":14,"type":"tomogram","subtype":"snapshot","fNotes":"testnotes"},{"software":"testacquisition","dname":"testfname.mp4","fname":"rawdata/testfname.mp4","fsize":14,"type":"tomogram","subtype":"tiltSeries","fNotes":"testfname.mp4"},{"fname":"keyimg_testseries_s.jpg","fsize":24,"type":"image","subtype":"thumbnail","cType":"image/jpeg"},{"fname":"keyimg_testseries.jpg","fsize":24,"type":"tomogram","subtype":"keyimg","cType":"image/jpeg"},{"fname":"keymov_testseries.mp4","fsize":1990544,"type":"tomogram","subtype":"keymov","cType":"video/mp4"},{"fname":"keymov_testseries.flv","fsize":6389760,"type":"tomogram","subtype":"keymov","cType":"video/x-flv"}]},"signature":"H+q6Czkh87izQEPAZ4bLYbGaGABi1wpqsDJz+zqdiPpYenJJYoFBcbL61/6Vt8H1QGA7TisFITpATZTlD5Rjvgk="}}}}

	sig, err := signMessage(config.FloAddress, preImage)
	
	//fmt.Println(sig)
	//IOLrpIR8IEqIQpzdyodn3mxFM9GZswxF73pGrZZp9yU/QXXf4CAKPcO3CNoFVa3Qhh70QXvKWmGaVoWMRQ/TeaA=
	//IOLrpIR8IEqIQpzdyodn3mxFM9GZswxF73pGrZZp9yU/QXXf4CAKPcO3CNoFVa3Qhh70QXvKWmGaVoWMRQ/TeaA=

	if err != nil {
		return "", err
	}

	meta := []string{p1, p2, config.FloAddress, reference, sig}
	floData := prefix + strings.Join(meta, ",") + suffix + data
	
	//fmt.Println(floData)
	//oip-mp(0,1,oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ,,H+WnAEyhafXh2yJuC3UT7nU80qsfUOfv0ABKEd3k9N2NKwCQrqsZ2fNFr/bwT/0nLwyyjtKcABtwG5KNcU5hljo=):json:{"oip042":{"publish":{"artifact":{"floAddress":"oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ","timestamp":1605498443,"type":"research","subtype":"tomogram","info":{"title":"testtitle","tags":"etdb,jensen.lab,tomogram,electron.tomography","description":"Auto imported from etdb"},"details":{"date":1577836800,"NCBItaxID":1,"artNotes":"Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n","scopeName":"testscope","speciesName":"testsname","strain":"teststrain","tiltSingleDual":1,"defocus":0.1,"dosage":0.3,"tiltMin":0.4,"tiltMax":2,"tiltStep":0.1,"magnification":0.2,"emdb":"testemdb","microscopist":"testuname","institution":"Caltech","lab":"Jensen Lab","sid":"testseries"},"storage":{"network":"ipfs","location":"QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx","files":[{"dname":"testfname.png","fname":"file_123/te
	//oip-mp(1,1,oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ,ba96c5abbfb24e7a57bcba5e1aae6fd431dcf1be03469211718f6ec462a9fc14,HwGhNeFDiD9fY4EJdkKBACrFtp2eZKfpvHIwe2WDOUTZaon5t/D1VvbfiY+J+odEY9Vnwa1socKZWUvkh0TTG30=):stfname.png","fsize":14,"type":"tomogram","subtype":"snapshot","fNotes":"testnotes"},{"software":"testacquisition","dname":"testfname.mp4","fname":"rawdata/testfname.mp4","fsize":14,"type":"tomogram","subtype":"tiltSeries","fNotes":"testfname.mp4"},{"fname":"keyimg_testseries_s.jpg","fsize":24,"type":"image","subtype":"thumbnail","cType":"image/jpeg"},{"fname":"keyimg_testseries.jpg","fsize":24,"type":"tomogram","subtype":"keyimg","cType":"image/jpeg"},{"fname":"keymov_testseries.mp4","fsize":1990544,"type":"tomogram","subtype":"keymov","cType":"video/mp4"},{"fname":"keymov_testseries.flv","fsize":6389760,"type":"tomogram","subtype":"keymov","cType":"video/x-flv"}]},"signature":"IBArDz8Zax4LetOvkeVkLxohP3/AppcnpS4hXq362PRZArmuQ6yQFh9FbYBdwhoqrq7Dxu/32foddj4bzY6WRHg="}}}}

	txid, err := sendToAddress(config.FloAddress, 0.1, floData)
	if err != nil {
		return "", err
	}

	return txid, nil
}
