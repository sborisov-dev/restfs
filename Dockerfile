FROM golang:1.9
EXPOSE 8080
VOLUME /tmp/restfs
WORKDIR /restfs
ADD . /go/src/rest-fs/
RUN go get github.com/julienschmidt/httprouter && \
    go build -o /restfs/app /go/src/rest-fs/main.go
CMD ["/restfs/app"]