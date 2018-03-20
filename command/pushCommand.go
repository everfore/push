package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/everfore/exc"
	cr "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toukii/goutils"
	"github.com/toukii/icat"
)

var Command = &cobra.Command{
	Use:   "push",
	Short: "git push",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		size := len(args)
		if size > 0 {
			viper.Set("commit", args[0])
		} else {
			viper.Set("commit", "the+same+as+the+last")
		}
		Excute()
	},
}

type Repo struct {
	Commition    string
	Branch       string
	RemoteBranch string
	Remote       string
	Tag          string
	git          *exc.CMD
}

func init() {
	Command.PersistentFlags().IntP("squash", "s", 0, "git squash")
	Command.PersistentFlags().StringP("add", "a", "", "git add && commit")
	Command.PersistentFlags().BoolP("quit", "q", false, "quit add all, just push the committed code \n\tgit push")
	Command.PersistentFlags().BoolP("force", "F", false, "-f: push --force")
	Command.PersistentFlags().BoolP("only_commit", "o", false, "-o only commit")
	Command.PersistentFlags().StringP("remote", "R", "origin", "-r $remote \n\tgit push $origin")
	Command.PersistentFlags().StringP("branch", "b", "", "-b $branch \n\tgit push $origin $branch:$remote_branch")
	Command.PersistentFlags().StringP("remote_branch", "r", "", "-rb $remote_branch \n\tgit push $origin $branch:$remote_branch")
	Command.PersistentFlags().StringP("tag", "t", "", "-t $tag \n\tgit tag -a $tag -m $tag;git push $origin --tags $tag:$tag")
	Command.PersistentFlags().BoolP("ignore_dlt_tag", "d", false, "-d: delete the tag after 50 seconds \n\tgit tag -d $tag;git push $origin --tags :$tag")

	viper.BindPFlag("squash", Command.PersistentFlags().Lookup("squash"))
	viper.BindPFlag("add", Command.PersistentFlags().Lookup("add"))
	viper.BindPFlag("quit", Command.PersistentFlags().Lookup("quit"))
	viper.BindPFlag("force", Command.PersistentFlags().Lookup("force"))
	viper.BindPFlag("force", Command.PersistentFlags().Lookup("force"))
	viper.BindPFlag("only_commit", Command.PersistentFlags().Lookup("only_commit"))
	viper.BindPFlag("remote", Command.PersistentFlags().Lookup("remote"))
	viper.BindPFlag("branch", Command.PersistentFlags().Lookup("branch"))
	viper.BindPFlag("remote_branch", Command.PersistentFlags().Lookup("remote_branch"))
	viper.BindPFlag("tag", Command.PersistentFlags().Lookup("tag"))
	viper.BindPFlag("ignore_dlt_tag", Command.PersistentFlags().Lookup("ignore_dlt_tag"))
}

func Excute() error {
	// exc.NewCMD("e icat /Users/toukii/toukii/images/github.svg").Exec(true)
	icat.GitHubSVGThree(time.Now().Unix())
	var repo Repo
	repo.Init()

	return repo.ExcuteGit()
}

func LocalBranch() string {
	bs, err := exc.NewCMD("git rev-parse --abbrev-ref HEAD").DoNoTime()
	if err != nil {
		fmt.Printf("no-HEAD on %s\n", cr.GreenString("* master"))
		return "master"
	}
	cb := string(bs[:len(bs)-1])
	fmt.Printf("git on %s\n", cr.GreenString("* %s", cb))
	return cb
}

func (r *Repo) Init() {
	r.git = exc.NewCMD("").Debug(true)

	r.Commition = viper.GetString("commit")
	branch := LocalBranch()

	vb := viper.GetString("branch")
	if vb != "" {
		if vb == "nil" {
			r.Branch = ""
		} else {
			r.Branch = vb
		}
	} else {
		r.Branch = branch
	}

	vrb := viper.GetString("remote_branch")
	if vrb != "" {
		r.RemoteBranch = vrb
	} else {
		r.RemoteBranch = branch
	}

	vr := viper.GetString("remote")
	if vr != "" {
		r.Remote = vr
	} else {
		r.Remote = "origin"
	}

	r.Tag = viper.GetString("tag")
}

func (r *Repo) ExcuteGit() error {
	r.Commit()
	if squash := viper.GetInt("squash"); squash >= 2 {
		// bs1, err := exc.Bash(fmt.Sprintf("git rebase -i HEAD~%d", squash)).Do()
		// fmt.Printf("%s", bs1)
		// return err

		bs, err := exc.Bash(fmt.Sprintf(`osascript <<EOF
tell application "System Events"
    tell process "iTerm2"
        # set frontmost to true
        keystroke "git rebase -i HEAD~%d"
        keystroke return
	end tell
end tell
EOF`, squash)).Exec(true).Do()

		fmt.Printf("%s", bs)
		return err

	}
	r.Push()
	return nil
}

func (r *Repo) Status() bool {
	bs, err := exc.NewCMD(r.status()).Do()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if strings.Contains(goutils.ToString(bs), "nothing to commit, working tree clean") {
		return false
	}
	// fmt.Printf("%s\n", bs)
	return true
}

func (r *Repo) status() string {
	return `git status`
}

func (r *Repo) Commit() {
	cstatus := r.Status()
	if viper.GetBool("quit") {
		return
	}

	if !cstatus {
		return
	}

	add := viper.GetString("add")
	fmt.Println("add:", add)
	if add != "" {
		r.git.Reset("git add " + add).Execute()
	} else {
		if !viper.GetBool("only_commit") {
			r.git.Reset("git add -A").Execute()
		}
	}

	r.git.Reset(r.commit()).Execute()
	r.Status()
}

func (r *Repo) commit() string {
	return fmt.Sprintf(`git commit -m %s`, r.Commition)
}

func (r *Repo) Push() {
	if viper.GetBool("only_commit") {
		return
	}
	r.git.Reset(r.push()).Execute()

	fmt.Println("push...")
	if pushtag, cmd := r.tag(); pushtag {
		r.git.Reset(cmd).Execute()
	} else {
		fmt.Println(pushtag, cmd)
	}
}

func (r *Repo) push() string {
	force := ""
	if viper.GetBool("force") {
		force = "-f "
	}
	return fmt.Sprintf("git push %s%s %s:%s", force, r.Remote, r.Branch, r.RemoteBranch)
}

func (r *Repo) tag() (bool, string) {
	if r.Tag == "" {
		return false, ""
	}
	_, err := r.git.Reset(fmt.Sprintf("git tag -a %s -m %s", r.Tag, r.Tag)).DoNoTime()
	if err != nil {
		return false, ""
	}
	return true, fmt.Sprintf("git push %s --tags %s:%s", r.Remote, r.Tag, r.Tag)
}
