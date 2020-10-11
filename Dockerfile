FROM golang:1.15 as build

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -v -o backupshq

FROM scratch

COPY --from=build /build/backupshq /usr/bin/backupshq

CMD ["/usr/bin/backupshq"]
