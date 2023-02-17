FROM golang:alpine

WORKDIR /app

EXPOSE 8082

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o app
CMD ["./app"]