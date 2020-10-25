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
//
// Root 函数返回指令的根指令
func (cmd *Command) Root() *Command {
	if cmd.parent != nil {
		return cmd.parent.Root()
	}
	return cmd
}

// runnable determines if the command is itself runnable.
func (cmd *Command) runnable() bool {
	return cmd.Run != nil || cmd.RunE != nil
}

func (cmd *Command) findNext(next string) *Command {
	for _, nextCmd := range cmd.commands {
		if nextCmd.Name() == next {
			return nextCmd
		}
	}
	return nil
}
