// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// periph-info prints out information about the loaded periph drivers.
package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"periph.io/x/periph"
)

// driverAfter is an optional function.
// TODO(maruel): Remove in v3.
type driverAfter interface {
	After() []string
}

func printDrivers(drivers []periph.DriverFailure) {
	if len(drivers) == 0 {
		os.Stdout.WriteString("  <none>\n")
		return
	}
	max := 0
	for _, f := range drivers {
		if m := len(f.D.String()); m > max {
			max = m
		}
	}
	for _, f := range drivers {
		os.Stdout.WriteString("- ")
		os.Stdout.WriteString(f.D.String()) //"%-*s")
		os.Stdout.WriteString(": ")
		os.Stdout.WriteString(f.Err.Error())
		os.Stdout.WriteString("\n")
	}
}

func mainImpl() error {
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)
	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}

	state, err := hostInit()
	if err != nil {
		return err
	}

	os.Stdout.WriteString("Drivers loaded and their dependencies, if any:\n")
	if len(state.Loaded) == 0 {
		os.Stdout.WriteString("  <none>\n")
	} else {
		max := 0
		for _, d := range state.Loaded {
			if m := len(d.String()); m > max {
				max = m
			}
		}
		for _, d := range state.Loaded {
			p := d.Prerequisites()
			var a []string
			if da, ok := d.(driverAfter); ok {
				a = da.After()
			}
			if len(p) == 0 && len(a) == 0 {
				os.Stdout.WriteString("- ")
				os.Stdout.WriteString(d.String())
				os.Stdout.WriteString("\n")
				continue
			}
			os.Stdout.WriteString("- ")
			os.Stdout.WriteString(d.String()) //"%-*s", max, d)
			os.Stdout.WriteString(":")
			if len(p) != 0 {
				os.Stdout.WriteString(" ")
				// TODO: os.Stdout.WriteString(p)
			}
			if len(a) != 0 {
				os.Stdout.WriteString(" optional: ")
				// TODO: os.Stdout.WriteString(a)
			}
			os.Stdout.WriteString("\n")
		}
	}

	os.Stdout.WriteString("Drivers skipped and the reason why:\n")
	printDrivers(state.Skipped)
	os.Stdout.WriteString("Drivers failed to load and the error:\n")
	printDrivers(state.Failed)
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		os.Stderr.WriteString("periph-info: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString(".\n")
		os.Exit(1)
	}
}
