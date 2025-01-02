package main

import (
	"bufio"
	"fmt"
	"log"
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

func getInput(inPath string) []string {
	file, err := os.Open(inPath)
	check(err)
	defer file.Close()

	input := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	check(scanner.Err())
	return input
}

func combineRules(rulesLines []string) map[string][]string {
	rules := make(map[string][]string)

	for _, line := range rulesLines {
		pair := strings.Split(line, `|`)
		key := pair[0]
		val := pair[1]

		// check for init
		if len(rules[key]) == 0 {
			rules[key] = make([]string, 0)
		}

		// don't dupe
		if slices.Contains(rules[key], val) {
			continue
		}

		// append if we got here
		rules[key] = append(rules[key], val)
	}

	return rules
}

func validateUpdate(update []string, rules map[string][]string) bool {
	for i, currVal := range update {
		pagesBefore := update[:i]

		for _, checkVal := range pagesBefore {
			if slices.Contains(rules[currVal], checkVal) {
				return false
			}
		}
	}
	return true
}

func categorizeUpdates(updatesLines []string, rules map[string][]string) ([][]string, [][]string) {

	validUpdates := [][]string{}
	invalidUpdates := [][]string{}

	for _, line := range updatesLines {
		update := strings.Split(line, `,`)
		valid := validateUpdate(update, rules)
		if !valid {
			fmt.Println(update, "invalid")
			invalidUpdates = append(invalidUpdates, update)
			continue
		}
		fmt.Println(update, "valid")
		validUpdates = append(validUpdates, update)
	}

	return validUpdates, invalidUpdates
}

func addMiddleNums(updates [][]string) int {
	result := 0

	for _, update := range updates {
		length := len(update)
		if length%2 == 0 {
			log.Panic(update, "not odd!")
		}
		middleStr := update[(length-1)/2]
		middleVal, err := strconv.Atoi(middleStr)
		check(err)
		fmt.Println("adding", middleVal)
		result += middleVal
	}

	return result
}

func fixInvalidUpdates(invalids [][]string, rules map[string][]string) [][]string {
	fixed := make([][]string, 0)

	for _, update := range invalids {
		//fmt.Println("fixing", update)

		// initialize
		fixedUpdate := make([]string, 1)

		for _, v := range update {

			for i := 0; i < len(fixedUpdate); i++ {
				fixedUpdate = slices.Insert(fixedUpdate, i, v)

				// check then re-insert in next slot if not valid
				isInPosition := validateUpdate(fixedUpdate, rules)
				if isInPosition {
					//fmt.Println("valid")
					break
				}
				//fmt.Println("invalid")
				fixedUpdate = slices.Delete(fixedUpdate, i, i+1)
			}
		}
		fixedUpdate = fixedUpdate[0 : len(fixedUpdate)-1]
		fixed = append(fixed, fixedUpdate)
	}
	fmt.Println(fixed)
	return fixed
}

func main() {
	start := time.Now()

	in := getInput("input.bak")

	rulesLines := []string{}
	updatesLines := []string{}

	for _, line := range in {
		if strings.Contains(line, `|`) {
			rulesLines = append(rulesLines, line)
		}
		if strings.Contains(line, `,`) {
			updatesLines = append(updatesLines, line)
		}
	}
	rules := combineRules(rulesLines)
	validUpdates, invalidUpdates := categorizeUpdates(updatesLines, rules)
	fmt.Println("part 1:", addMiddleNums(validUpdates))

	fixedUpdates := fixInvalidUpdates(invalidUpdates, rules)
	fmt.Println("part 2:", addMiddleNums(fixedUpdates))

	fmt.Println("completed in", time.Since(start))
}
