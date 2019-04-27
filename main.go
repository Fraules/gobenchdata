package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bobheadxi/gobenchdata/bench"
	"github.com/spf13/pflag"
)

// Version is the version of gobenchdata
var Version string

var (
	jsonOut   = pflag.String("json", "", "output as json to file")
	appendOut = pflag.BoolP("append", "a", false, "append to output file")

	version = pflag.StringP("version", "v", "", "version to tag in your benchmark output")
	date    = pflag.StringP("date", "d", time.Now().UTC().String(), "date of this run, defaults to UTC time.Now()")
	tags    = pflag.StringArrayP("tag", "t", nil, "array of tags to include in result")
)

func main() {
	pflag.Parse()
	if len(pflag.Args()) > 0 {
		switch cmd := pflag.Args()[0]; cmd {
		case "version":
			if Version == "" {
				println("gobenchdata version unknown")
			} else {
				println("gobenchdata " + Version)
			}
		case "help":
			showHelp()
		case "merge":
			args := pflag.Args()[1:]
			if len(args) < 1 {
				panic("no merge targets provided")
			}
			merge(args...)
		default:
			showHelp()
			os.Exit(1)
		}
		return
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	} else if fi.Mode()&os.ModeNamedPipe == 0 {
		showHelp()
		panic("gobenchdata should be used with a pipe")
	}

	parser := bench.NewParser(bufio.NewReader(os.Stdin))
	suites, err := parser.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("detected %d benchmark suites\n", len(suites))

	// set up results
	results := []Run{{
		Version: *version,
		Date:    *date,
		Tags:    *tags,
		Suites:  suites,
	}}
	if *appendOut {
		if *jsonOut == "" {
			panic("file output needs to be set (try '--json')")
		}
		b, err := ioutil.ReadFile(*jsonOut)
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		} else if !os.IsNotExist(err) {
			var runs []Run
			if err := json.Unmarshal(b, &runs); err != nil {
				panic(err)
			}
			results = append(results, runs...)
		} else {
			fmt.Printf("could not find specified output file '%s' - creating a new file\n", *jsonOut)
		}
	}

	output(results)
}

// Run denotes one run of gobenchdata, useful for grouping benchmark records
type Run struct {
	Version string `json:",omitempty"`
	Date    string
	Tags    []string `json:",omitempty"`
	Suites  []bench.Suite
}