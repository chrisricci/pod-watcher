FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/redhat/podwatcher-controller
COPY src/ .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/podwatcher-controller

FROM scratch
COPY --from=builder /go/bin/podwatcher-controller /go/bin/podwatcher-controller
USER 1001
ENTRYPOINT ["/go/bin/podwatcher-controller"]