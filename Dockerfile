FROM golang:1.20-alpine3.16 as build
ENV CGO_ENABLED 0
WORKDIR $GOPATH/src/github.com/frantjc/sneasler
ARG semver=0.0.0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w -X github.com/frantjc/sneasler.Semver=${semver}" -o /assets/sneasler ./cmd/sneasler

FROM alpine:3.16
ENTRYPOINT ["sneasler"]
COPY --from=build /assets /usr/local/bin
