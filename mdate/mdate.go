package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/watarukura/OpenUspTukubaiGolang/util"
)

const usageText = `
Usage of %s:
   %s -y <yyyymmdd> | [-e] <yyyymmdd>/+-<diff> | [-e] <yyyymmdd> <yyyymmdd> | [-e] <yyyymm>m/+-<diff> | [-e] <yyyymm>m <yyyymm>m | -ly <yyyymm>m
`

// 曜日をintにする
type dayOfWeekNum int

// 曜日をintにする
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
	isDiffMode      bool // true: 差の数値を出力する false: 差分を加算した日付を出力する
	isSequenceMode  bool
	isLastYearMode  bool
	isMonthMode     bool
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
	flags := flag.NewFlagSet("mdate", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	option := &option{isDayOfWeekMode: false, isDiffMode: false, isSequenceMode: false, isLastYearMode: false, isMonthMode: false}
	flags.BoolVar(&option.isDayOfWeekMode, "y", false, "day of week")
	flags.BoolVar(&option.isSequenceMode, "e", false, "days seq")
	flags.BoolVar(&option.isLastYearMode, "ly", false, "last year")

	if err := flags.Parse(args[1:]); err != nil {
		return util.ExitCodeParseFlagErr
	}
	param := flags.Args()
	// param, err := shellwords.Parse(strings.Join(args[1:], " "))
	// if err != nil {
	// 	util.Fatal(err, util.ExitCodeFlagErr)
	// }
	// fmt.Println(param)
	// option := &option{isDayOfWeekMode: false, isDiffMode: false, isSequenceMode: false, isLastYearMode: false, isMonthMode: false}

	firstDate, lastDate, firstMonth, lastMonth := validateParam(param, c.inStream, option)
	//fmt.Println(option)
	//fmt.Println(firstMonth)
	//fmt.Println(lastMonth)
	switch {
	case option.isDayOfWeekMode:
		mdateDayOfWeek(firstDate, option, c.outStream)
	case option.isLastYearMode:
		mdateLastYear(firstMonth, option, c.outStream)
	case option.isSequenceMode && !option.isMonthMode:
		mdateDiffSeq(firstDate, lastDate, option, c.outStream)
	case !option.isDiffMode && !option.isSequenceMode && !option.isMonthMode:
		mdateDiff(firstDate, lastDate, option, c.outStream)
	case option.isDiffMode && !option.isSequenceMode && !option.isMonthMode:
		mdateDiffLastDate(lastDate, option, c.outStream)
	case option.isSequenceMode && option.isMonthMode:
		mdateDiffSeqMonth(firstMonth, lastMonth, option, c.outStream)
	case !option.isDiffMode && !option.isSequenceMode && option.isMonthMode:
		//fmt.Println(1)
		mdateDiffMonth(firstMonth, lastMonth, option, c.outStream)
	case option.isDiffMode && !option.isSequenceMode && option.isMonthMode:
		mdateDiffLastMonth(lastMonth, option, c.outStream)
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
		// if strings.HasPrefix(p, "-y") {
		// 	if len(p) == 2 {
		// 		opt.isDayOfWeekMode = true
		// 		continue
		// 	} else {
		// 		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		// 	}
		// }
		// if strings.HasPrefix(p, "-e") {
		// 	if len(p) == 2 {
		// 		opt.isSequenceMode = true
		// 		continue
		// 	} else {
		// 		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		// 	}
		// }
		// if strings.HasPrefix(p, "-ly") {
		// 	if len(p) == 3 {
		// 		opt.isLastYearMode = true
		// 		continue
		// 	} else {
		// 		util.Fatal(errors.New("failed to read param"), util.ExitCodeFlagErr)
		// 	}
		// }
		if strings.Contains(p, "m") {
			// month mode
			opt.isMonthMode = true

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
					signStr := p[8:9]
					deltaStr := p[9:]
					delta, err := strconv.Atoi(deltaStr)
					if err != nil {
						util.Fatal(err, util.ExitCodeFlagErr)
					}
					if signStr == "+" {
						lastMonth = addMonth(firstMonth, delta)
					} else if signStr == "-" {
						lastMonth = addMonth(firstMonth, -1*delta)
					}
					opt.isDiffMode = true
				}
			}
		} else {
			// month modeでない場合
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
					opt.isDiffMode = true
				}
			}
		}
	}

	return firstDate, lastDate, firstMonth, lastMonth
}

