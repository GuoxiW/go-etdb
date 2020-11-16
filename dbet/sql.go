// sql 相关的函数。

package main

import (
	"database/sql" // sql 包提供了保证SQL或类SQL数据库的泛用接口
	"errors"
	"io/ioutil"
	"regexp" // 正则表达式
	"strconv"
	"strings"
	"github.com/go-sql-driver/mysql" // mysql
	"github.com/jmoiron/sqlx"        // sql 通用拓展
)

var (
	dbh *sqlx.DB
	// sql queries loaded from files
	selectTiltSeriesSummarySql string
	selectDataFilesSql         string
	selectThreeDFilesSql       string
	selectFilterSql            string
)

func init() {
	buf, err := ioutil.ReadFile("./sql/selectTiltSeriesSummary.sql")
	if err != nil {
		panic(err)
	}
	selectTiltSeriesSummarySql = string(buf) // sql文件原文读进来

	buf, err = ioutil.ReadFile("./sql/selectDataFiles.sql")
	if err != nil {
		panic(err)
	}
	selectDataFilesSql = string(buf) // sql文件原文读进来

	buf, err = ioutil.ReadFile("./sql/selectThreeDFiles.sql")
	if err != nil {
		panic(err)
	}
	selectThreeDFilesSql = string(buf) // sql文件原文读进来

	buf, err = ioutil.ReadFile("./sql/filter.sql")
	if err != nil {
		panic(err)
	}
	selectFilterSql = string(buf) // sql文件原文读进来

	conf := mysql.NewConfig()
	conf.User = config.DatabaseConfiguration.User       // 读取 config.go 中的数据
	conf.Passwd = config.DatabaseConfiguration.Password // 读取 config.go 中的数据
	conf.Net = config.DatabaseConfiguration.Net         // 读取 config.go 中的数据
	conf.Addr = config.DatabaseConfiguration.Address    // 读取 config.go 中的数据
	conf.DBName = config.DatabaseConfiguration.Name     // 读取 config.go 中的数据

	newDb, err := sqlx.Connect("mysql", conf.FormatDSN()) // 连接 sql
	//fmt.Println(newDb) //&{0xc42010caa0 mysql false 0xc4201baff0}
	//fmt.Println(err) //<nil>
	if err != nil {
		panic(err)
	}
	dbh = newDb // 将连接到的publicdb命名为dbh
	//fmt.Println(dbh) //&{0xc42010caa0 mysql false 0xc4201baff0}
}

type tiltSeriesRow struct {
	TiltSeriesID        sql.NullString  `db:"tiltSeriesID"`
	Title               sql.NullString  `db:"title"`
	TomoDate            mysql.NullTime  `db:"tomo_date"`
	TsdTXTNotes         sql.NullString  `db:"tsd_TXT_notes"`
	Scope               sql.NullString  `db:"scope"`
	Roles               sql.NullString  `db:"roles"`
	ScdTXTNotes         sql.NullString  `db:"scd_TXT_notes"`
	SpeciesName         sql.NullString  `db:"SpeciesName"`
	SpdTXTNotes         sql.NullString  `db:"spd_TXT_notes"`
	Strain              sql.NullString  `db:"strain"`
	TaxId               sql.NullInt64   `db:"tax_id"`
	SingleDual          sql.NullInt64   `db:"single_dual"`
	Defocus             sql.NullFloat64 `db:"defocus"`
	Magnification       sql.NullFloat64 `db:"magnification"`
	Dosage              sql.NullFloat64 `db:"dosage"`
	TiltConstant        sql.NullFloat64 `db:"tilt_constant"`
	TiltMin             sql.NullFloat64 `db:"tilt_min"`
	TiltMax             sql.NullFloat64 `db:"tilt_max"`
	TiltStep            sql.NullString  `db:"tilt_step"`
	SoftwareAcquisition sql.NullString  `db:"software_acquisition"`
	SoftwareProcess     sql.NullString  `db:"software_process"`
	Emdb                sql.NullString  `db:"emdb"`
	KeyImg              sql.NullInt64   `db:"keyimg"`
	KeyMov              sql.NullInt64   `db:"keymov"`
	FullName            sql.NullString  `db:"fullname"`
}

