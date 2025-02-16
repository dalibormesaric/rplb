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

ENV RPLB_F=
ENV RPLB_B=
ENV RPLB_A=$RPLB_A

CMD ["sh", "-c", "./app -f \"${RPLB_F}\" -b \"${RPLB_B}\" -a \"${RPLB_A}\""]
