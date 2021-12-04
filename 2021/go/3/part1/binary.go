package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Diagnostic struct {
	𝛾, ε uint
}

func (d Diagnostic) PowerConsumption() uint {
	return d.ε * d.𝛾
}

func (d *Diagnostic) CalcRates(s [][]string) {
	counts := make(map[int]map[string]int)
	for _, line := range s {
		for colInd, digit := range line {
			mvalue, ok := counts[colInd]
			if !ok {
				counts[colInd] = map[string]int{
					digit: 1,
				}
				continue
			}
			_, ok = mvalue[digit]
			if !ok {
				mvalue[digit] = 1
				continue
			}
			mvalue[digit] += 1
		}
	}

	var s𝛾 string
	var sε string
	for index := range s[0] {
		count := counts[index]
		if count["0"] > count["1"] {
			s𝛾 += "0"
			sε += "1"
		} else {
			s𝛾 += "1"
			sε += "0"
		}
	}

	r, err := strconv.ParseUint(s𝛾, 2, 64)
	if err != nil {
		panic("invalid data")
	}
	d.𝛾 = uint(r)

	r, err = strconv.ParseUint(sε, 2, 64)
	if err != nil {
		panic("invalid data")
	}
	d.ε = uint(r)
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("{%b %b}", d.𝛾, d.ε)
}

func main() {
	data := data()
	var diagnostic Diagnostic
	start := time.Now()
	diagnostic.CalcRates(data)
	elapsed := time.Since(start)
	fmt.Printf("𝛾: %b, %v\n", diagnostic.𝛾, diagnostic.𝛾)
	fmt.Printf("ε: %b, %v\n", diagnostic.ε, diagnostic.ε)
	fmt.Printf("Power Consumption: %b, %v\n", diagnostic.PowerConsumption(), diagnostic.PowerConsumption())
	fmt.Printf("Elapsed time: %v", elapsed)
}

func data() [][]string {
	data, err := os.ReadFile("../input")
	if err != nil {
		panic("no data available")
	}
	lines := strings.Split(string(data), "\n")
	var matrix [][]string
	for _, line := range lines {
		matrix = append(matrix, strings.Split(line, ""))
	}
	return matrix
}
