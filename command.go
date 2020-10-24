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
// 这个结构是 cobra.Command 结构的简化版本，提供了记录子命令、父命令、帮助描述等相关字段。
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
	// RunE: Run but returns an error.
	RunE func(cmd *Command, args []string) error

	// helpFunc is help func defined by user.
	helpFunc func(*Command, []string)
	// helpCommand is command with usage 'help'. If it's not defined by user,
	// cobra uses default help command.
	helpCommand *Command

	// args is actual args parsed from flags.
	args []string
	// flags is full secom/warpmatrix/simple-cobra"t of flags.
	// flags *flag.FlagSet

	// parent is a parent command for this command.
	parent *Command
	// commands is the list of commands supported by this program.
	commands []*Command
}

// Execute executes the command.
func (cmd *Command) Execute() error {
	cmd.InitDefaultHelpCmd()
	args := cmd.args
	// Workaround FAIL with "go test -v" or "cobra.test -test.v", see #155
	if cmd.args == nil && filepath.Base(os.Args[0]) != "cobra.test" {
		args = os.Args[1:]
	}
	targetCmd, flags, err := cmd.Find(args)
	if err != nil {
		return err
	}
	err = targetCmd.execute(flags)
	if err != nil {
		// Always show help if requested, even if SilenceErrors is in effect
		if err == flag.ErrHelp {
			targetCmd.HelpFunc()(targetCmd, args)
			return nil
		}
	}
	return err
}

func (cmd *Command) execute(a []string) (err error) {
	if cmd == nil {
		return fmt.Errorf("Called Execute() on a nil Command")
	}
	if !cmd.Runnable() {
		return flag.ErrHelp
	}
	if cmd.RunE != nil {
		err := cmd.RunE(cmd, a)
		return err
	}
	cmd.Run(cmd, a)
	return nil
}

// Find the target command given the args and command tree
// Meant to be run on the highest node. Only searches down.
//
// Find 函数为用户输入的 cmd 指令，找到并返回最终执行的子指令，并将剩余的参数作为 flags 返回
func (cmd *Command) Find(args []string) (*Command, []string, error) {
	var innerfind func(*Command, []string) (*Command, []string)
	innerfind = func(cmd *Command, innerArgs []string) (*Command, []string) {
		// args without flags
		argsWOflags := innerArgs
		if len(argsWOflags) == 0 {
			return cmd, innerArgs
		}
		nextSubCmd := argsWOflags[0]
		targetCmd := cmd.findNext(nextSubCmd)
		if targetCmd != nil {
			return innerfind(targetCmd, argsMinusFirstX(innerArgs, nextSubCmd))
		}
		return cmd, innerArgs
	}
	commandFound, flags := innerfind(cmd, args)
	return commandFound, flags, nil
}

func (cmd *Command) findNext(next string) *Command {
	for _, cmd := range cmd.commands {
		if cmd.Name() == next {
			return cmd
		}
	}
	return nil
}
