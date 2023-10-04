package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

var peak float64
var offPeak float64

// ---------------------

type consumptionData struct {
	ts          time.Time
	consumption float64
}

func newComsumptionData(line string) consumptionData {
	cd := consumptionData{}
	t, err := time.ParseInLocation("2006-01-02 15:04", line[1:17], time.UTC)
	if err != nil {
		panic(err)
	}
	cd.ts = t.Add(-30 * time.Minute) // move time from end of period to start

	index := strings.Index(line, ",\"")
	u := line[index+2 : len(line)-1]
	f, err := strconv.ParseFloat(u, 32)
	if err != nil {
		fmt.Printf("Unable to parse %s\n", u)
	} else {
		cd.consumption = f
	}

	return cd
}

func (cd consumptionData) OffPeak() bool {
	return cd.ts.Local().Hour() < 5
}

func (cd consumptionData) WithinDates(start time.Time, end time.Time) bool {
	return cd.ts.After(start) && cd.ts.Before(end)
}

func (cd consumptionData) String() string {
	return fmt.Sprintf("%s, %f", cd.ts.Local().String(), cd.consumption)
}

// ---------------------

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

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage main.go '20230127' '20230321'")
		return
	}

	startDate, err := time.ParseInLocation("20060102", os.Args[1], time.Local)
	if err != nil {
		panic(err)
	}
	endDate, err := time.ParseInLocation("20060102", os.Args[2], time.Local)
	if err != nil {
		panic(err)
	}
	// adjust dates to make them inclusive
	startDate = startDate.Add(-1 * time.Minute)
	endDate = endDate.Add(24 * time.Hour)

	files := findByExtension(".", ".csv")
	if len(files) == 0 {
		fmt.Println("No .csv files found")
		return
	}

	slices.Sort(files)
	slices.Reverse(files)
	file := files[0]
	fmt.Printf("%v\n", files)
	fmt.Printf("Using %s\n", file)

	readFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan() // read header line
	for fileScanner.Scan() {
		line := fileScanner.Text()
		cd := newComsumptionData(line)
		if cd.WithinDates(startDate, endDate) {
			if cd.OffPeak() {
				fmt.Println(cd)
				offPeak += cd.consumption
			} else {
				peak += cd.consumption
			}

		}
	}
	fmt.Printf("\nOff peak consumption: %f kW\n", offPeak)
	fmt.Printf("Peak consumption: %f kW\n", peak)
	fmt.Printf("Total consumption: %f kW\n\n", offPeak+peak)

	fmt.Println("Start time is " + startDate.String())
	fmt.Println("End time is " + endDate.String())

}
