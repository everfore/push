package main

import (
	"flag"
	"fmt"

	"github.com/everfore/exc"
)

var (
	quit           = false
	featured       = false
	commit         = ""
	only_commit    = ""
	remote         = ""
	branch         = ""
	remote_branch  = ""
	tag            = ""
	status         = false
	ignore_dlt_tag = false // delete the tag after 25 sec
)

func init() {
	flag.BoolVar(&quit, "q", false, "-q: quit add all, just push the committed code \n\tgit push")
	flag.BoolVar(&featured, "f", false, "-f: push feature/branch")
	flag.StringVar(&only_commit, "om", "", "-om $commit,only commit no push \n\tgit commit $commit")
	flag.BoolVar(&status, "s", false, "git status")
	flag.StringVar(&commit, "m", "", "-m $commit, commit and push \n\tgit add -A;git commit $commit;git push")
	flag.StringVar(&remote, "r", "origin", "-r $remote \n\tgit push $origin")
	flag.StringVar(&branch, "b", "develop", "-b $branch \n\tgit push $origin $branch:$remote_branch")
	flag.StringVar(&remote_branch, "rb", "develop", "-rb $remote_branch \n\tgit push $origin $branch:$remote_branch")
	flag.StringVar(&tag, "t", "", "-t $tag \n\tgit tag -a $tag -m $tag;git push $origin --tags $tag:$tag")
	flag.BoolVar(&ignore_dlt_tag, "d", false, "-d: delete the tag after 50 seconds \n\tgit tag -d $tag;git push $origin --tags :$tag")
}

func main() {
	flag.Parse()
	first := "git add ."
	var git *exc.CMD
	git = exc.NewCMD(first).Wd().Debug()
	var err error
	if status {
		git.Reset("git status").Execute()
		return
	}
	branch = currentBranch()
	remote_branch = branch
	if featured {
		branch = fmt.Sprintf("feature/%s", branch)
		remote_branch = fmt.Sprintf("feature/%s", remote_branch)
	}

	if quit {
		git.Reset(fmt.Sprintf("git push %s %s:%s", remote, branch, remote_branch)).Execute()
		goto TAG
	}
	if len(only_commit) > 0 {
		fmt.Println("only_commit", only_commit)
		only_commit = fmt.Sprintf(`git commit -m %s`, only_commit)
		_ = git.Reset(only_commit).Execute()
		return
	}
	_ = git.Execute()
	if len(commit) > 0 {
		commit = fmt.Sprintf(`git commit -m %s`, commit)
	} else {
		if len(tag) > 0 {
			commit = fmt.Sprintf("git commit -m %s", tag)
		} else {
			commit = "git commit -m same_as_the_last_commit"
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
		fmt.Println("ignore_dlt_tag:", ignore_dlt_tag)
		if ignore_dlt_tag {
			return
		}

		fmt.Printf("%d seconds later...\n", 25)
		git.Reset(fmt.Sprintf("git tag -d %s", tag)).ExecuteAfter(25)
		git.Reset(fmt.Sprintf("git push %s --tag :%s", remote, tag)).Execute()
	}
}

func currentBranch() string {
	bs, err := exc.NewCMD("git rev-parse --abbrev-ref HEAD").DoNoTime()
	if err != nil {
		panic(err)
	}
	cb := string(bs[:len(bs)-1])
	fmt.Printf("* %s\n", cb)
	return cb
}

func checkerr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
