FROM scratch

WORKDIR /app

ADD templates /app/templates
ADD public /app/public
COPY sb-web /app/

EXPOSE 8080

ENTRYPOINT ["./sb-web"]