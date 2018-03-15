# Build stage
FROM golang:1.10 AS build-env
ADD . $GOPATH/src/github.com/jonasf/systembolaget-beer-releases
WORKDIR $GOPATH/src/github.com/jonasf/systembolaget-beer-releases

## Install dependencies
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only

RUN cd cmd/systembolaget-beer-releases && CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o ../../build/systembolaget-beer-releases

# Package stage
FROM scratch

WORKDIR /app
# NOTE: hard coded $GOPATH
COPY --from=build-env /go/src/github.com/jonasf/systembolaget-beer-releases/cmd/systembolaget-beer-releases/templates /app/templates
COPY --from=build-env /go/src/github.com/jonasf/systembolaget-beer-releases/cmd/systembolaget-beer-releases/public /app/public
COPY --from=build-env /go/src/github.com/jonasf/systembolaget-beer-releases/build/systembolaget-beer-releases /app/

EXPOSE 8080

ENTRYPOINT ["./systembolaget-beer-releases"]
