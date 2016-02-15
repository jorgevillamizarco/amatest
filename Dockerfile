FROM golang:onbuild
 
ADD . /go/src/github.com/jorgevillamizarco/amatest
RUN go install /go/src/github.com/jorgevillamizarco/amatest
ENTRYPOINT /go/bin/amatest_server
 
EXPOSE 8080