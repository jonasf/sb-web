# Build stage
FROM golang:1.10 AS build-env
ADD . $GOPATH/src/build
WORKDIR $GOPATH/src/build

## Install dependencies
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only

RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o sb-web

# Package stage
FROM scratch

WORKDIR /app
# NOTE: hard coded $GOPATH
COPY --from=build-env /go/src/build/templates /app/templates
COPY --from=build-env /go/src/build/public /app/public
COPY --from=build-env /go/src/build/sb-web /app/

EXPOSE 8080

ENTRYPOINT ["./sb-web"]
