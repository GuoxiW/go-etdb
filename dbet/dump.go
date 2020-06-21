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
	ids, err := GetFilterIdList()  // 获得 id 列表
	if err != nil { // := 是声明并赋值
		panic(err)
	}

	for _, id := range ids {
		if _, ok := history[id]; ok {
			fmt.Printf("Tilt %s already published\n", id)  // 判断是否已经发布
			continue
		}

		pt, err := tiltIdToPublishTomogram(id)  // 发布
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
			ids, err := sendToBlockchain("json:" + string(min))
			if err != nil {
				fmt.Println(ids)
				panic(err)
			} else {
				history[id] = ids
				PrettyPrint(ids) // json 格式处理
			}
		}

		err = saveHistory()
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
	s, err := ipfsPinPath("/services/tomography/data/"+row.Id, row.Id)
	if err != nil {
		return h, err
	}
	h.Data = s

	km := "keymov_" + row.Id
	if row.KeyMov > 0 && row.KeyMov <= 4 {
		flv := "/services/tomography/data/" + row.Id + "/" + km + ".flv"
		mp4 := "/services/tomography/data/Videos/" + km + ".mp4"

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
	if err != nil {
		panic(err)
	}

	//PrettyPrint(tsr)
	var pt oip042.PublishTomogram

	hash, ok := ipfsHashes[tiltSeriesId]
	emptyDir := false
	if ok {
		emptyDir, err = containsEmptyFolder(hash.Data)  // 判断是否空文件夹
		if err != nil {
			return pt, err
		}
	}
	if !ok || hash.Data == "" || hash.KeyMov == "" || hash.Combined == "" || emptyDir {
		hash, err = processFiles(tsr)
		if err != nil {
			return pt, err
		}
		ipfsHashes[tiltSeriesId] = hash
		saveIpfsHashes()
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
		fName := strings.TrimPrefix(df.FilePath, "/services/tomography/data/"+tsr.Id+"/")
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
			fName =  "AutoCaps/" + strings.TrimPrefix(df.FilePath, "/services/tomography/data/Caps/")
		}

		fi, err := os.Stat(df.FilePath)
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
			Fname:    strings.TrimPrefix(tdf.FilePath, "/services/tomography/data/"+tsr.Id+"/"),
			Software: tdf.Software,
		}
		pt.Storage.Files = append(pt.Storage.Files, af)
	}

	if tsr.KeyImg > 0 && tsr.KeyImg <= 4 {
		kif := "keyimg_" + tsr.Id + "_s.jpg"
		fi, err := os.Stat("/services/tomography/data/" + tsr.Id + "/" + kif)
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
		fi, err = os.Stat("/services/tomography/data/" + tsr.Id + "/" + kif)
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
		fi, err := os.Stat("/services/tomography/data/Videos/" + kmf)
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
		fi, err = os.Stat("/services/tomography/data/" + tsr.Id + "/" + kmf)
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
