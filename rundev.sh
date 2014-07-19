#!/bin/sh

CMD="echo 'trying to kill old version...' && ps -A | grep main | grep -v grep | grep -v watchmedo | awk '/.*/ {print \$1}' | xargs kill; echo 'building...' && go build main.go && echo 'running...' && ./main --env=dev"

eval $CMD &
watchmedo shell-command \
	--patterns="*.go;*.ini" \
	--recursive \
	--command="$CMD"
