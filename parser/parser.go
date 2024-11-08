package parser

import (
	"bufio"
	"encoding/csv"
	"os"
	"regexp"
	"time"

	"github.com/daniyalumer/repeated-ip-report/models"
)

var (
	timeRegex    = regexp.MustCompile(`(\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})`)
	ipRegex      = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	keywordRegex = regexp.MustCompile(`keyword=([^ ]+)`)
	userRegex    = regexp.MustCompile(`"([^"]+)"$`)
	urlRegex     = regexp.MustCompile(`\s+(\S+)\s+HTTP/1.1`)
)

func ParseFile(file *os.File, uniqueRedirectLogs map[string][]models.RedirectLog, filteredLogs map[string][]models.RedirectLog) (map[string][]models.RedirectLog, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		redirectLog, err := ParseLine(line)
		if err != nil {
			return nil, err
		}
		if _, exists := uniqueRedirectLogs[redirectLog.Ip]; !exists {
			uniqueRedirectLogs[redirectLog.Ip] = []models.RedirectLog{redirectLog}
		} else {
			lastLog := uniqueRedirectLogs[redirectLog.Ip][len(uniqueRedirectLogs[redirectLog.Ip])-1]
			if lastLog.Keyword != redirectLog.Keyword && redirectLog.Timestamp.Sub(lastLog.Timestamp) <= 5*time.Second {
				if !ContainsLog(filteredLogs[redirectLog.Ip], lastLog) {
					filteredLogs[redirectLog.Ip] = append(filteredLogs[redirectLog.Ip], lastLog)
				}
				uniqueRedirectLogs[redirectLog.Ip] = append(uniqueRedirectLogs[redirectLog.Ip], redirectLog)
				filteredLogs[redirectLog.Ip] = append(filteredLogs[redirectLog.Ip], redirectLog)
			}
		}
	}
	err := scanner.Err()
	return filteredLogs, err
}

func ParseLine(line string) (models.RedirectLog, error) {
	timeStampStr := timeRegex.FindString(line)

	timeStamp, err := time.Parse("02/Jan/2006:15:04:05 +0000", timeStampStr)

	ip := ipRegex.FindString(line)

	keyword := keywordRegex.FindString(line)

	userAgent := userRegex.FindString(line)

	urlString := urlRegex.FindString(line)

	redirectLog := models.RedirectLog{
		Timestamp: timeStamp,
		Ip:        ip,
		Keyword:   keyword,
		UserAgent: userAgent,
		Url:       urlString,
	}

	return redirectLog, err
}

func ContainsLog(logs []models.RedirectLog, log models.RedirectLog) bool {
	for _, l := range logs {
		if l == log {
			return true
		}
	}
	return false
}

func CreateCsv() (*os.File, *csv.Writer, error) {
	f, err := os.Create("output.csv")
	if err != nil {
		return nil, nil, err
	}

	w := csv.NewWriter(f)

	w.Write([]string{"Date", "IP", "Keyword", "URL"})

	return f, w, nil
}

func WriteToCsv(w *csv.Writer, returnLogs map[string][]models.RedirectLog) error {
	for _, logs := range returnLogs {
		for _, lg := range logs {
			err := w.Write([]string{string(lg.Timestamp.Format(time.RFC3339)), lg.Ip, lg.Keyword, lg.Url})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
