FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE  8000
CMD ["go", "run","."]