package cmd

import (
	"fmt"

	"igen/cmd/command"
)

func Usage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tigen command [options]")
	fmt.Println("")
	fmt.Println("Commands:")

	for _, c := range command.CMDs {
		fmt.Printf("\t%-16s%s\n", c.Name(), c.Short)
	}
	//fmt.Printf("\t%-16s%s\n", "all", "init both model and controller")
	fmt.Println("")
	fmt.Println("Use igen help [command] for more information about a command")
	fmt.Println("")
}

func Help(args []string) {
	if len(args) != 1 {
		Usage()
		return
	}
	for _, c := range command.CMDs {
		if args[0] == c.Name() {
			c.Usage()
		}
	}
}
