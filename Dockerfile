# use scratch for K8S
ARG BASEIMAGE=scratch

FROM golang:1.12 as builder
COPY kaput /kaput
WORKDIR /kaput
RUN STATIC=1 make build

# NOTE: cf requires more than scratch
# while K8S is fine with it.
# build image for cf with 
# docker build -t <iamge name> --build-arg BASEIMAGE=alpine:3.9 .

FROM $BASEIMAGE
COPY --from=builder /kaput/bin/kaput /
ENV PORT=8080
EXPOSE 8080
CMD [ "/kaput -s" ]
