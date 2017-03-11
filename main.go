// Package rclone-json streams an rclone sync activity as a json object stream.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mh-cbon/rclone-json/rclone"
)

// Version of the program
var Version = "0.0.0"

func main() {

	version := false
	help := false
	cmd := rclone.New("", "")

	flag.BoolVar(&version, "version", false, "Show version")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&cmd.BinPath, "rclone", "rclone", "")
	flag.StringVar(&cmd.BwLimit, "bwlimit", "", "")
	flag.StringVar(&cmd.Stats, "stats", "", "")
	flag.StringVar(&cmd.Checkers, "checkers", "", "")
	flag.StringVar(&cmd.TransferLimit, "transfers", "", "")
	flag.Parse()

	if version {
		ver()
		os.Exit(0)
		return
	}

	if help {
		usage("")
		os.Exit(0)
		return
	}

	if flag.NArg() != 3 {
		usage(fmt.Sprint(`
Wrong usage: Missing source or dest arguments.
It should be:
rclone-json stats [options] src/ dst/
`))
		os.Exit(1)
	}

	// rclone \
	// sync -vv --stats 1s --bwlimit 500k --checkers 2 \
	// --transfers 20 \
	// ~/test/rclone-v1.35-linux-amd64/source/ \
	// ~/test/rclone-v1.35-linux-amd64/dest/

	// _ := flag.Arg(0) // the cmd to run like stats or ls, not needed so far as only stats is implemented
	cmd.Src = flag.Arg(1)
	cmd.Dst = flag.Arg(2)

	cmd.Stdout = os.Stderr

	mustNotErr(cmd.Start())
	mustNotErr(cmd.ConvertTo(json.NewEncoder(os.Stdout)))
	mustNotErr(cmd.Wait())
}

func ver() {
	fmt.Fprintf(os.Stderr, "rclone-json - %v\n", Version)
}

func usage(err string) {
	ver()
	flag.Usage()
	if err != "" {
		fmt.Fprintln(os.Stderr, err)
	}
}

func mustNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
