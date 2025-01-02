package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addLinetoSlices(x []int, y []int, s string) ([]int, []int) {
	strSlice := strings.Split(s, "   ")

	num0, err := strconv.Atoi(string(strSlice[0]))
	check(err)
	num1, err := strconv.Atoi(string(strSlice[1]))
	check(err)

	return append(x, num0), append(y, num1)
}

func diffSlices(x []int, y []int) int {
	total := 0

	for i := 0; i < len(x); i++ {
		total += int(math.Abs(float64(x[i] - y[i])))
	}

	return total
}

func calcSimScore(x []int, y []int) int {
	total := 0

	// loop over each value of first slice
	for _, i := range x {
		multiplier := 0
		// loop until not found in second slice
		for {
			n, found := slices.BinarySearch(y, i)
			if found != true {
				break
			}
			multiplier += 1
			y = slices.Delete(y, n, n+1)
		}
		total += multiplier * i
	}

	return total
}

func main() {
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	s1 := []int{}
	s2 := []int{}
	for scanner.Scan() {
		s1, s2 = addLinetoSlices(s1, s2, scanner.Text())
	}

	check(scanner.Err())

	slices.Sort(s1)
	slices.Sort(s2)

	// fmt.Println("part 1 diff: ", diffSlices(s1, s2))
	fmt.Println("part 2 score:", calcSimScore(s1, s2))
}
