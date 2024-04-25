FROM golang:latest as build
WORKDIR /app

COPY . .
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o out/main

FROM alpine:latest as final-stage
WORKDIR /app
COPY --from=build /app/out/main /app/main
COPY --from=build /app/.env /app/.env
EXPOSE 8080

ENTRYPOINT ["/app/main"]