type dataFileRow struct {
	Filetype        sql.NullString `db:"filetype"`
	Filename        sql.NullString `db:"filename"`
	Notes           sql.NullString `db:"TXT_notes"`
	ThreeDFileImage sql.NullString `db:"ThreeDFileImage"`
	DefId           sql.NullInt64  `db:"DEF_id"`
	Auto            sql.NullInt64  `db:"auto"`
}

type threeDFileRow struct {
	Classify sql.NullString `db:"classify"`
	Notes    sql.NullString `db:"TXT_notes"`
	Filename sql.NullString `db:"filename"`
	DefId    sql.NullInt64  `db:"DEF_id"`
}

var extractTiltStepRe = regexp.MustCompile(`^[0-9.]+`) // 正则表达式匹配 ^ 表示补集、 + 表示一次或多次匹配

func GetTiltSeriesById(tiltSeriesId string) (ts TiltSeries, err error) { // 通过id获得ts数据。
	var tsr tiltSeriesRow
	err = dbh.Get(&tsr, selectTiltSeriesSummarySql, 1, tiltSeriesId) // 给定tsid, 用sql语句提取tsrow结构体
	//fmt.Println(tsr) //{{testseries true} {testtitle true} {2020-01-01 00:00:00 +0000 UTC true} {testnotes true} {testscope true} {testroles true} {testnotes true} {testsname true} {testnotes true} {teststrain true} {1 true} {1 true} {0.1 true} {0.2 true} {0.3 true} {0 true} {0.4 true} {2 true} {0.1 true} {testacquisition true} {testprocess true} {testemdb true} {1 true} {1 true} {testuname true}}
	//PrettyPrint(tsr)
	//{
	//  "TiltSeriesID": {
	//    "String": "testseries",
	//    "Valid": true
	//  },
	//  "Title": {
	//    "String": "testtitle",
	//    "Valid": true
	//  },
	//  "TomoDate": {
	//    "Time": "2020-01-01T00:00:00Z",
	//    "Valid": true
	//  },
	//  "TsdTXTNotes": {
	//    "String": "testnotes",
	//    "Valid": true
	//  },
	//  "Scope": {
	//    "String": "testscope",
	//    "Valid": true
	//  },
	//  "Roles": {
	//    "String": "testroles",
	//    "Valid": true
	//  },
	//  "ScdTXTNotes": {
	//    "String": "testnotes",
	//    "Valid": true
	//  },
	//  "SpeciesName": {
	//    "String": "testsname",
	//    "Valid": true
	//  },
	//  "SpdTXTNotes": {
	//    "String": "testnotes",
	//    "Valid": true
	//  },
	//  "Strain": {
	//    "String": "teststrain",
	//    "Valid": true
	//  },
	//  "TaxId": {
	//    "Int64": 1,
	//    "Valid": true
	//  },
	//  "SingleDual": {
	//    "Int64": 1,
	//    "Valid": true
	//  },
	//  "Defocus": {
	//    "Float64": 0.1,
	//    "Valid": true
	//  },
	//  "Magnification": {
	//    "Float64": 0.2,
	//    "Valid": true
	//  },
	//  "Dosage": {
	//    "Float64": 0.3,
	//    "Valid": true
	//  },
	//  "TiltConstant": {
	//    "Float64": 0,
	//    "Valid": true
	//  },
	//  "TiltMin": {
	//    "Float64": 0.4,
	//    "Valid": true
	//  },
	//  "TiltMax": {
	//    "Float64": 2,
	//    "Valid": true
	//  },
	//  "TiltStep": {
	//    "String": "0.1",
	//    "Valid": true
	//  },
	//  "SoftwareAcquisition": {
	//    "String": "testacquisition",
	//    "Valid": true
	//  },
	//  "SoftwareProcess": {
	//    "String": "testprocess",
	//    "Valid": true
	//  },
	//  "Emdb": {
	//    "String": "testemdb",
	//    "Valid": true
	//  },
	//  "KeyImg": {
	//    "Int64": 1,
	//    "Valid": true
	//  },
	//  "KeyMov": {
	//    "Int64": 1,
	//    "Valid": true
	//  },
	//  "FullName": {
	//    "String": "testuname",
	//    "Valid": true
	//  }
	//}
	if err != nil {
		return
	}

	// 命名函数返回的 ts
	if tsr.TiltSeriesID.Valid { // testseries
		ts.Id = tsr.TiltSeriesID.String // 检验后命名Id
		//fmt.Println(ts.Id)
	} else {
		return ts, errors.New("tiltSeriesId returned no result")
	}
	if tsr.Title.Valid { // testtitle
		ts.Title = tsr.Title.String // 检验后命名title
		//fmt.Println(ts.Title)
	}
	if tsr.SpeciesName.Valid { // testsname
		ts.SpeciesName = tsr.SpeciesName.String // 检验后命名speciesname
		//fmt.Println(ts.SpeciesName)
	}
	if len(ts.Title) == 0 { // testtitle 所以不会运行
		ts.Title = ts.SpeciesName // 如果没有title, 用speciesname代替title
		//fmt.Println(ts.Title)
	}
	if tsr.TomoDate.Valid { // 2020-01-01 00:00:00 +0000 UTC
		ts.Date = tsr.TomoDate.Time // 检验后命名date
		//fmt.Println(ts.Date)
	}
	if tsr.TsdTXTNotes.Valid { // testnotes
		ts.TiltSeriesNotes = tsr.TsdTXTNotes.String // 检验后命名tsnotes
		//fmt.Println(ts.TiltSeriesNotes)
	}
	if tsr.Scope.Valid { // testscope
		ts.ScopeName = tsr.Scope.String // 检验后命名scopename
		//fmt.Println(ts.ScopeName)
	}
	if tsr.Roles.Valid { // testroles
		ts.Roles = tsr.Roles.String // 检验后命名roles
		//fmt.Println(ts.Roles)
	}
	if tsr.ScdTXTNotes.Valid { // testnotes
		ts.ScopeNotes = tsr.ScdTXTNotes.String // 检验后命名scopenotes
		//fmt.Println(ts.ScopeNotes)
	}
	if tsr.SpdTXTNotes.Valid { // testnotes
		ts.SpeciesNotes = tsr.SpdTXTNotes.String // 检验后命名speciesnotes
		//fmt.Println(ts.SpeciesNotes)
	}
	if tsr.Strain.Valid { // teststrain
		ts.SpeciesStrain = tsr.Strain.String // 检验后命名speciesstrain
		//fmt.Println(ts.SpeciesStrain)
	}
	if tsr.TaxId.Valid { // 1
		ts.SpeciesTaxId = tsr.TaxId.Int64 // 检验后命名speciestaxid
		//fmt.Println(ts.SpeciesTaxId)
	}
	if tsr.SingleDual.Valid { // 1
		ts.SingleDual = tsr.SingleDual.Int64 // 检验后命名singledual
		//fmt.Println(ts.SingleDual)
	}
	if tsr.Defocus.Valid { // 0.1
		ts.Defocus = tsr.Defocus.Float64 // 检验后命名defocus
		//fmt.Println(ts.Defocus)
	}
	if tsr.Magnification.Valid { // 0.2
		ts.Magnification = tsr.Magnification.Float64 // 检验后命名magnification
		//fmt.Println(ts.Magnification)
	}
	if tsr.Dosage.Valid { // 0.3
		ts.Dosage = tsr.Dosage.Float64 // 检验后命名dosage
		//fmt.Println(ts.Dosage)
	}
	if tsr.TiltConstant.Valid { // 0
		ts.TiltConstant = tsr.TiltConstant.Float64 // 检验后命名tiltconstant
		//fmt.Println(ts.TiltConstant)
	}
	if tsr.TiltMin.Valid { // 0.4
		ts.TiltMin = tsr.TiltMin.Float64 // 检验后命名tiltmin
		//fmt.Println(ts.TiltMin)
	}
	if tsr.TiltMax.Valid { // 2.0
		ts.TiltMax = tsr.TiltMax.Float64 // 检验后命名tiltmax
		//fmt.Println(ts.TiltMax)
	}
	if tsr.TiltStep.Valid {
		tss := tsr.TiltStep.String                                                 // 0.1
		ts.TiltStep, _ = strconv.ParseFloat(extractTiltStepRe.FindString(tss), 64) // 转换为 float64 // 0.1
		//fmt.Println(tss)
		//fmt.Println(ts.TiltStep)
	}
	if tsr.SoftwareAcquisition.Valid { // testacquisition
		ts.SoftwareAcquisition = tsr.SoftwareAcquisition.String // 检验后命名softwareacquisition
		//fmt.Println(ts.SoftwareAcquisition)
	}
	if tsr.SoftwareProcess.Valid { // testprocess
		ts.SoftwareProcess = tsr.SoftwareProcess.String // 检验后命名softwareprocess
		//fmt.Println(ts.SoftwareProcess)
	}
	if tsr.Emdb.Valid { // testemdb
		ts.Emdb = tsr.Emdb.String // 检验后命名emdb
		//fmt.Println(ts.Emdb)
	}
	if tsr.KeyMov.Valid { // 1
		ts.KeyMov = tsr.KeyMov.Int64 // 检验后命名keymov
		//fmt.Println(ts.KeyMov)
	}
	if tsr.KeyImg.Valid { // 1
		ts.KeyImg = tsr.KeyImg.Int64 // 检验后命名keyimg
		//fmt.Println(ts.KeyImg)
	}
	if tsr.FullName.Valid { // testuname
		ts.Microscopist = tsr.FullName.String // 检验后命名microscopist
		//fmt.Println(ts.Microscopist)
	}

	rows, err := dbh.Queryx(selectDataFilesSql, tiltSeriesId) // datafile 的 sql语句来查询 tsid
	//fmt.Println(rows) //&{0xc420254280 false 0xc4201b8ff0 false [] []}

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var dfr dataFileRow
		err = rows.StructScan(&dfr) // 扫描结构体
		//fmt.Println(dfr) //{{2dimage true} {testfname.png true} {testnotes true} {testtdimage true} {123 true} {0 true}}
		//PrettyPrint(dfr)
		//{
		//  "Filetype": {
		//    "String": "2dimage",
		//    "Valid": true
		//  },
		//  "Filename": {
		//    "String": "testfname.png",
		//    "Valid": true
		//  },
		//  "Notes": {
		//    "String": "testnotes",
		//    "Valid": true
		//  },
		//  "ThreeDFileImage": {
		//    "String": "testtdimage",
		//    "Valid": true
		//  },
		//  "DefId": {
		//    "Int64": 123,
		//    "Valid": true
		//  },
		//  "Auto": {
		//    "Int64": 0,
		//    "Valid": true
		//  }
		//}
		if err != nil {
			return
		}
		df := DataFile{}
		if dfr.Filename.Valid { // testfname.png
			df.Filename = dfr.Filename.String // 检验命名filename
			//fmt.Println(df.Filename)
			if len(strings.TrimSpace(df.Filename)) == 0 {
				// No file name, no file...
				continue
			}
		} else {
			// No file name, no file...
			continue
		}
		if dfr.Filetype.Valid { // 2dimage
			df.Filetype = dfr.Filetype.String // 检验命名filetype
			//fmt.Println(df.Filetype)
		}
		if dfr.ThreeDFileImage.Valid { // testtdimage
			df.ThreeDFileImage = dfr.ThreeDFileImage.String // 检验命名threedfileimage
			//fmt.Println(df.ThreeDFileImage)
		}
		if dfr.Notes.Valid { // testnotes
			df.Notes = dfr.Notes.String // 检验命名notes
			//fmt.Println(df.Notes)
		}
		if dfr.DefId.Valid { // 123
			df.DefId = dfr.DefId.Int64 // 检验命名defid
			//fmt.Println(df.DefId)
		}
		if dfr.Auto.Valid { // 0
			df.Auto = dfr.Auto.Int64 // 检验命名auto
			//fmt.Println(df.Auto)
		}
		df.Type = "tomogram"
		//fmt.Println(df.Filetype) //2dimage
		switch df.Filetype { // 根据filetype选择
		case "2dimage":
			df.SubType = "snapshot"
			if df.Auto == 2 {
				df.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/Caps/" + df.Filename
			} else {
				df.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tiltSeriesId + "/file_" + strconv.FormatInt(df.DefId, 10) + "/" + df.Filename
			}
		case "movie":
			df.SubType = "preview"
			df.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tiltSeriesId + "/file_" + strconv.FormatInt(df.DefId, 10) + "/" + df.Filename
		case "other":
			df.SubType = "other"
			df.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tiltSeriesId + "/file_" + strconv.FormatInt(df.DefId, 10) + "/" + df.Filename
		default:
			panic("Unknown new DataFile.FileType " + df.Filetype + " from DEF_id " + strconv.FormatInt(df.DefId, 10))
		}
		//PrettyPrint(ts.DataFiles) null
		ts.DataFiles = append(ts.DataFiles, df)
		//PrettyPrint(ts.DataFiles)
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
	}
	err = rows.Err()
	if err != nil {
		return
	}

	rows, err = dbh.Queryx(selectThreeDFilesSql, tiltSeriesId) // 3dfiles 的 sql语句来查询 tsid
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tdfr threeDFileRow
		err = rows.StructScan(&tdfr) // 扫描结构体
		//PrettyPrint(tdfr)
		//{
		//  "Classify": {
		//    "String": "rawdata",
		//    "Valid": true
		//  },
		//  "Notes": {
		//    "String": "testnotes",
		//    "Valid": true
		//  },
		//  "Filename": {
		//    "String": "testfname.mp4",
		//    "Valid": true
		//  },
		//  "DefId": {
		//    "Int64": 123,
		//    "Valid": true
		//  }
		//}
		if err != nil {
			return
		}
		tdf := ThreeDFile{}
		if tdfr.Filename.Valid { // testfname.mp4
			tdf.Filename = tdfr.Filename.String // 检验后命名filename
			//fmt.Println(tdf.Filename)
		}
		if tdfr.Classify.Valid { // rawdata
			tdf.Classify = tdfr.Classify.String // 检验后命名classify  文件类型
			//fmt.Println(tdf.Classify)
		}
		if tdfr.Notes.Valid { // testfname.mp4
			tdf.Notes = tdfr.Filename.String // 检验后命名notes
			//fmt.Println(tdf.Notes)
		}
		if tdfr.DefId.Valid { // 123
			tdf.DefId = tdfr.DefId.Int64 // 检验后命名defid
			//fmt.Println(tdf.DefId)
		}
		tdf.Type = "tomogram"
		//fmt.Println(tdf.Classify) //rawdata
		switch tdf.Classify {
		case "rawdata":
			tdf.SubType = "tiltSeries"
			//fmt.Println(tdf.Software) //空
			if !strings.Contains(ts.SoftwareAcquisition, ",") {
				tdf.Software = ts.SoftwareAcquisition
				//fmt.Println(tdf.Software) //testacquisition
			}
			tdf.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tiltSeriesId + "/rawdata/" + tdf.Filename
		case "reconstruction":
			tdf.SubType = "reconstruction"
			if !strings.Contains(ts.SoftwareProcess, ",") {
				tdf.Software = ts.SoftwareProcess
			}
			tdf.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tiltSeriesId + "/3dimage_" + strconv.FormatInt(tdf.DefId, 10) + "/" + tdf.Filename
		case "subvolume": // 子卷
			fallthrough // fallthrough 会强制执行后面 case 的代码,不管 case 是 true 还是 false, 就是 other 会默认执行。
		case "other":
			tdf.SubType = tdf.Classify
			tdf.FilePath = "/home/guoxi/snap/ipfs/blockchain/tomography/data/" + tiltSeriesId + "/3dimage_" + strconv.FormatInt(tdf.DefId, 10) + "/" + tdf.Filename
		default:
			panic("Unknown new DataFile.FileType " + tdf.Classify + " from DEF_id " + strconv.FormatInt(tdf.DefId, 10))
		}

		ts.ThreeDFiles = append(ts.ThreeDFiles, tdf)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func GetFilterIdList() ([]string, error) {
	var ids []string

	err := dbh.Select(&ids, selectFilterSql) // 根据 filter.sql 中的设置选择id
	//fmt.Println(ids) //[testseries]
	if err != nil {
		return nil, err
	}

	return ids, nil
}
