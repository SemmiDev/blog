FROM alpine:3.13

WORKDIR /app

COPY bin/main .

COPY app.env .

EXPOSE 3030

RUN CGO_ENABLED=0

CMD ["/app/main"]