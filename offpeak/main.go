package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var times = []string{
	"23:30",
	"00:00",
	"00:30",
	"01:00",
	"01:30",
	"02:00",
	"02:30",
	"03:00",
	"03:30",
	"04:00",
	"04:30",
	"05:00",
}

var gmt = times[2:11]
var bst = times[0:9]
var usage float64

func findByExtension(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func isOffpeak(times []string, line string) bool {
	for _, s := range times {
		if strings.Contains(line, s) {
			return true
		}
	}
	return false
}

func withinDates(startDate string, endDate string, line string) bool {
	date := line[1:17]
	if date < startDate {
		return false
	}
	if date > endDate {
		return false
	}
	return true
}

func addUsage(line string) {
	index := strings.Index(line, ",\"")
	u := line[index+2 : len(line)-1]
	f, err := strconv.ParseFloat(u, 32)
	if err != nil {
		fmt.Printf("Unable to parse %s\n", u)
	} else {
		usage = usage + f
	}
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage main.go '2023-01-27 23:00' '2023-03-21 05:00'")
		return
	}

	startDate := os.Args[1]
	endDate := os.Args[2]

	var offpeak []string
	if time.Now().IsDST() {
		offpeak = bst
	} else {
		offpeak = gmt
	}

	files := findByExtension(".", ".csv")
	if len(files) == 0 {
		fmt.Println("No .csv files found")
		return
	}

	file := files[0]
	fmt.Printf("Using %s\n", file)

	readFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if isOffpeak(offpeak, line) && withinDates(startDate, endDate, line) {
			fmt.Println(line)
			addUsage(line)
		}
	}
	fmt.Printf("Usage: %f kWH\n", usage)

}
