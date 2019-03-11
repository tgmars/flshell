package main

import "regexp"

func checkReErr(e error) {
	if e != nil {
		writeStringToFile("RegexErrors.txt", e.Error())
	}
}

func inodeMatcher(line string) string {
	re, err := regexp.Compile(`(?m).\/.\s(\*?\s*.*):\s.*`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func dirMatcher(line string) string {
	re, err := regexp.Compile(`(?m)(.\/.)\s\*?\s*.*:\s.*`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func nameMatcher(line string) string {
	re, err := regexp.Compile(`(?m):\s+(.+)`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func newlineCounter(input string) int {
	re, err := regexp.Compile(`(?m)\n`)
	checkReErr(err)
	return len(re.FindAllStringIndex(input, -1))
}
