FROM golang:latest
WORKDIR /app
COPY . /app
RUN go get github.com/joho/godotenv
RUN go get github.com/streadway/amqp
RUN go get github.com/sirupsen/logrus
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/zbindenren/logrus_mail
CMD ["go", "run", "main.go", "hub.go", "client.go", "socket.go", "query.go"]
EXPOSE 80

