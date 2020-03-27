package converter

import (
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
)

type Xls2CsvJob struct {
	sheetName   string
	sheetNumber int
	xls         *excelize.File
	csv         *csv.Writer
}

func (j *Xls2CsvJob) SheetName() string {
	if j.sheetName != "" {
		if sheetIndex := j.xls.GetSheetIndex(j.sheetName); sheetIndex > 0 {
			return j.sheetName
		}
		fmt.Printf("[!] Sheet with name [%s] has not been found, the active sheet will be used\n", j.sheetName)
	}

	if j.sheetNumber > 0 {
		if number, ok := j.xls.GetSheetMap()[j.sheetNumber]; ok {
			return number
		}
		fmt.Printf("[!] Sheet with number [%d] has not been found, the active sheet will be used\n", j.sheetNumber)
	}

	return j.xls.GetSheetMap()[j.xls.GetActiveSheetIndex()]
}

func NewXls2CsvJob(input io.Reader, output io.Writer, sheetName string, sheetNumber int, separator rune) *Xls2CsvJob {
	xls, err := excelize.OpenReader(input)
	if err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(output)
	csvWriter.Comma = separator

	return &Xls2CsvJob{
		sheetName:   sheetName,
		sheetNumber: sheetNumber,
		xls:         xls,
		csv:         csvWriter,
	}
}

func Xls2Csv(job *Xls2CsvJob) error {
	return convertXls2Csv(job.xls, job.csv, job.SheetName())
}

func convertXls2Csv(xls *excelize.File, csvWriter *csv.Writer, sheetName string) error {
	rows, err := xls.Rows(sheetName)
	if err != nil {
		return err
	}

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		if err = csvWriter.Write(columns); err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}
