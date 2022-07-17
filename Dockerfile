FROM golang:alpine as builder

#copy all main folder
COPY . /root/src/library-demo

#set GO variables
ENV GOPATH /root/
ENV GOROOT /usr/local/go
ENV GOBIN /root/bin
ENV PATH /usr/local/go/bin:$PATH
ENV GO111MODULE=off

RUN set -eux; \
	apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go git

RUN cd /root/ &&\
	mkdir bin pkg &&\
	cd /root/src/library-demo &&\
	pwd && ls -ll &&\
	go install &&\
	go clean 

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /root/src/library-demo/index.html .
COPY --from=builder /root/src/library-demo/.env .
COPY --from=builder /root/bin/library-demo .
ENTRYPOINT ["/root/library-demo"]
