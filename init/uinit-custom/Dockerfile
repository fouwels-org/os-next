FROM golang:1-alpine

WORKDIR /build
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
RUN go mod download

COPY . /build
RUN go build -o uinit-custom .
CMD ["cp", "/build/uinit-custom", "/out/uinit-custom"]