package cobra

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

// func stripFlags(args []string, cmd *Command) []string {
// 	if len(args) == 0 {
// 		return args
// 	}

// 	commands := []string{}
// 	flags := cmd.Flags()

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

// // Flags returns the complete FlagSet that applies
// // to this command (local and persistent declared here and by all parents).
// func (cmd *Command) Flags() *flag.FlagSet {
// 	if cmd.flags == nil {
// 		cmd.flags = flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
// 	}
// 	return cmd.flags
// }
