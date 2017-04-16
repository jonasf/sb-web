FROM scratch

WORKDIR /app

ADD templates /app/templates
COPY sb-web /app/

EXPOSE 8080

ENTRYPOINT ["./sb-web"]