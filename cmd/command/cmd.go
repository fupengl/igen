package command

import (
	"flag"
	"fmt"
	"strings"

	"igen/cmd/util"
)

type CMD struct {
	UsageLine string

	Short   string // usage短说明
	Long    string // usage长说明
	Example string

	Run func(a *util.Arg)

	Flag flag.FlagSet
	Arg  *util.Arg
}

var CMDs = []*CMD{}

func (c *CMD) Name() string {
	a := strings.Split(c.UsageLine, " ")
	return a[0]
}

func (c *CMD) Usage() {
	fmt.Println("")
	fmt.Printf("%s\n\n", c.Short)
	fmt.Println("Usage:")
	fmt.Printf("\tigen %s\n", c.UsageLine)
	if c.Long != "" {
		fmt.Printf("\n%s\n", c.Long)
	}
	fmt.Println("")
}

func (c *CMD) Init() {
	c.Arg, c.Flag = initArg()
}
