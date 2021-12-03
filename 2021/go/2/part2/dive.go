/*
I tried 2 different solutions just to see which one is more performant
since I'm new to Go language :P
*/
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Direction string

type Command struct {
	direction Direction
	units     uint8
}

type Plan []Command

type Position struct {
	horizontal int
	depth      int
	aim        int
}

type Course []Position

const (
	Forward Direction = "forward"
	Down    Direction = "down"
	Up      Direction = "up"
)

func main() {
	initial := new(Position)
	plan := plan()
	start := time.Now()
	course := follow(initial, plan)
	elapsed := time.Since(start)
	fmt.Println(course)
	last := &course[len(course)-1]
	fmt.Printf("Tracking the course results:\n"+
		"time elapsed: %v.\n"+
		"last position: %v.\n"+
		"result: %v.",
		elapsed, last, last.depth*last.horizontal)
	start = time.Now()
	last = new(Position)
	follow2(last, plan)
	elapsed = time.Since(start)
	fmt.Printf("\n\nNot tracking the course results:\n"+
		"time elapsed: %v.\n"+
		"last position: %v.\n"+
		"result: %v.",
		elapsed, last, last.depth*last.horizontal)
}

// follow applies a Plan to an initial position, tracking each step by
// creating fresh positions, and returns a Course
func follow(position *Position, plan Plan) Course {
	course := make(Course, len(plan)+1)
	currentPosition := position
	course[0] = *currentPosition
	for index, command := range plan {
		currentPosition = submit(currentPosition, &command)
		course[index+1] = *currentPosition
	}
	return course
}

func submit(position *Position, command *Command) *Position {
	switch command.direction {
	case Forward:
		return &Position{
			position.horizontal + int(command.units),
			position.aim*int(command.units) + position.depth,
			position.aim,
		}
	case Down:
		return &Position{
			position.horizontal,
			position.depth,
			position.aim + int(command.units),
		}
	case Up:
		return &Position{
			position.horizontal,
			position.depth,
			position.aim - int(command.units),
		}
	default:
		return position
	}
}

// follow2 applies a Plan to an initial position modifying the position
func follow2(position *Position, plan Plan) {
	for _, command := range plan {
		submit2(position, &command)
	}
}

func submit2(position *Position, command *Command) {
	switch command.direction {
	case Forward:
		position.horizontal += int(command.units)
		position.depth = position.aim*int(command.units) + position.depth
	case Down:
		position.aim += int(command.units)
	case Up:
		position.aim -= int(command.units)
	}
}

func plan() Plan {
	rawdata, err := os.ReadFile("../input")
	if err != nil {
		panic("no input file")
	}
	scommands := strings.Split(string(rawdata), "\n")
	commands := make(Plan, len(scommands))
	for index, scommand := range scommands {
		directionUnits := strings.Split(scommand, " ")
		if len(directionUnits) < 2 {
			panic(fmt.Sprintf("invalid input value: %v", scommand))
		}
		units, err := strconv.ParseUint(directionUnits[1], 10, 8)
		if err != nil {
			panic("invalid input value")
		}
		commands[index] = Command{
			Direction(directionUnits[0]),
			uint8(units),
		}
	}
	return commands
}
