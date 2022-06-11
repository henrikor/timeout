#!/bin/bash

if $1 -gt 30
then
    echo "Argument må være 30 eller mer"
    exit 1
fi

echo timeout: $1


sov=`expr $1 - 30`
sleep $sov
afplay beep-01a.mp3

sleep 26
afplay beep-01a.mp3
sleep 1
afplay beep-01a.mp3
sleep 1
afplay beep-01a.mp3
sleep 2
afplay cow.mp3
