package main

import (
	"github.com/everfore/exc"
	"os"
)

func checkerr(err error) bool {
	if err != nil {
		return true
	}
	return false
}
func main() {
	git := exc.NewCMD("git add -A")
	_, err := git.Debug().Do()
	if checkerr(err) {
		os.Exit(-1)
	}
	_, err = git.Reset(`git commit -m "auto"`).Wd().Do()
	if checkerr(err) {
		// os.Exit(-1)
	}
	git.Reset("git push origin master").Execute()
}
