FROM golang:1.23-alpine AS build

WORKDIR /usr/src/app

ARG VERSION

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags "-X github.com/dalibormesaric/rplb/internal/config.Version=$VERSION" -v -o /usr/local/bin/app ./cmd/rplb

FROM alpine

WORKDIR /root

COPY --from=build /usr/local/bin/app ./

EXPOSE 8000
EXPOSE 8080

ENV FE=
ENV BE=
ENV ALGO=sticky

CMD ["sh", "-c", "./app -f \"${FE}\" -b \"${BE}\" -a \"${ALGO}\""]
