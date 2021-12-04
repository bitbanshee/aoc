package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Diagnostic struct {
	𝛾, ε, oxygen, co2 uint
}

func (d Diagnostic) PowerConsumption() uint {
	return d.ε * d.𝛾
}

func (d Diagnostic) LifeSupport() uint {
	return d.oxygen * d.co2
}

func (d *Diagnostic) CalcRates(s [][]string) {
	counts := countDigits(s)

	// 𝛾 and ε
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

	d.𝛾 = parseUint(s𝛾)
	d.ε = parseUint(sε)

	// oxygen and co2
	var oxygens [][]string
	var co2s [][]string
	for cIndex := range s[0] {
		// extracts lines from s only for column 0
		if cIndex == 0 {
			count := counts[cIndex]
			mostCommon := mostCommonDigit(count)
			for lIndex, line := range s {
				if mostCommon == "" {
					if s[lIndex][cIndex] == "1" {
						oxygens = append(oxygens, line)
					} else {
						co2s = append(co2s, line)
					}
				} else {
					if s[lIndex][cIndex] == mostCommon {
						oxygens = append(oxygens, line)
					} else {
						co2s = append(co2s, line)
					}
				}
			}
			continue
		}

		if len(oxygens) > 1 {
			_oxygens := filterBy(true, cIndex, oxygens)
			if len(_oxygens) > 0 {
				oxygens = _oxygens
			} else {
				oxygens = oxygens[:1]
			}
		}

		if len(co2s) > 1 {
			_co2s := filterBy(false, cIndex, co2s)
			if len(_co2s) > 0 {
				co2s = _co2s
			} else {
				co2s = co2s[:1]
			}
		}

		if len(oxygens) == 1 && len(co2s) == 1 {
			break
		}
	}

	d.oxygen = parseUint(strings.Join(oxygens[0], ""))
	d.co2 = parseUint(strings.Join(co2s[0], ""))
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
	fmt.Printf("Oxygen: %b, %v\n", diagnostic.oxygen, diagnostic.oxygen)
	fmt.Printf("CO2: %b, %v\n", diagnostic.co2, diagnostic.co2)
	fmt.Printf("Life Support: %b, %v\n", diagnostic.LifeSupport(), diagnostic.LifeSupport())
	fmt.Printf("Elapsed time: %v", elapsed)
}

func parseUint(s string) uint {
	r, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		panic("invalid data for " + s)
	}
	return uint(r)
}

func filterBy(mostCommon bool, column int, m [][]string) [][]string {
	counts := countDigits(m)
	mostCommonDigit := mostCommonDigit(counts[column])
	var _m [][]string
	for _, line := range m {
		if mostCommon {
			if mostCommonDigit == "" &&
				line[column] == "1" ||
				line[column] == mostCommonDigit {
				_m = append(_m, line)
			}
		} else {
			if mostCommonDigit == "" {
				if line[column] == "0" {
					_m = append(_m, line)
				}
			} else {
				if line[column] != mostCommonDigit {
					_m = append(_m, line)
				}
			}
		}
	}
	return _m
}

func mostCommonDigit(m map[string]int) (mostCommon string) {
	switch {
	case m["0"] > m["1"]:
		mostCommon = "0"
	case m["1"] > m["0"]:
		mostCommon = "1"
	}
	return
}

func countDigits(s [][]string) map[int]map[string]int {
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
	return counts
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
