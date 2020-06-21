// sql 相关的函数。

package main

import (
	"io/ioutil"
	"database/sql" // sql 包提供了保证SQL或类SQL数据库的泛用接口

	"github.com/jmoiron/sqlx" // sql 通用拓展
	"github.com/go-sql-driver/mysql" // mysql
	"errors"
	"strings"
	"strconv"
	"regexp" // 正则表达式
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
	selectTiltSeriesSummarySql = string(buf)

	buf, err = ioutil.ReadFile("./sql/selectDataFiles.sql")
	if err != nil {
		panic(err)
	}
	selectDataFilesSql = string(buf)

	buf, err = ioutil.ReadFile("./sql/selectThreeDFiles.sql")
	if err != nil {
		panic(err)
	}
	selectThreeDFilesSql = string(buf)

	buf, err = ioutil.ReadFile("./sql/filter.sql")
	if err != nil {
		panic(err)
	}
	selectFilterSql = string(buf)

	conf := mysql.NewConfig()
	conf.User = config.DatabaseConfiguration.User
	conf.Passwd = config.DatabaseConfiguration.Password
	conf.Net = config.DatabaseConfiguration.Net
	conf.Addr = config.DatabaseConfiguration.Address
	conf.DBName = config.DatabaseConfiguration.Name

	newDb, err := sqlx.Connect("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}
	dbh = newDb
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

var extractTiltStepRe = regexp.MustCompile(`^[0-9.]+`)

func GetTiltSeriesById(tiltSeriesId string) (ts TiltSeries, err error) {
	var tsr tiltSeriesRow
	err = dbh.Get(&tsr, selectTiltSeriesSummarySql, 1, tiltSeriesId)
	if err != nil {
		return
	}

	if tsr.TiltSeriesID.Valid {
		ts.Id = tsr.TiltSeriesID.String
	} else {
		return ts, errors.New("tiltSeriesId returned no result")
	}
	if tsr.Title.Valid {
		ts.Title = tsr.Title.String
	}
	if tsr.SpeciesName.Valid {
		ts.SpeciesName = tsr.SpeciesName.String
	}
	if len(ts.Title) == 0 {
		ts.Title = ts.SpeciesName
	}
	if tsr.TomoDate.Valid {
		ts.Date = tsr.TomoDate.Time
	}
	if tsr.TsdTXTNotes.Valid {
		ts.TiltSeriesNotes = tsr.TsdTXTNotes.String
	}
	if tsr.Scope.Valid {
		ts.ScopeName = tsr.Scope.String
	}
	if tsr.Roles.Valid {
		ts.Roles = tsr.Roles.String
	}
	if tsr.ScdTXTNotes.Valid {
		ts.ScopeNotes = tsr.ScdTXTNotes.String
	}
	if tsr.SpdTXTNotes.Valid {
		ts.SpeciesNotes = tsr.SpdTXTNotes.String
	}
	if tsr.Strain.Valid {
		ts.SpeciesStrain = tsr.Strain.String
	}
	if tsr.TaxId.Valid {
		ts.SpeciesTaxId = tsr.TaxId.Int64
	}
	if tsr.SingleDual.Valid {
		ts.SingleDual = tsr.SingleDual.Int64
	}
	if tsr.Defocus.Valid {
		ts.Defocus = tsr.Defocus.Float64
	}
	if tsr.Magnification.Valid {
		ts.Magnification = tsr.Magnification.Float64
	}
	if tsr.Dosage.Valid {
		ts.Dosage = tsr.Dosage.Float64
	}
	if tsr.TiltConstant.Valid {
		ts.TiltConstant = tsr.TiltConstant.Float64
	}
	if tsr.TiltMin.Valid {
		ts.TiltMin = tsr.TiltMin.Float64
	}
	if tsr.TiltMax.Valid {
		ts.TiltMax = tsr.TiltMax.Float64
	}
	if tsr.TiltStep.Valid {
		tss := tsr.TiltStep.String
		ts.TiltStep, _ = strconv.ParseFloat(extractTiltStepRe.FindString(tss), 64)
	}
	if tsr.SoftwareAcquisition.Valid {
		ts.SoftwareAcquisition = tsr.SoftwareAcquisition.String
	}
	if tsr.SoftwareProcess.Valid {
		ts.SoftwareProcess = tsr.SoftwareProcess.String
	}
	if tsr.Emdb.Valid {
		ts.Emdb = tsr.Emdb.String
	}
	if tsr.KeyMov.Valid {
		ts.KeyMov = tsr.KeyMov.Int64
	}
	if tsr.KeyImg.Valid {
		ts.KeyImg = tsr.KeyImg.Int64
	}
	if tsr.FullName.Valid {
		ts.Microscopist = tsr.FullName.String
	}

	rows, err := dbh.Queryx(selectDataFilesSql, tiltSeriesId)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var dfr dataFileRow
		err = rows.StructScan(&dfr)
		if err != nil {
			return
		}
		df := DataFile{}
		if dfr.Filename.Valid {
			df.Filename = dfr.Filename.String
			if len(strings.TrimSpace(df.Filename)) == 0 {
				// No file name, no file...
				continue
			}
		} else {
			// No file name, no file...
			continue
		}
		if dfr.Filetype.Valid {
			df.Filetype = dfr.Filetype.String
		}
		if dfr.ThreeDFileImage.Valid {
			df.ThreeDFileImage = dfr.ThreeDFileImage.String
		}
		if dfr.Notes.Valid {
			df.Notes = dfr.Notes.String
		}
		if dfr.DefId.Valid {
			df.DefId = dfr.DefId.Int64
		}
		if dfr.Auto.Valid {
			df.Auto = dfr.Auto.Int64
		}
		df.Type = "tomogram"
		switch df.Filetype {
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

	rows, err = dbh.Queryx(selectThreeDFilesSql, tiltSeriesId)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tdfr threeDFileRow
		err = rows.StructScan(&tdfr)
		if err != nil {
			return
		}
		tdf := ThreeDFile{}
		if tdfr.Filename.Valid {
			tdf.Filename = tdfr.Filename.String
		}
		if tdfr.Classify.Valid {
			tdf.Classify = tdfr.Classify.String
		}
		if tdfr.Notes.Valid {
			tdf.Notes = tdfr.Filename.String
		}
		if tdfr.DefId.Valid {
			tdf.DefId = tdfr.DefId.Int64
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
		case "subvolume":
			fallthrough
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

	err := dbh.Select(&ids, selectFilterSql)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
