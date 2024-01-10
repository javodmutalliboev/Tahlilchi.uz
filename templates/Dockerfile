FROM golang:1.21.5

LABEL maintainer="Javod Mutalliboev <javodmutalliboev@gmail.com>"

WORKDIR /home/javod/Desktop/Projects/Tahlilchi.uz/applications/server

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]