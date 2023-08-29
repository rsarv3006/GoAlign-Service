#!/bin/bash

# Get the number of times to perform the action
read -p "Enter number of times: " num

# Perform the action in a loop
for ((i=1; i<=$num; i++)); do
  ./test.sh &
  disown
done
