package main

import (
	"bufio"
	"log"
	"os"
	"time"

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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		redirectLog := parser.ParseLine(line)
		if _, exists := uniqueRedirectLogs[redirectLog.Ip]; !exists {
			uniqueRedirectLogs[redirectLog.Ip] = []RedirectLog{redirectLog}
		} else {
			lastLog := uniqueRedirectLogs[redirectLog.Ip][len(uniqueRedirectLogs[redirectLog.Ip])-1]
			if lastLog.Keyword != redirectLog.Keyword && lastLog.Timestamp.Sub(redirectLog.Timestamp) <= 5*time.Second {
				if !parser.ContainsLog(returnLogs[redirectLog.Ip], lastLog) {
					returnLogs[redirectLog.Ip] = append(returnLogs[redirectLog.Ip], lastLog)
				}
				uniqueRedirectLogs[redirectLog.Ip] = append(uniqueRedirectLogs[redirectLog.Ip], redirectLog)
				returnLogs[redirectLog.Ip] = append(returnLogs[redirectLog.Ip], redirectLog)
			}
		}
	}
	defer file.Close()

	f, w := parser.CreateCsv()
	defer f.Close()

	parser.WriteToCsv(w, returnLogs)

	w.Flush()

}
