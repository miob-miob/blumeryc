FROM golang:1.19.1-bullseye AS build

WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY *.go ./

RUN go build index.go 
# -o /docker-gs-ping
RUN mkdir output && mv index output/index

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /app/output /app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/index"]