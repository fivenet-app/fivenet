# syntax=docker/dockerfile:1
FROM docker.io/library/golang:1.20 AS gobuilder
WORKDIR /go/src/github.com/galexrt/arpanet/
COPY . ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o arpanet .

FROM docker.io/library/node:16.19-alpine3.16 AS nodebuilder
WORKDIR /app
COPY . ./
RUN yarn && yarn build

FROM docker.io/library/alpine:3.17.2
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=gobuilder /go/src/github.com/galexrt/arpanet/arpanet ./
COPY --from=nodebuilder /app/dist ./
CMD ["./arpanet"]
