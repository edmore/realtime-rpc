FROM golang:1.21-alpine as builder
WORKDIR /app/
COPY . /app/

RUN ls -alh
RUN CGO_ENABLED=0 go build -o /app/jit-server .

FROM scratch
COPY --from=builder /app/jit-server /jit-server
EXPOSE 50051
CMD ["/jit-server"]