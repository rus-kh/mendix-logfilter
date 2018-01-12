package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/juju/ratelimit"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	drops    int
	flagRate int64
	flagKeep bool
	flagFile string
)

func init() {
	rate, _ := strconv.Atoi(os.Getenv("LOG_RATELIMIT"))
	if rate == 0 {
		rate = 1000
	}
	flag.Int64Var(&flagRate, "r", int64(rate), "limit to r messages per second (drops those exceeding the limit)")
	flag.BoolVar(&flagKeep, "k", false, "keep the messages instead of dropping them")
	flag.StringVar(&flagFile, "f", "", "define a file as input")
	flag.Parse()
}

func reset() {
	if drops > 0 {
		fmt.Fprintf(os.Stderr, "Rate-limiting to %d loglines/second. Suppressed %d.\n", flagRate, drops)
		drops = 0
	}
}

func parseMendixLogline(line string) string {
	re := regexp.MustCompile("^[0-9-]+\\s[0-9:\\.]+\\s")
	return re.ReplaceAllString(line, "")
}

func openFile(file string) (*os.File, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open \"%s\" (file does not exist)", file))
	}

	input, err := os.Open(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open \"%s\" (%s)", file, err))
	}
	return input, nil
}

func main() {
	var input *os.File
	if flagFile != "" {
		var err error
		input, err = openFile(flagFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		input = os.Stdin
	}
	scanner := bufio.NewScanner(input)
	c := make(chan string)
	done := make(chan bool)
	go func(c chan string, done chan bool) {
		for scanner.Scan() {
			c <- scanner.Text()
		}
		close(done)
	}(c, done)
	l := ratelimit.NewBucket(time.Second/time.Duration(flagRate), flagRate)
	timer := time.NewTicker(time.Second)
	for {
		select {
		case line := <-c:
			if flagKeep {
				l.Wait(1)
				fmt.Println(parseMendixLogline(line))
			} else {
				if l.TakeAvailable(1) > 0 {
					fmt.Println(parseMendixLogline(line))
					reset()
				} else {
					drops++
				}
			}
		case <-timer.C:
			if l.Available() == flagRate {
				reset()
			}
		case <-done:
			reset()
			return
		}
	}
}
