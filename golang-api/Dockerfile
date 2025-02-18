# Gunakan base image untuk Golang
FROM golang:1.21.1 AS builder

# Set environment variable
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Buat working directory dalam container
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu
COPY go.mod go.sum ./

# Unduh semua dependency sebelum menyalin kode sumber
RUN go mod download

# Copy semua file dari direktori lokal ke dalam container
COPY . .

# Pastikan go.mod dan go.sum sudah di dalam direktori
RUN ls -la /app

# Download dependencies
RUN go mod tidy

# Compile aplikasi
RUN go build -o go-saas-api main.go

# Gunakan lightweight image untuk runtime
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary dari tahap build ke stage runtime
COPY --from=builder /app/go-saas-api .

# Expose port 8080 untuk HTTP
EXPOSE 8080

# Jalankan aplikasi
CMD ["./go-saas-api"]
