package util

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"time"
)

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Red(str string) string {
	return textColor(textRed, str)
}

func Green(str string) string {
	return textColor(textGreen, str)
}

func Yellow(str string) string {
	return textColor(textYellow, str)
}

func textColor(c int, s string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", c, s)
}

func NewKvDisplayTable() table.Writer {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"key", "value"})
	tbl.Style().Format.Header = text.FormatLower
	return tbl
}

func NewNormalDisplayTable(titles []interface{}) table.Writer {
	tbl := table.NewWriter()
	tbl.Style().Format.Header = text.FormatLower
	tbl.SetOutputMirror(os.Stdout)
	t := table.Row{}
	for _, v := range titles {
		t = append(t, v)
	}
	tbl.AppendHeader(t)
	return tbl
}

func EmptyResult(t time.Time) {
	fmt.Printf("Empty set (%0.2f sec)\n", Duration(t))
}

func Render(tbl table.Writer) {
	fmt.Println()
	tbl.Render()
	fmt.Println()
}

func NRowsAffected(n int, t time.Time) {
	fmt.Printf("Query OK, %d rows affected (%0.2f sec)\n", n, Duration(t))
}

func QueryOkNRows(n int, t time.Time) {
	fmt.Printf("Query OK, %d rows affected (%0.2f sec)\n", n, Duration(t))
}

func NRowsInSet(n int, t time.Time) {
	fmt.Printf("%d rows in set (%0.2f sec)\n", n, Duration(t))
}

func Duration(t time.Time) float64 {
	return time.Now().Sub(t).Seconds()
}

func IsNil(i interface{}) interface{} {
	if i == nil {
		return "NULL"
	}
	return i
}

func Status(s string) string {
	if s == "normal" {
		return Green(s)
	}
	return Red(s)
}

func RunSuccess(s string) {
	fmt.Println(fmt.Sprintf("Run `%s` success.", s))
}
