# syntax=docker/dockerfile:1
FROM docker.io/library/node:16.19-alpine3.16 AS nodebuilder
WORKDIR /app
COPY . ./
ENV VITE_BASE="/dist"
RUN yarn && yarn build
COPY ./public/tiles ./dist

FROM docker.io/library/golang:1.20 AS gobuilder
WORKDIR /go/src/github.com/galexrt/arpanet/
COPY . ./
COPY --from=nodebuilder /app/dist ./dist
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o arpanet .

FROM docker.io/library/alpine:3.17.2
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=gobuilder /go/src/github.com/galexrt/arpanet/arpanet /usr/local/bin

EXPOSE 8080/tcp
EXPOSE 9090/tcp

CMD ["arpanet", "server"]
