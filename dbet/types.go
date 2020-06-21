// 定义了 TiltSeries、DataFile、ThreeDFile 三个结构体。

package main

import (
	"time"
)

type TiltSeries struct {
	Id                  string
	Title               string
	Date                time.Time
	TiltSeriesNotes     string
	ScopeName           string
	Roles               string
	ScopeNotes          string
	SpeciesName         string
	SpeciesNotes        string
	SpeciesStrain       string
	SpeciesTaxId        int64
	SingleDual          int64
	Defocus             float64
	Magnification       float64
	Dosage              float64
	TiltConstant        float64
	TiltMin             float64
	TiltMax             float64
	TiltStep            float64
	SoftwareAcquisition string
	SoftwareProcess     string
	Emdb                string
	KeyImg              int64
	KeyMov              int64
	Microscopist        string
	Institution         string
	Lab                 string
	DataFiles           []DataFile
	ThreeDFiles         []ThreeDFile
}

type DataFile struct {
	Filetype        string
	Filename        string
	Notes           string
	ThreeDFileImage string
	Type            string
	SubType         string
	FilePath        string
	DefId           int64
	Auto            int64
	Software        string
}

type ThreeDFile struct {
	Classify string
	Notes    string
	Filename string
	Type     string
	SubType  string
	FilePath string
	DefId    int64
	Software string
}
