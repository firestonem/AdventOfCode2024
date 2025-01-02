package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func AddEquations(slice []string) int {
	total := 0
	for _, text := range slice {
		text = strings.TrimLeft(text, "mul(")
		text = strings.TrimRight(text, ")")
		numsStr := strings.Split(text, ",")

		num0, err := strconv.Atoi(numsStr[0])
		check(err)
		num1, err := strconv.Atoi(numsStr[1])
		check(err)

		total += num0 * num1
	}

	return total
}

func stripDontMuls(s []string) []string {
	slice := make([]string, len(s))
	copy(slice, s)
	enabled := true

	for i, val := range slice {
		switch val {
		case "do()":
			enabled = true
			slice[i] = "mul(0,0)"
		case "don't()":
			enabled = false
			slice[i] = "mul(0,0)"
		default:
			if !enabled {
				slice[i] = "mul(0,0)"
			}
		}
	}
	return slice
}

func main() {
	start := time.Now()
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	var inputString string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputString += scanner.Text()
	}

	r, _ := regexp.Compile(`(?:mul\([0-9]{1,3},[0-9]{1,3}\))`)

	equations := r.FindAllString(inputString, -1)
	result := AddEquations(equations)
	elapsed1 := time.Since(start)
	fmt.Println(result, "time:", elapsed1)

	start = time.Now()
	r, _ = regexp.Compile(`(?:mul\([0-9]{1,3},[0-9]{1,3}\))|(?:do\(\))|(?:don\'t\(\))`)

	equationsWithFlags := r.FindAllString(inputString, -1)
	equations = stripDontMuls(equationsWithFlags)
	result = AddEquations(equations)
	elapsed2 := time.Since(start)
	fmt.Println(result, "time:", elapsed2)
}
