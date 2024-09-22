FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o /usr/local/bin/app ./main.go

EXPOSE 8888

ENV DELAY=

CMD ["sh", "-c", "app -delay \"${DELAY}\""]