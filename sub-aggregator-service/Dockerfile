FROM golang:1.26.1-alpine
WORKDIR /app/sub-aggregator-service
COPY . .
RUN go build -o ./bin/main ./cmd
EXPOSE 7000

CMD ["./bin/main"]
