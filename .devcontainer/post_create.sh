#!/bin/bash
set -e

apt-get update && apt-get install -y vim git

# tools
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/go-delve/delve/cmd/dlv@latest
echo export PATH="$PATH:$(go env GOPATH)/bin" >> ~/.bashrc

git config --local core.editor vim
git config --local pull.rebase false
echo "source /usr/share/bash-completion/completions/git" >> ~/.bashrc

# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
