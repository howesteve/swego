// +build IGNORE

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	log.SetPrefix("")
	flag.Parse()

	var header string

	const defaultHeader = "../swisseph/sweph.h"
	if _, err := os.Stat(defaultHeader); err == nil {
		header = defaultHeader
	}

	if header == "" {
		if flag.NArg() == 0 {
			log.Fatal("Provide path to sweph.h.")
		}

		header = flag.Arg(0)
	}

	re, err := regexp.Compile(`#define SE_VERSION\s+"(.*)"`)
	if err != nil {
		log.Fatalln("Error compiling regexp:", err)
	}

	data, err := ioutil.ReadFile(header)
	if err != nil {
		log.Fatalln("Error reading header file:", err)
	}

	match := re.FindSubmatch(data)
	version := string(match[1])

	parts := strings.Split(version, ".")
	trimmed := make([]string, len(parts))
	for i, p := range parts {
		trimmed[i] = strings.TrimPrefix(p, "0")
	}

	if len(trimmed) == 2 {
		parts = append(parts, "00")
		trimmed = append(trimmed, "0")
	}

	f, err := os.Create("sweversion.h")
	if err != nil {
		log.Fatalln("Error creating sweversion.h:", err)
	}

	fmt.Fprintln(f, "// DO NOT EDIT - generated by running:")
	fmt.Fprint(f, "// \tgo run genversion.go")

	if header == defaultHeader {
		fmt.Fprintln(f)
	} else {
		fmt.Fprintf(f, " %s\n", header)
	}

	fmt.Fprintln(f, "//")
	fmt.Fprintln(f, "// Re-run this command when the Swiss Ephemeris is updated.")
	fmt.Fprintln(f, "//")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "#define SWEX_VERSION %s\n", strings.Join(parts, ""))
	fmt.Fprintf(f, "#define SWEX_VERSION_MAJOR %s\n", trimmed[0])
	fmt.Fprintf(f, "#define SWEX_VERSION_MINOR %s\n", trimmed[1])
	fmt.Fprintf(f, "#define SWEX_VERSION_PATCH %s\n", trimmed[2])

	f.Close()
}
