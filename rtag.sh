#!/bin/sh

for tag in `git tag`;do
echo $tag
git tag -d $tag
git push origin --tag :$tag
done
