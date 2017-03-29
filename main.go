// Package rclone-json streams an rclone sync activity as a json object stream.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mh-cbon/rclone-json/rclone"
)

// VERSION of the program
var VERSION = "0.0.0"

type cliargs struct {
	version       bool
	help          bool
	binPath       string
	bwlimit       string
	stats         string
	checker       string
	transferlimit string
}

func main() {

	args := cliargs{}

	flag.BoolVar(&args.version, "version", false, "Show version")
	flag.BoolVar(&args.help, "help", false, "Show help")
	flag.StringVar(&args.binPath, "rclone", "rclone", "")

	flag.StringVar(&args.bwlimit, "bwlimit", "", "")
	flag.StringVar(&args.stats, "stats", "", "")
	flag.StringVar(&args.checker, "checkers", "", "")
	flag.StringVar(&args.transferlimit, "transfers", "", "")
	flag.Parse()

	if args.version {
		ver()
		os.Exit(0)
		return
	}

	if args.help {
		usage("")
		os.Exit(0)
		return
	}

	if flag.NArg() < 1 {
		usage(fmt.Sprint(`
Wrong usage: Missing subcommand.
Excpected one of check|size|sync
`))
		os.Exit(1)
	}
	subcmd := flag.Arg(0) // the cmd to run like sync / size / check

	if subcmd == "sync" {

		if flag.NArg() != 3 {
			usage(fmt.Sprint(`
Wrong usage: Missing source or dest arguments.
It should be:
rclone-json sync [options] src/ dst/
			`))
			os.Exit(1)
		}

		/*
			rclone \
			sync -vv --stats 1s --bwlimit 500k --checkers 2 \
			--transfers 20 \
			test/source/ \
			test/dest/
		*/

		cmd := rclone.New(flag.Arg(1), flag.Arg(2))

		cmd.BinPath = args.binPath
		cmd.Verbose = true
		cmd.BwLimit = args.bwlimit
		cmd.Stats = args.stats
		cmd.Checkers = args.checker
		cmd.TransferLimit = args.transferlimit

		cmd.Stdout = os.Stderr

		mustNotErr(cmd.Start())
		mustNotErr(cmd.ConvertTo(json.NewEncoder(os.Stdout)))
		mustNotErr(cmd.Wait())

	} else if subcmd == "size" {

		if flag.NArg() != 2 {
			usage(fmt.Sprint(`
Wrong usage: Missing source arguments.
It should be:
rclone-json size [options] src/
			`))
			os.Exit(1)
		}

		cmd := rclone.NewSize(flag.Arg(1))

		cmd.BinPath = args.binPath
		mustNotErr(cmd.Exec())

		b, err := json.Marshal(cmd.Res)
		mustNotErr(err)
		os.Stdout.Write(b)

	} else if subcmd == "check" {

		if flag.NArg() != 3 {
			usage(fmt.Sprint(`
Wrong usage: Missing source or dest arguments.
It should be:
rclone-json check [options] src/ dst/
			`))
			os.Exit(1)
		}

		cmd := rclone.NewSize(flag.Arg(1))

		cmd.BinPath = args.binPath
		mustNotErr(cmd.Exec())

		b, err := json.Marshal(cmd.Res)
		mustNotErr(err)
		os.Stdout.Write(b)

	} else {
		usage(fmt.Sprintf(`
		Wrong usage: Uknown sub command %q
		Excpected one of check|size|sync
		`, subcmd))
		os.Exit(1)
	}
}

func ver() {
	fmt.Fprintf(os.Stderr, "rclone-json - %v\n", VERSION)
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
