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
	flag.BoolVar(&quit, "q", false, "-q: quit add all, just push the committed code \n[git push]")
	flag.StringVar(&commit, "m", "", "-m: commit content, \n[git add -A;git commit $commit]")
	flag.StringVar(&remote, "r", "origin", "-r origin \n[git push $origin]")
	flag.StringVar(&tag, "t", "", "-t: tag \n[git tag -a $tag -m $tag;git push $origin --tags $tag:$tag]")
	flag.BoolVar(&dlt_tag, "d", true, "-d: delete the tag after 50 seconds \n[git tag -D $tag;git push $origin --tags :$tag]")
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
	_ = git.Debug().Execute()
	if len(commit) > 0 {
		commit = fmt.Sprintf(`git commit -m %s`, commit)
	} else {
		if len(tag) > 0 {
			commit = fmt.Sprintf("git commit -m %s", tag)
		} else {
			commit = `git commit -m auto`
		}
	}
	_ = git.Reset(commit).Execute()
	git.Reset(fmt.Sprintf("git push %s master", remote)).Execute()

TAG:
	if len(tag) > 0 {
		_, err = git.Reset(fmt.Sprintf("git tag -a %s -m %s", tag, tag)).DoNoTime()
		if checkerr(err) {
			return
		}
		//		git.Reset(fmt.Sprintf("git push %s master --tag %s:%s", remote, tag, tag)).Execute()
		if !dlt_tag {
			return
		}
		fmt.Printf("%d seconds later...\n", 50)
		git.Reset(fmt.Sprintf("git tag -d %s", tag)).ExecuteAfter(5)
		git.Reset(fmt.Sprintf("git push %s --tag :%s", remote, tag)).Execute()
	}
}

func checkerr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
