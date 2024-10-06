FROM --platform=$BUILDPLATFORM golang:1.23.1-bookworm AS vscode
WORKDIR /app
COPY . /app
# RUN  ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime


FROM --platform=$BUILDPLATFORM golang:1.23.1-bookworm AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bot main.go


FROM --platform=$BUILDPLATFORM gcr.io/distroless/base-debian12:latest
ARG USERNAME=nonroot
ARG GROUPNAME=nonroot
# ENV TZ Asia/Tokyo
WORKDIR /app
COPY --chown=${USERNAME}:${GROUPNAME} --chmod=100  --from=build /app/bot /app/bot
USER ${USERNAME}
ENTRYPOINT [ "/app/bot" ]
