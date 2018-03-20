#!/bin/sh

osascript <<EOF
tell application "iTerm2"
    activate
    set t to (make new terminal)
    tell t
        tell (make new session at the end of sessions)
            exec command "cd Downloads"
            exec command "clear"
            exec command "pwd"
        end tell
    end tell
end tell
EOF

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
