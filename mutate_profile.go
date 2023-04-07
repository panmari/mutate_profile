package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/pprof/profile"
)

var profileFile = flag.String("profile", "", "The profile to mutate.")
var outFile = flag.String("output", "", "The mutated profile file to write. By default, appends '.mutated' to input profile.")

func parse() (*profile.Profile, error) {
	f, err := os.Open(*profileFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", *profileFile, err)
	}
	return profile.Parse(f)
}

func write(prof *profile.Profile) error {
	fName := *outFile
	if fName == "" {
		fName = *profileFile + ".mutated"
	}
	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	return prof.Write(f)
}

func main() {
	flag.Parse()

	prof, err := parse()
	if err != nil {
		log.Fatal(err)
	}

	// Mutations on the profile.
	for _, f := range prof.Function {
		f.Name = strings.ReplaceAll(f.Name, "[...]", "")
		f.Name = strings.ReplaceAll(f.Name, "*filter", "*Filter")
	}

	if err := write(prof); err != nil {
		log.Fatal(err)
	}
}
