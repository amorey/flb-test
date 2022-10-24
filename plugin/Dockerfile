FROM golang:1.19.2 AS gobuilder

RUN mkdir plugin
WORKDIR /plugin

# install dependencies (for cache)
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy code and compile
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o out_plugin.so .

# -----------------------------------------------------------------------------

FROM fluent/fluent-bit:1.9.9-debug

COPY --from=gobuilder /plugin/out_plugin.so /fluent-bit/bin/

ENTRYPOINT ["/fluent-bit/bin/fluent-bit", "--plugin", "/fluent-bit/bin/out_plugin.so"]
