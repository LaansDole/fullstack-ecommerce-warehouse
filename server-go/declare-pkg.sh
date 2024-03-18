#!/bin/bash

dir=$1

for file in $dir/*.go
do
    if ! grep -q "^package " $file
    then
        echo "package $dir" | cat - $file > temp && mv temp $file
    fi
done
