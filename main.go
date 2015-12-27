package main

import (
	"flag"
	"fmt"
	"github.com/everfore/exc"
)

var (
	quit    = false
	commit  = ""
	remote  = ""
	tag     = ""
	dlt_tag = true
)

func init() {
	flag.BoolVar(&quit, "q", false, "-q: quit add all")
	flag.StringVar(&commit, "m", "", "-m: commit content")
	flag.StringVar(&remote, "r", "origin", "-r origin")
	flag.StringVar(&tag, "t", "no_tag", "-t: tag")
	flag.BoolVar(&dlt_tag, "d", true, "-d: delete the tag after 50 seconds")
}

func main() {
	flag.Parse()
	first := "git add -A"
	var git *exc.CMD
	var err error
	if quit {
		git = exc.NewCMD(fmt.Sprintf("git push %s master", remote)).Debug()
		git.Execute()
		goto TAG
	}
	git = exc.NewCMD(first).Wd()
	_, _ = git.Debug().DoNoTime()
	if len(commit) > 0 {
		commit = fmt.Sprintf(`git commit -m %s`, commit)
	} else {
		if len(tag) > 0 {
			commit = fmt.Sprintf("git commit -m %s", tag)
		} else {
			commit = `git commit -m auto`
		}
	}
	_, _ = git.Reset(commit).DoNoTime()
	git.Reset(fmt.Sprintf("git push %s master", remote)).Execute()

TAG:
	if "no_tag" != tag {
		_, err = git.Reset(fmt.Sprintf("git tag -a %s", tag)).DoNoTime()
		if checkerr(err) {
			return
		}
		git.Reset(fmt.Sprintf("git push %s master --tag %s:%s", remote, tag, tag)).Execute()
		if !dlt_tag {
			return
		}
		fmt.Printf("%d seconds later...\n", 50)
		git.Reset(fmt.Sprintf("git tag -d %s", tag)).ExecuteAfter(50)
		git.Reset(fmt.Sprintf("git push %s --tag :%s", remote, tag)).Execute()
	}
}

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
