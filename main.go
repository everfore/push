package main

import (
	"flag"
	"fmt"

	"github.com/everfore/exc"
)

var (
	quit          = false
	commit        = ""
	remote        = ""
	branch        = ""
	remote_branch = ""
	tag           = ""
	no_dlt_tag    = false // no delete the tag
)

func init() {
	flag.BoolVar(&quit, "q", false, "-q: quit add all, just push the committed code \n\tgit push")
	flag.StringVar(&commit, "m", "", "-m: commit content, \n\tgit add -A;git commit $commit")
	flag.StringVar(&remote, "r", "origin", "-r origin \n\tgit push $origin")
	flag.StringVar(&branch, "b", "master", "-b master \n\tgit push $origin $branch:$remote_branch")
	flag.StringVar(&remote_branch, "rb", "master", "-rb master \n\tgit push $origin $branch:$remote_branch")
	flag.StringVar(&tag, "t", "", "-t: tag \n\tgit tag -a $tag -m $tag;git push $origin --tags $tag:$tag")
	flag.BoolVar(&no_dlt_tag, "d", false, "-d: delete the tag after 50 seconds \n\tgit tag -d $tag;git push $origin --tags :$tag")
}

func main() {
	flag.Parse()
	first := "git add -A"
	var git *exc.CMD
	var err error
	if quit {
		git = exc.NewCMD(fmt.Sprintf("git push %s %s:%s", remote, branch, remote_branch)).Debug()
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
	git.Reset(fmt.Sprintf("git push %s %s:%s", remote, branch, remote_branch)).Execute()

TAG:
	if len(tag) > 0 {
		_, err = git.Reset(fmt.Sprintf("git tag -a %s -m %s", tag, tag)).DoNoTime()
		if checkerr(err) {
			return
		}
		git.Reset(fmt.Sprintf("git push %s --tag %s:%s", remote, tag, tag)).Execute()
		fmt.Println("no_del_tag:", no_dlt_tag)
		if no_dlt_tag {
			return
		}

		fmt.Printf("%d seconds later...\n", 50)
		git.Reset(fmt.Sprintf("git tag -d %s", tag)).ExecuteAfter(15)
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
