FROM golang:1.20-alpine3.16 AS go
ENV CGO_ENABLED 0
WORKDIR $GOPATH/src/github.com/frantjc/sneasler
ARG semver=0.0.0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w -X github.com/frantjc/sneasler.Semver=${semver}" -o /assets/sneasler ./cmd/sneasler

FROM node:19-alpine3.16 AS remix
WORKDIR /src/github.com/frantjc/sneasler
COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile
COPY . .
RUN yarn build

FROM node:19-alpine3.16
ENTRYPOINT ["sneasler"]
ENV SNEASLER_JS_ENTRYPOINT /app/index.js
COPY public/ /app/public
COPY --from=go /assets /usr/local/bin
COPY --from=remix /src/github.com/frantjc/sneasler/dist /app
