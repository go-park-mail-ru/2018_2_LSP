FROM golang:alpine

RUN apk add --no-cache git

ADD . /go/src/github.com/go-park-mail-ru/2018_2_LSP

RUN cd /go/src/github.com/go-park-mail-ru/2018_2_LSP && go get ./...

RUN go install github.com/go-park-mail-ru/2018_2_LSP

ENTRYPOINT /go/bin/2018_2_LSP

EXPOSE 8080