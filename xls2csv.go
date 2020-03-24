package main

import (
	"fmt"
	"github.com/igorrius/xls2csv/action"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Version: "1.0.1",
		Name:    "XLS to CSV converter",
		Usage:   "xls2csv -i in.xls -o out.csv",
		Action:  action.Xls2Csv(),
		Flags:   action.Xls2CsvFlags(),
		Before: func(c *cli.Context) error {
			return action.Xls2CsvFlagsValidation(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("\n[!] Application error: " + err.Error())
	}
}
