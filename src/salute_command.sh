#!/bin/bash
CHOICES="goodmornin goodafternoon goodnight"
# shellcheck disable=SC2086
CHOSEN=$(gum choose $CHOICES)

if [ -z "$1" ] 
then 
  USER=$(gum input --placeholder "Insert user to salute")
else
  USER="$1"
fi

echo "$CHOSEN, $USER"