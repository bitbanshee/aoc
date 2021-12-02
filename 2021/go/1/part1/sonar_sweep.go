/*
I tried 3 different solutions just to see which one is more performant
since I'm new to the Go language :P
*/
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ClosureIterator func() (item int, any bool)

func main() {
	report := buildReport()
	start := time.Now()
	increased, decreased := Res1(report)
	elapsed := time.Since(start)
	fmt.Printf("Resolution 1: 'for' recursion with 'range' iterator, time elapsed: %v\nincreased: %v, decreased: %v", elapsed, increased, decreased)
	start = time.Now()
	increased, decreased = Res2(report)
	elapsed = time.Since(start)
	fmt.Printf("\nResolution 2: 'func' recursion and a 'chan' with goroutine based iterator, time elapsed: %v\nincreased: %v, decreased: %v", elapsed, increased, decreased)
	start = time.Now()
	increased, decreased = Res3(report)
	elapsed = time.Since(start)
	fmt.Printf("\nResolution 3: 'func' recursion and a closure based iterator, time elapsed: %v\nincreased: %v, decreased: %v", elapsed, increased, decreased)
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

// Res1 uses simple for recursion with range iterator
func Res1(report []int) (increased, decreased int) {
	if len(report) < 2 {
		return
	}
	for last, index := report[0], 1; index < len(report); last, index = report[index], index+1 {
		switch {
		case last < report[index]:
			increased++
		case last > report[index]:
			decreased++
		}
	}
	return
}

// Res2 uses func recursion and a chan with goroutine based iterator
func Res2(report []int) (increased, decreased int) {
	if len(report) < 2 {
		return
	}
	return _res2(iterator(report[2:]), report[0], report[1])
}

func _res2(iterator <-chan int, last, current int) (increased, decreased int) {
	switch {
	case last < current:
		increased++
	case last > current:
		decreased++
	}
	next, ok := <-iterator
	if !ok {
		return
	}
	inc, dec := _res2(iterator, current, next)
	return increased + inc, decreased + dec
}

func iterator(slice []int) <-chan int {
	ch := make(chan int)
	go func() {
		for _, item := range slice {
			ch <- item
		}
		close(ch)
	}()
	return (<-chan int)(ch)
}

// Res3 uses func recursion and a closure based iterator
func Res3(report []int) (increased, decreased int) {
	if len(report) < 2 {
		return
	}
	return _res3(iterator2(report[2:]), report[0], report[1])
}

func _res3(iterator ClosureIterator, last, current int) (increased, decreased int) {
	switch {
	case last < current:
		increased++
	case last > current:
		decreased++
	}
	next, ok := iterator()
	if !ok {
		return
	}
	inc, dec := _res3(iterator, current, next)
	return increased + inc, decreased + dec
}

func iterator2(slice []int) ClosureIterator {
	index := 0
	return func() (item int, any bool) {
		if index < len(slice) {
			// we can do the index incrementation using a deferred function,
			// but it's (a little bit) slower
			index++
			return slice[index-1], true
		}
		return 0, false
	}
}
