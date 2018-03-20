#!/bin/sh

function psh(){
name=$1
$(git add .)
$(git commit -m $name)
osascript <<EOF
tell application "System Events"
    tell process "iTerm2"
        set frontmost to true
        keystroke "git rebase -i HEAD~2"
        keystroke return
	end tell
end tell
EOF
return $?
}
