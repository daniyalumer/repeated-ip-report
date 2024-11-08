package main

import (
	"log"
	"os"

	"github.com/daniyalumer/repeated-ip-report/models"
	"github.com/daniyalumer/repeated-ip-report/parser"
)

func main() {
	file, err := os.Open("logs_1feb.txt")
	if err != nil {
		log.Fatal(err)
	}

	uniqueRedirectLogs := make(map[string][]models.RedirectLog)
	filteredLogs := make(map[string][]models.RedirectLog)

	filteredLogs, err = parser.ParseFile(file, uniqueRedirectLogs, filteredLogs)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	f, w := parser.CreateCsv()

	defer f.Close()

	parser.WriteToCsv(w, filteredLogs)

	w.Flush()
}
