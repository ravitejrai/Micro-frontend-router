package cmd

import (
	exec "golang.org/x/sys/execabs"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)



func IsGitRepo() bool {
	out, _ := execCommand("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	return strings.TrimSpace(string(out)) == "true"
}

func PrintGitVersion() {
	RunInTerminalWithColor("git", []string{"--version"})
}

func execCommand(name string, arg ...string) *exec.Cmd {
	log.Debug().Msg(name + " " + strings.Join(arg, " "))
	c := exec.Command(name, arg...)
	c.Env = os.Environ()

	if name == "git" {
		// exec commands are parsed by bit without getting printed.
		// parsing git assumes english
		c.Env = append(c.Env, "LANG=C")
	}
	return c
}
