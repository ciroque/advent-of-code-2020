#!/usr/bin/env bash

die () {
    echo >&2 "$@"
    exit 1
}

[ "$#" -eq 1 ] || die "1 argument required, $# provided"
echo $1 | grep -E -q '^[a-zA-Z0-9-]+$' || die "Directory name argument required, '$1' provided"

DIR_NAME=$1

mkdir $DIR_NAME

cp ./template.go $DIR_NAME/solution.go

touch $DIR_NAME/puzzle.txt
touch $DIR_NAME/example-input.dat
touch $DIR_NAME/puzzle-input.dat

git add $DIR_NAME/
