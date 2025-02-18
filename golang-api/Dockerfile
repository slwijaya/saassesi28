# Gunakan base image untuk Golang
FROM golang:1.21.1 as builder

# Set environment variable
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Buat working directory dalam container
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk caching dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file dari direktori lokal ke dalam container
COPY . .

# Unduh kembali dependency jika ada perubahan
RUN go mod tidy

# Compile aplikasi
RUN go build -o go-saas-api main.go

# Gunakan lightweight image untuk runtime
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary dari tahap build ke stage runtime
COPY --from=builder /app/go-saas-api .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./go-saas-api"]
