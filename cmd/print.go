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

func printTable(printers map[int]utils.PaperCutPrinter) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "name", "location"})
	table.SetRowLine(true)

	for _, p := range printers {
		table.Append(p.ToListStrings())
	}

	table.Render()
}

/*
Prompts user to select number of copies.
*/
func selectCopies() int {
	var copies string
	fmt.Print("Number of copies: ")
	fmt.Scanln(&copies)

	// check if printerID is an int and in range
	num, err := strconv.Atoi(copies)
	if err != nil {
		fmt.Println("Not a valid number of copies!")
		os.Exit(1)
	} else if num < 1 {
		fmt.Println("Not a valid number of copies!")
		os.Exit(1)
	}
	return num
}

/*
Prompts user to select printer.
Pass in map of printers.
Returns selected printer.
*/
func selectPrinter(printers map[int]utils.PaperCutPrinter) utils.PaperCutPrinter {

	// If only one printer, return that printer
	if len(printers) == 1 {
		return printers[0]
	}

	var printerID string
	fmt.Print("Select a printer ID: ")
	fmt.Scanln(&printerID)

	// check if printerID is an int and a valid ID
	id, err := strconv.Atoi(printerID)
	if err != nil {
		fmt.Println("Not a valid ID!")
		os.Exit(1)
	} else if _, ok := printers[id]; !ok {
		fmt.Println("Not a valid ID!")
		os.Exit(1)
	}

	return printers[id]
}

/*
Handles login with user. Exits if failed login.
Returns credentials object.
*/
func login() *utils.PaperCutCredentials {
	var username string
	fmt.Print("Username for 'https://guprint.gonzaga.edu': ")
	fmt.Scanln(&username)
	password, _ := speakeasy.Ask("Password for 'https://guprint.gonzaga.edu': ")

	credentials := utils.CreatePaperCutCredentials(username, password)

	if !credentials.IsLoggedIn() {
		fmt.Println("Could not connect to Gonzaga Print Services")
		os.Exit(1)
	}

	return credentials
}

/*
Extracts the filePath from command line args. Exits if no file path.
*/
func getFilePath() string {

	if len(os.Args) < 3 {
		fmt.Println("Need to specify a file to print!")
		os.Exit(1)
	}

	return os.Args[2]
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

		filePath := getFilePath()
		credentials := login()
		printers := utils.GetPaperCutPrinters(credentials)
		printTable(printers)
		printer := selectPrinter(printers)
		copies := selectCopies()
		utils.CreatePrintJob(credentials, &printer, copies, filePath)

		fmt.Println("Printing " + strconv.Itoa(copies) + " copies of " +
			filePath + " to printer " + printer.GetName() + ".")

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
