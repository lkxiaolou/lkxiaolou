package toolx

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	SEQ  = "%seq"
	READ = "%read"
	DATE = "%date"
)

func readFile(filePath string) (string, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeFile(filePath string, content string) error {
	newFile, err := os.Create(filePath)
	defer newFile.Close()

	if err != nil {
		return err
	}
	_, err = newFile.WriteString(content)
	return err
}

func handleOneTab(tab string) (string, map[string]string) {

	tabContent := strings.Builder{}

	context := make(map[string]string)

	if strings.Contains(tab, "http") {

		urls := parseUrls(tab)
		icoMap := getIcoMap(urls)

		if len(icoMap) == 0 {
			return tab, context
		} else {
			icoIndex := IcoMin
			subs := strings.Split(tab, " ")
			for _, sub := range subs {
				trimedSub := strings.TrimSpace(sub)
				if len(trimedSub) == 0 {
					tabContent.WriteString(sub)
				} else {
					if strings.HasPrefix(trimedSub, "http") {
						for icoIndex <= IcoMax {
							if ico, ok := icoMap[icoIndex]; ok {
								tabContent.WriteString(ico)
								icoIndex++
								break
							} else {
								icoIndex++
							}
						}
					} else {
						tabContent.WriteString(sub)
					}
				}
				tabContent.WriteString(" ")
			}
		}
		tab = tabContent.String()

		readCount := getReadCount(urls)
		context[READ] = strconv.Itoa(readCount)

		date := getDate(urls)
		context[DATE] = date
	}

	return tab, context
}

func handleOneLine(line string, lastSeq int) (string, int) {

	lineContent := &strings.Builder{}

	trimedLine := strings.TrimSpace(line)
	var tabStr string
	var context map[string]string
	var tmpContext map[string]string

	if strings.HasPrefix(trimedLine, "|") {
		tabs := strings.Split(trimedLine, "|")
		for i, tab := range tabs {
			tabStr, tmpContext = handleOneTab(tab)
			if len(tmpContext) > 0 {
				context = tmpContext
			}
			lineContent.WriteString(tabStr)
			if i+1 < len(tabs) {
				lineContent.WriteString("|")
			}
		}
	} else {
		lineContent.WriteString(line)
	}

	lineStr := lineContent.String()
	if strings.Contains(lineStr, SEQ) {
		lastSeq = lastSeq + 1
		lineStr = strings.ReplaceAll(lineStr, SEQ, strconv.Itoa(lastSeq))
	}
	if date, ok := context[DATE]; ok {
		lineStr = strings.ReplaceAll(lineStr, DATE, date)
	}
	if read, ok := context[READ]; ok {
		lineStr = strings.ReplaceAll(lineStr, READ, read)
	}

	return lineStr, lastSeq
}

func FormatReadMe(templateFile string, outputFile string) error {
	content, err := readFile(templateFile)
	if err != nil {
		return err
	}

	finalContent := strings.Builder{}

	lines := strings.Split(content, "\n")
	seq := 0
	var lineStr string

	for _, line := range lines {
		lineStr, seq = handleOneLine(line, seq)
		finalContent.WriteString(lineStr)
		finalContent.WriteString("\n")
	}

	return writeFile(outputFile, finalContent.String())
}

func parseUrls(tab string) []string {
	ss := strings.Split(tab, " ")
	urls := make([]string, 0)
	for _, s := range ss {
		st := strings.TrimSpace(s)
		if strings.HasPrefix(st, "http") {
			urls = append(urls, st)
		}
	}
	return urls
}

func getIcoMap(urls []string) map[int]string {
	icoMap := make(map[int]string)
	for _, u := range urls {
		ico, icoType := GetIconLink(u)
		if icoType != IcoUnknown {
			icoMap[icoType] = ico
		}
	}
	return icoMap
}

func getReadCount(urls []string) int {
	return 0
}

func getDate(urls []string) string {
	for _, u := range urls {
		if strings.Contains(u, wechatHost) {
			//content, err := HttpGetWithCache(u)
			//if err == nil {
			//	fmt.Println(content)
			//}
		}
	}
	return ""
}
