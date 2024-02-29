##Builder Image
FROM golang:1.22.0-alpine3.19 as builder
ENV GO111MODULE=on
RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates
COPY . /risk-management
WORKDIR /risk-management/cmd
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/application

#s Run Image
FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /risk-management/cmd/bin/application application
EXPOSE 8888
ENTRYPOINT ["./application"]