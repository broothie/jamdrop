FROM golang:1.14-alpine AS bin_builder
WORKDIR /go/src/github.com/broothie/queuecumber
COPY . .
RUN apk add --update ca-certificates
RUN go build cmd/server/main.go

FROM node:14 AS bundle_builder
WORKDIR /usr/src/app
COPY package.json .
COPY yarn.lock .
COPY frontend frontend
RUN yarn
RUN yarn build

FROM alpine:3.7
COPY --from=bin_builder /go/src/github.com/broothie/queuecumber/main main
COPY --from=bundle_builder /usr/src/app/public public
CMD ./main
