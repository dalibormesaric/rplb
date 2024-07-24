FROM golang:1.22-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/rplb

ARG VERSION
ENV VERSION=$VERSION

EXPOSE 8000
EXPOSE 8080

CMD ["sh", "-c", "app -f \"${FE}\" -b \"${BE}\" -v ${VERSION}"]
