# Build image (toolchain Go hanya dipakai di tahap build)
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY . .

# Cache modul, build cache, dan file intermediate semuanya ditaruh di tmpfs (RAM).
# Tidak ada yang perlu bertahan di layer image, sehingga peak pemakaian disk saat
# build hanya binary (~26 MB) alih-alih ~450 MB.
# -p=2 sesuai board ini (4x A53, RAM terbatas).
RUN --mount=type=tmpfs,target=/tmp \
    --mount=type=tmpfs,target=/go/pkg/mod \
    GOCACHE=/tmp/gocache \
    sh -c 'go mod download && go build -p=2 -o /main .'

# Runtime image: hanya binary hasil build + libc
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /main /app/main

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./main"]