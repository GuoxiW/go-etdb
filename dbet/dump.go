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
	fmt.Println(string(out))
	if err != nil && !strings.HasSuffix(string(out), "already exists. Exiting.\n") {
		return err
	}
	return nil
}

func processFiles(row TiltSeries) (ipfsHash, error) {
	h := ipfsHash{}
	//fmt.Println(h) // {   }
	//fmt.Println(row) // {testseries testtitle 2020-01-01 00:00:00 +0000 UTC testnotes testscope testroles testnotes testsname testnotes teststrain 1 1 0.1 0.2 0.3 0 0.4 2 0.1 testacquisition testprocess testemdb 0 0 testuname   [{2dimage testfname testnotes testtdimage tomogram snapshot /home/guoxi/blockchain/tomography/data/testseries/file_123/testfname 123 0 }] [{rawdata testfname testfname tomogram tiltSeries /home/guoxi/blockchain/tomography/data/testseries/rawdata/testfname 123 testacquisition}]}
	//fmt.Println(row.Id) // testseries
	s, err := ipfsPinPath("/home/guoxi/blockchain/tomography/data/"+row.Id, row.Id)
	if err != nil {
		return h, err
	}
	h.Data = s

	km := "keymov_" + row.Id
	if row.KeyMov > 0 && row.KeyMov <= 4 {
		flv := "/home/guoxi/blockchain/tomography/data/" + row.Id + "/" + km + ".flv"
		mp4 := "/home/guoxi/blockchain/tomography/data/Videos/" + km + ".mp4"

		err := convertVideo(flv, mp4)
		if err != nil {
			return h, err
		}
		s, err := ipfsPinPath(mp4, km+".mp4")
		if err != nil {
			return h, err
		}
		h.KeyMov = s
	} else {
		h.KeyMov = "n/a"
	}

	if h.KeyMov == "n/a" {
		h.Combined = h.Data
	} else {
		nh, err := ipfsAddLink(h.Data, km+".mp4", h.KeyMov)
		if err != nil {
			return h, err
		}
		h.Combined = nh
	}

	return h, nil
}

