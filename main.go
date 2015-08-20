package main

import (
	"flag"
	"fmt"
	"github.com/everfore/exc"
	"os"
)

var (
	quit   = false
	commit = ""
)

func init() {
	flag.BoolVar(&quit, "q", false, "-q: quit add all")
	flag.StringVar(&commit, "m", "", "-m: commit content")
}

func main() {
	flag.Parse()
	fmt.Println(quit)
	first := "git add -A"
	if quit {
		first = "git status"
	}
	git := exc.NewCMD(first).Wd()
	_, err := git.Debug().Do()
	if checkerr(err) {
		os.Exit(-1)
	}
	fmt.Println(commit)
	if len(commit) > 0 {
		commit = fmt.Sprintf(`git commit -m "%s"`, commit)
	} else {
		commit = `git commit -m "auto"`
	}
	_, err = git.Reset(commit).Do()
	if checkerr(err) {
		os.Exit(-1)
	}
	git.Reset("git status").Execute()
	return
	git.Reset("git push origin master").Execute()
}

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
