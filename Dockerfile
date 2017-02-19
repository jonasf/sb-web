FROM scratch

WORKDIR /app

COPY sb-web /app/

EXPOSE 8080

ENTRYPOINT ["./sb-web"]