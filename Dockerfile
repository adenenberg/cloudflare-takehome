FROM golang as builder
WORKDIR /build/api
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o api

# post build stage
FROM alpine
WORKDIR /root
COPY --from=builder /build/api/api .
EXPOSE 8080
CMD ["./api"]

# FROM golang:latest

# ENV GO111MODULE=on
# ENV PORT=8080
# WORKDIR /app
# COPY go.mod /app
# COPY go.sum /app

# RUN go mod download
# RUN go install -mod=mod github.com/githubnemo/CompileDaemon
# COPY . /app
# ENTRYPOINT CompileDaemon --build="go build -o main" --command=./main