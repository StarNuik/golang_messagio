FROM golang:1.22-alpine3.20 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ARG SERVICE_SRC=./cmd/messages/**

COPY $SERVICE_SRC ./
COPY ./internal/** ./internal/
# COPY *.go ./
# COPY internal ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build

FROM alpine:3.20 AS final
WORKDIR /app
COPY --from=build /app/build .

EXPOSE 8080
CMD ["/app/build"]
