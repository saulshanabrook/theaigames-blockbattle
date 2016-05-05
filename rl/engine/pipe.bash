#!/bin/bash

# called with <read_from> <write_to>
# All stdin is written to the <write_to> file
# <read_from> is catted out

READ_FROM="$1"
WRITE_TO="$2"

# so that tail is cleaned up
# http://stackoverflow.com/a/2173421/907060
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

tail -f $READ_FROM &
cat /dev/stdin > $WRITE_TO
