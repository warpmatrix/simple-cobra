package cobra

import (
	"reflect"
	"testing"
)

func TestHelpFunc(t *testing.T) {
	rootCmd := &Command{helpFunc: nil}
	if rootCmd.HelpFunc() == nil {
		t.Error("c.HelpFunc() want: not nil, get: nil")
	}
	befPtr := reflect.ValueOf(rootCmd.HelpFunc()).Pointer()
	helpFunc := func(*Command, []string) {}
	rootCmd.helpFunc = helpFunc
	helpFuncPtr := reflect.ValueOf(helpFunc).Pointer()
	rootHelpFuncPtr := reflect.ValueOf(rootCmd.HelpFunc()).Pointer()
	if befPtr == helpFuncPtr {
		t.Error("default help function should be different with self-defined help function")
	}
	if helpFuncPtr != rootHelpFuncPtr {
		t.Errorf("c.HelpFunc() want: %v, get: %v", helpFuncPtr, rootHelpFuncPtr)
	}
	childCmd := &Command{helpFunc: nil}
	rootCmd.AddCommand(childCmd)
	childHelpFuncPtr := reflect.ValueOf(childCmd.HelpFunc()).Pointer()
	if rootHelpFuncPtr != childHelpFuncPtr {
		t.Errorf("c.HelpFunc() want: %v, get: %v", rootHelpFuncPtr, childHelpFuncPtr)
	}
}

func TestSetHelpFunc(t *testing.T) {
	cmd := &Command{helpFunc: nil}
	f := func(*Command, []string) {}
	cmd.SetHelpFunc(f)
	if cmd.helpFunc == nil {
		t.Errorf("c.helpFunc want: %p, get: nil", f)
	}
	cmd.SetHelpFunc(nil)
	if cmd.helpFunc != nil {
		t.Errorf("c.helpFunc want: nil, get: %p", cmd.helpFunc)
	}
}

func TestInitDefaultHelpCmd(t *testing.T) {
	rootCmd := &Command{Use: "root"}
	childCmd := &Command{Use: "child"}
	rootCmd.AddCommand(childCmd)
	rootCmd.initDefaultHelpCmd()
	if rootCmd.helpCommand == nil {
		t.Error("cmd.initDefaultHelpCmd failed, helpCommand got nil")
	}
	isFound := false
	for _, cmd := range rootCmd.commands {
		if cmd == rootCmd.helpCommand {
			isFound = true
			break
		}
	}
	if !isFound {
		t.Error("helpCommand is not in commands of the rootCmd")
	}
	if rootCmd.helpCommand.parent != rootCmd {
		t.Error("helpCommand's parent is not rootCmd")
	}
}
