FROM golang:1.15-alpine AS builder

ADD . /src

WORKDIR /src

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/server main.go

FROM scratch

COPY --from=builder /bin/server /bin/server
ADD ./migrations /migrations

ENTRYPOINT [ "/bin/server" ]