FROM golang:alpine3.20

RUN apk update && apk add --no-cache git
RUN apk add --update alpine-sdk

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go env -w CGO_ENABLED=1
RUN go build -o binary
VOLUME ["/app/form_config"]

ENTRYPOINT ["/app/binary"]
