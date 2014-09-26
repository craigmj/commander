package commander

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var ErrUnrecognizedCommand = errors.New("No command executed")

// Command wraps together the short command name, the description
// for a command, the commands Flags and the function that will handle
// the command.
type Command struct {
	Command string
	Description string
	FlagSet *flag.FlagSet
	F func(args []string) error
}

// NewCommand creates a new comandeer Command struct with the given parameters.
func NewCommand(cmd, description string, flagset *flag.FlagSet, f func(args []string) error) *Command {
	return &Command{ cmd, description, flagset, f }
}

// CommandFunction returns a command
type CommandFunction func() *Command

// Execute takes an args array, and executes the appropriate command from the 
// array of commandFunctions. If nil is passed as the args array, os.Args is used
// by default.
func Execute(args []string, commandFns ...CommandFunction) error {
	if nil==args {
		args = os.Args[1:]
	}
	commands := make(map[string]*Command, len(commandFns))
	for _, c := range commandFns {
		cmd := c()
		commands[strings.ToLower(cmd.Command)] = cmd
	}

	if 0==len(args) || strings.ToLower(args[0])=="help" {
		if 1<len(args) {
			for _, c := range args[1:] {
				cmd, ok := commands[strings.ToLower(c)]
				if !ok {
					fmt.Println("Unrecognized sub-command: ", cmd)
					continue
				}
				if nil!=cmd.FlagSet {
					cmd.FlagSet.PrintDefaults()
				} else {
					fmt.Printf("%s takes no arguments: %s", cmd.Command, cmd.Description)
				}
			}
			return nil
		}
		fmt.Println(`Commands are:`)
		for _, c := range commands {
			fmt.Printf("%s\t\t%s\n", c.Command, c.Description)			
		}
		return nil
	}

	c, ok := commands[strings.ToLower(args[0])]
	if !ok {
		return ErrUnrecognizedCommand
	}
	args = args[1:]
	if nil!=c.FlagSet {
		c.FlagSet.Parse(args)
	}
	return c.F(args)
}
