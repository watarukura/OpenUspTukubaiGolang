package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"

	util "github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
日付    mdate -y <yyyymmdd>                     : 曜日
        mdate -e <yyyywwdd>/±<dif>              : dif 日先までの日付を連続出力
        mdate -e <yyyymmdd1> <yyyymmdd2>        : 日付の範囲を連続出力
        mdate <yyyywwdd>/±<dif>                 : dif 日先の日付
        mdate <yyyymmdd1> <yyyymmdd2>           : 日付の差
        mdate <yyyymm>m/±<dif>                  : dif 月先の月
        mdate -e <yyyymm>m/±<dif>               : dif 月先までの月を連続出力
        mdate <yyyymm1>m <yyyymm2>m             : 月の差
        mdate -ly <yyyymm>m                     : 前年月
`

type dayOfWeekNum int

const (
	Monday dayOfWeekNum = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type option struct {
	isDayOfWeekMode bool
	isDiffMode      bool
	isSequenceMode  bool
	isLastYearMode  bool
}

type cli struct {
	outStream, errStream io.Writer
	inStream             io.Reader
}

var (
	layoutMonth = "200601"
	layoutDate  = "20060102"
)

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr, inStream: os.Stdin}
	os.Exit(cli.run(os.Args))
}

func (c *cli) run(args []string) int {
	param, err := shellwords.Parse(strings.Join(args[1:], " "))
	if err != nil {
		util.Fatal(err, util.ExitCodeFlagErr)
	}
	// fmt.Println(param)
	option := &option{isDayOfWeekMode: false, isDiffMode: false, isSequenceMode: false, isLastYearMode: false}

	firstDate, lastDate, firstMonth, lastMonth := validateParam(param, c.inStream, option)
	switch {
	case option.isDayOfWeekMode:
		mdateDayOfWeek(firstDate, option, c.outStream)
	case option.isLastYearMode:
		mdateLastYear(firstMonth, option, c.outStream)
	case option.isDiffMode && option.isSequenceMode:
		mdateDayOfWeek(lastDate, option, c.outStream)
	case option.isDiffMode && !option.isSequenceMode:
		mdateDayOfWeek(firstMonth, option, c.outStream)
	case !option.isDiffMode && option.isSequenceMode:
		mdateDayOfWeek(lastMonth, option, c.outStream)
	default:
		return util.ExitCodeNG
	}
	return util.ExitCodeOK
}

// func validateParam(param []string, inStream io.Reader, opt *option) (firstDate, lastDate, firstMonth, lastMonth time.Time) {
func validateParam(param []string, inStream io.Reader, opt *option) (firstDate, lastDate, firstMonth, lastMonth time.Time) {
	if len(param) < 1 || len(param) > 5 {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
	}

	firstDateStr := ""
	lastDateStr := ""
	firstMonthStr := ""
	lastMonthStr := ""
	var err error
	for _, p := range param {
		// fmt.Print(i)
		// fmt.Println(": " + p)
		// fmt.Println("prev: " + prev)
		if strings.HasPrefix(p, "-y") {
			if len(p) == 2 {
				opt.isDayOfWeekMode = true
				continue
			} else {
				util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
			}
		}
		if strings.HasPrefix(p, "-e") {
			if len(p) == 2 {
				opt.isSequenceMode = true
				continue
			} else {
				util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
			}
		}
		if strings.HasPrefix(p, "-ly") {
			if len(p) == 3 {
				opt.isLastYearMode = true
				continue
			} else {
				util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
			}
		}
		if strings.HasSuffix(p, "m") {
			if len(p) == 7 {
				if firstMonthStr == "" {
					firstMonthStr = p[:6]
					firstMonth, err = time.Parse(layoutMonth, firstMonthStr)
					if err != nil {
						util.Fatal(err, util.ExitCodeFlagErr)
					}
				} else {
					lastMonthStr = p[:6]
					lastMonth, err = time.Parse(layoutMonth, lastMonthStr)
					if err != nil {
						util.Fatal(err, util.ExitCodeFlagErr)
					}
				}
			}
			if strings.Contains(p, "/") {
				if firstMonthStr == "" {
					firstMonthStr = p[0:6]
					firstMonth, err = time.Parse(layoutMonth, firstMonthStr)
					if err != nil {
						util.Fatal(err, util.ExitCodeFlagErr)
					}
					signStr := p[7:8]
					deltaStr := p[8:]
					delta, err := strconv.Atoi(deltaStr)
					if err != nil {
						util.Fatal(err, util.ExitCodeFlagErr)
					}
					if signStr == "+" {
						lastMonth = addMonth(firstMonth, delta)
					} else if signStr == "-" {
						lastMonth = addMonth(firstMonth, -1*delta)
					}
				}
			}
		}
		if len(p) == 8 {
			if firstDateStr == "" {
				firstDateStr = p
				firstDate, err = time.Parse(layoutDate, firstDateStr)
				if err != nil {
					util.Fatal(err, util.ExitCodeFlagErr)
				}
			} else {
				lastDateStr = p
				lastDate, err = time.Parse(layoutDate, lastDateStr)
				if err != nil {
					util.Fatal(err, util.ExitCodeFlagErr)
				}
			}
		}
		if strings.Contains(p, "/") {
			if firstDateStr == "" {
				firstDateStr = p[0:8]
				firstDate, err = time.Parse(layoutDate, firstDateStr)
				if err != nil {
					util.Fatal(err, util.ExitCodeFlagErr)
				}
				signStr := p[9:10]
				deltaStr := p[10:]
				delta, err := strconv.Atoi(deltaStr)
				if err != nil {
					util.Fatal(err, util.ExitCodeFlagErr)
				}
				if signStr == "+" {
					lastDate = firstDate.AddDate(0, 0, delta)
				} else if signStr == "-" {
					lastDate = firstDate.AddDate(0, 0, -1*delta)
				}
			}
		}
	}

	return firstDate, lastDate, firstMonth, lastMonth
}

func addMonth(t time.Time, d_month int) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	newMonth := int(month) + d_month
	newLastDay := getLastDay(year, newMonth)
	var newDay int
	if day > newLastDay {
		newDay = newLastDay
	} else {
		newDay = day
	}

	return time.Date(year, time.Month(newMonth), newDay, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

}

// その月の最終日を求める
func getLastDay(year, month int) int {
	t := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local)
	t = t.AddDate(0, 0, -1)
	return t.Day()
}

func mdateDayOfWeek(firstDate time.Time, opt *option, outStream io.Writer) {
	dayOfWeek := firstDate.Weekday()
	fmt.Fprint(outStream, dayOfWeekNum(dayOfWeek))
}

func mdateLastYear(firstMonth time.Time, opt *option, outStream io.Writer) {
	lastYearMonth := addMonth(firstMonth, -12)
	fmt.Fprint(outStream, lastYearMonth.Format(layoutMonth))
}
