FROM golang:1.16 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags "-X 'DatabaseCamp/router.BuildTime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`'\
    -X 'DatabaseCamp/router.BuildCommit=`git rev-parse --short HEAD`'"\
    -o DatabaseCamp .

FROM alpine:3.13



RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime
RUN echo "Asia/Bangkok" >  /etc/timezone

WORKDIR /usr/src/app

COPY --from=builder /src/DatabaseCamp /usr/src/app/DatabaseCamp
COPY --from=builder /src/.env /usr/src/app/.env
COPY --from=builder /src/service_account.json /usr/src/app/service_account.json

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

EXPOSE 8080
CMD ["./DatabaseCamp"]