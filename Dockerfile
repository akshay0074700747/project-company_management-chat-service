FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o chat-service

FROM busybox:latest

WORKDIR /chat-service

COPY --from=build /app/cmd/chat-service .

COPY --from=build /app/cmd/.env .

EXPOSE 50006

CMD [ "./chat-service" ]