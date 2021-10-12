FROM golang:1.16 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -o DatabaseCamp .

FROM alpine:3.13

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime
RUN echo "Asia/Bangkok" >  /etc/timezone

WORKDIR /usr/src/app

COPY --from=builder /src/DatabaseCamp /usr/src/app/DatabaseCamp
COPY --from=builder /src/.env /usr/src/app/.env

EXPOSE 8080
CMD ["./DatabaseCamp"]