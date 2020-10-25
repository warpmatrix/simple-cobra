// Package cobra 实现了支持简单带子命令的命令行程序开发的功能。
// 支持导入子命令目录，为开发的应用程序增加需要的自定义子命令。
// 所有的子命令文件夹中，一个 go 文件对应应用程序中的一个子命令。
// 可以通过增加或删除 go 文件，对直接对应用程序中的子命令进行改动。
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
	// Use 字段描述指令的使用格式，支持的格式如下：
	//   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.（可选参数）
	//   ... indicates that you can specify multiple values for the previous argument.（多值参数）
	//   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
	//       argument to the right of the separator. You cannot use both arguments in a single use of the command.（互斥参数）
	//   { } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
	//       optional, they are enclosed in brackets ([ ]).）（互斥分隔府）
	// Example: add [-F file | -D dir]... [-f format] profile
	Use string
	// Short is the short description shown in the 'help' output.
	// Short 字段用于 help 中的简短描述
	Short string
	// Long is the long message shown in the 'help <this-command>' output.
	// Long 字段用于 help 中的长信息描述
	Long string

	// Run: Typically the actual work function. Most commands will only implement this.
	// Run 函数由用户提供，描述执行指令的行为（无返回值）
	Run func(cmd *Command, args []string)
	// RunE: Run but returns an error.
	// Run 函数由用户提供，描述执行指令的行为，返回可能产生的错误
	RunE func(cmd *Command, args []string) error

	// helpFunc is help func defined by user.
	// helpFunc 字段描述 help 指令执行的函数，系统会提供默认函数，可以被用户重定义
	helpFunc func(*Command, []string)
	// helpCommand is command with usage 'help'. If it's not defined by user,
	// cobra uses default help command.
	// helpCommand 字段系统为指令添加默认的 help 子指令，可以被用户重定义
	helpCommand *Command

	// args is actual args parsed from flags.
	// args 字段可以通过配置文件读取指令默认的参数（尚未实现）
	args []string

	// parent is a parent command for this command.
	// parent 字段记录该指令的父指令
	parent *Command
	// commands is the list of commands supported by this program.
	// commands 字段记录该指令的子指令集合
	commands []*Command
}

// Execute executes the command.
//
// Execute 函数执行用户输入的指令，会解析出具体执行的指令并执行
func (cmd *Command) Execute() error {
	cmd.initDefaultHelpCmd()
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
	if !cmd.runnable() {
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

// AddCommand adds one or more commands to this parent command.
//
// AddCommand 函数为 cmd 指令增加子指令
func (cmd *Command) AddCommand(subCmds ...*Command) {
	for i, subCmd := range subCmds {
		if subCmds[i] == cmd {
			panic("Command can't be a child of itself")
		}
		subCmds[i].parent = cmd
		cmd.commands = append(cmd.commands, subCmd)
	}
}

// RemoveCommand removes one or more commands from a parent command.
//
// RemoveCommand 函数为 cmd 指令移除子命令
func (cmd *Command) RemoveCommand(rmCmds ...*Command) {
	commands := []*Command{}
main:
	for _, command := range cmd.commands {
		for _, rmCmd := range rmCmds {
			if command == rmCmd {
				command.parent = nil
				continue main
			}
		}
		commands = append(commands, command)
	}
	cmd.commands = commands
}
