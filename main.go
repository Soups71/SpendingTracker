package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

var VERSION = "0.0.1"

type transType string

const (
	Credit transType = "Credit"
	Debit  transType = "Debit"
)

type Transaction struct {
	RawDate  string    `csv:"Booking Date"`
	Amount   float32   `csv:"Amount"`
	Tt       transType `csv:"Credit Debit Indicator"`
	Category string    `csv:"Category"`
	Date     time.Time `csv:"-"`
}

func getFiles(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, path+"/"+file.Name())
		}
	}
	return fileNames

}
func readCSV(filepath string) []*Transaction {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("[-] The CSV file, %s, provided could not be read.\n", filepath)
		os.Exit(-1)
	}

	defer file.Close()
	var records []*Transaction
	if err := gocsv.UnmarshalFile(file, &records); err != nil {
		fmt.Printf("[-] The CSV file, %s,  could not be parsed.\n", filepath)
		os.Exit(-1)
	}
	// Define the date layout
	layout := "01/02/2006"

	for _, row := range records {
		// Parse the date string into a time.Time value
		parsedTime, err := time.Parse(layout, row.RawDate)
		if err != nil {
			fmt.Printf("[-] A date value could not be parsed.\n")
			continue
		}

		// Update the date value in the entry
		row.Date = parsedTime
	}
	return records
}

func MonthlyTrans(transactions []*Transaction, month time.Month, year int) []*Transaction {
	var filteredTransactions []*Transaction

	for _, transaction := range transactions {
		if transaction.Date.Month() == month && transaction.Date.Year() == year {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	return filteredTransactions

}
func TransBreakDown(transactions []*Transaction) map[string]float64 {
	totalAmountPerCategory := make(map[string]float64)
	for _, trans := range transactions {
		if trans.Tt == Debit {
			totalAmountPerCategory[trans.Category] += float64(trans.Amount)
		} else {
			totalAmountPerCategory[trans.Category] -= float64(trans.Amount)
		}

	}
	return totalAmountPerCategory
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <Folder> <Month (i.e. 3)> <Year (i.e. 2024)>\n", os.Args[0])
		os.Exit(1)
	}
	filePath := os.Args[1]

	month, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("[-] The month value provided was not valid.")
		os.Exit(1)
	}

	year, _ := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("[-] The month value provided was not valid.")
		os.Exit(1)
	}

	filenames := getFiles(filePath)

	for _, filename := range filenames {
		// Open the CSV file
		records := readCSV(filename)
		fmt.Printf("Transaction Details for: %s\n", strings.Title(filename[strings.Index(filename, "/")+1:strings.Index(filename, ".")]))

		curr_trans := MonthlyTrans(records, time.Month(month), year)
		var total_debt float32 = 0.0
		var total_earned float32 = 0.0
		for _, row := range curr_trans {
			if row.Tt == Debit {
				total_debt += (row.Amount)
			} else {
				total_earned += (row.Amount)
			}
		}
		fmt.Printf("\tThe total net for %s in %d: %.2f\n\t\tThe total earned: %.2f\n\t\tThe total spent: %.2f\n", time.Month(month), year, total_earned-total_debt, total_earned, total_debt)

		//
		breakdown := TransBreakDown(curr_trans)
		var earnedoutput string
		var spentoutput string
		for category, amount := range breakdown {
			if amount < 0 {
				earnedoutput += fmt.Sprintf("\tAmount earned on %s was: %.02f\n", category, amount)
			} else {

				spentoutput += fmt.Sprintf("\tAmount spent on %s was: %.02f\n", category, amount)
			}
		}
		fmt.Print(earnedoutput + spentoutput)
	}

}
