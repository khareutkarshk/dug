FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go build -o backend .

EXPOSE 3001

CMD ["./backend"]