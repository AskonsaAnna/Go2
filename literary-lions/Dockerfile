FROM golang:1.22.2
WORKDIR /app

COPY . .

RUN go get
RUN go build -o bin .

ENTRYPOINT ["/app/bin"]
EXPOSE 8080