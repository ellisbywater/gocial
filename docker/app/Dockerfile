FROM golang:1.22-bookworm


RUN apt install direnv


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN direnv allow .

RUN go build -o ./main ./cmd/api/*.go

EXPOSE 8080

CMD [ "./main" ]