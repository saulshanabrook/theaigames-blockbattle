#!/bin/bash

# called with <read_from> <write_to>
# All stdin is written to the <write_to> file
# <read_from> is catted out

READ_FROM="$1"
WRITE_TO="$2"

tail -f $READ_FROM &
cat /dev/stdin > $WRITE_TO
