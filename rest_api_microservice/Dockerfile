FROM golang:1.22.2-alpine AS build-env

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=build-env /app/app .

RUN chmod +x app

CMD [ "./app" ]
