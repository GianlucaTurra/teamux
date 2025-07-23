tmux new -s "JMRP" -d -n "Back-End" -c /home/giatur/git/jmrp/
tmux neww -d -n "Front-End" -c /home/giatur/git/jmrp-ui/ -t "JMRP"
tmux neww -d -n "Other" -c "$HOME" -t "JMRP"
tmux split-window -t "JMRP:Back-End" -h -l 25% -c /home/giatur/.config/jmrp
tmux split-window -t "JMRP:Front-End" -v -l 23% -c /home/giatur/git/jmrp-ui/
