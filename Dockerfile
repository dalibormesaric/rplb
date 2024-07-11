FROM golang:1.22-alpine

LABEL org.opencontainers.image.source=https://github.com/dalibormesaric/rplb
LABEL org.opencontainers.image.description RPLB

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server

EXPOSE 8000
EXPOSE 8080

CMD ["sh", "-c", "app -f ${FE} -b ${BE}"]