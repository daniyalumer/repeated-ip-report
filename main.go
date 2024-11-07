package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"
)

type RedirectLog struct {
	Timestamp time.Time
	Ip        string
	Keyword   string
	UserAgent string
	Url       string
}

func parseLine(line string) RedirectLog {
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

func containsLog(logs []RedirectLog, log RedirectLog) bool {
	for _, l := range logs {
		if l == log {
			return true
		}
	}
	return false
}

func main() {
	file, err := os.Open("logs-testing.txt")
	if err != nil {
		log.Fatal(err)
	}

	uniqueRedirectLogs := make(map[string][]RedirectLog)
	returnLogs := make(map[string][]RedirectLog)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		redirectLog := parseLine(line)
		if _, exists := uniqueRedirectLogs[redirectLog.Ip]; !exists {
			uniqueRedirectLogs[redirectLog.Ip] = []RedirectLog{redirectLog}
		} else {
			lastLog := uniqueRedirectLogs[redirectLog.Ip][len(uniqueRedirectLogs[redirectLog.Ip])-1]
			if lastLog.Keyword != redirectLog.Keyword && lastLog.Timestamp.Sub(redirectLog.Timestamp) <= 5*time.Second {
				if !containsLog(returnLogs[redirectLog.Ip], lastLog) {
					returnLogs[redirectLog.Ip] = append(returnLogs[redirectLog.Ip], lastLog)
				}
				uniqueRedirectLogs[redirectLog.Ip] = append(uniqueRedirectLogs[redirectLog.Ip], redirectLog)
				returnLogs[redirectLog.Ip] = append(returnLogs[redirectLog.Ip], redirectLog)
			}
		}
	}
	defer file.Close()

	log.Printf("returnLogs: %+v", returnLogs)
}
