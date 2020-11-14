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
	//fmt.Println(ids) //[testseries]
	//fmt.Println(history) //map[]


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
			//{
			//  "floAddress": "ofMvqGLqxjdJr784cVGRquV3edJA5jEykd",
			//  "timestamp": 1605166706,
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
			//    "artNotes": "Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n",
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
			//    "location": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx",
			//    "files": [
			//      {
			//        "dname": "testfname.png",
			//        "fname": "file_123/testfname.png",
			//        "fsize": 14,
			//        "type": "tomogram",
			//        "subtype": "snapshot",
			//        "fNotes": "testnotes"
			//      },
			//      {
			//        "software": "testacquisition",
			//        "dname": "testfname.mp4",
			//        "fname": "rawdata/testfname.mp4",
			//        "fsize": 14,
			//        "type": "tomogram",
			//        "subtype": "tiltSeries",
			//        "fNotes": "testfname.mp4"
			//      },
			//      {
			//        "fname": "keyimg_testseries_s.jpg",
			//        "fsize": 24,
			//        "type": "image",
			//        "subtype": "thumbnail",
			//        "cType": "image/jpeg"
			//      },
			//      {
			//        "fname": "keyimg_testseries.jpg",
			//        "fsize": 24,
			//        "type": "tomogram",
			//        "subtype": "keyimg",
			//        "cType": "image/jpeg"
			//      },
			//      {
			//        "fname": "keymov_testseries.mp4",
			//        "fsize": 1990544,
			//        "type": "tomogram",
			//        "subtype": "keymov",
			//        "cType": "video/mp4"
			//      },
			//      {
			//        "fname": "keymov_testseries.flv",
			//        "fsize": 6389760,
			//        "type": "tomogram",
			//        "subtype": "keymov",
			//        "cType": "video/x-flv"
			//      }
			//    ]
			//  },
			//  "signature": "Hx4Lf+StYF01XvPvULBk8jgqBerF51Bi6aqZKL3pCht+UcLKkDqHMgCVzfcBDls/1iGYnVZy/NPa0G6VFTF+JlQ="
			//}

			min, err := json.Marshal(rWrap{OipPublish{OipArtifact{pt}}}) // 将数据编码成json字符串
			//PrettyPrint(OipArtifact{pt})
			//{
			//  "artifact": {
			//    "floAddress": "ofMvqGLqxjdJr784cVGRquV3edJA5jEykd",
			//    "timestamp": 1605167128,
			//    "type": "research",
			//    "subtype": "tomogram",
			//    "info": {
			//      "title": "testtitle",
			//      "tags": "etdb,jensen.lab,tomogram,electron.tomography",
			//      "description": "Auto imported from etdb"
			//    },
			//    "details": {
			//      "date": 1577836800,
			//      "NCBItaxID": 1,
			//      "artNotes": "Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n",
			//      "scopeName": "testscope",
			//      "speciesName": "testsname",
			//      "strain": "teststrain",
			//      "tiltSingleDual": 1,
			//      "defocus": 0.1,
			//      "dosage": 0.3,
			//      "tiltMin": 0.4,
			//      "tiltMax": 2,
			//      "tiltStep": 0.1,
			//      "magnification": 0.2,
			//      "emdb": "testemdb",
			//      "microscopist": "testuname",
			//      "institution": "Caltech",
			//      "lab": "Jensen Lab",
			//      "sid": "testseries"
			//    },
			//    "storage": {
			//      "network": "ipfs",
			//      "location": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx",
			//      "files": [
			//        {
			//          "dname": "testfname.png",
			//          "fname": "file_123/testfname.png",
			//          "fsize": 14,
			//          "type": "tomogram",
			//          "subtype": "snapshot",
			//          "fNotes": "testnotes"
			//        },
			//        {
			//          "software": "testacquisition",
			//          "dname": "testfname.mp4",
			//          "fname": "rawdata/testfname.mp4",
			//          "fsize": 14,
			//          "type": "tomogram",
			//          "subtype": "tiltSeries",
			//          "fNotes": "testfname.mp4"
			//        },
			//        {
			//          "fname": "keyimg_testseries_s.jpg",
			//          "fsize": 24,
			//          "type": "image",
			//          "subtype": "thumbnail",
			//          "cType": "image/jpeg"
			//        },
			//        {
			//          "fname": "keyimg_testseries.jpg",
			//          "fsize": 24,
			//          "type": "tomogram",
			//          "subtype": "keyimg",
			//          "cType": "image/jpeg"
			//        },
			//        {
			//          "fname": "keymov_testseries.mp4",
			//          "fsize": 1990544,
			//          "type": "tomogram",
			//          "subtype": "keymov",
			//          "cType": "video/mp4"
			//        },
			//        {
			//          "fname": "keymov_testseries.flv",
			//          "fsize": 6389760,
			//          "type": "tomogram",
			//          "subtype": "keymov",
			//          "cType": "video/x-flv"
			//        }
			//      ]
			//    },
			//    "signature": "Hx4Lf+StYF01XvPvULBk8jgqBerF51Bi6aqZKL3pCht+UcLKkDqHMgCVzfcBDls/1iGYnVZy/NPa0G6VFTF+JlQ="
			//  }
			//}

			//PrettyPrint(OipPublish{OipArtifact{pt}})
			//{
			//  "publish": {
			//    "artifact": {
			//      "floAddress": "ofMvqGLqxjdJr784cVGRquV3edJA5jEykd",
			//      "timestamp": 1605167302,
			//      "type": "research",
			//      "subtype": "tomogram",
			//      "info": {
			//        "title": "testtitle",
			//        "tags": "etdb,jensen.lab,tomogram,electron.tomography",
			//        "description": "Auto imported from etdb"
			//      },
			//      "details": {
			//        "date": 1577836800,
			//        "NCBItaxID": 1,
			//        "artNotes": "Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n",
			//        "scopeName": "testscope",
			//        "speciesName": "testsname",
			//        "strain": "teststrain",
			//        "tiltSingleDual": 1,
			//        "defocus": 0.1,
			//        "dosage": 0.3,
			//        "tiltMin": 0.4,
			//        "tiltMax": 2,
			//        "tiltStep": 0.1,
			//        "magnification": 0.2,
			//        "emdb": "testemdb",
			//        "microscopist": "testuname",
			//        "institution": "Caltech",
			//        "lab": "Jensen Lab",
			//        "sid": "testseries"
			//      },
			//      "storage": {
			//        "network": "ipfs",
			//        "location": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx",
			//        "files": [
			//          {
			//            "dname": "testfname.png",
			//            "fname": "file_123/testfname.png",
			//            "fsize": 14,
			//            "type": "tomogram",
			//            "subtype": "snapshot",
			//            "fNotes": "testnotes"
			//          },
			//          {
			//            "software": "testacquisition",
			//            "dname": "testfname.mp4",
			//            "fname": "rawdata/testfname.mp4",
			//            "fsize": 14,
			//            "type": "tomogram",
			//            "subtype": "tiltSeries",
			//            "fNotes": "testfname.mp4"
			//          },
			//          {
			//            "fname": "keyimg_testseries_s.jpg",
			//            "fsize": 24,
			//            "type": "image",
			//            "subtype": "thumbnail",
			//            "cType": "image/jpeg"
			//          },
			//          {
			//            "fname": "keyimg_testseries.jpg",
			//            "fsize": 24,
			//            "type": "tomogram",
			//            "subtype": "keyimg",
			//            "cType": "image/jpeg"
			//          },
			//          {
			//            "fname": "keymov_testseries.mp4",
			//            "fsize": 1990544,
			//            "type": "tomogram",
			//            "subtype": "keymov",
			//            "cType": "video/mp4"
			//          },
			//          {
			//            "fname": "keymov_testseries.flv",
			//            "fsize": 6389760,
			//            "type": "tomogram",
			//            "subtype": "keymov",
			//            "cType": "video/x-flv"
			//          }
			//        ]
			//      },
			//      "signature": "Hx4Lf+StYF01XvPvULBk8jgqBerF51Bi6aqZKL3pCht+UcLKkDqHMgCVzfcBDls/1iGYnVZy/NPa0G6VFTF+JlQ="
			//    }
			//  }
			//}

			//PrettyPrint(rWrap{OipPublish{OipArtifact{pt}}})
			//{
			//  "oip042": {
			//    "publish": {
			//      "artifact": {
			//        "floAddress": "ofMvqGLqxjdJr784cVGRquV3edJA5jEykd",
			//        "timestamp": 1605167444,
			//        "type": "research",
			//        "subtype": "tomogram",
			//        "info": {
			//          "title": "testtitle",
			//          "tags": "etdb,jensen.lab,tomogram,electron.tomography",
			//          "description": "Auto imported from etdb"
			//        },
			//        "details": {
			//          "date": 1577836800,
			//          "NCBItaxID": 1,
			//          "artNotes": "Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n",
			//          "scopeName": "testscope",
			//          "speciesName": "testsname",
			//          "strain": "teststrain",
			//          "tiltSingleDual": 1,
			//          "defocus": 0.1,
			//          "dosage": 0.3,
			//          "tiltMin": 0.4,
			//          "tiltMax": 2,
			//          "tiltStep": 0.1,
			//          "magnification": 0.2,
			//          "emdb": "testemdb",
			//          "microscopist": "testuname",
			//          "institution": "Caltech",
			//          "lab": "Jensen Lab",
			//          "sid": "testseries"
			//        },
			//        "storage": {
			//          "network": "ipfs",
			//          "location": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx",
			//          "files": [
			//            {
			//              "dname": "testfname.png",
			//              "fname": "file_123/testfname.png",
			//              "fsize": 14,
			//              "type": "tomogram",
			//              "subtype": "snapshot",
			//              "fNotes": "testnotes"
			//            },
			//            {
			//              "software": "testacquisition",
			//              "dname": "testfname.mp4",
			//              "fname": "rawdata/testfname.mp4",
			//              "fsize": 14,
			//              "type": "tomogram",
			//              "subtype": "tiltSeries",
			//              "fNotes": "testfname.mp4"
			//            },
			//            {
			//              "fname": "keyimg_testseries_s.jpg",
			//              "fsize": 24,
			//              "type": "image",
			//              "subtype": "thumbnail",
			//              "cType": "image/jpeg"
			//            },
			//            {
			//              "fname": "keyimg_testseries.jpg",
			//              "fsize": 24,
			//              "type": "tomogram",
			//              "subtype": "keyimg",
			//              "cType": "image/jpeg"
			//            },
			//            {
			//              "fname": "keymov_testseries.mp4",
			//              "fsize": 1990544,
			//              "type": "tomogram",
			//              "subtype": "keymov",
			//              "cType": "video/mp4"
			//            },
			//            {
			//              "fname": "keymov_testseries.flv",
			//              "fsize": 6389760,
			//              "type": "tomogram",
			//              "subtype": "keymov",
			//              "cType": "video/x-flv"
			//            }
			//          ]
			//        },
			//        "signature": "Hx4Lf+StYF01XvPvULBk8jgqBerF51Bi6aqZKL3pCht+UcLKkDqHMgCVzfcBDls/1iGYnVZy/NPa0G6VFTF+JlQ="
			//      }
			//    }
			//  }
			//}

			//fmt.Println(json.Marshal(rWrap{OipPublish{OipArtifact{pt}}}))
			//[123 34 111 105 112 48 52 50 34 58 123 34 112 117 98 108 105 115 104 34 58 123 34 97 114 116 105 102 97 99 116 34 58 123 34 102 108 111 65 100 100 114 101 115 115 34 58 34 111 102 77 118 113 71 76 113 120 106 100 74 114 55 56 52 99 86 71 82 113 117 86 51 101 100 74 65 53 106 69 121 107 100 34 44 34 116 105 109 101 115 116 97 109 112 34 58 49 54 48 53 49 54 55 53 50 54 44 34 116 121 112 101 34 58 34 114 101 115 101 97 114 99 104 34 44 34 115 117 98 116 121 112 101 34 58 34 116 111 109 111 103 114 97 109 34 44 34 105 110 102 111 34 58 123 34 116 105 116 108 101 34 58 34 116 101 115 116 116 105 116 108 101 34 44 34 116 97 103 115 34 58 34 101 116 100 98 44 106 101 110 115 101 110 46 108 97 98 44 116 111 109 111 103 114 97 109 44 101 108 101 99 116 114 111 110 46 116 111 109 111 103 114 97 112 104 121 34 44 34 100 101 115 99 114 105 112 116 105 111 110 34 58 34 65 117 116 111 32 105 109 112 111 114 116 101 100 32 102 114 111 109 32 101 116 100 98 34 125 44 34 100 101 116 97 105 108 115 34 58 123 34 100 97 116 101 34 58 49 53 55 55 56 51 54 56 48 48 44 34 78 67 66 73 116 97 120 73 68 34 58 49 44 34 97 114 116 78 111 116 101 115 34 58 34 83 99 111 112 101 32 110 111 116 101 115 58 32 116 101 115 116 110 111 116 101 115 92 110 83 112 101 99 105 101 115 32 110 111 116 101 115 58 32 116 101 115 116 110 111 116 101 115 92 110 84 105 108 116 32 115 101 114 105 101 115 32 110 111 116 101 115 58 32 116 101 115 116 110 111 116 101 115 92 110 34 44 34 115 99 111 112 101 78 97 109 101 34 58 34 116 101 115 116 115 99 111 112 101 34 44 34 115 112 101 99 105 101 115 78 97 109 101 34 58 34 116 101 115 116 115 110 97 109 101 34 44 34 115 116 114 97 105 110 34 58 34 116 101 115 116 115 116 114 97 105 110 34 44 34 116 105 108 116 83 105 110 103 108 101 68 117 97 108 34 58 49 44 34 100 101 102 111 99 117 115 34 58 48 46 49 44 34 100 111 115 97 103 101 34 58 48 46 51 44 34 116 105 108 116 77 105 110 34 58 48 46 52 44 34 116 105 108 116 77 97 120 34 58 50 44 34 116 105 108 116 83 116 101 112 34 58 48 46 49 44 34 109 97 103 110 105 102 105 99 97 116 105 111 110 34 58 48 46 50 44 34 101 109 100 98 34 58 34 116 101 115 116 101 109 100 98 34 44 34 109 105 99 114 111 115 99 111 112 105 115 116 34 58 34 116 101 115 116 117 110 97 109 101 34 44 34 105 110 115 116 105 116 117 116 105 111 110 34 58 34 67 97 108 116 101 99 104 34 44 34 108 97 98 34 58 34 74 101 110 115 101 110 32 76 97 98 34 44 34 115 105 100 34 58 34 116 101 115 116 115 101 114 105 101 115 34 125 44 34 115 116 111 114 97 103 101 34 58 123 34 110 101 116 119 111 114 107 34 58 34 105 112 102 115 34 44 34 108 111 99 97 116 105 111 110 34 58 34 81 109 102 74 120 119 69 66 67 98 102 101 53 83 82 81 112 80 49 84 49 106 97 74 114 67 76 77 117 66 83 119 80 56 70 103 112 101 86 53 52 114 115 80 76 120 34 44 34 102 105 108 101 115 34 58 91 123 34 100 110 97 109 101 34 58 34 116 101 115 116 102 110 97 109 101 46 112 110 103 34 44 34 102 110 97 109 101 34 58 34 102 105 108 101 95 49 50 51 47 116 101 115 116 102 110 97 109 101 46 112 110 103 34 44 34 102 115 105 122 101 34 58 49 52 44 34 116 121 112 101 34 58 34 116 111 109 111 103 114 97 109 34 44 34 115 117 98 116 121 112 101 34 58 34 115 110 97 112 115 104 111 116 34 44 34 102 78 111 116 101 115 34 58 34 116 101 115 116 110 111 116 101 115 34 125 44 123 34 115 111 102 116 119 97 114 101 34 58 34 116 101 115 116 97 99 113 117 105 115 105 116 105 111 110 34 44 34 100 110 97 109 101 34 58 34 116 101 115 116 102 110 97 109 101 46 109 112 52 34 44 34 102 110 97 109 101 34 58 34 114 97 119 100 97 116 97 47 116 101 115 116 102 110 97 109 101 46 109 112 52 34 44 34 102 115 105 122 101 34 58 49 52 44 34 116 121 112 101 34 58 34 116 111 109 111 103 114 97 109 34 44 34 115 117 98 116 121 112 101 34 58 34 116 105 108 116 83 101 114 105 101 115 34 44 34 102 78 111 116 101 115 34 58 34 116 101 115 116 102 110 97 109 101 46 109 112 52 34 125 44 123 34 102 110 97 109 101 34 58 34 107 101 121 105 109 103 95 116 101 115 116 115 101 114 105 101 115 95 115 46 106 112 103 34 44 34 102 115 105 122 101 34 58 50 52 44 34 116 121 112 101 34 58 34 105 109 97 103 101 34 44 34 115 117 98 116 121 112 101 34 58 34 116 104 117 109 98 110 97 105 108 34 44 34 99 84 121 112 101 34 58 34 105 109 97 103 101 47 106 112 101 103 34 125 44 123 34 102 110 97 109 101 34 58 34 107 101 121 105 109 103 95 116 101 115 116 115 101 114 105 101 115 46 106 112 103 34 44 34 102 115 105 122 101 34 58 50 52 44 34 116 121 112 101 34 58 34 116 111 109 111 103 114 97 109 34 44 34 115 117 98 116 121 112 101 34 58 34 107 101 121 105 109 103 34 44 34 99 84 121 112 101 34 58 34 105 109 97 103 101 47 106 112 101 103 34 125 44 123 34 102 110 97 109 101 34 58 34 107 101 121 109 111 118 95 116 101 115 116 115 101 114 105 101 115 46 109 112 52 34 44 34 102 115 105 122 101 34 58 49 57 57 48 53 52 52 44 34 116 121 112 101 34 58 34 116 111 109 111 103 114 97 109 34 44 34 115 117 98 116 121 112 101 34 58 34 107 101 121 109 111 118 34 44 34 99 84 121 112 101 34 58 34 118 105 100 101 111 47 109 112 52 34 125 44 123 34 102 110 97 109 101 34 58 34 107 101 121 109 111 118 95 116 101 115 116 115 101 114 105 101 115 46 102 108 118 34 44 34 102 115 105 122 101 34 58 54 51 56 57 55 54 48 44 34 116 121 112 101 34 58 34 116 111 109 111 103 114 97 109 34 44 34 115 117 98 116 121 112 101 34 58 34 107 101 121 109 111 118 34 44 34 99 84 121 112 101 34 58 34 118 105 100 101 111 47 120 45 102 108 118 34 125 93 125 44 34 115 105 103 110 97 116 117 114 101 34 58 34 72 120 52 76 102 43 83 116 89 70 48 49 88 118 80 118 85 76 66 107 56 106 103 113 66 101 114 70 53 49 66 105 54 97 113 90 75 76 51 112 67 104 116 43 85 99 76 75 107 68 113 72 77 103 67 86 122 102 99 66 68 108 115 47 49 105 71 89 110 86 90 121 47 78 80 97 48 71 54 86 70 84 70 43 74 108 81 61 34 125 125 125 125] <nil>
			if err != nil {
				panic(err)
			}
			//fmt.Println(string(min))
			//{"oip042":{"publish":{"artifact":{"floAddress":"ofMvqGLqxjdJr784cVGRquV3edJA5jEykd","timestamp":1605167645,"type":"research","subtype":"tomogram","info":{"title":"testtitle","tags":"etdb,jensen.lab,tomogram,electron.tomography","description":"Auto imported from etdb"},"details":{"date":1577836800,"NCBItaxID":1,"artNotes":"Scope notes: testnotes\nSpecies notes: testnotes\nTilt series notes: testnotes\n","scopeName":"testscope","speciesName":"testsname","strain":"teststrain","tiltSingleDual":1,"defocus":0.1,"dosage":0.3,"tiltMin":0.4,"tiltMax":2,"tiltStep":0.1,"magnification":0.2,"emdb":"testemdb","microscopist":"testuname","institution":"Caltech","lab":"Jensen Lab","sid":"testseries"},"storage":{"network":"ipfs","location":"QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx","files":[{"dname":"testfname.png","fname":"file_123/testfname.png","fsize":14,"type":"tomogram","subtype":"snapshot","fNotes":"testnotes"},{"software":"testacquisition","dname":"testfname.mp4","fname":"rawdata/testfname.mp4","fsize":14,"type":"tomogram","subtype":"tiltSeries","fNotes":"testfname.mp4"},{"fname":"keyimg_testseries_s.jpg","fsize":24,"type":"image","subtype":"thumbnail","cType":"image/jpeg"},{"fname":"keyimg_testseries.jpg","fsize":24,"type":"tomogram","subtype":"keyimg","cType":"image/jpeg"},{"fname":"keymov_testseries.mp4","fsize":1990544,"type":"tomogram","subtype":"keymov","cType":"video/mp4"},{"fname":"keymov_testseries.flv","fsize":6389760,"type":"tomogram","subtype":"keymov","cType":"video/x-flv"}]},"signature":"Hx4Lf+StYF01XvPvULBk8jgqBerF51Bi6aqZKL3pCht+UcLKkDqHMgCVzfcBDls/1iGYnVZy/NPa0G6VFTF+JlQ="}}}}
			ids, err := sendToBlockchain("json:" + string(min))  // 发送
			if err != nil {
				fmt.Println(ids)
				panic(err)
			} else {
				history[id] = ids
				PrettyPrint(ids) // json 格式处理
			}
		}

		//err = saveHistory()  // 记录历史
		//if err != nil {
		//	panic(err)
		//}
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
	//fmt.Println(tsr) //{testseries testtitle 2020-01-01 00:00:00 +0000 UTC testnotes testscope testroles testnotes testsname testnotes teststrain 1 1 0.1 0.2 0.3 0 0.4 2 0.1 testacquisition testprocess testemdb 1 1 testuname   [{2dimage testfname.png testnotes testtdimage tomogram snapshot /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png 123 0 }] [{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}]}
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
	//  "KeyImg": 1,
	//  "KeyMov": 1,
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

	var pt oip042.PublishTomogram
	//fmt.Println(tiltSeriesId) //testseries
	hash, ok := ipfsHashes[tiltSeriesId]  // 算它的 ipfs 哈希值
	//fmt.Println(hash) //{   } 新输入文件这里可能就是空的
	//fmt.Println(ok) //false
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
		//saveIpfsHashes()  // 本目录下保存 ipfs 哈希值
	}

	ts := time.Now().Unix()
	//fmt.Println(ts) //1605354951
	floAddress := config.FloAddress
	//fmt.Println(floAddress) //oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ
	//fmt.Println(tsr.Title) //testtitle
	//fmt.Println(hash.Combined) //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx

	//fmt.Println(tsr.Microscopist) //testuname
	//fmt.Println(tsr.Id) //testseries
	//fmt.Println(tsr.Magnification) //0.2
	//fmt.Println(tsr.Defocus) //0.1
	//fmt.Println(tsr.Dosage) //0.3
	//fmt.Println(tsr.TiltConstant) //0
	//fmt.Println(tsr.TiltMin) //0.4
	//fmt.Println(tsr.TiltMax) //2
	//fmt.Println(tsr.TiltStep) //0.1
	//fmt.Println(tsr.SpeciesStrain) //teststrain
	//fmt.Println(tsr.SpeciesName) //testsname
	//fmt.Println(tsr.ScopeName) //testscope
	//fmt.Println(tsr.Date.Unix()) //1577836800
	//fmt.Println(tsr.Emdb) //testemdb
	//fmt.Println(tsr.SingleDual) //1
	//fmt.Println(tsr.SpeciesTaxId) //1
	//fmt.Println(tsr.SpeciesName) //testsname



	pt = oip042.PublishTomogram{
		PublishArtifact: oip042.PublishArtifact{
			Type:       "research",
			SubType:    "tomogram",
			Timestamp:  ts, //1605354951
			FloAddress: floAddress, //oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ
			Info: &oip042.ArtifactInfo{
				Title:       tsr.Title, //testtitle
				Description: "Auto imported from etdb",
				Tags:        "etdb,jensen.lab,tomogram,electron.tomography",
			},
			Storage: &oip042.ArtifactStorage{
				Network:  "ipfs",
				Location: hash.Combined, //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx
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
	//  "floAddress": "oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ",
	//  "timestamp": 1605361515,
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
	//    "location": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx"
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
	//PrettyPrint(tsr.DataFiles)
	//[
	//  {
	//    "Filetype": "2dimage",
	//    "Filename": "testfname.png",
	//    "Notes": "testnotes",
	//    "ThreeDFileImage": "testtdimage",
	//    "Type": "tomogram",
	//    "SubType": "snapshot",
	//    "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/file_123/testfname.png",
	//    "DefId": 123,
	//    "Auto": 0,
	//    "Software": ""
	//  }
	//]
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
		//fmt.Println(fi) //&{testfname.png 14 420 {789376539 63739467815 0x913da0} {2049 5781644 1 33188 1000 1000 0 0 14 4096 8 {1605325738 610129418} {1603871015 789376539} {1603871015 789376539} [0 0 0]}}
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
		//saveIpfsHashes()
	}

	//fmt.Println(tsr.ThreeDFiles) //[{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}]
	//PrettyPrint(tsr.ThreeDFiles)
	//[
	//  {
	//    "Classify": "rawdata",
	//    "Notes": "testfname.mp4",
	//    "Filename": "testfname.mp4",
	//    "Type": "tomogram",
	//    "SubType": "tiltSeries",
	//    "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4",
	//    "DefId": 123,
	//    "Software": "testacquisition"
	//  }
	//]
	for _, tdf := range tsr.ThreeDFiles {
		//fmt.Println(tdf) //{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}
		//PrettyPrint(tdf)
		//{
		//  "Classify": "rawdata",
		//  "Notes": "testfname.mp4",
		//  "Filename": "testfname.mp4",
		//  "Type": "tomogram",
		//  "SubType": "tiltSeries",
		//  "FilePath": "/home/guoxi/snap/ipfs/blockchain/tomography/data/testseries/rawdata/testfname.mp4",
		//  "DefId": 123,
		//  "Software": "testacquisition"
		//}
		fi, err := os.Stat(tdf.FilePath)
		//fmt.Println(fi) //&{testfname.mp4 14 420 {113280763 63739467862 0x913da0} {2049 5781646 1 33188 1000 1000 0 0 14 4096 8 {1605325738 654129649} {1603871062 113280763} {1603871062 117280756} [0 0 0]}}
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
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes}]
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
		pt.Storage.Files = append(pt.Storage.Files, af)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4}]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  }
		//]
	}

	if tsr.KeyImg > 0 && tsr.KeyImg <= 4 {
		kif := "keyimg_" + tsr.Id + "_s.jpg"
		//fmt.Println(kif) //keyimg_testseries_s.jpg
		fi, err := os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tsr.Id + "/" + kif)
		//fmt.Println(fi) //&{keyimg_testseries_s.jpg 24 420 {628002886 63740354107 0x914da0} {2049 5784432 1 33188 1000 1000 0 0 24 4096 8 {1605325738 610129418} {1604757307 628002886} {1604757307 628002886} [0 0 0]}}
		//fmt.Println(err) //<nil>
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
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4}]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  }
		//]
		pt.Storage.Files = append(pt.Storage.Files, ki)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  }
		//]

		kif = "keyimg_" + tsr.Id + ".jpg"
		//fmt.Println(kif) //keyimg_testseries.jpg
		fi, err = os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tsr.Id + "/" + kif)
		//fmt.Println(fi) //&{keyimg_testseries.jpg 24 420 {628002000 63740354107 0x913da0} {2049 5780637 1 33188 1000 1000 0 0 24 4096 8 {1605325738 610129418} {1604757307 628002000} {1604760331 869782938} [0 0 0]}}
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
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  }
		//]
		pt.Storage.Files = append(pt.Storage.Files, ki)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keyimg_testseries.jpg",
		//    "fsize": 24,
		//    "type": "tomogram",
		//    "subtype": "keyimg",
		//    "cType": "image/jpeg"
		//  }
		//]
	}
	if tsr.KeyMov > 0 && tsr.KeyMov <= 4 {
		kmf := "keymov_" + tsr.Id + ".mp4"
		//fmt.Println(kmf) //keymov_testseries.mp4
		fi, err := os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/Videos/" + kmf)  // 获取文件属性
		//fmt.Println(fi) //&{keymov_testseries.mp4 1990544 436 {429177064 63740164999 0x915da0} {2049 5770299 1 33204 1000 1000 0 0 1990544 4096 3888 {1605325738 918131039} {1604568199 429177064} {1604568199 429177064} [0 0 0]}}
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
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keyimg_testseries.jpg",
		//    "fsize": 24,
		//    "type": "tomogram",
		//    "subtype": "keyimg",
		//    "cType": "image/jpeg"
		//  }
		//]
		pt.Storage.Files = append(pt.Storage.Files, km)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg } { false  0 keymov_testseries.mp4 1990544 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/mp4 }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keyimg_testseries.jpg",
		//    "fsize": 24,
		//    "type": "tomogram",
		//    "subtype": "keyimg",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keymov_testseries.mp4",
		//    "fsize": 1990544,
		//    "type": "tomogram",
		//    "subtype": "keymov",
		//    "cType": "video/mp4"
		//  }
		//]

		kmf = "keymov_" + tsr.Id + ".flv"
		//fmt.Println(kmf) //keymov_testseries.flv
		fi, err = os.Stat("/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tsr.Id + "/" + kmf)  // 获取文件属性
		//fmt.Println(fi) //&{keymov_testseries.flv 6389760 436 {948048000 63740164231 0x915da0} {2049 5781462 1 33204 1000 1000 0 0 6389760 4096 12480 {1605325738 610129418} {1604567431 948048000} {1604567541 21644465} [0 0 0]}}
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
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg } { false  0 keymov_testseries.mp4 1990544 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/mp4 }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keyimg_testseries.jpg",
		//    "fsize": 24,
		//    "type": "tomogram",
		//    "subtype": "keyimg",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keymov_testseries.mp4",
		//    "fsize": 1990544,
		//    "type": "tomogram",
		//    "subtype": "keymov",
		//    "cType": "video/mp4"
		//  }
		//]
		pt.Storage.Files = append(pt.Storage.Files, km)
		//fmt.Println(pt.Storage.Files) //[{ false testfname.png 0 file_123/testfname.png 14 0 0 0 0 0 0 0 tomogram  false 0 0 snapshot  testnotes} {testacquisition false testfname.mp4 0 rawdata/testfname.mp4 14 0 0 0 0 0 0 0 tomogram  false 0 0 tiltSeries  testfname.mp4} { false  0 keyimg_testseries_s.jpg 24 0 0 0 0 0 0 0 image  false 0 0 thumbnail image/jpeg } { false  0 keyimg_testseries.jpg 24 0 0 0 0 0 0 0 tomogram  false 0 0 keyimg image/jpeg } { false  0 keymov_testseries.mp4 1990544 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/mp4 } { false  0 keymov_testseries.flv 6389760 0 0 0 0 0 0 0 tomogram  false 0 0 keymov video/x-flv }]
		//PrettyPrint(pt.Storage.Files)
		//[
		//  {
		//    "dname": "testfname.png",
		//    "fname": "file_123/testfname.png",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "snapshot",
		//    "fNotes": "testnotes"
		//  },
		//  {
		//    "software": "testacquisition",
		//    "dname": "testfname.mp4",
		//    "fname": "rawdata/testfname.mp4",
		//    "fsize": 14,
		//    "type": "tomogram",
		//    "subtype": "tiltSeries",
		//    "fNotes": "testfname.mp4"
		//  },
		//  {
		//    "fname": "keyimg_testseries_s.jpg",
		//    "fsize": 24,
		//    "type": "image",
		//    "subtype": "thumbnail",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keyimg_testseries.jpg",
		//    "fsize": 24,
		//    "type": "tomogram",
		//    "subtype": "keyimg",
		//    "cType": "image/jpeg"
		//  },
		//  {
		//    "fname": "keymov_testseries.mp4",
		//    "fsize": 1990544,
		//    "type": "tomogram",
		//    "subtype": "keymov",
		//    "cType": "video/mp4"
		//  },
		//  {
		//    "fname": "keymov_testseries.flv",
		//    "fsize": 6389760,
		//    "type": "tomogram",
		//    "subtype": "keymov",
		//    "cType": "video/x-flv"
		//  }
		//]
	}

	loc := hash.Combined
	//fmt.Println(loc) //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx
	//PrettyPrint(hash)
	//{
	//  "d": "QmfDatAAh1KzPuFVsAaHo1VS7b4WPeHo4CpwAXS3U2xsWE",
	//  "k": "QmQ7bMwCYZtvorCgaQsYa1hRRtC9xvsJ81jGGpGBG3kZG9",
	//  "c": "QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx"
	//}
	//fmt.Println(capDir) //空
	if capDir != "" {
		loc = hash.Caps
	}
	//fmt.Println(loc) //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx
	//fmt.Println(floAddress) //oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ
	//fmt.Println(strconv.FormatInt(ts, 10)) //1605364418
	v := []string{loc, floAddress, strconv.FormatInt(ts, 10)}
	//fmt.Println(v) //[QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ 1605364418]
	preImage := strings.Join(v, "-")
	//fmt.Println(preImage) //QmfJxwEBCbfe5SRQpP1T1jaJrCLMuBSwP8FgpeV54rsPLx-oYidNmhhCZ76BEKyZJdiqu7YFAURRfSbCJ-1605364418
	signature, err := signMessage(floAddress, preImage)
	//fmt.Println(signature) //{IMkFK0rmu/jyfSe68eCroSPoCELj/MWa94L8EPgsQLZGLH9H7CH4HeaTTwpcJe5Cilv/oqZkuAcf472S2hI4rac= <nil> 0xc42023f360}
	//fmt.Println(err) //<nil>
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
