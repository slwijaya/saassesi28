# Gunakan base image untuk Golang versi 1.21.1
FROM golang:1.21.1 as builder

# Set environment variable untuk Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Buat direktori kerja dalam container
WORKDIR /app

# Copy file go.mod dan go.sum dari folder golang-api untuk caching dependency
COPY golang-api/go.mod golang-api/go.sum ./

# Unduh semua dependency yang diperlukan sebelum menyalin kode sumber
RUN go mod download

# Copy seluruh kode sumber dari folder golang-api ke dalam container
COPY golang-api/ .

# Unduh ulang dependency jika ada perubahan
RUN go mod tidy

# Compile aplikasi menjadi binary
RUN go build -o go-saas-api main.go

# Gunakan lightweight image untuk runtime
FROM alpine:latest

# Set direktori kerja di dalam container runtime
WORKDIR /root/

# Copy binary yang telah dikompilasi dari tahap build ke tahap runtime
COPY --from=builder /app/go-saas-api .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./go-saas-api"]