func tiltIdToPublishTomogram(tiltSeriesId string) (oip042.PublishTomogram, error) {
	tsr, err := GetTiltSeriesById(tiltSeriesId)  // 获取序列
	//fmt.Println(tsr) //{testseries testtitle 2020-01-01 00:00:00 +0000 UTC testnotes testscope testroles testnotes testsname testnotes teststrain 1 1 0.1 0.2 0.3 0 0.4 2 0.1 testacquisition testprocess testemdb 0 0 testuname   [{2dimage testfname.png testnotes testtdimage tomogram snapshot /home/guoxi/blockchain/tomography/data/testseries/file_123/testfname.png 123 0 }] [{rawdata testfname.mp4 testfname.mp4 tomogram tiltSeries /home/guoxi/blockchain/tomography/data/testseries/rawdata/testfname.mp4 123 testacquisition}]}
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
	//      "FilePath": "/home/guoxi/blockchain/tomography/data/testseries/file_123/testfname.png",
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
	//      "FilePath": "/home/guoxi/blockchain/tomography/data/testseries/rawdata/testfname.mp4",
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
	floAddress := config.FloAddress

	pt = oip042.PublishTomogram{
		PublishArtifact: oip042.PublishArtifact{
			Type:       "research",
			SubType:    "tomogram",
			Timestamp:  ts,
			FloAddress: floAddress,
			Info: &oip042.ArtifactInfo{
				Title:       tsr.Title,
				Description: "Auto imported from etdb",
				Tags:        "etdb,jensen.lab,tomogram,electron.tomography",
			},
			Storage: &oip042.ArtifactStorage{
				Network:  "ipfs",
				Location: hash.Combined,
				Files:    []oip042.ArtifactFiles{},
			},
			Payment: nil, // it's free
		},
		TomogramDetails: oip042.TomogramDetails{
			Microscopist:   tsr.Microscopist,
			Institution:    "Caltech",
			Lab:            "Jensen Lab",
			Sid:            tsr.Id,
			Magnification:  tsr.Magnification,
			Defocus:        tsr.Defocus,
			Dosage:         tsr.Dosage,
			TiltConstant:   tsr.TiltConstant,
			TiltMin:        tsr.TiltMin,
			TiltMax:        tsr.TiltMax,
			TiltStep:       tsr.TiltStep,
			Strain:         tsr.SpeciesStrain,
			SpeciesName:    tsr.SpeciesName,
			ScopeName:      tsr.ScopeName,
			Date:           tsr.Date.Unix(),
			Emdb:           tsr.Emdb,
			TiltSingleDual: tsr.SingleDual,
			NCBItaxID:      tsr.SpeciesTaxId,
			// ToDo: Needs database cleanup before publishing Roles
			//Roles:        tsr.Roles,
		},
	}

	if len(tsr.ScopeNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Scope notes: " + tsr.ScopeNotes + "\n"
	}
	if len(tsr.SpeciesNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Species notes: " + tsr.SpeciesNotes + "\n"
	}
	if len(tsr.TiltSeriesNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Tilt series notes: " + tsr.TiltSeriesNotes + "\n"
	}

	capDir := ""
	for _, df := range tsr.DataFiles {
		fName := strings.TrimPrefix(df.FilePath, "/home/guoxi/blockchain/tomography/data/"+tsr.Id+"/")  // 返回不含前缀字符的 df.FilePath
		if df.Auto == 2 {
			if capDir == "" {
				capDir, err = ipfsNewUnixFsDir()
				if err != nil {
					return pt, err
				}
			}
			h, err := ipfsPinPath(df.FilePath, df.Filename)
			if err != nil {
				return pt, err
			}
			capDir, err = ipfsAddLink(capDir, df.Filename, h)
			if err != nil {
				return pt, err
			}
			fName =  "AutoCaps/" + strings.TrimPrefix(df.FilePath, "/home/guoxi/blockchain/tomography/data/Caps/")  // 返回不含前缀字符的 df.FilePath
		}

		fi, err := os.Stat(df.FilePath)  // 获取文件属性
		if err != nil {
			return pt, err
		}
		af := oip042.ArtifactFiles{
			Type:    df.Type,
			SubType: df.SubType,
			FNotes:  df.Notes,
			Fsize:   fi.Size(),
			Dname:   df.Filename,
			Fname:   fName,
		}
		pt.Storage.Files = append(pt.Storage.Files, af)
	}

	if capDir != "" {
		hash.Caps, err = ipfsAddLink(hash.Combined, "AutoCaps", capDir)
		if err != nil {
			return pt, err
		}
		pt.Storage.Location = hash.Caps
		ipfsHashes[tsr.Id] = hash
		saveIpfsHashes()
	}

	for _, tdf := range tsr.ThreeDFiles {
		fi, err := os.Stat(tdf.FilePath)
		if err != nil {
			return pt, err
		}
		af := oip042.ArtifactFiles{
			Type:     tdf.Type,
			SubType:  tdf.SubType,
			FNotes:   tdf.Notes,
			Fsize:    fi.Size(),
			Dname:    tdf.Filename,
			Fname:    strings.TrimPrefix(tdf.FilePath, "/home/guoxi/blockchain/tomography/data/"+tsr.Id+"/"),
			Software: tdf.Software,
		}
		pt.Storage.Files = append(pt.Storage.Files, af)
	}

	if tsr.KeyImg > 0 && tsr.KeyImg <= 4 {
		kif := "keyimg_" + tsr.Id + "_s.jpg"
		fi, err := os.Stat("/home/guoxi/blockchain/tomography/data/" + tsr.Id + "/" + kif)
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
		pt.Storage.Files = append(pt.Storage.Files, ki)

		kif = "keyimg_" + tsr.Id + ".jpg"
		fi, err = os.Stat("/home/guoxi/blockchain/tomography/data/" + tsr.Id + "/" + kif)
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
		pt.Storage.Files = append(pt.Storage.Files, ki)
	}
	if tsr.KeyMov > 0 && tsr.KeyMov <= 4 {
		kmf := "keymov_" + tsr.Id + ".mp4"
		fi, err := os.Stat("/home/guoxi/blockchain/tomography/data/Videos/" + kmf)  // 获取文件属性
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
		pt.Storage.Files = append(pt.Storage.Files, km)
		kmf = "keymov_" + tsr.Id + ".flv"
		fi, err = os.Stat("/home/guoxi/blockchain/tomography/data/" + tsr.Id + "/" + kmf)  // 获取文件属性
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
		pt.Storage.Files = append(pt.Storage.Files, km)
	}

	loc := hash.Combined
	if capDir != "" {
		loc = hash.Caps
	}
	v := []string{loc, floAddress, strconv.FormatInt(ts, 10)}
	preImage := strings.Join(v, "-")
	signature, err := signMessage(floAddress, preImage)
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
