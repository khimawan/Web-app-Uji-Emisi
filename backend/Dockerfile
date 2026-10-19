# Build image (toolchain Go hanya dipakai di tahap build)
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build cache dan $WORK intermediate dipindah ke tmpfs (RAM) supaya tidak
#Didakan lewat layer image pada eMMC yang hanya 6.5 GB.
# -p=2 sesuai board ini (4x A53, RAM terbatas).
RUN --mount=type=tmpfs,target=/tmp \
    GOCACHE=/tmp/gocache \
    go build -p=2 -o /main .

# Runtime image: hanya binary hasil build + libc
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /main /app/main

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./main"]