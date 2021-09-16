FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN GOCGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o /build/main
FROM scratch
COPY --from=build /build/main /main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY data /data
CMD ["/main"]