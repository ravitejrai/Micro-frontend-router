// Package main is complete tool for the go command line
package cmd

import (
	"fmt"
	"github.com/chriswalz/complete/v3"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Bitcomplete() {
	fmt.Println("Inside bitcomplete")
	compLine := os.Getenv("COMP_LINE")
	compPoint := os.Getenv("COMP_POINT")
	doInstall := os.Getenv("COMP_INSTALL") == "1"
	doUninstall := os.Getenv("COMP_UNINSTALL") == "1"

	bitcompletionNotNeeded := compLine == "" && compPoint == "" && !doInstall && !doUninstall
	if bitcompletionNotNeeded {
		return
	}

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	bitcomplete, _ := CreateSuggestionMap(BitCmd)

	complete.Complete("harmony", bitcomplete)
	//bitcomplete.Complete("bit")
}
