package output

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

func maxKeyLen(data []SummaryItem) int {
	maxLen := 0
	for _, i := range data {
		if len(i.Key) > maxLen {
			maxLen = len(i.Key)
		}
	}
	return maxLen
}

type SummaryItem struct {
	Key   string
	Value string
}

func SummaryTable(data []SummaryItem) string {
	keyFormat := fmt.Sprintf("%%-%ds", maxKeyLen(data)+1)
	var out strings.Builder
	for _, i := range data {
		out.WriteString(fmt.Sprintf("%s %s\n", text.Bold.Sprintf(keyFormat, i.Key+":"), i.Value))
	}

	return out.String()
}

func Failed(failureCount int) string {
	return text.Colors{text.Bold, text.FgRed}.Sprintf("%d failed", failureCount)
}

func Passed(successCount int) string {
	return text.Colors{text.Bold, text.FgGreen}.Sprintf("%d passed", successCount)
}

func Total(count int) string {
	return fmt.Sprintf("%d total", count)
}
