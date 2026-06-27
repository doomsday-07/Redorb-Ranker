FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod main.go ./
RUN go build -o ranker .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/ranker .
COPY sample_100.jsonl ./candidates.jsonl
ENTRYPOINT ["/app/ranker", "candidates.jsonl", "/app/submission.csv"]
