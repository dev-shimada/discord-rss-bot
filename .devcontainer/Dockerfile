FROM --platform=$BUILDPLATFORM golang:1.23.1-bookworm
WORKDIR /app
COPY . /app
# RUN  ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN <<EOF
apt-get update && apt-get install -y vim git
# tools
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/go-delve/delve/cmd/dlv@latest
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
EOF
