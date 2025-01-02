package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addLineToMatrix(m [][]int, line string) [][]int {
	strSlice := strings.Split(line, " ")

	intSlice := []int{}
	for _, val := range strSlice {
		i, err := strconv.Atoi(val)
		check(err)
		intSlice = append(intSlice, i)
	}

	return append(m, intSlice)
}

func isSorted(s []int) bool {
	slice := make([]int, len(s))
	copy(slice, s)
	isIncreasing := slices.IsSorted(slice)

	slices.Reverse(slice)
	isDecreasing := slices.IsSorted(slice)

	return (isIncreasing || isDecreasing)
}

func isSmallDiff(s []int) bool {
	for i := 0; i < len(s)-1; i++ {
		diff := int(math.Abs(float64(s[i] - s[i+1])))
		if diff > 3 || diff == 0 {
			return false
		}
	}
	return true
}

func isSafeWithDampener(s []int) bool {
	for i := 0; i < len(s); i++ {
		slice := make([]int, len(s))
		copy(slice, s)
		slice = append(slice[:i], slice[i+1:]...)
		if isSorted(slice) && isSmallDiff(slice) {
			return true
		}
	}
	return false
}

func checkSafety(m [][]int) (int, [][]int) {
	numSafe := 0
	retries := [][]int{}

	for _, line := range m {
		if !isSorted(line) || !isSmallDiff(line) {
			retries = append(retries, line)
			continue
		}
		numSafe += 1
	}

	return numSafe, retries
}

func checkSafetyWithDampener(m [][]int) int {
	numSafe := 0

	for _, line := range m {
		if !isSafeWithDampener(line) {
			continue
		}
		numSafe += 1
	}

	return numSafe
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	m := [][]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m = addLineToMatrix(m, scanner.Text())
	}
	check(scanner.Err())

	safe1, retries := checkSafety(m)
	elapsed1 := time.Since(start)
	safe2 := checkSafetyWithDampener(retries)
	elapsed2 := time.Since(start)
	fmt.Println("safe #1:", safe1, "time:", elapsed1.Milliseconds(), "ms")
	fmt.Println("safe #2:", safe1+safe2, "time:", elapsed2.Milliseconds(), "ms")
}
