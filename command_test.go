package cobra

import (
	"fmt"
	"reflect"
	"testing"

	flag "github.com/spf13/pflag"
)

func TestExecute(t *testing.T) {
	args := []string{}
	execCases := []struct {
		cmd  *Command
		want error
	}{
		{&Command{}, flag.ErrHelp},
		{nil, fmt.Errorf("Called Execute() on a nil Command")},
		{&Command{RunE: func(*Command, []string) error {
			return fmt.Errorf("error")
		}}, fmt.Errorf("error")},
		{&Command{Run: func(*Command, []string) {}}, nil},
	}
	for _, c := range execCases {
		err := c.cmd.execute(args)
		if err == nil && c.want == nil {
			continue
		}
		if err == nil || c.want == nil || err.Error() != c.want.Error() {
			t.Errorf("cmd.execute() return error want: %v, got: %v", c.want, err)
		}
	}

	ExecCases := []struct {
		cmd  *Command
		want error
	}{
		{&Command{}, nil},
		{&Command{RunE: func(*Command, []string) error {
			return fmt.Errorf("error")
		}}, fmt.Errorf("error")},
		{&Command{Run: func(*Command, []string) {}}, nil},
	}
	for _, c := range ExecCases {
		err := c.cmd.Execute()
		if err == nil && c.want == nil {
			continue
		}
		if err == nil || c.want == nil || err.Error() != c.want.Error() {
			t.Errorf("cmd.Execute() return error want: %v, got: %v", c.want, err)
		}
	}
}

func TestFind(t *testing.T) {
	var leafCmd = &Command{Use: "leaf", commands: []*Command{}}
	var childCmd = &Command{Use: "child", commands: []*Command{leafCmd}}
	var rootCmd = &Command{Use: "root", commands: []*Command{childCmd}}

	cases := []struct {
		cmd       *Command
		args      []string
		targetCmd *Command
		flags     []string
		err       error
	}{
		{rootCmd, []string{childCmd.Name(), leafCmd.Name()}, leafCmd, []string{}, nil},
		{childCmd, []string{leafCmd.Name(), "arg1", "arg2"}, leafCmd, []string{"arg1", "arg2"}, nil},
		{leafCmd, []string{}, leafCmd, []string{}, nil},
	}

	for _, c := range cases {
		targetCmd, flags, err := c.cmd.Find(c.args)
		if targetCmd != c.targetCmd {
			t.Errorf("want target command: %v, got target command: %v", c.targetCmd.Name(), targetCmd.Name())
		}
		if !reflect.DeepEqual(flags, c.flags) {
			t.Errorf("want target flags: %v, got target flags: %v", c.flags, flags)
		}
		if err != nil {
			t.Errorf("error want: nil, got: %v", err)
		}
	}

}

func TestAddCommand(t *testing.T) {
	rootCmd := &Command{Use: "root"}
	childCmd := &Command{Use: "child"}
	rootCmd.AddCommand(childCmd)
	isFound := false
	for _, cmd := range rootCmd.commands {
		if cmd == childCmd {
			isFound = true
			break
		}
	}
	if !isFound {
		t.Error("childCmd is not in commands of the rootCmd")
	}
	if childCmd.parent != rootCmd {
		t.Error("childCmd's parent is not rootCmd")
	}
}

func TestRemoveCommand(t *testing.T) {
	rootCmd := &Command{Use: "root"}
	childCmd := &Command{Use: "child"}
	rootCmd.AddCommand(childCmd)
	rootCmd.RemoveCommand(childCmd)
	for _, cmd := range rootCmd.commands {
		if cmd == childCmd {
			t.Error("childCmd is in commands of the rootCmd")
		}
	}
	if childCmd.parent == rootCmd {
		t.Error("childCmd's parent is rootCmd, expected nil")
	}
}
