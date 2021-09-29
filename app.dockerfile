FROM golang:1.17 as builder
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
# todo: this should cache dependencies to avoid downloading them
# on every rebuild
RUN go mod download
COPY . .
RUN go build -v -o booker ./cmd/booker/booker.go

FROM debian:stable-slim
COPY --from=builder /go/src/booker /app/booker
# todo: consider switching to go:embed
COPY --from=builder /go/src/booker/sql /app/booker/sql
ARG PORT=8080
ENV PORT=${PORT}
EXPOSE ${PORT}
CMD ["/app/booker"]
