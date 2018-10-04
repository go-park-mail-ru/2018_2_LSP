FROM golang:alpine

COPY . /go/src/lsp

RUN go install lsp

EXPOSE 8080

CMD lsp