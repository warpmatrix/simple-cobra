package cobra

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

// Command is just that, a command for your application.
// E.g.  'go run ...' - 'run' is the command. Cobra requires
// you to define the usage and description as part of your command
// definition to ensure usability.
//
// Command 结构定义了应用程序所使用的指令，提供了 usage 以及一些其他描述性的接口。
// 这个结构是 cobra.Command 结构的简化版本，提供了记录子命令、父命令等相关字段
type Command struct {
	// Use is the one-line usage message.
	// Recommended syntax is as follow:
	//   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
	//   ... indicates that you can specify multiple values for the previous argument.
	//   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
	//       argument to the right of the separator. You cannot use both arguments in a single use of the command.
	//   { } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
	//       optional, they are enclosed in brackets ([ ]).
	// Example: add [-F file | -D dir]... [-f format] profile
	Use string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	// Run: Typically the actual work function. Most commands will only implement this.
	Run func(cmd *Command, args []string)

	// args is actual args parsed from flags.
	args []string
}

// Execute executes the command.
func (c *Command) Execute() error {
	args := c.args
	// Workaround FAIL with "go test -v" or "cobra.test -test.v", see #155
	if c.args == nil && filepath.Base(os.Args[0]) != "cobra.test" {
		args = os.Args[1:]
	}
	// TODO: c.Find(args)
	// cmd, _, err := c.Find(args)
	// if err != nil {
	// 	// If found parse to a subcommand and then failed, talk about the subcommand
	// 	if cmd != nil {
	// 		c = cmd
	// 	}
	// 	return err
	// }
	// err = c.execute(args)
	err := c.execute(args)
	if err != nil {
		// Always show help if requested, even if SilenceErrors is in effect
		if err == flag.ErrHelp {
			// cmd.HelpFunc()(cmd, args)
			return nil
		}
	}
	return err
}

func (c *Command) execute(a []string) (err error) {
	if c == nil {
		return fmt.Errorf("Called Execute() on a nil Command")
	}
	c.Run(c, a)
	return nil
}

// // Find the target command given the args and command tree
// // Meant to be run on the highest node. Only searches down.
// func (c *Command) Find(args []string) (*Command, []string, error) {
// 	var innerfind func(*Command, []string) (*Command, []string)

// 	innerfind = func(c *Command, innerArgs []string) (*Command, []string) {
// 		argsWOflags := stripFlags(innerArgs, c)
// 		if len(argsWOflags) == 0 {
// 			return c, innerArgs
// 		}
// 		nextSubCmd := argsWOflags[0]

// 		cmd := c.findNext(nextSubCmd)
// 		if cmd != nil {
// 			return innerfind(cmd, argsMinusFirstX(innerArgs, nextSubCmd))
// 		}
// 		return c, innerArgs
// 	}

// 	commandFound, a := innerfind(c, args)
// 	if commandFound.Args == nil {
// 		return commandFound, a, legacyArgs(commandFound, stripFlags(a, commandFound))
// 	}
// 	return commandFound, a, nil
// }

// func stripFlags(args []string, c *Command) []string {
// 	if len(args) == 0 {
// 		return args
// 	}

// 	commands := []string{}
// 	flags := c.Flags()

// Loop:
// 	for len(args) > 0 {
// 		s := args[0]
// 		args = args[1:]
// 		switch {
// 		case s == "--":
// 			// "--" terminates the flags
// 			break Loop
// 		case strings.HasPrefix(s, "--") && !strings.Contains(s, "=") && !hasNoOptDefVal(s[2:], flags):
// 			// If '--flag arg' then
// 			// delete arg from args.
// 			fallthrough // (do the same as below)
// 		case strings.HasPrefix(s, "-") && !strings.Contains(s, "=") && len(s) == 2 && !shortHasNoOptDefVal(s[1:], flags):
// 			// If '-f arg' then
// 			// delete 'arg' from args or break the loop if len(args) <= 1.
// 			if len(args) <= 1 {
// 				break Loop
// 			} else {
// 				args = args[1:]
// 				continue
// 			}
// 		case s != "" && !strings.HasPrefix(s, "-"):
// 			commands = append(commands, s)
// 		}
// 	}

// 	return commands
// }

// func hasNoOptDefVal(name string, fs *flag.FlagSet) bool {
// 	flag := fs.Lookup(name)
// 	if flag == nil {
// 		return false
// 	}
// 	return flag.NoOptDefVal != ""
// }

// func shortHasNoOptDefVal(name string, fs *flag.FlagSet) bool {
// 	if len(name) == 0 {
// 		return false
// 	}

// 	flag := fs.ShorthandLookup(name[:1])
// 	if flag == nil {
// 		return false
// 	}
// 	return flag.NoOptDefVal != ""
// }
