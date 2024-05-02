FROM golang:1.21-alpine as builder
WORKDIR /app/
COPY . /app/

RUN ls -alh
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/jit-client .

FROM scratch
COPY --from=builder /app/jit-client /jit-client

CMD ["/jit-client"]