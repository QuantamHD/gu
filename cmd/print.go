// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"os"
	"strconv"

	"github.com/bgentry/speakeasy"
	"github.com/olekukonko/tablewriter"
	"github.com/quantamhd/gu/utils"
	"github.com/spf13/cobra"
)

func printTable(printers []*utils.PaperCutPrinter) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "name", "location"})
	table.SetRowLine(true)

	for _, p := range printers {
		table.Append(p.ToListStrings())
	}

	table.Render()
}

func selectPrinter(printers []*utils.PaperCutPrinter) *utils.PaperCutPrinter {
	var id string
	fmt.Print("Select Printer id(i.e. 1,2 or 4, ...): ")
	fmt.Scanln(&id)

	idInt, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Please input an integer.")
		return selectPrinter(printers)
	}

	var idExists bool
	var selectedPrinter *utils.PaperCutPrinter

	for _, p := range printers {
		if p.GetID() == idInt {
			idExists = true
			selectedPrinter = p
		}
	}

	if !idExists {
		fmt.Println("That id is not in the system.")
		return selectPrinter(printers)
	}

	return selectedPrinter
}

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print <document file>",
	Short: "Prints your document at the selected location",
	Long: `This command connects to the Gonzaga Print system and sends your
document off to the printer.

Examples

gu print homework.pdf
gu print concreteReport.docx
gu print /home/family_photo.jpg

Supported Document Types

+-------------------------+-------------------------------------------+
| Application / File Type | File Extension(s)                         |
+-------------------------+-------------------------------------------+
| Microsoft Excel         | xlam, xls, xlsb, xlsm, xlsx, xltm, xltx   |
+-------------------------+-------------------------------------------+
| Microsoft PowerPoint    | pot, potm, potx, ppam, pps, ppsm,         |
|                         | ppsx, ppt, pptm, pptx                     |
+-------------------------+-------------------------------------------+
| Microsoft Word          | doc, docm, docx, dot, dotm, dotx, rtf     |
+-------------------------+-------------------------------------------+
| PDF                     | pdf                                       |
+-------------------------+-------------------------------------------+
| Picture Files           | bmp, dib, gif, jfif, jif, jpe, jpeg, jpg, |
|                         | png, tif, tiff                            |
+-------------------------+-------------------------------------------+
| XPS                     | xps                                       |
+-------------------------+-------------------------------------------+
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		var username string
		fmt.Print("Username for 'https://guprint.gonzaga.edu': ")
		fmt.Scanln(&username)
		password, _ := speakeasy.Ask("Password for 'https://guprint.gonzaga.edu': ")

		credentials := utils.CreatePaperCutCredentials(username, password)

		if !credentials.IsLoggedIn() {
			fmt.Println("Could not connect to Gonzaga Print Services")
			os.Exit(1)
		}

		printers := utils.GetPaperCutPrinters(credentials)
		printTable(printers)
		selectedPrinter := selectPrinter(printers)
		utils.CreatePrintJob(credentials, selectedPrinter, 1, "")
	},
}

func init() {
	RootCmd.AddCommand(printCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
