# use scratch for K8S production
ARG BASEIMAGE=scratch

FROM golang:1.12 as builder
WORKDIR /go/src/github.com/frnksgr/kaput/kaput
COPY kaput .
RUN CGO_ENABLED=0 GOOS=linux go install ./...

# NOTE: cf requires more than scratch
# while K8S is fine with it.
# build image for cf with 
# docker build -t frnksgr/fibo-cf --build-arg BASEIMAGE=alpine:3.9 .

FROM $BASEIMAGE
COPY --from=builder /go/bin/kaput-server /kaput-server
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/kaput-server"]
