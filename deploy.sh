#!/bin/bash

# Get Go install.
wget https://go.dev/dl/go1.17.3.linux-amd64.tar.gz

# Unzip Go installation.
sudo tar -C /usr/local/ -xzf go1.17.3.linux-amd64.tar.gz

# Export Go to PATH env var.
export PATH=$PATH:/usr/local/go/bin

# Turn off Go modules.
export GO111MODULE=off

# Change to VoteTo directory.
cd go/src/github.com/NathanielRand/VoteTo/

# Pull most recent updates from Github repo.
git clone github.com/NathanielRand/VoteTo

# Build go program.
go build

# Get previously running background process PID and assign to var PID
pid=$(pgrep VoteTo)

# Kill previously running background process.
kill pid

# Run and detach updated go program into a new process.
nohup ./VoteTo &