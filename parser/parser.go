package parser

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/daniyalumer/repeated-ip-report/models"
)

type RedirectLog = models.RedirectLog

func ParseFile(file *os.File, uniqueRedirectLogs map[string][]RedirectLog, returnLogs map[string][]RedirectLog) map[string][]RedirectLog {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		redirectLog := ParseLine(line)
		if _, exists := uniqueRedirectLogs[redirectLog.Ip]; !exists {
			uniqueRedirectLogs[redirectLog.Ip] = []RedirectLog{redirectLog}
		} else {
			lastLog := uniqueRedirectLogs[redirectLog.Ip][len(uniqueRedirectLogs[redirectLog.Ip])-1]
			if lastLog.Keyword != redirectLog.Keyword && redirectLog.Timestamp.Sub(lastLog.Timestamp) <= 5*time.Second {
				if !ContainsLog(returnLogs[redirectLog.Ip], lastLog) {
					returnLogs[redirectLog.Ip] = append(returnLogs[redirectLog.Ip], lastLog)
				}
				uniqueRedirectLogs[redirectLog.Ip] = append(uniqueRedirectLogs[redirectLog.Ip], redirectLog)
				returnLogs[redirectLog.Ip] = append(returnLogs[redirectLog.Ip], redirectLog)
			}
		}
	}
	return returnLogs
}

func ParseLine(line string) RedirectLog {
	timeRegex := regexp.MustCompile(`(\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})`)
	timeStampStr := timeRegex.FindString(line)

	timeStamp, err := time.Parse("02/Jan/2006:15:04:05 +0000", timeStampStr)
	if err != nil {
		log.Fatal(err)
	}

	ipRegex := regexp.MustCompile(`(\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3})`)
	ip := ipRegex.FindString(line)

	keywordRegex := regexp.MustCompile(`keyword=([^ ]+)`)
	keyword := keywordRegex.FindString(line)

	userRegex := regexp.MustCompile(`"([^"]+)"$`)
	userAgent := userRegex.FindString(line)

	urlRegex := regexp.MustCompile(`\s+(\S+)\s+HTTP/1.1`)
	urlString := urlRegex.FindString(line)

	redirectLog := RedirectLog{
		Timestamp: timeStamp,
		Ip:        ip,
		Keyword:   keyword,
		UserAgent: userAgent,
		Url:       urlString,
	}

	return redirectLog
}

func ContainsLog(logs []RedirectLog, log RedirectLog) bool {
	for _, l := range logs {
		if l == log {
			return true
		}
	}
	return false
}

func CreateCsv() (*os.File, *csv.Writer) {
	f, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(f)

	w.Write([]string{"Date", "IP", "Keyword", "URL"})

	return f, w
}

func WriteToCsv(w *csv.Writer, returnLogs map[string][]RedirectLog) {
	for _, logs := range returnLogs {
		for _, lg := range logs {
			err := w.Write([]string{fmt.Sprintf("%s", lg.Timestamp), lg.Ip, lg.Keyword, lg.Url})
			if err != nil {
				log.Println(err)
			}
		}
	}
}