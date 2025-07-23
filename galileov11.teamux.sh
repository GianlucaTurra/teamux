#!/bin/sh

tmux new -s "V11" -d -n "Back-End" -c /home/giatur/git/galileo//
tmux neww -d -n "Front-End" -c /home/giatur/git/galileo-ui/ -t "V11"
tmux neww -d -n "Libs" -c /home/giatur/git/libs/ -t "V11"
tmux neww -d -n "Other" -c "$HOME" -t "V11"
tmux split-window -t "V11:Back-End" -h -l 25% -c /home/giatur/.config/galileoV11/
tmux split-window -t "V11:Front-End" -v -l 23% -c /home/giatur/git/galileo-ui/
