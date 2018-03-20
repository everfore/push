#!/bin/sh

git add .
git commit -m "auto"
echo "======="
osascript <<EOF
tell application "iTerm 2"
	tell application "System Events" 
         display dialog "bal bal" 
    end tell
end tell
EOF
echo "======="