func addMonth(t time.Time, dMonth int) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	newMonth := int(month) + dMonth
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
	if dayOfWeek == time.Sunday {
		fmt.Fprintln(outStream, 7)
	} else {
		fmt.Fprintln(outStream, dayOfWeekNum(dayOfWeek))
	}
}

func mdateLastYear(firstMonth time.Time, opt *option, outStream io.Writer) {
	lastYearMonth := addMonth(firstMonth, -12)
	fmt.Fprintln(outStream, lastYearMonth.Format(layoutMonth))
}

func mdateDiffSeq(firstDate, lastDate time.Time, opt *option, outStream io.Writer) {
	sign := firstDate.Sub(lastDate)
	var lowerDate, higherDate time.Time
	if sign > 0 {
		lowerDate = lastDate
		higherDate = firstDate
	} else {
		lowerDate = firstDate
		higherDate = lastDate
	}

	fmt.Fprint(outStream, lowerDate.Format(layoutDate))
	date := lowerDate.AddDate(0, 0, 1)
	for {
		duration := higherDate.Sub(date)
		if duration < 0 {
			fmt.Fprint(outStream, "\n")
			break
		}
		fmt.Fprint(outStream, " "+date.Format(layoutDate))
		date = date.AddDate(0, 0, 1)
	}
}

func mdateDiff(firstDate, lastDate time.Time, opt *option, outStream io.Writer) {
	duration := firstDate.Sub(lastDate)
	hours := int(duration.Hours())
	days := hours / 24
	fmt.Fprintln(outStream, days)
}

func mdateDiffSeqMonth(firstMonth, lastMonth time.Time, opt *option, outStream io.Writer) {
	sign := firstMonth.Sub(lastMonth)
	var lowerMonth, higherMonth time.Time
	if sign > 0 {
		lowerMonth = lastMonth
		higherMonth = firstMonth
	} else {
		lowerMonth = firstMonth
		higherMonth = lastMonth
	}

	fmt.Fprint(outStream, lowerMonth.Format(layoutMonth))
	month := addMonth(lowerMonth, 1)
	for {
		duration := higherMonth.Sub(month)
		if duration < 0 {
			fmt.Fprint(outStream, "\n")
			break
		}
		fmt.Fprint(outStream, " "+month.Format(layoutMonth))
		month = addMonth(month, 1)
	}
}

func mdateDiffMonth(firstMonth, lastMonth time.Time, opt *option, outStream io.Writer) {
	duration := firstMonth.Sub(lastMonth)
	count := 0
	//fmt.Println(duration)
	if duration > 0 {
		for {
			month := addMonth(lastMonth, count)
			//fmt.Println(month)
			if month.Sub(firstMonth) >= 0 {
				break
			}
			count++
			//fmt.Println(count)
		}
	} else {
		for {
			month := addMonth(lastMonth, count)
			//fmt.Println(month)
			if month.Sub(firstMonth) <= 0 {
				break
			}
			count--
			//fmt.Println(count)
		}
	}
	fmt.Fprintln(outStream, count)
}

func mdateDiffLastDate(lastDate time.Time, opt *option, outStream io.Writer) {
	fmt.Fprintln(outStream, lastDate.Format(layoutDate))
}

func mdateDiffLastMonth(lastMonth time.Time, opt *option, outStream io.Writer) {
	fmt.Fprintln(outStream, lastMonth.Format(layoutMonth))
}
