package main

import (
	"flag"
	"fmt"
	"github.com/everfore/exc"
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
	first := "git add -A"
	var git *exc.CMD
	if quit {
		exc.NewCMD("git push origin master").Debug().Execute()
		return
	}
	git = exc.NewCMD(first).Wd()
	_, err := git.Debug().Do()
	if checkerr(err) {
		return
	}
	if len(commit) > 0 {
		commit = fmt.Sprintf(`git commit -m "%s"`, commit)
	} else {
		commit = `git commit -m "auto"`
	}
	_, err = git.Reset(commit).Do()
	if checkerr(err) {
		return
	}
	git.Reset("git push origin master").Execute()
}

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
