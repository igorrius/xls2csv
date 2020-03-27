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
	InputFilenameFlag      = "input"
	OutputFilenameFlag     = "output"
	SheetNameFlag          = "sheet-name"
	SheetNumberFlag        = "sheet-number"
	SeparatorCharacterFlag = "separator"
)

var (
	ErrEmptyInputFileName                          = errors.New("input filename can not be empty")
	ErrInvalidInputParametersSheetNameAndNumberSet = errors.New("invalid input parameters: sheet name and number set at the same time")
	ErrInvalidInputParametersSeparatorLength       = errors.New("invalid input parameters: separator length must be one character")
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
		// CSV file separator character
		&cli.StringFlag{Name: SeparatorCharacterFlag, Aliases: []string{"sep"}, Usage: "CSV separator character", Value: ",", DefaultText: ","},
	}
}

func Xls2CsvFlagsValidation(c *cli.Context) error {
	cli.ShowVersion(c)

	if separator := c.String(SeparatorCharacterFlag); len(separator) != 1 {
		return ErrInvalidInputParametersSeparatorLength
	}

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
		defer func() {
			if err := outputFile.Close(); err != nil {
				_, _ = fmt.Fprintf(c.App.Writer, "[!] output file save error: %s", err.Error())
			}
		}()

		sep := []rune(c.String(SeparatorCharacterFlag))

		job := converter.NewXls2CsvJob(inputFile, outputFile, c.String(SheetNameFlag), c.Int(SheetNumberFlag), sep[0])
		_, _ = fmt.Fprintf(c.App.Writer,
			"[+] converting %s to %s, using sheet with name [%s], using CSV separator %c\n",
			inputFileName, outputFileName, job.SheetName(), sep)

		return converter.Xls2Csv(job)
	}
}

func outputFileNameFromInputFile(file *os.File) string {
	filename := file.Name()
	return strings.TrimSuffix(filename, path.Ext(filename)) + ".csv"
}
