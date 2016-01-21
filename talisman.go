package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/thoughtworks/talisman/git_repo"
)

var (
	fdebug  bool
	githook string
)

const (
	PrePush   = "pre-push"
	PreCommit = "pre-commit"
)

func init() {
	log.SetOutput(os.Stderr)
}

//Logger is the default log device, set to emit at the Error level by default
func main() {
	flag.BoolVar(&fdebug, "debug", false, "enable debug mode (warning: very verbose)")
	flag.BoolVar(&fdebug, "d", false, "short form of debug (warning: very verbose)")
	flag.StringVar(&githook, "githook", PrePush, "either pre-push or pre-commit")
	os.Exit(run(os.Stdin))
}

func run(stdin io.Reader) (returnCode int) {
	flag.Parse()
	if fdebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	if githook == "" {
		githook = PrePush
	}

	log.Info("Running %s hook", githook)

	var additions []git_repo.Addition
	if githook == PreCommit {
		preCommitHook := NewPreCommitHook()
		additions = preCommitHook.GetRepoAdditions()
	} else {
		prePushHook := NewPrePushHook(readRefAndSha(stdin))
		additions = prePushHook.GetRepoAdditions()
	}

	return NewRunner(additions).RunWithoutErrors()
}

func readRefAndSha(file io.Reader) (string, string, string, string) {
	text, _ := bufio.NewReader(file).ReadString('\n')
	refsAndShas := strings.Split(strings.Trim(string(text), "\n"), " ")
	if len(refsAndShas) < 4 {
		return EmptySha, EmptySha, "", ""
	}
	return refsAndShas[0], refsAndShas[1], refsAndShas[2], refsAndShas[3]
}
