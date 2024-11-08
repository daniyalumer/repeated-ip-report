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
	if filteredLogs == nil {
		log.Fatal("filteredLogs is nil")
	}

	defer file.Close()

	f, w, err := parser.CreateCsv()
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	err = parser.WriteToCsv(w, filteredLogs)
	if err != nil {
		log.Fatal(err)
	}

	w.Flush()
}
