# Build stage
FROM golang:1.9 AS build-env
ADD . /src

# TODO: Replace with depencdency manager
RUN go get golang.org/x/net/context
RUN go get gopkg.in/olivere/elastic.v5
RUN go get github.com/gorilla/mux

RUN cd /src && CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o sb-web

# Package stage
FROM scratch

WORKDIR /app
COPY --from=build-env src/templates /app/templates
COPY --from=build-env src/public /app/public
COPY --from=build-env /src/sb-web /app/

EXPOSE 8080

ENTRYPOINT ["./sb-web"]
