#!/bin/bash

FILE_PATH=$1
TARGET_FILE=$2

./stackcollapse-go.pl $FILE_PATH | ./flamegraph.pl > $TARGET_FILE
