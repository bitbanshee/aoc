package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const inputFile string = "../input"

type Board [][]Cell

func (b Board) Mark(n uint8) {
	for lIndex := 0; lIndex < len(b); lIndex++ {
		for cIndex := 0; cIndex < len(b[lIndex]); cIndex++ {
			if b[lIndex][cIndex].value == n {
				b[lIndex][cIndex].marked = true
				return
			}
		}
	}
}

func (b Board) Check() bool {
	colsWin := make([]bool, len(b[0]))
	for i := range colsWin {
		colsWin[i] = true
	}
	for _, line := range b {
		lineWin := true
		for cIndex, cell := range line {
			if !cell.marked {
				lineWin = false
				colsWin[cIndex] = false
			}
		}
		if lineWin {
			return true
		}
	}
	for _, win := range colsWin {
		if win {
			return true
		}
	}
	return false
}

type Cell struct {
	value  uint8
	marked bool
}

type WinnerBoard struct {
	value uint8
	board Board
}

func main() {
	drawn, boards := data()
	drawCh := make(chan uint8)
	go func() {
		for _, n := range drawn {
			drawCh <- n
		}
		close(drawCh)
	}()
	drawChs := multiplex(drawCh, len(boards))
	var winChannels []<-chan WinnerBoard
	for index, board := range boards {
		winChannel := bingo(drawChs[index], board)
		winChannels = append(winChannels, winChannel)
	}
	winChannel := merge(winChannels...)
	winnerBoard := <-winChannel
	sum := sumNotMarked(winnerBoard.board)
	fmt.Println(sum * int(winnerBoard.value))
}

func sumNotMarked(b Board) (sum int) {
	for _, line := range b {
		for _, cell := range line {
			if !cell.marked {
				sum += int(cell.value)
			}
		}
	}
	return
}

func multiplex(ch chan uint8, n int) []chan uint8 {
	var channels []chan uint8
	for i := 0; i < n; i++ {
		channels = append(channels, make(chan uint8))
	}
	go func() {
		for n := range ch {
			for _, _ch := range channels {
				_ch <- n
			}
		}

		for _, _ch := range channels {
			close(_ch)
		}
	}()
	return channels
}

func merge(chs ...<-chan WinnerBoard) <-chan WinnerBoard {
	ch := make(chan WinnerBoard)
	for _, c := range chs {
		go func(c <-chan WinnerBoard) {
			for w := range c {
				ch <- w
			}
		}(c)
	}
	return ch
}

func bingo(drawn <-chan uint8, board Board) <-chan WinnerBoard {
	win := make(chan WinnerBoard)
	go func() {
		// optimization
		alreadyDrawn := 0
		for n := range drawn {
			alreadyDrawn++
			board.Mark(n)
			if alreadyDrawn > 5 && board.Check() {
				win <- WinnerBoard{n, board}
				break
			}
		}
		close(win)
	}()
	return win
}

func data() (drawn []uint8, boards []Board) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	sdata := string(data)
	splitted := strings.SplitN(sdata, "\n", 2)
	sdrawn := strings.Split(splitted[0], ",")
	for _, sn := range sdrawn {
		dvalue, err := strconv.ParseUint(sn, 10, 8)
		if err != nil {
			panic(err)
		}
		drawn = append(drawn, uint8(dvalue))
	}
	rboards := strings.Split(splitted[1], "\n")
	blanks := regexp.MustCompile(`\s+`)
	sblanks := regexp.MustCompile(`^\s+`)
	var currentBoard Board
	for _, line := range rboards {
		if len(line) == 0 {
			continue
		}

		rcells := strings.Split(
			sblanks.ReplaceAllString(
				blanks.ReplaceAllString(line, " "), ""), " ")

		var cells []Cell
		for _, rcell := range rcells {
			cvalue, err := strconv.ParseUint(rcell, 10, 8)
			if err != nil {
				panic(err)
			}
			cells = append(cells, Cell{value: uint8(cvalue)})
		}
		currentBoard = append(currentBoard, cells)

		if len(currentBoard) == 5 {
			boards = append(boards, currentBoard)
			currentBoard = nil
		}
	}
	return
}
