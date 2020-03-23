package action

import (
	"errors"
	"fmt"
	"github.com/igorrius/xls2csv/converter"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strings"
)

const (
	InputFilenameFlag  = "input"
	OutputFilenameFlag = "output"
	SheetNameFlag      = "sheet-name"
	SheetNumberFlag    = "sheet-number"
)

var (
	ErrEmptyInputFileName                          = errors.New("input filename can not be empty")
	ErrInvalidInputParametersSheetNameAndNumberSet = errors.New("invalid input parameters: sheet name and number set at the same time")
)

func Xls2CsvFlags() []cli.Flag {
	return []cli.Flag{
		// Input (source) file name flag
		&cli.StringFlag{Name: InputFilenameFlag, Aliases: []string{"i"}, Required: true, Usage: "Input file name"},
		// Output (result) file name flag
		// This flag is not mandatory. If this flag has omitted using filename from source
		&cli.StringFlag{Name: OutputFilenameFlag, Aliases: []string{"o"}, Usage: "Output file name"},
		// If processing sheet name or processing file name has not been set up will use active sheet
		// Pay attention that if both parameters have been set this throw a validation error
		// Processing sheet name in input XLS file
		&cli.StringFlag{Name: SheetNameFlag, Aliases: []string{"sname"}, Usage: "Processing sheet name"},
		// Processing sheet number in input XLS file
		&cli.IntFlag{Name: SheetNumberFlag, Aliases: []string{"snum"}, Usage: "Processing sheet number"},
	}
}

func Xls2CsvFlagsValidation(c *cli.Context) error {
	cli.ShowVersion(c)

	name := c.String(SheetNameFlag)
	number := c.Int(SheetNumberFlag)
	if name != "" && number > 0 {
		return ErrInvalidInputParametersSheetNameAndNumberSet
	}

	return nil
}

func Xls2Csv() cli.ActionFunc {
	return func(c *cli.Context) error {
		inputFileName := c.String(InputFilenameFlag)
		if inputFileName == "" {
			return ErrEmptyInputFileName
		}

		inputFile, err := os.Open(inputFileName)
		if err != nil {
			return err
		}
		defer func() { _ = inputFile.Close() }()

		outputFileName := c.String(OutputFilenameFlag)
		if outputFileName == "" {
			outputFileName = outputFileNameFromInputFile(inputFile)
		}

		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return err
		}
		defer func() { _ = outputFile.Close() }()

		job := converter.NewXls2CsvJob(inputFile, outputFile, c.String(SheetNameFlag), c.Int(SheetNumberFlag))
		_, _ = fmt.Fprintf(c.App.Writer,
			"[+] converting %s to %s, using sheet with name [%s]\n",
			inputFileName, outputFileName, job.SheetName())

		return converter.Xls2Csv(job)
	}
}

func outputFileNameFromInputFile(file *os.File) string {
	filename := file.Name()
	return strings.TrimSuffix(filename, path.Ext(filename)) + ".csv"
}
