FROM golang:alpine

# WORKDIR /usr/src/app
WORKDIR /app

# EXPOSE 8090

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o app
CMD ["./app"]


