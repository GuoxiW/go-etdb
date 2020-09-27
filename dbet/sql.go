// sql 相关的函数。

package main

import (
	"database/sql" // sql 包提供了保证SQL或类SQL数据库的泛用接口
	"io/ioutil"

	"errors"
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
	if err != nil {
		panic(err)
	}
	dbh = newDb // 将连接到的publicdb命名为dbh
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
	if err != nil {
		return
	}

	// 命名函数返回的 ts
	if tsr.TiltSeriesID.Valid {
		ts.Id = tsr.TiltSeriesID.String  // 检验后命名Id
	} else {
		return ts, errors.New("tiltSeriesId returned no result")
	}
	if tsr.Title.Valid {
		ts.Title = tsr.Title.String // 检验后命名title
	}
	if tsr.SpeciesName.Valid {
		ts.SpeciesName = tsr.SpeciesName.String // 检验后命名speciesname
	}
	if len(ts.Title) == 0 {
		ts.Title = ts.SpeciesName // 如果没有title, 用speciesname代替title
	}
	if tsr.TomoDate.Valid {
		ts.Date = tsr.TomoDate.Time // 检验后命名date
	}
	if tsr.TsdTXTNotes.Valid {
		ts.TiltSeriesNotes = tsr.TsdTXTNotes.String // 检验后命名tsnotes
	}
	if tsr.Scope.Valid {
		ts.ScopeName = tsr.Scope.String // 检验后命名scopename
	}
	if tsr.Roles.Valid {
		ts.Roles = tsr.Roles.String // 检验后命名roles
	}
	if tsr.ScdTXTNotes.Valid {
		ts.ScopeNotes = tsr.ScdTXTNotes.String // 检验后命名scopenotes
	}
	if tsr.SpdTXTNotes.Valid {
		ts.SpeciesNotes = tsr.SpdTXTNotes.String // 检验后命名speciesnotes
	}
	if tsr.Strain.Valid {
		ts.SpeciesStrain = tsr.Strain.String // 检验后命名speciesstrain
	}
	if tsr.TaxId.Valid {
		ts.SpeciesTaxId = tsr.TaxId.Int64 // 检验后命名speciestaxid
	}
	if tsr.SingleDual.Valid {
		ts.SingleDual = tsr.SingleDual.Int64 // 检验后命名singledual
	}
	if tsr.Defocus.Valid {
		ts.Defocus = tsr.Defocus.Float64 // 检验后命名defocus
	}
	if tsr.Magnification.Valid {
		ts.Magnification = tsr.Magnification.Float64 // 检验后命名magnification
	}
	if tsr.Dosage.Valid {
		ts.Dosage = tsr.Dosage.Float64 // 检验后命名dosage
	}
	if tsr.TiltConstant.Valid {
		ts.TiltConstant = tsr.TiltConstant.Float64 // 检验后命名tiltconstant
	}
	if tsr.TiltMin.Valid {
		ts.TiltMin = tsr.TiltMin.Float64 // 检验后命名tiltmin
	}
	if tsr.TiltMax.Valid {
		ts.TiltMax = tsr.TiltMax.Float64 // 检验后命名tiltmax
	}
	if tsr.TiltStep.Valid {
		tss := tsr.TiltStep.String
		ts.TiltStep, _ = strconv.ParseFloat(extractTiltStepRe.FindString(tss), 64) // 转换为 float64
	}
	if tsr.SoftwareAcquisition.Valid {
		ts.SoftwareAcquisition = tsr.SoftwareAcquisition.String // 检验后命名softwareacquisition
	}
	if tsr.SoftwareProcess.Valid {
		ts.SoftwareProcess = tsr.SoftwareProcess.String // 检验后命名softwareprocess
	}
	if tsr.Emdb.Valid {
		ts.Emdb = tsr.Emdb.String // 检验后命名emdb
	}
	if tsr.KeyMov.Valid {
		ts.KeyMov = tsr.KeyMov.Int64 // 检验后命名keymov
	}
	if tsr.KeyImg.Valid {
		ts.KeyImg = tsr.KeyImg.Int64 // 检验后命名keyimg
	}
	if tsr.FullName.Valid {
		ts.Microscopist = tsr.FullName.String // 检验后命名microscopist
	}

	rows, err := dbh.Queryx(selectDataFilesSql, tiltSeriesId) // datafile 的 sql语句来查询 tsid
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var dfr dataFileRow
		err = rows.StructScan(&dfr) // 扫描结构体
		if err != nil {
			return
		}
		df := DataFile{}
		if dfr.Filename.Valid {
			df.Filename = dfr.Filename.String // 检验命名filename
			if len(strings.TrimSpace(df.Filename)) == 0 {
				// No file name, no file...
				continue
			}
		} else {
			// No file name, no file...
			continue
		}
		if dfr.Filetype.Valid {
			df.Filetype = dfr.Filetype.String // 检验命名filetype
		}
		if dfr.ThreeDFileImage.Valid {
			df.ThreeDFileImage = dfr.ThreeDFileImage.String // 检验命名threedfileimage
		}
		if dfr.Notes.Valid {
			df.Notes = dfr.Notes.String // 检验命名notes
		}
		if dfr.DefId.Valid {
			df.DefId = dfr.DefId.Int64 // 检验命名defid
		}
		if dfr.Auto.Valid {
			df.Auto = dfr.Auto.Int64 // 检验命名auto
		}
		df.Type = "tomogram"
		switch df.Filetype { // 根据filetype选择
		case "2dimage":
			df.SubType = "snapshot"
			if df.Auto == 2 {
				df.FilePath = "/services/tomography/data/Caps/" + df.Filename
			} else {
				df.FilePath = "/services/tomography/data/" + tiltSeriesId + "/file_" + strconv.FormatInt(df.DefId, 10) + "/" + df.Filename
			}
		case "movie":
			df.SubType = "preview"
			df.FilePath = "/services/tomography/data/" + tiltSeriesId + "/file_" + strconv.FormatInt(df.DefId, 10) + "/" + df.Filename
		case "other":
			df.SubType = "other"
			df.FilePath = "/services/tomography/data/" + tiltSeriesId + "/file_" + strconv.FormatInt(df.DefId, 10) + "/" + df.Filename
		default:
			panic("Unknown new DataFile.FileType " + df.Filetype + " from DEF_id " + strconv.FormatInt(df.DefId, 10))
		}
		ts.DataFiles = append(ts.DataFiles, df)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	rows, err = dbh.Queryx(selectThreeDFilesSql, tiltSeriesId)  // 3dfiles 的 sql语句来查询 tsid
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tdfr threeDFileRow
		err = rows.StructScan(&tdfr) // 扫描结构体
		if err != nil {
			return
		}
		tdf := ThreeDFile{}
		if tdfr.Filename.Valid {
			tdf.Filename = tdfr.Filename.String  // 检验后命名filename
		}
		if tdfr.Classify.Valid {
			tdf.Classify = tdfr.Classify.String  // 检验后命名classify  文件类型
		}
		if tdfr.Notes.Valid {
			tdf.Notes = tdfr.Filename.String  // 检验后命名notes
		}
		if tdfr.DefId.Valid {
			tdf.DefId = tdfr.DefId.Int64  // 检验后命名defid
		}
		tdf.Type = "tomogram"
		switch tdf.Classify {
		case "rawdata":
			tdf.SubType = "tiltSeries"
			if !strings.Contains(ts.SoftwareAcquisition, ",") {
				tdf.Software = ts.SoftwareAcquisition
			}
			tdf.FilePath = "/services/tomography/data/" + tiltSeriesId + "/rawdata/" + tdf.Filename
		case "reconstruction":
			tdf.SubType = "reconstruction"
			if !strings.Contains(ts.SoftwareProcess, ",") {
				tdf.Software = ts.SoftwareProcess
			}
			tdf.FilePath = "/services/tomography/data/" + tiltSeriesId + "/3dimage_" + strconv.FormatInt(tdf.DefId, 10) + "/" + tdf.Filename
		case "subvolume":  // 子卷
			fallthrough  // fallthrough 会强制执行后面 case 的代码,不管 case 是 true 还是 false, 就是 other 会默认执行。
		case "other":
			tdf.SubType = tdf.Classify
			tdf.FilePath = "/services/tomography/data/" + tiltSeriesId + "/3dimage_" + strconv.FormatInt(tdf.DefId, 10) + "/" + tdf.Filename
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

	err := dbh.Select(&ids, selectFilterSql)  // 根据 filter.sql 中的设置选择id
	if err != nil {
		return nil, err
	}

	return ids, nil
}
