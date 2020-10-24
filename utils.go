package cobra

import "strings"

// argsMinusFirstX removes only the first x from args.  Otherwise, commands that look like
// openshift admin policy add-role-to-user admin my-user, lose the admin argument (arg[4]).
func argsMinusFirstX(args []string, x string) []string {
	for i, y := range args {
		if x == y {
			ret := []string{}
			ret = append(ret, args[:i]...)
			ret = append(ret, args[i+1:]...)
			return ret
		}
	}
	return args
}

// Name returns the command's name: the first word in the use line.
//
// Name 函数通过 Command.Use 字段返回 command 的名字（use 的首单词）
func (cmd *Command) Name() string {
	name := cmd.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Root finds root command.
func (cmd *Command) Root() *Command {
	if cmd.parent != nil {
		return cmd.parent.Root()
	}
	return cmd
}

// Runnable determines if the command is itself runnable.
func (cmd *Command) Runnable() bool {
	return cmd.Run != nil || cmd.RunE != nil
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
