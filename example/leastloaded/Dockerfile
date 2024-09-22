FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY . .

EXPOSE 8888

ENV DELAY=1

CMD ["sh", "-c", "go run main.go -delay \"${DELAY}\""]