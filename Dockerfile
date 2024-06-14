FROM golang:1.22-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /doc-app

CMD ["/doc-app"]
