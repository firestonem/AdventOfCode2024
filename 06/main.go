package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

// define runes for easy references
var exitRune = []rune("E")[0]
var guardRune = []rune("^")[0]
var obstacleRune = []rune("#")[0]
var openRune = []rune(".")[0]

type Coordinates struct {
	x int
	y int
}

type Part2Record struct {
	position    Coordinates
	orientation string
}

type Guard struct {
	orientation         string
	startingOrientation string
	position            Coordinates
	startingPosition    Coordinates
	visited             []Coordinates
	fullHistory         []Part2Record
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// pointer receiver to turn guard's direction
func (g *Guard) turn90Degrees() {
	switch g.orientation {
	case "N":
		g.orientation = "E"
	case "E":
		g.orientation = "S"
	case "S":
		g.orientation = "W"
	case "W":
		g.orientation = "N"
	}
}

func (g Guard) getNextPosition() Coordinates {
	nextPosition := Coordinates{
		x: 0,
		y: 0,
	}
	switch g.orientation {
	case "N":
		nextPosition.y -= 1
	case "S":
		nextPosition.y += 1
	case "E":
		nextPosition.x += 1
	case "W":
		nextPosition.x -= 1
	}

	nextPosition.x += g.position.x
	nextPosition.y += g.position.y

	return nextPosition
}

// eval next spot on grid
func (g Guard) getNextSymbol(m [][]rune) rune {
	nextPosition := g.getNextPosition()
	// fmt.Println("next position:", nextPosition)

	// check for exit
	yAxisLen := len(m)
	xAxisLen := len(m[0])
	if nextPosition.x < 0 || nextPosition.y < 0 || nextPosition.x == xAxisLen || nextPosition.y == yAxisLen {
		return exitRune
	}

	return m[nextPosition.y][nextPosition.x]
}

// guard takes a step forward
func (g *Guard) takeStep(m [][]rune) [][]rune {
	nextPosition := g.getNextPosition()

	// move guard's rune on the matrix
	m[g.position.y][g.position.x] = openRune
	// m[nextPosition.y][nextPosition.x] = guardRune

	// update guard's position
	g.position.x = nextPosition.x
	g.position.y = nextPosition.y

	return m
}

func (g *Guard) updateHistories() {

	if !slices.Contains(g.visited, g.position) {
		g.visited = append(g.visited, g.position)
	}

	g.fullHistory = append(g.fullHistory, Part2Record{
		position:    g.position,
		orientation: g.orientation,
	})
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	input := []string{}

	guard := Guard{
		orientation:         "N",
		startingOrientation: "N",
		position: Coordinates{
			x: 0,
			y: 0,
		},
		startingPosition: Coordinates{
			x: 0,
			y: 0,
		},
		visited: []Coordinates{},
	}

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())

		// set guard's starting position while reading inputs
		if strings.Contains(scanner.Text(), `^`) {
			guard.position.y = lineCount
			guard.position.x = strings.Index(scanner.Text(), `^`)
			guard.visited = append(guard.visited, guard.position)
			guard.startingPosition = guard.position
		}
		lineCount++
	}

	check(scanner.Err())

	// turn input into matrix to make working with it easier
	inputCols := len(input[0])
	inputRows := len(input)

	m := make([][]rune, inputRows)
	for i := range inputRows {
		m[i] = make([]rune, inputCols)
		for j := range inputCols {
			m[i][j] = rune(input[i][j])
		}
	}

part1:
	for {
		// fmt.Println("(", guard.position.x, ",", guard.position.y, ")")
		nextSymbol := guard.getNextSymbol(m)
		// fmt.Println("next symbol:", string(nextSymbol))
		switch nextSymbol {
		case openRune:
			// fmt.Println("step")
			m = guard.takeStep(m)
		case obstacleRune:
			// fmt.Println("turn")
			guard.turn90Degrees()
		case exitRune:
			// fmt.Println("exit")
			break part1
		}
		guard.updateHistories()
	}
	fmt.Println("Guard visited", len(guard.visited), "positions")
	fmt.Println("part1 finished in", time.Since(start))

	start = time.Now()

	m[guard.startingPosition.y][guard.startingPosition.x] = openRune
	m[guard.visited[len(guard.visited)-1].y][guard.visited[len(guard.visited)-1].x] = openRune

	// remove starting position
	guard.visited = guard.visited[1:]
	fmt.Println("Testing", len(guard.visited), "obstacle placements")
	count := 0

	for _, testObstacle := range guard.visited {

		// create new matrix
		matrix := make([][]rune, inputRows)
		copy(matrix, m)
		matrix[testObstacle.y][testObstacle.x] = obstacleRune

		guardClone := Guard{
			orientation:         "N",
			startingOrientation: "N",
			position:            guard.startingPosition,
			startingPosition:    guard.startingPosition,
			visited:             []Coordinates{},
			fullHistory:         []Part2Record{},
		}

		guardClone.visited = slices.Delete(guardClone.visited, 0, len(guardClone.visited))
		guardClone.fullHistory = slices.Delete(guardClone.fullHistory, 0, len(guardClone.fullHistory))

		infinite := false

	addedObstacleLoop:
		for {
			// fmt.Println("(", guard.position.x, ",", guard.position.y, ")")
			nextSymbol := guardClone.getNextSymbol(matrix)
			// fmt.Println("next symbol:", string(nextSymbol))
			switch nextSymbol {
			case openRune:
				// fmt.Println("step")
				matrix = guardClone.takeStep(matrix)
			case obstacleRune:
				// fmt.Println("turn")
				guardClone.turn90Degrees()
			case exitRune:
				// fmt.Println("exit")

				// restore runes when complete
				matrix[testObstacle.y][testObstacle.x] = openRune
				// matrix[guard.startingPosition.y][guard.startingPosition.x] = guardRune
				// matrix[guardClone.visited[len(guardClone.visited)-1].y][guardClone.visited[len(guardClone.visited)-1].x] = openRune
				// fmt.Println("--- not infinite")

				break addedObstacleLoop
			}

			record := Part2Record{
				position:    guardClone.position,
				orientation: guardClone.orientation,
			}
			// if this record is already in the history, it is a loop
			if slices.Contains(guardClone.fullHistory, record) {
				infinite = true

				// restore runes when complete
				matrix[testObstacle.y][testObstacle.x] = openRune
				// matrix[guard.startingPosition.y][guard.startingPosition.x] = guardRune
				// matrix[guardClone.visited[len(guardClone.visited)].y][guardClone.visited[len(guardClone.visited)-1].x] = openRune

				// fmt.Println("--- infinite at ", guardClone.position, guardClone.orientation)
				break addedObstacleLoop
			}

			guardClone.updateHistories()
		}

		if infinite == true {
			count += 1
			// fmt.Println("+1 to count")
		}
		// fmt.Println(guardClone.fullHistory)
	}

	fmt.Println("created", count, "infinite routes")
	fmt.Println("part2 finished in", time.Since(start))

}
