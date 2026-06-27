FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -o /out/server ./cmd/server

FROM alpine:3.20
COPY --from=build /out/server /server
EXPOSE 8000
CMD ["/server"]
