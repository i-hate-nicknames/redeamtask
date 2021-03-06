FROM golang:1.17 as builder
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
# todo: this should cache dependencies to avoid downloading them
# on every rebuild
RUN go mod download
COPY . .
RUN make build

FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /go/src/booker /app/booker
COPY --from=builder /go/src/sql /app/sql
ARG PORT=8080
EXPOSE ${PORT}
CMD ["/app/booker"]
