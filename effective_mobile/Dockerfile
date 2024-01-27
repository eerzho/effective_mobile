FROM golang:1.21

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

ENTRYPOINT ["air", "-c", ".air.toml"]
