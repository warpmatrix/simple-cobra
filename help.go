package cobra

import "fmt"

// HelpFunc returns either the function set by SetHelpFunc for this command
// or a parent, or it returns a function with default help behavior.
//
// HelpFunc 函数返回执行 help 指令的函数，可以被用户自定义或继承实现
func (cmd *Command) HelpFunc() func(*Command, []string) {
	if cmd.helpFunc != nil {
		return cmd.helpFunc
	}
	if cmd.parent != nil {
		return cmd.parent.HelpFunc()
	}
	return func(cmd *Command, a []string) {
		if cmd.Long != "" {
			fmt.Println(cmd.Long)
		}
		if cmd.Use != "" {
			fmt.Println("Usage:")
			fmt.Printf("  %s\n\n", cmd.Use)
		}
		if len(cmd.commands) != 0 {
			fmt.Println("Available Commands:")
			for _, subCmd := range cmd.commands {
				fmt.Printf("  %-10s %s\n", subCmd.Name(), subCmd.Short)
			}
			fmt.Println("")
			fmt.Printf("Use \"%s help [command] \" for more information about a command.\n", cmd.Name())
		}
	}
}

// initDefaultHelpCmd adds default help command to c.
// It is called automatically by executing the c or by calling help and usage.
// If c already has help command or c has no subcommands, it will do nothing.
func (cmd *Command) initDefaultHelpCmd() {
	if len(cmd.commands) == 0 {
		return
	}

	if cmd.helpCommand == nil {
		cmd.helpCommand = &Command{
			Use:   "help [command]",
			Short: "Help about any command",
			Long: `Help provides help for any command in the application.
Simply type ` + cmd.Name() + ` help [path to command] for full details.`,
			Run: func(c *Command, args []string) {
				targetCmd, _, e := c.Root().Find(args)
				if targetCmd == nil || e != nil {
					fmt.Printf("Unknown help topic %#q\n", args)
					c.HelpFunc()
				} else {
					targetCmd.HelpFunc()(targetCmd, []string{})
				}
			},
		}
	}
	cmd.RemoveCommand(cmd.helpCommand)
	cmd.AddCommand(cmd.helpCommand)
}

// SetHelpFunc sets help function. Can be defined by Application.
//
// SetHelpFunc 函数可以让用户自定义 helpFunc 字段
func (cmd *Command) SetHelpFunc(f func(*Command, []string)) {
	cmd.helpFunc = f
}
