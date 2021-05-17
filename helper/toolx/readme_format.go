package toolx

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func FormatReadMe(templateFile string, outputFile string) error {
	file, err := os.OpenFile(templateFile, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	finalContent := strings.Builder{}
	totalCount := 0
	subTotal := 0
	replaceCount := 0

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "%d") {
			replaceCount++
		}
	}

	totals := make([]int, replaceCount)
	totalIndex := 0
	curIndex := 0
	endCount := false

	for _, line := range lines {
		trimedLine := strings.TrimSpace(line)

		if strings.Contains(line, "%d") {
			if curIndex != totalIndex {
				if endCount {
					totals[curIndex] = subTotal
					curIndex++
					endCount = false
				}
				if !endCount {
					endCount = true
				}
				subTotal = 0
			} else {
				curIndex++
			}
		}

		if strings.HasPrefix(trimedLine, "-") {
			icoMap := handleLine(trimedLine)
			if len(icoMap) == 0 {
				finalContent.WriteString(line)
			} else {
				totalCount++
				subTotal++
				icoIndex := IcoMin
				subs := strings.Split(line, " ")
				for _, sub := range subs {
					trimedSub := strings.TrimSpace(sub)
					if len(trimedSub) == 0 {
						finalContent.WriteString(sub)
					} else {
						if strings.HasPrefix(trimedSub, "http") {
							for icoIndex <= IcoMax {
								if ico, ok := icoMap[icoIndex]; ok {
									finalContent.WriteString(ico)
									icoIndex++
									break
								} else {
									icoIndex++
								}
							}
						} else {
							finalContent.WriteString(sub)
						}
					}
					finalContent.WriteString(" ")
				}
			}
		} else {
			finalContent.WriteString(line)
		}
		finalContent.WriteString("\n")
	}

	if endCount {
		totals[curIndex] = subTotal
	}

	totals[totalIndex] = totalCount
	finalContentStr := finalContent.String()
	for _, total := range totals {
		finalContentStr = strings.Replace(finalContentStr, "%d", strconv.Itoa(total), 1)
	}

	// 写入文件
	newFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	_, err = newFile.WriteString(finalContentStr)
	return err
}

func handleLine(line string) map[int]string {
	ss := strings.Split(line, " ")

	icoMap := make(map[int]string)

	for _, s := range ss {
		st := strings.TrimSpace(s)
		if strings.HasPrefix(st, "http") {
			ico, icoType := GetIconLink(st)
			if icoType != IcoUnknown {
				icoMap[icoType] = ico
			}
		}
	}

	return icoMap
}
