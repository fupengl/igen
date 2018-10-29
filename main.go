package main

import (
	"flag"
	"fmt"

	"igen/cmd"
	"igen/cmd/command"
	"log"
)

func main() {
	flag.Usage = cmd.Usage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		cmd.Usage()
		return
	}

	if args[0] == "help" {
		cmd.Help(args[1:])
		return
	}

	for _, c := range command.CMDs {
		if c.Name() == args[0] {
			c.Flag.Usage = func() { c.Usage() }
			err := c.Flag.Parse(args[1:])
			if err != nil {
				log.Fatal(err)
				return
			}

			if c.Run != nil {
				c.Run(c.Arg)
			}
			return
		}
	}

	fmt.Printf("\nunknow command %s\n\n", args[0])
}
