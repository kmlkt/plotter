FROM golang:latest AS builder
WORKDIR /
COPY ./ ./
RUN go mod download
RUN cd /cmd && CGO_ENABLED=0 go build -o ./main


FROM scratch
WORKDIR /
COPY --from=builder /cmd/main /main
EXPOSE 80
CMD ["/main"]
