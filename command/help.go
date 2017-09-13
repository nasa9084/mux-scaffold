package command

import "fmt"

const (
	helpDescription = `mux-scaffold help command prints help about the given command.`
	helpHelp        = `Usage:
        mux-scaffold help <command name>

Description:
        %s

Example:
        mux-scaffold help init
`
)

// Help represents help command
type Help struct {
}

// Exec to implement Command interface
func (c *Help) Exec(args []string) int {
	if len(args) == 0 {
		c.Help()
		return 1
	}

	target, ok := List[args[0]]

	if !ok {
		c.Help()
		return 1
	}

	target.Help()

	return 0
}

// Description to implement Command interface
func (c *Help) Description() string {
	return helpDescription
}

// Help to implement Command interface
func (c *Help) Help() {
	fmt.Printf(helpHelp, helpDescription)
}
