package main

import (
	"flag"
	"fmt"
	"github.com/everfore/exc"
)

var (
	quit   = false
	commit = ""
	tag    = ""
)

func init() {
	flag.BoolVar(&quit, "q", false, "-q: quit add all")
	flag.StringVar(&tag, "t", "no_tag", "-t: tag")
	flag.StringVar(&commit, "m", "", "-m: commit content")
}

func main() {
	flag.Parse()
	first := "git add -A"
	var git *exc.CMD
	var err error
	if quit {
		exc.NewCMD("git push origin master").Debug().Execute()
		goto TAG
	}
	git = exc.NewCMD(first).Wd()
	_, err = git.Debug().Do()
	if checkerr(err) {
		goto TAG
	}
	if len(commit) > 0 {
		commit = fmt.Sprintf(`git commit -m %s`, commit)
	} else {
		commit = `git commit -m auto`
	}
	_, err = git.Reset(commit).Do()
	if checkerr(err) {
		goto TAG
	}
	// git.Reset("git push origin master").Execute()
TAG:
	if "no_tag" == tag {
		_, err = git.Reset(fmt.Sprintf("git tag -a %s -m %s", tag, tag)).Do()
		if checkerr(err) {
			return
		}
		git.Reset(fmt.Sprintf("git push origin master --tag %s:%s", tag, tag)).Execute()
	}
}

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
