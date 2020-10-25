package cobra

import "testing"

func TestName(t *testing.T) {
	rootCmd := &Command{Use: "root   "}
	childCmd := &Command{Use: "child arg   "}
	rootCmd.AddCommand(childCmd)
	cases := []struct {
		cmd  *Command
		want string
	}{
		{rootCmd, "root"},
		{childCmd, "child"},
	}
	for _, c := range cases {
		name := c.cmd.Name()
		if name != c.want {
			t.Errorf("cmd.Name() get: %v, want: %v", name, c.want)
		}
	}
}

func TestRunnable(t *testing.T) {
	withRunCmd := &Command{Run: func(cmd *Command, args []string) {}}
	withRunECmd := &Command{RunE: func(cmd *Command, args []string) error { return nil }}
	noRunCmd := &Command{}
	cases := []struct {
		cmd  *Command
		want bool
	}{
		{withRunCmd, true},
		{withRunECmd, true},
		{noRunCmd, false},
	}
	for _, c := range cases {
		runnable := c.cmd.runnable()
		if runnable != c.want {
			t.Errorf("cmd.runnable() get: %v, want: %v", runnable, c.want)
		}
	}
}

func TestRoot(t *testing.T) {
	rootCmd := &Command{Use: "root"}
	childCmd := &Command{Use: "child"}
	leafCmd := &Command{Use: "leaf"}
	rootCmd.AddCommand(childCmd)
	childCmd.AddCommand(leafCmd)
	cases := []struct {
		cmd  *Command
		want *Command
	}{
		{rootCmd, rootCmd},
		{childCmd, rootCmd},
		{leafCmd, rootCmd},
	}
	for _, c := range cases {
		root := c.cmd.Root()
		if root != c.want {
			t.Errorf("cmd.Root() get: %p, want: %p", root, c.want)
		}
	}
}

func TestFindNext(t *testing.T) {
	childCmd := &Command{Use: "child", commands: []*Command{}}
	rootCmd := &Command{Use: "root", commands: []*Command{childCmd}}
	cases := []struct {
		cmd  *Command
		next string
		want *Command
	}{
		{rootCmd, childCmd.Name(), childCmd},
		{rootCmd, "nonExistCmd", nil},
	}
	for _, c := range cases {
		nextCmd := c.cmd.findNext(c.next)
		if nextCmd != c.want {
			t.Errorf("nextCmd want: %v, got: %v", nextCmd, c.want)
		}
	}
}
