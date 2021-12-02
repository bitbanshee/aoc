package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// If you change the windowSize to 1, Res will give you the answer of the first part
const windowSize int = 3

func main() {
	report := buildReport()
	start := time.Now()
	increased, decreased := Res(report)
	elapsed := time.Since(start)
	fmt.Printf("Resolution: 'for' recursion and a sum function, time elapsed: %v\nincreased: %v, decreased: %v", elapsed, increased, decreased)
}

func buildReport() []int {
	rawdata, err := os.ReadFile("../input")
	if err != nil {
		panic("no input file")
	}
	sreport := strings.Split(string(rawdata), "\n")
	report := make([]int, len(sreport))
	for index, item := range sreport {
		report[index], err = strconv.Atoi(item)
		if err != nil {
			panic("invalid input value")
		}
	}
	return report
}

// Res uses simple for recursion and a sum function
func Res(report []int) (increased, decreased int) {
	if len(report) < (windowSize + 1) {
		return
	}
	maxIndex := len(report) - windowSize + 1
	for index, lastSum := 1, sumSliceItems(report, windowSize); index < maxIndex; index++ {
		currentSum := sumSliceItems(report[index:], windowSize)
		switch {
		case lastSum < currentSum:
			increased++
		case lastSum > currentSum:
			decreased++
		}
		lastSum = currentSum
	}
	return
}

func sumSliceItems(slice []int, n int) (sum int) {
	for _, item := range slice[:n] {
		sum += item
	}
	return
}
