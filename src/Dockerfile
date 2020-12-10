FROM golang:buster as builder
RUN update-ca-certificates
WORKDIR /go/src/app
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/parklakes-viberbot .

FROM gcr.io/distroless/base
COPY --from=builder /go/bin/parklakes-viberbot /app/parklakes-viberbot
ENTRYPOINT ["/app/parklakes-viberbot"]