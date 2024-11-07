package main

import (
	"log"
	"os"

	"github.com/daniyalumer/repeated-ip-report/parser"
)

type RedirectLog = parser.RedirectLog

func main() {
	file, err := os.Open("logs_1feb.txt")
	if err != nil {
		log.Fatal(err)
	}

	uniqueRedirectLogs := make(map[string][]RedirectLog)
	returnLogs := make(map[string][]RedirectLog)

	parser.ParseFile(file, uniqueRedirectLogs, returnLogs)

	defer file.Close()

	f, w := parser.CreateCsv()

	defer f.Close()

	parser.WriteToCsv(w, returnLogs)

	w.Flush()
}
