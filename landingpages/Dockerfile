# Gunakan image Nginx untuk melayani file statis
FROM nginx:latest

# Copy file HTML, CSS, JS ke dalam direktori default Nginx
COPY . /usr/share/nginx/html

# Expose port 80 untuk HTTP
EXPOSE 80

# Jalankan Nginx
CMD ["nginx", "-g", "daemon off;"]
