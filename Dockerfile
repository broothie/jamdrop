FROM golang:1.14-alpine AS builder
WORKDIR /go/src/github.com/broothie/queuecumber
COPY . .
RUN apk add --update ca-certificates
RUN go build cmd/server/main.go

FROM alpine:3.7
COPY --from=builder /go/src/github.com/broothie/queuecumber/main main
COPY views views
COPY public public
CMD ./main
