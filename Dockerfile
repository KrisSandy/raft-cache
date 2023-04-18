FROM golang:1.20.3 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o raft-cache-app ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o raftctl ./cmd/raftctl


FROM alpine:3.17.3 as raft-cache-app
WORKDIR /
COPY --from=builder /app/raft-cache-app raft-cache-app
COPY --from=builder /app/raftctl /bin/raftctl
CMD ["/raft-cache-app"]