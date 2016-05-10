#!/bin/bash

# called with <read_from> <write_to>
# All stdin is written to the <write_to> file
# <read_from> is catted out

READ_FROM="$1"
WRITE_TO="$2"

# so that tail is cleaned up
#http://stackoverflow.com/a/21807140/907060
trap '[[ -n $tailPid ]] && kill $tailPid 2>/dev/null' EXIT

tail -f $READ_FROM & tailPid=$!

cat /dev/stdin > $WRITE_TO
