FROM golang:1.26-alpine

ENV DB_HOST=host.docker.internal DB_PASSWORD=rahasia DB_USER=postgres DB_NAME=wallet_db DB_PORT=5432 DB_SSLMODE=disable SECRET_KEY=kuncirahasia-sangat-rahasia-32ch

WORKDIR /app

COPY . .

RUN go build -o app

CMD ./app