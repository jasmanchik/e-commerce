FROM --platform='linux/arm64/v8' golang:1.20
WORKDIR /app
RUN go install github.com/cosmtrek/air@v1.43.0
COPY ./src .
COPY ./.air.toml /.air.toml
RUN go mod tidy
RUN go build -o main .
CMD ["air", "-c", "/.air.toml"]