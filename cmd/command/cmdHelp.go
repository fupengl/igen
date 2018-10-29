package command

var cmdHelp = &CMD{
	UsageLine: "help [command]",
	Short:     "",
}

func init() {
	CMDs = append(CMDs, cmdHelp)
}
