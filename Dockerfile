FROM golang:1.14-alpine AS bin_builder
WORKDIR /go/src/jamdrop
COPY . .
RUN apk add --update ca-certificates
RUN scripts/build.sh

FROM node:14 AS bundle_builder
WORKDIR /usr/src/app
COPY package.json .
COPY yarn.lock .
COPY frontend frontend
RUN yarn
RUN yarn build

FROM alpine:3.7
COPY --from=bin_builder /go/src/jamdrop/main main
COPY --from=bundle_builder /usr/src/app/dist dist
CMD ./main
