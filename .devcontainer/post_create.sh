#!/bin/bash
set -e

echo export PATH="$PATH:$(go env GOPATH)/bin" >> ~/.bashrc
echo "source /usr/share/bash-completion/completions/git" >> ~/.bashrc
