// 发布数据到区块链。

package main

import (
	"fmt" //格式化输出
	"encoding/json"
	//"github.com/oipwg/media-protocol/oip042"
	"github.com/GuoxiW/media-protocol/oip042"//  OIP-042 JSON 标准
	"strings"
	"time"
	"strconv" // 数据类型转换
	"os"
	"os/exec" // 执行外部命令
)

type OipArtifact struct {
	Pt oip042.PublishTomogram `json:"artifact"` // artifact 人工添加的信息 https://github.com/GuoxiW/media-protocol
}
type OipPublish struct {
	OipArtifact `json:"publish"`
}
type rWrap struct {
	OipPublish `json:"oip042"`
}

func main() {
	ids, err := GetFilterIdList()  // 根据filter.sql的设定获得 id 列表
	if err != nil { // := 是声明并赋值
		panic(err)
	}
	//fmt.Println(ids) [testseries]
	//fmt.Println(history) map[]


	for _, id := range ids {  // 列表进行循环
		fmt.Println(id)
		if _, ok := history[id]; ok {
			fmt.Printf("Tilt %s already published\n", id)  // 判断是否已经发布
			continue
		}

		pt, err := tiltIdToPublishTomogram(id)  // 发布单个id
		//fmt.Println(pt) // {{ 0   <nil> [] <nil> <nil> } {0 0 0      0 0 0 0 0 0 0 0     }}
		if err != nil {
			fmt.Println("Unable to obtain " + id)
			fmt.Println(err)
		} else {
			fmt.Println("---------")
			//PrettyPrint(pt)

			min, err := json.Marshal(rWrap{OipPublish{OipArtifact{pt}}}) // 将数据编码成json字符串
			if err != nil {
				panic(err)
			}
			ids, err := sendToBlockchain("json:" + string(min))  // 发送
			if err != nil {
				fmt.Println(ids)
				panic(err)
			} else {
				history[id] = ids
				PrettyPrint(ids) // json 格式处理
			}
		}

		err = saveHistory()  // 记录历史
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func convertVideo(flv string, mp4 string) error { // 转换数据格式
	fmt.Println("Converting " + flv + " -> " + mp4)
	bin := "ffmpeg"
	args := []string{"-i", flv, "-movflags", "faststart", "-nostats",
		"-n", "-vcodec", "libx264", "-pix_fmt", "yuv420p", "-vf",
		"pad=width=ceil(iw/2)*2:height=ceil(ih/2)*2", mp4}
	ial := exec.Command(bin, args...)
	out, err := ial.CombinedOutput()
	//fmt.Println(string(out))
	//ffmpeg version 3.4.8-0ubuntu0.2 Copyright (c) 2000-2020 the FFmpeg developers
	//  built with gcc 7 (Ubuntu 7.5.0-3ubuntu1~18.04)
	//  configuration: --prefix=/usr --extra-version=0ubuntu0.2 --toolchain=hardened --libdir=/usr/lib/x86_64-linux-gnu --incdir=/usr/include/x86_64-linux-gnu --enable-gpl --disable-stripping --enable-avresample --enable-avisynth --enable-gnutls --enable-ladspa --enable-libass --enable-libbluray --enable-libbs2b --enable-libcaca --enable-libcdio --enable-libflite --enable-libfontconfig --enable-libfreetype --enable-libfribidi --enable-libgme --enable-libgsm --enable-libmp3lame --enable-libmysofa --enable-libopenjpeg --enable-libopenmpt --enable-libopus --enable-libpulse --enable-librubberband --enable-librsvg --enable-libshine --enable-libsnappy --enable-libsoxr --enable-libspeex --enable-libssh --enable-libtheora --enable-libtwolame --enable-libvorbis --enable-libvpx --enable-libwavpack --enable-libwebp --enable-libx265 --enable-libxml2 --enable-libxvid --enable-libzmq --enable-libzvbi --enable-omx --enable-openal --enable-opengl --enable-sdl2 --enable-libdc1394 --enable-libdrm --enable-libiec61883 --enable-chromaprint --enable-frei0r --enable-libopencv --enable-libx264 --enable-shared
	//  libavutil      55. 78.100 / 55. 78.100
	//  libavcodec     57.107.100 / 57.107.100
	//  libavformat    57. 83.100 / 57. 83.100
	//  libavdevice    57. 10.100 / 57. 10.100
	//  libavfilter     6.107.100 /  6.107.100
	//  libavresample   3.  7.  0 /  3.  7.  0
	//  libswscale      4.  8.100 /  4.  8.100
	//  libswresample   2.  9.100 /  2.  9.100
	//  libpostproc    54.  7.100 / 54.  7.100
	//Input #0, flv, from '/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/keymov_testseries.flv':
	//  Metadata:
	//    encoder         : Lavf55.11.101
	//  Duration: 00:00:06.88, start: 0.000000, bitrate: 7435 kb/s
	//    Stream #0:0: Video: flv1, yuv420p, 512x356, 200 kb/s, 24 fps, 24 tbr, 1k tbn
	//Stream mapping:
	//  Stream #0:0 -> #0:0 (flv1 (flv) -> h264 (libx264))
	//Press [q] to stop, [?] for help
	//[libx264 @ 0x55d0124e3c00] using cpu capabilities: MMX2 SSE2Fast SSSE3 SSE4.2 AVX FMA3 BMI2 AVX2 AVX512
	//[libx264 @ 0x55d0124e3c00] profile High, level 2.1
	//[libx264 @ 0x55d0124e3c00] 264 - core 152 r2854 e9a5903 - H.264/MPEG-4 AVC codec - Copyleft 2003-2017 - http://www.videolan.org/x264.html - options: cabac=1 ref=3 deblock=1:0:0 analyse=0x3:0x113 me=hex subme=7 psy=1 psy_rd=1.00:0.00 mixed_ref=1 me_range=16 chroma_me=1 trellis=1 8x8dct=1 cqm=0 deadzone=21,11 fast_pskip=1 chroma_qp_offset=-2 threads=11 lookahead_threads=1 sliced_threads=0 nr=0 decimate=1 interlaced=0 bluray_compat=0 constrained_intra=0 bframes=3 b_pyramid=2 b_adapt=1 b_bias=0 direct=1 weightb=1 open_gop=0 weightp=2 keyint=250 keyint_min=24 scenecut=40 intra_refresh=0 rc_lookahead=40 rc=crf mbtree=1 crf=23.0 qcomp=0.60 qpmin=0 qpmax=69 qpstep=4 ip_ratio=1.40 aq=1:1.00
	//Output #0, mp4, to '/home/guoxi/snap/ipfs/blockchain/tomography/data/Videos/keymov_testseries.mp4':
	//  Metadata:
	//    encoder         : Lavf57.83.100
	//    Stream #0:0: Video: h264 (libx264) (avc1 / 0x31637661), yuv420p, 512x356, q=-1--1, 24 fps, 12288 tbn, 24 tbc
	//    Metadata:
	//      encoder         : Lavc57.107.100 libx264
	//    Side data:
	//      cpb: bitrate max/min/avg: 0/0/0 buffer size: 0 vbv_delay: -1
	//[flv @ 0x55d0124eadc0] illegal ac vlc code at 1x1
	//[flv @ 0x55d0124eadc0] Error at MB: 34
	//[flv @ 0x55d0124eadc0] concealing 736 DC, 736 AC, 736 MV errors in P frame
	//[mp4 @ 0x55d0124eb280] Starting second pass: moving the moov atom to the beginning of the file
	//frame=   79 fps=0.0 q=-1.0 Lsize=    1944kB time=00:00:03.16 bitrate=5028.6kbits/s speed=5.57x    
	//video:1943kB audio:0kB subtitle:0kB other streams:0kB global headers:0kB muxing overhead: 0.059919%
	//[libx264 @ 0x55d0124e3c00] frame I:1     Avg QP:27.96  size: 54123
	//[libx264 @ 0x55d0124e3c00] frame P:77    Avg QP:27.86  size: 25122
	//[libx264 @ 0x55d0124e3c00] frame B:1     Avg QP:31.00  size:   152
	//[libx264 @ 0x55d0124e3c00] consecutive B-frames: 97.5%  2.5%  0.0%  0.0%
	//[libx264 @ 0x55d0124e3c00] mb I  I16..4:  0.0% 77.7% 22.3%
	//[libx264 @ 0x55d0124e3c00] mb P  I16..4:  0.0%  0.0%  0.8%  P16..4: 68.8% 12.6% 17.8%  0.0%  0.0%    skip: 0.0%
	//[libx264 @ 0x55d0124e3c00] mb B  I16..4:  0.0%  0.0%  0.0%  B16..8:  4.1%  0.0%  0.0%  direct: 0.0%  skip:95.9%  L0: 0.0% L1:96.7% BI: 3.3%
	//[libx264 @ 0x55d0124e3c00] 8x8 transform intra:47.4% inter:68.2%
	//[libx264 @ 0x55d0124e3c00] coded y,uvDC,uvAC intra: 79.2% 0.0% 0.0% inter: 98.0% 0.0% 0.0%
	//[libx264 @ 0x55d0124e3c00] i8 v,h,dc,ddl,ddr,vr,hd,vl,hu: 12%  0% 41% 12%  7% 10%  2% 12%  3%
	//[libx264 @ 0x55d0124e3c00] i4 v,h,dc,ddl,ddr,vr,hd,vl,hu: 63%  1% 10%  5%  5%  6%  2%  5%  3%
	//[libx264 @ 0x55d0124e3c00] i8c dc,h,v,p: 100%  0%  0%  0%
	//[libx264 @ 0x55d0124e3c00] Weighted P-Frames: Y:24.7% UV:0.0%
	//[libx264 @ 0x55d0124e3c00] ref P L0: 88.7% 10.9%  0.3%  0.0%  0.0%
	//[libx264 @ 0x55d0124e3c00] kb/s:4833.20
	if err != nil && !strings.HasSuffix(string(out), "already exists. Exiting.\n") {
		return err
	}
	return nil
}

func processFiles(row TiltSeries) (ipfsHash, error) {
	h := ipfsHash{}
	//fmt.Println(h) //{   }
	//fmt.Println(row) // {testseries testtitle 2020-01-01 00:00:00 +0000 UTC testnotes testscope testroles testnotes testsname testnotes teststrain 1 1 0.1 0.2 0.3 0 0.4 2 0.1 testacquisition testprocess testemdb 0 0 testuname   [{2dimage testfname testnotes testtdimage tomogram snapshot /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname 123 0 }] [{rawdata testfname testfname tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname 123 testacquisition}]}
	//fmt.Println(row.Id) // testseries
	s, err := ipfsPinPath("/home/guoxi/snap/ipfs/blockchain/tomography/data/"+row.Id, row.Id)
	//fmt.Println(s)  //QmeDS7vWdhPQ6tL1kgUGdqYfy4CxDpxbrt7AbPweKDedKy 
	if err != nil {
		return h, err
	}
	h.Data = s
	//fmt.Println(h) //{QmeDS7vWdhPQ6tL1kgUGdqYfy4CxDpxbrt7AbPweKDedKy   }

	km := "keymov_" + row.Id
	//fmt.Println(km) //keymov_testseries
	//PrettyPrint(row)
	//{
	//  "Id": "testseries",
	//  "Title": "testtitle",
	//  "Date": "2020-01-01T00:00:00Z",
	//  "TiltSeriesNotes": "testnotes",
	//  "ScopeName": "testscope",
	//  "Roles": "testroles",
	//  "ScopeNotes": "testnotes",
	//  "SpeciesName": "testsname",
	//  "SpeciesNotes": "testnotes",
	//  "SpeciesStrain": "teststrain",
	//  "SpeciesTaxId": 1,
	//  "SingleDual": 1,
	//  "Defocus": 0.1,
	//  "Magnification": 0.2,
	//  "Dosage": 0.3,
	//  "TiltConstant": 0,
	//  "TiltMin": 0.4,
	//  "TiltMax": 2,
	//  "TiltStep": 0.1,
	//  "SoftwareAcquisition": "testacquisition",
	//  "SoftwareProcess": "testprocess",
	//  "Emdb": "testemdb",
	//  "KeyImg": 0,
	//  "KeyMov": 0,
	//  "Microscopist": "testuname",
	//  "Institution": "",
	//  "Lab": "",
	//  "DataFiles": [
	//    {
	//      "Filetype": "2dimage",
	//      "Filename": "testfname.png",
	//      "Notes": "testnotes",
	//      "ThreeDFileImage": "testtdimage",
	//      "Type": "tomogram",
	//      "SubType": "snapshot",
	//      "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png",
	//      "DefId": 123,
	//      "Auto": 0,
	//      "Software": ""
	//    }
	//  ],
	//  "ThreeDFiles": [
	//    {
	//      "Classify": "rawdata",
	//      "Notes": "testfname.mp4",
	//      "Filename": "testfname.mp4",
	//      "Type": "tomogram",
	//      "SubType": "tiltSeries",
	//      "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4",
	//      "DefId": 123,
	//      "Software": "testacquisition"
	//    }
	//  ]
	//}
	if row.KeyMov > 0 && row.KeyMov <= 4 {
		flv := "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + row.Id + "/" + km + ".flv"
		mp4 := "/home/guoxi/snap/ipfs/blockchain/tomography/data/Videos/" + km + ".mp4"
		// 只有一个flv格式的视频,转换成根目录/videos/下的.mp4视频
		err := convertVideo(flv, mp4)
		if err != nil {
			return h, err
		}
		s, err := ipfsPinPath(mp4, km+".mp4")
		//fmt.Println(s) // QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9
		if err != nil {
			return h, err
		}
		h.KeyMov = s
		//PrettyPrint(h)
		//{
		//  "d": "QmYwcsnkqo7qbeLHE4UkNWK7UezqrssBA6gudSFRT2eZ8T",
		//  "k": "QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9"
		//}
		//fmt.Println(h) //{QmYwcsnkqo7qbeLHE4UkNWK7UezqrssBA6gudSFRT2eZ8T QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9  }
	} else {
		h.KeyMov = "n/a"
	}
	//fmt.Println(h.Combined) // 空
	//fmt.Println(h.Data) // QmYwcsnkqo7qbeLHE4UkNWK7UezqrssBA6gudSFRT2eZ8T

	if h.KeyMov == "n/a" {
		h.Combined = h.Data
	} else {
		nh, err := ipfsAddLink(h.Data, km+".mp4", h.KeyMov) //为指定对象加入一个新的链接,为指定的一个数据链接到keymov。
		//fmt.Println(nh) // QmVqYs6amKxvrb4VETbkrKbCi4KvbX2M64ad6nASZ4Xoyq
		if err != nil {
			return h, err
		}
		h.Combined = nh
	}

	return h, nil
}

func tiltIdToPublishTomogram(tiltSeriesId string) (oip042.PublishTomogram, error) {
	tsr, err := GetTiltSeriesById(tiltSeriesId)  // 获取序列
	//fmt.Println(tsr) //{testseries testtitle 2020-01-01 00:00:00 +0000 UTC testnotes testscope testroles testnotes testsname testnotes teststrain 1 1 0.1 0.2 0.3 0 0.4 2 0.1 testacquisition testprocess testemdb 0 0 testuname   [{2dimage testfname.png testnotes testtdimage tomogram snapshot /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 123 0 }] [{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}]}
	//PrettyPrint(tsr)
	//{
	//  "Id": "testseries",
	//  "Title": "testtitle",
	//  "Date": "2020-01-01T00:00:00Z",
	//  "TiltSeriesNotes": "testnotes",
	//  "ScopeName": "testscope",
	//  "Roles": "testroles",
	//  "ScopeNotes": "testnotes",
	//  "SpeciesName": "testsname",
	//  "SpeciesNotes": "testnotes",
	//  "SpeciesStrain": "teststrain",
	//  "SpeciesTaxId": 1,
	//  "SingleDual": 1,
	//  "Defocus": 0.1,
	//  "Magnification": 0.2,
	//  "Dosage": 0.3,
	//  "TiltConstant": 0,
	//  "TiltMin": 0.4,
	//  "TiltMax": 2,
	//  "TiltStep": 0.1,
	//  "SoftwareAcquisition": "testacquisition",
	//  "SoftwareProcess": "testprocess",
	//  "Emdb": "testemdb",
	//  "KeyImg": 0,
	//  "KeyMov": 0,
	//  "Microscopist": "testuname",
	//  "Institution": "",
	//  "Lab": "",
	//  "DataFiles": [
	//    {
	//      "Filetype": "2dimage",
	//      "Filename": "testfname.png",
	//      "Notes": "testnotes",
	//      "ThreeDFileImage": "testtdimage",
	//      "Type": "tomogram",
	//      "SubType": "snapshot",
	//      "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png",
	//      "DefId": 123,
	//      "Auto": 0,
	//      "Software": ""
	//    }
	//  ],
	//  "ThreeDFiles": [
	//    {
	//      "Classify": "rawdata",
	//      "Notes": "testfname.mp4",
	//      "Filename": "testfname.mp4",
	//      "Type": "tomogram",
	//      "SubType": "tiltSeries",
	//      "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4",
	//      "DefId": 123,
	//      "Software": "testacquisition"
	//    }
	//  ]
	//}

	if err != nil {
		panic(err)
	}

	//PrettyPrint(tsr)
	var pt oip042.PublishTomogram
	//fmt.Println(tiltSeriesId)
	hash, ok := ipfsHashes[tiltSeriesId]  // 算它的 ipfs 哈希值
	//fmt.Println(hash) // {   } 新输入文件这里可能就是空的
	//fmt.Println(ok) // false
	emptyDir := false
	if ok {  // 如果hash不是空的，意味着不是新输入文件
		emptyDir, err = containsEmptyFolder(hash.Data)  // 判断是否空文件夹

		if err != nil {
			return pt, err
		}
	}
	if !ok || hash.Data == "" || hash.KeyMov == "" || hash.Combined == "" || emptyDir {  // 或关系，满足一个即可
		hash, err = processFiles(tsr)
		if err != nil {
			return pt, err
		}
		ipfsHashes[tiltSeriesId] = hash
		saveIpfsHashes()  // 本目录下保存 ipfs 哈希值
	}

	ts := time.Now().Unix()
	//fmt.Println(ts) //1604628499
	floAddress := config.FloAddress
	//fmt.Println(floAddress) //FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse
	//fmt.Println(tsr.Title)
	//fmt.Println(hash.Combined)

	//fmt.Println("reach here")
	//fmt.Println(tsr.Microscopist)
	//fmt.Println(tsr.Id)
	//fmt.Println(tsr.Magnification)
	//fmt.Println(tsr.Defocus)
	//fmt.Println(tsr.Dosage)
	//fmt.Println("reach here")
	//fmt.Println(tsr.TiltConstant)
	//fmt.Println(tsr.TiltMin)
	//fmt.Println(tsr.TiltMax)
	//fmt.Println(tsr.TiltStep)
	//fmt.Println(tsr.SpeciesStrain)
	//fmt.Println(tsr.SpeciesName)
	//fmt.Println("reach here")
	//fmt.Println(tsr.ScopeName)
	//fmt.Println(tsr.Date.Unix())
	//fmt.Println(tsr.Emdb)
	//fmt.Println(tsr.SingleDual)
	//fmt.Println(tsr.SpeciesTaxId)
	//fmt.Println(tsr.SpeciesName)



	pt = oip042.PublishTomogram{
		PublishArtifact: oip042.PublishArtifact{
			Type:       "research",
			SubType:    "tomogram",
			Timestamp:  ts,
			FloAddress: floAddress,
			Info: &oip042.ArtifactInfo{
				Title:       tsr.Title, //testtitle
				Description: "Auto imported from etdb",
				Tags:        "etdb,jensen.lab,tomogram,electron.tomography",
			},
			Storage: &oip042.ArtifactStorage{
				Network:  "ipfs",
				Location: hash.Combined, //QmVqYs6amKxvrb4VETbkrKbCi4KvbX2M64ad6nASZ4Xoyq
				Files:    []oip042.ArtifactFiles{},
			},
			Payment: nil, // it's free
		},
		TomogramDetails: oip042.TomogramDetails{
			Microscopist:   tsr.Microscopist, //testuname
			Institution:    "Caltech",
			Lab:            "Jensen Lab",
			Sid:            tsr.Id, //testseries
			Magnification:  tsr.Magnification, //0.2
			Defocus:        tsr.Defocus, //0.1
			Dosage:         tsr.Dosage, //0.3
			TiltConstant:   tsr.TiltConstant, //0
			TiltMin:        tsr.TiltMin, //0.4
			TiltMax:        tsr.TiltMax, //2
			TiltStep:       tsr.TiltStep, //0.1
			Strain:         tsr.SpeciesStrain, //teststrain
			SpeciesName:    tsr.SpeciesName, //testsname
			ScopeName:      tsr.ScopeName, //testscope
			Date:           tsr.Date.Unix(), // 1577836800
			Emdb:           tsr.Emdb, //testemdb
			TiltSingleDual: tsr.SingleDual, //1
			NCBItaxID:      tsr.SpeciesTaxId, //1
			// ToDo: Needs database cleanup before publishing Roles
			//Roles:        tsr.Roles,
		},
	}

	//PrettyPrint(pt)
	//{
	//  "floAddress": "FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse",
	//  "timestamp": 1604628871,
	//  "type": "research",
	//  "subtype": "tomogram",
	//  "info": {
	//    "title": "testtitle",
	//    "tags": "etdb,jensen.lab,tomogram,electron.tomography",
	//    "description": "Auto imported from etdb"
	//  },
	//  "details": {
	//    "date": 1577836800,
	//    "NCBItaxID": 1,
	//    "scopeName": "testscope",
	//    "speciesName": "testsname",
	//    "strain": "teststrain",
	//    "tiltSingleDual": 1,
	//    "defocus": 0.1,
	//    "dosage": 0.3,
	//    "tiltMin": 0.4,
	//    "tiltMax": 2,
	//    "tiltStep": 0.1,
	//    "magnification": 0.2,
	//    "emdb": "testemdb",
	//    "microscopist": "testuname",
	//    "institution": "Caltech",
	//    "lab": "Jensen Lab",
	//    "sid": "testseries"
	//  },
	//  "storage": {
	//    "network": "ipfs",
	//    "location": "QmVqYs6amKxvrb4VETbkrKbCi4KvbX2M64ad6nASZ4Xoyq"
	//  }
	//}
	//fmt.Println(tsr.ScopeNotes) //testnotes
	//fmt.Println(tsr.SpeciesNotes) //testnotes
	//fmt.Println(tsr.TiltSeriesNotes) //testnotes
	//fmt.Println(pt.TomogramDetails.ArtNotes) //空
	if len(tsr.ScopeNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Scope notes: " + tsr.ScopeNotes + "\n"
	}
	//fmt.Println(pt.TomogramDetails.ArtNotes)
	//Scope notes: testnotes
	if len(tsr.SpeciesNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Species notes: " + tsr.SpeciesNotes + "\n"
	}
	//fmt.Println(pt.TomogramDetails.ArtNotes)
	//Scope notes: testnotes
	//Species notes: testnotes
	if len(tsr.TiltSeriesNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Tilt series notes: " + tsr.TiltSeriesNotes + "\n"
	}
	//fmt.Println(pt.TomogramDetails.ArtNotes)
	//Scope notes: testnotes
	//Species notes: testnotes
	//Tilt series notes: testnotes

	capDir := ""
	//fmt.Println(tsr.DataFiles) //[{2dimage testfname.png testnotes testtdimage tomogram snapshot /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 123 0 }]
	for _, df := range tsr.DataFiles {
		//fmt.Println(df) //{2dimage testfname.png testnotes testtdimage tomogram snapshot /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 123 0 }
		fName := strings.TrimPrefix(df.FilePath, "/home/guoxi/snap/ipfs/blockchain/tomography/data/"+tsr.Id+"/")  // 返回不含前缀字符的 df.FilePath
		//fmt.Println(df.FilePath) ///home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png
		//fmt.Println(fName) //file_123/testfname.png
		//fmt.Println(df.Auto) //0
		if df.Auto == 2 {
			if capDir == "" {
				capDir, err = ipfsNewUnixFsDir()
				//fmt.Println(capDir) //QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn
				//fmt.Println(err) //<nil>
				if err != nil {
					return pt, err
				}
			}
			//fmt.Println(df.FilePath) ///home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png
			//fmt.Println(df.Filename) //testfname.png
			h, err := ipfsPinPath(df.FilePath, df.Filename)
			//fmt.Println(h) //bafkreig65botpnfaoyaqw6y4fum42mtpjv7uwpn5jzmug52mxuhwnnr26m
			if err != nil {
				return pt, err
			}
			capDir, err = ipfsAddLink(capDir, df.Filename, h)
			//fmt.Println(capDir) //QmPt7pF3tW5ED86Gp6cVsgfZZfMSFmsmLWUhpVbpNzz2fY
			if err != nil {
				return pt, err
			}
			fName =  "AutoCaps/" + strings.TrimPrefix(df.FilePath, "/home/guoxi/snap/ipfs/blockchain/tomography/data/Caps/")  // 返回不含前缀字符的 df.FilePath
			//fmt.Println(fName) //AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png
		}

		fi, err := os.Stat(df.FilePath)  // 获取文件属性
		//fmt.Println(fi) //&{testfname.png 14 420 {789376539 63739467815 0x913d80} {2049 5781644 1 33188 1000 1000 0 0 14 4096 8 {1604628498 361266611} {1603871015 789376539} {1603871015 789376539} [0 0 0]}}
		if err != nil {
			return pt, err
		}
		af := oip042.ArtifactFiles{
			Type:    df.Type, //"tomogram"
			SubType: df.SubType, //"snapshot"
			FNotes:  df.Notes, //"testnotes"
			Fsize:   fi.Size(), //14
			Dname:   df.Filename, //"testfname.png"
			Fname:   fName, //"file_123/testfname.png"
		}
		//PrettyPrint(af)
		//{
		//  "dname": "testfname.png",
		//  "fname": "file_123/testfname.png",
		//  "fsize": 14,
		//  "type": "tomogram",
		//  "subtype": "snapshot",
		//  "fNotes": "testnotes"
		//}
		//PrettyPrint(pt.Storage.Files)
		//[]
		pt.Storage.Files = append(pt.Storage.Files, af)
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  }
		//]		
	}

	//fmt.Println(capDir) //auto不等于2时为空
	if capDir != "" {
		hash.Caps, err = ipfsAddLink(hash.Combined, "AutoCaps", capDir)
		//fmt.Println(hash) // {QmYwcsnkqo7qbeLHE4UkNWK7UezqrssBA6gudSFRT2eZ8T QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9 QmVqYs6amKxvrb4VETbkrKbCi4KvbX2M64ad6nASZ4Xoyq QmbHY3oEcDaiEg3yMgDAFS5gc4FfsXNFi3PhtgmQvoaqAc}
		//PrettyPrint(hash)
		//{
		//  "d": "QmYwcsnkqo7qbeLHE4UkNWK7UezqrssBA6gudSFRT2eZ8T",
		//  "k": "QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9",
		//  "c": "QmVqYs6amKxvrb4VETbkrKbCi4KvbX2M64ad6nASZ4Xoyq",
		//  "caps": "QmbHY3oEcDaiEg3yMgDAFS5gc4FfsXNFi3PhtgmQvoaqAc"
		//}
		//fmt.Println(hash.Caps) //QmbHY3oEcDaiEg3yMgDAFS5gc4FfsXNFi3PhtgmQvoaqAc
		if err != nil {
			return pt, err
		}
		pt.Storage.Location = hash.Caps
		ipfsHashes[tsr.Id] = hash
		//PrettyPrint(ipfsHashes)
		//{
		//  "testseries": {
		//    "d": "QmYwcsnkqo7qbeLHE4UkNWK7UezqrssBA6gudSFRT2eZ8T",
		//    "k": "QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9",
		//    "c": "QmVqYs6amKxvrb4VETbkrKbCi4KvbX2M64ad6nASZ4Xoyq",
		//    "caps": "QmbHY3oEcDaiEg3yMgDAFS5gc4FfsXNFi3PhtgmQvoaqAc"
		//  }
		//}		
		saveIpfsHashes()
	}

	//fmt.Println(tsr.ThreeDFiles) //[{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}]
	for _, tdf := range tsr.ThreeDFiles {
		//fmt.Println(tdf) //{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}
		fi, err := os.Stat(tdf.FilePath)
		//fmt.Println(fi) //&{testfname.mp4 14 420 {113280763 63739467862 0x913d80} {2049 5781646 1 33188 1000 1000 0 0 14 4096 8 {1604755839 284226878} {1603871062 113280763} {1603871062 117280756} [0 0 0]}}
		//fmt.Println(err) //<nil>
		if err != nil {
			return pt, err
		}
		af := oip042.ArtifactFiles{
			Type:     tdf.Type,
			SubType:  tdf.SubType,
			FNotes:   tdf.Notes,
			Fsize:    fi.Size(),
			Dname:    tdf.Filename,
			Fname:    strings.TrimPrefix(tdf.FilePath, "/home/guoxi/snap/ipfs/blockchain/tomography/data/"+tsr.Id+"/"),
			Software: tdf.Software,
		}
		//PrettyPrint(af)
		//{
		//  "software": "testacquisition",
		//  "dname": "testfname.mp4",
		//  "fname": "rawdata/testfname.mp4",
		//  "fsize": 14,
		//  "type": "tomogram",
		//  "subtype": "tiltSeries",
		//  "fNotes": "testfname.mp4"
		//}
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes}]
		pt.Storage.Files = append(pt.Storage.Files, af)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4}]
	}

	if tsr.KeyImg > 0 && tsr.KeyImg <= 4 {
		kif := "keyimg_" + tsr.Id + "_s.jpg"
		//fmt.Println(kif) //keyimg_testseries_s.jpg
		fi, err := os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tsr.Id + "/" + kif)
		//fmt.Println(fi) &{keyimg_testseries_s.jpg 24 420 {628002886 63740354107 0x913d80} {2049 5784432 1 33188 1000 1000 0 0 24 4096 8 {1604757308 64007574} {1604757307 628002886} {1604757307 628002886} [0 0 0]}}
		//fmt.Println(err)
		if err != nil {
			return pt, err
		}
		ki := oip042.ArtifactFiles{
			Type:    "image",
			SubType: "thumbnail",
			CType:   "image/jpeg",
			Fsize:   fi.Size(),
			Fname:   kif,
		}
		//PrettyPrint(ki)
		//{
		//  "fname": "keyimg_testseries_s.jpg",
		//  "fsize": 24,
		//  "type": "image",
		//  "subtype": "thumbnail",
		//  "cType": "image/jpeg"
		//}
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4}]
		pt.Storage.Files = append(pt.Storage.Files, ki)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg }]

		kif = "keyimg_" + tsr.Id + ".jpg"
		//fmt.Println(kif) //keyimg_testseries.jpg
		fi, err = os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tsr.Id + "/" + kif)
		//fmt.Println(fi) //&{keyimg_testseries.jpg 24 420 {628002000 63740354107 0x913d80} {2049 5780637 1 33188 1000 1000 0 0 24 4096 8 {1604760324 429539000} {1604757307 628002000} {1604760331 869782938} [0 0 0]}}
		//fmt.Println(err) //<nil>
		if err != nil {
			return pt, err
		}
		ki = oip042.ArtifactFiles{
			Type:    "tomogram",
			SubType: "keyimg",
			CType:   "image/jpeg",
			Fsize:   fi.Size(),
			Fname:   "keyimg_" + tsr.Id + ".jpg",
		}
		//PrettyPrint(ki)
		//{
		//  "fname": "keyimg_testseries.jpg",
		//  "fsize": 24,
		//  "type": "tomogram",
		//  "subtype": "keyimg",
		//  "cType": "image/jpeg"
		//}
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg }]
		pt.Storage.Files = append(pt.Storage.Files, ki)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg }]
	}
	if tsr.KeyMov > 0 && tsr.KeyMov <= 4 {
		kmf := "keymov_" + tsr.Id + ".mp4"
		fi, err := os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/Videos/" + kmf)  // 获取文件属性
		//fmt.Println(fi) //&{keymov_testseries.mp4 1990544 436 {429177064 63740164999 0x913d80} {2049 5770299 1 33204 1000 1000 0 0 1990544 4096 3888 {1604755839 500229195} {1604568199 429177064} {1604568199 429177064} [0 0 0]}}
		if err != nil {
			return pt, err
		}
		km := oip042.ArtifactFiles{
			Type:    "tomogram",
			SubType: "keymov",
			CType:   "video/mp4",
			Fsize:   fi.Size(),
			Fname:   kmf,
		}
		//PrettyPrint(km)
		//{
		//  "fname": "keymov_testseries.mp4",
		//  "fsize": 1990544,
		//  "type": "tomogram",
		//  "subtype": "keymov",
		//  "cType": "video/mp4"
		//}
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg }]
		pt.Storage.Files = append(pt.Storage.Files, km)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg } { false  0 keymov_testseries.mp4 1990544 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/mp4 }]

		kmf = "keymov_" + tsr.Id + ".flv"
		//fmt.Println(kmf) //keymov_testseries.flv
		fi, err = os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tsr.Id + "/" + kmf)  // 获取文件属性
		//fmt.Println(fi) //&{keymov_testseries.flv 6389760 436 {948048000 63740164231 0x913d80} {2049 5781462 1 33204 1000 1000 0 0 6389760 4096 12480 {1604907635 634991295} {1604567431 948048000} {1604567541 21644465} [0 0 0]}}
		//fmt.Println(err) //<nil>
		if err != nil {
			return pt, err
		}
		km = oip042.ArtifactFiles{
			Type:    "tomogram",
			SubType: "keymov",
			CType:   "video/x-flv",
			Fsize:   fi.Size(),
			Fname:   kmf,
		}
		//PrettyPrint(km)
		//{
		//  "fname": "keymov_testseries.flv",
		//  "fsize": 6389760,
		//  "type": "tomogram",
		//  "subtype": "keymov",
		//  "cType": "video/x-flv"
		//}
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg } { false  0 keymov_testseries.mp4 1990544 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/mp4 }]
		pt.Storage.Files = append(pt.Storage.Files, km)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 AutoCaps//home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg } { false  0 keymov_testseries.mp4 1990544 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/mp4 } { false  0 keymov_testseries.flv 6389760 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/x-flv }]
	}

	loc := hash.Combined
	//fmt.Println(loc) //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx
	//PrettyPrint(hash)
	//{
	//  "d": "QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE",
	//  "k": "QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9",
	//  "c": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx",
	//  "caps": "QmZ8P5BLHRxKEg3pA91J7nktqGdeikvvapWiMz2Tzkskv3"
	//}
	if capDir != "" {
		loc = hash.Caps
	}
	//fmt.Println(loc) //QmZ8P5BLHRxKEg3pA91J7nktqGdeikvvapWiMz2Tzkskv3
	//fmt.Println(floAddress) //FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse
	//fmt.Println(strconv.FormatInt(ts, 10)) //1604908393
	v := []string{loc, floAddress, strconv.FormatInt(ts, 10)}
	//fmt.Println(v) //[QmZ8P5BLHRxKEg3pA91J7nktqGdeikvvapWiMz2Tzkskv3 FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse 1604908393]
	preImage := strings.Join(v, "-")
	//fmt.Println(preImage) //QmZ8P5BLHRxKEg3pA91J7nktqGdeikvvapWiMz2Tzkskv3-FS267EyRkAEyLNuSaTdbtkgcjjsXEBYhse-1604908393
	signature, err := signMessage(floAddress, preImage)
	//fmt.Println(signature) //{<nil> <nil> <nil>} Error getting json reply: Error sending json message: Post http://<username>:<password>@localhost:17317: dial tcp 127.0.0.1:17317: connect: connection refused
	fmt.Println(err)
	if err != nil {
		return pt, err
	}

	pt.Signature = signature

	return pt, nil
}

func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ") // json 格式处理
	fmt.Println(string(b))
}
