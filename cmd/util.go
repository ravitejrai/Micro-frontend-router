package cmd

import (
	"fmt"
	exec "golang.org/x/sys/execabs"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/c-bata/go-prompt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func RunInTerminalWithColor(cmdName string, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	return RunInTerminalWithColorInDir(cmdName, dir, args)
}

func RunInTerminalWithColorInDir(cmdName string, dir string, args []string) error {
	log.Debug().Msg(cmdName + " " + strings.Join(args, " "))

	_, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runtime.GOOS != "windows" {
		cmd.ExtraFiles = []*os.File{w}
	}

	err = cmd.Run()
	log.Debug().Err(err).Send()
	return err
}

func AskConfirm(q string) bool {
	ans := false
	survey.AskOne(&survey.Confirm{
		Message: q,
	}, &ans)
	return ans
}

func CobraCommandToSuggestions(cmds []*cobra.Command) []prompt.Suggest {
	var suggestions []prompt.Suggest
	for _, branch := range cmds {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        branch.Use,
			Description: branch.Short,
		})
	}
	return suggestions
}

func CobraCommandToName(cmds []*cobra.Command) []string {
	var ss []string
	for _, cmd := range cmds {
		ss = append(ss, cmd.Use)
	}
	return ss
}

func CobraCommandToDesc(cmds []*cobra.Command) []string {
	var ss []string
	for _, cmd := range cmds {
		ss = append(ss, cmd.Short)
	}
	return ss
}

type Branch struct {
	Author       string
	FullName     string
	RelativeDate string
	AbsoluteDate string
}

type FileChange struct {
	Name   string
	Status string
}

type PromptTheme struct {
	PrefixTextColor             prompt.Color
	SelectedSuggestionBGColor   prompt.Color
	SuggestionBGColor           prompt.Color
	SuggestionTextColor         prompt.Color
	SelectedSuggestionTextColor prompt.Color
	DescriptionBGColor          prompt.Color
	DescriptionTextColor        prompt.Color
}

var DefaultTheme = PromptTheme{
	PrefixTextColor:             prompt.Yellow,
	SelectedSuggestionBGColor:   prompt.Yellow,
	SuggestionBGColor:           prompt.Yellow,
	SuggestionTextColor:         prompt.DarkGray,
	SelectedSuggestionTextColor: prompt.Blue,
	DescriptionBGColor:          prompt.Black,
	DescriptionTextColor:        prompt.White,
}

var InvertedTheme = PromptTheme{
	PrefixTextColor:             prompt.Blue,
	SelectedSuggestionBGColor:   prompt.LightGray,
	SelectedSuggestionTextColor: prompt.White,
	SuggestionBGColor:           prompt.Blue,
	SuggestionTextColor:         prompt.White,
	DescriptionBGColor:          prompt.LightGray,
	DescriptionTextColor:        prompt.Black,
}

var MonochromeTheme = PromptTheme{}

func SuggestionPrompt(prefix string, completer func(d prompt.Document) []prompt.Suggest) string {
	theme := DefaultTheme
	themeName := os.Getenv("BIT_THEME")
	if strings.EqualFold(themeName, "inverted") {
		theme = InvertedTheme
	}
	if strings.EqualFold(themeName, "monochrome") {
		theme = MonochromeTheme
	}
	result := prompt.Input(prefix, completer,
		prompt.OptionTitle(""),
		prompt.OptionHistory([]string{""}),
		prompt.OptionPrefixTextColor(theme.PrefixTextColor), // fine
		prompt.OptionSelectedSuggestionBGColor(theme.SelectedSuggestionBGColor),
		prompt.OptionSuggestionBGColor(theme.SuggestionBGColor),
		prompt.OptionSuggestionTextColor(theme.SuggestionTextColor),
		prompt.OptionSelectedSuggestionTextColor(theme.SelectedSuggestionTextColor),
		prompt.OptionDescriptionBGColor(theme.DescriptionBGColor),
		prompt.OptionDescriptionTextColor(theme.DescriptionTextColor),
		// prompt.OptionPreviewSuggestionBGColor(prompt.Yellow),
		// prompt.OptionPreviewSuggestionTextColor(prompt.Yellow),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionOnDown(),
		prompt.OptionSwitchKeyBindMode(prompt.EmacsKeyBind),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn:  exit,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x62},
			Fn:        prompt.GoLeftWord,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x66},
			Fn:        prompt.GoRightWord,
		}),
	)
	branchName := strings.TrimSpace(result)
	if strings.HasPrefix(result, "origin/") {
		branchName = branchName[7:]
	}
	return branchName
}

type Exit int

func exit(_ *prompt.Buffer) {
	panic(Exit(0))
}

func HandleExit() {
	switch v := recover().(type) {
	case nil:
		return
	case Exit:
		os.Exit(int(v))
	default:
		fmt.Println(v)
		fmt.Println(string(debug.Stack()))
		fmt.Println("OS:", runtime.GOOS, runtime.GOARCH)
		fmt.Println("bit version " + GetVersion())

	}
}

func AllBitSubCommands(rootCmd *cobra.Command) ([]*cobra.Command, map[string]*cobra.Command) {
	bitCmds := rootCmd.Commands()
	bitCmdMap := map[string]*cobra.Command{}
	for _, bitCmd := range bitCmds {
		bitCmdMap[bitCmd.Name()] = bitCmd
	}
	return bitCmds, bitCmdMap
}

func AllBitAndGitSubCommands(rootCmd *cobra.Command) (cc []*cobra.Command) {
	commonCommands := CommonCommandsList()
	return concatCopyPreAllocate([][]*cobra.Command{commonCommands})
}

func concatCopyPreAllocate(slices [][]*cobra.Command) []*cobra.Command {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]*cobra.Command, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func CommonCommandsList() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "status",
			Short: "Show the working tree status",
		},
		{
			Use:   "pull --rebase origin master",
			Short: "Rebase on origin master branch",
		},
		{
			Use:   "push --force-with-lease",
			Short: "force push with a safety net",
		},
		{
			Use:   "stash pop",
			Short: "Use most recently stashed changes",
		},
		{
			Use:   "commit -am \"",
			Short: "Commit all tracked files",
		},
		{
			Use:   "commit -a --amend --no-edit",
			Short: "Amend most recent commit with new changes",
		},
		{
			Use:   "commit --amend --no-edit",
			Short: "Amend most recent commit with added changes",
		},
		{
			Use:   "merge --squash",
			Short: "Squash and merge changes from a specified branch",
		},
		{
			Use:   "release bump",
			Short: "Commit unstaged changes, bump minor tag, push",
		},
		{
			Use:   "log --oneline",
			Short: "Display one commit per line",
		},
		{
			Use:   "diff --cached",
			Short: "Shows all staged changes",
		},
	}
}

func isBranchCompletionCommand(command string) bool {
	return command == "checkout" || command == "switch" || command == "co" || command == "pr" || command == "merge" || command == "rebase"
}

func isBranchChangeCommand(command string) bool {
	return command == "checkout" || command == "switch" || command == "co" || command == "pr"
}

func Find(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
