#!/bin/sh

git add .
git commit -m "auto"
echo "======="
osascript <<EOF
tell application "System Events"
    tell process "iTerm2"
        set frontmost to true
	end tell
end tell	
EOF
# osascript xxx
echo "======="