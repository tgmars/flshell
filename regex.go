package main

import "regexp"

func checkReErr(e error) {
	if e != nil {
		writeStringToFile("RegexErrors.txt", e.Error())
	}
}

func inodeMatcher(line string) string {
	//re, err := regexp.Compile(`(?m).\/.\s(\*?\s*.*):\s.*`)
	// While go is interpreting static input with a \t literal rather than
	// space, regex is modified to catch this.
	re, err := regexp.Compile(`(?m).\/.\s(\*?\s*.*):\s.*`)
	checkReErr(err)
	return re.FindStringSubmatch(line)[1]
}

func dirMatcher(line string) string {
	// Modified regex to match on characters inserted before the fls output
	// Only applicable for println version where command to be run is specified by a
	// character insert before the fls string to be targeted.
	re, err := regexp.Compile(`(?m).*(.\/.)\s\*?\s*.*:.*`)
	//re, err := regexp.Compile(`(?m).*(.\/.)\s\*?\s*.*:\s.*`)
	//re, err := regexp.Compile(`(?m)(.\/.)\s\*?\s*.*:\s.*`)
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
