package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// takes arranged matrix and outputs slice
func convertToSlice(m [][]byte) []string {
	output := []string{}

	for _, line := range m {
		output = append(output, string(line))
	}
	return output
}

// takes pre-arranged slice of strings to count XMASes
func countXmas(x []string) int {
	r, _ := regexp.Compile(`(?:XMAS)`)
	r2, _ := regexp.Compile(`(?:SAMX)`)
	count := 0
	for _, line := range x {
		fmt.Print(line)
		finds := r.FindAllString(line, -1)
		finds = append(finds, r2.FindAllString(line, -1)...)
		if len(finds) > 0 {
			fmt.Print(" - found ", len(finds))
		}
		fmt.Println()
		count += len(finds)
	}
	fmt.Print("\n")
	fmt.Println("---")
	return count
}

func transposeMatrix(orig [][]byte) [][]byte {
	rows, cols := len(orig), len(orig[0])
	matrix := make([][]byte, cols)

	for i := range cols {
		matrix[i] = make([]byte, rows)
		for j := range rows {
			matrix[i][j] = orig[j][i]
		}
	}

	return matrix
}

func reverseArray(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// https://stackoverflow.com/a/8664879
func rotateMatrix(orig [][]byte) [][]byte {
	rot := transposeMatrix(orig)

	for i := range rot {
		rot[i] = reverseArray(rot[i])
	}

	return rot
}

// top-left to bottom-right diagonals
func rotate45Degrees(m [][]byte) [][]byte {
	rows, cols := len(m), len(m[0])
	output := make([][]byte, cols)

	// start at bottom and move upward with each iteration,
	// starting positions from [0,n] to [0,0]
	for j := rows - 1; j >= 0; j-- {
		line := make([]byte, rows)

		for i := 0; i <= cols-1; i++ {
			if i+j > cols-1 {
				continue
			}
			line[i] = m[j+i][i]
		}

		output = append(output, line)

	}

	// move left and grab each diagonal
	// starting positions from (0,n) to (0,1)
	for i := len(m[0]) - 1; i >= 1; i-- {
		line := []byte{}
		for j := 0; j <= len(m[0])-1-i; j++ {
			// if i+j > len(m[j])-1 {
			// 	continue
			// }
			line = append(line, m[j][i+j])
		}
		output = append(output, line)
	}

	return output
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	input := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	check(scanner.Err())

	// turn input into byte matrix to make working with it easier
	inputCols := len(input[0])
	inputRows := len(input)

	m := make([][]byte, inputRows)
	for i := range inputRows {
		m[i] = make([]byte, inputCols)
		for j := range inputCols {
			m[i][j] = byte(input[i][j])
		}
	}

	count := 0

	// count L to R
	//    && R to L
	count += countXmas(convertToSlice(m))

	// count top-left to bottom-right
	//    && bottom-right to top-left
	count += countXmas(convertToSlice(rotate45Degrees(m)))

	// count up-to-down
	//    && down-to-up
	rotatedMatrix := rotateMatrix(m)
	count += countXmas(convertToSlice(rotatedMatrix))

	// count top-right to bottom-left
	//    && bottom-left to top-right
	count += countXmas(convertToSlice(rotate45Degrees(rotatedMatrix)))

	fmt.Println(count, "in", time.Since(start))
}
