# syntax=docker/dockerfile:1
FROM docker.io/library/golang:1.20 AS builder
WORKDIR /go/src/github.com/galexrt/arpanet/
COPY . ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o arpanet .

FROM docker.io/library/alpine:3.17.2
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/galexrt/arpanet/arpanet ./
CMD ["./arpanet"]
