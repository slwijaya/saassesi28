document.addEventListener("DOMContentLoaded", function () {
    // Fungsi untuk menangani register
    document.getElementById("registerForm").addEventListener("submit", async function (event) {
        event.preventDefault(); // Mencegah form submit default

        // Ambil data dari form register
        const fullName = document.getElementById("registerName").value;
        const email = document.getElementById("registerEmail").value;
        const password = document.getElementById("registerPassword").value;
        const confirmPassword = document.getElementById("confirmPassword").value;

        // Validasi password
        if (password !== confirmPassword) {
            alert("Password dan konfirmasi password tidak cocok.");
            return;
        }

        // Kirim data ke API Register
        try {
            const response = await fetch("http://localhost:8080/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ full_name: fullName, email: email, password: password })
            });

            const result = await response.json();

            if (response.ok) {
                alert("Registrasi berhasil! Silakan login.");
                document.getElementById("registerForm").reset();
                new bootstrap.Modal(document.getElementById('registerModal')).hide();
            } else {
                alert(result.error || "Registrasi gagal.");
            }
        } catch (error) {
            console.error("Error:", error);
            alert("Terjadi kesalahan pada server.");
        }
    });

    // Fungsi untuk menangani login
    document.getElementById("loginForm").addEventListener("submit", async function (event) {
        event.preventDefault(); // Mencegah form submit default

        // Ambil data dari form login
        const email = document.getElementById("loginEmail").value;
        const password = document.getElementById("loginPassword").value;

        // Kirim data ke API Login
        try {
            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email: email, password: password })
            });

            const result = await response.json();

            if (response.ok) {
                alert("Login berhasil!");
                localStorage.setItem("token", result.token); // Simpan token di Local Storage
                document.getElementById("loginForm").reset();
                new bootstrap.Modal(document.getElementById('loginModal')).hide();
            } else {
                alert(result.error || "Login gagal.");
            }
        } catch (error) {
            console.error("Error:", error);
            alert("Terjadi kesalahan pada server.");
        }
    });
});

let allProducts = []; // Menyimpan semua produk

// Fetch produk dari API saat halaman dimuat
function fetchProducts() {
    fetch("http://localhost:8080/search-products")
        .then(response => response.json())
        .then(data => {
            allProducts = data; // Simpan semua produk ke dalam array
            displayProducts(allProducts); // Tampilkan semua produk
        })
        .catch(error => console.error("Error fetching product data:", error));
}

// Fungsi untuk menampilkan produk dalam HTML
function displayProducts(products) {
    const productCards = document.getElementById("productCards");
    productCards.innerHTML = "";

    products.forEach(product => {
        const card = `
            <div class="col-md-3 product-card">
                <div class="card h-100">
                    <img src="${product.image}" class="card-img-top" alt="${product.title}">
                    <div class="card-body">
                        <h6 class="card-title">${product.title}</h6>
                        <p class="card-text"><strong>Rp ${product.price.toLocaleString()}</strong></p>
                        <small class="text-muted">⭐ ${product.rating.rate} | ${product.rating.count} reviews</small>
                    </div>
                </div>
            </div>
        `;
        productCards.innerHTML += card;
    });
}

// Fungsi untuk mencari produk berdasarkan input pengguna
function searchProduct() {
    const query = document.getElementById("searchInput").value.toLowerCase();
    
    if (query === "") {
        displayProducts(allProducts); // Jika kosong, tampilkan semua produk
        return;
    }

    const filteredProducts = allProducts.filter(product =>
        product.title.toLowerCase().includes(query)
    );

    displayProducts(filteredProducts);
}

// Fetch produk saat halaman dimuat
fetchProducts();

// document.getElementById("checkoutForm").addEventListener("submit", function (e) {
//     e.preventDefault();

//     let productElement = document.getElementById("product");
//     let selectedProduct = productElement.options[productElement.selectedIndex];
//     let amount = selectedProduct.dataset.price;
//     let email = document.getElementById("email").value;

//     // Kirim permintaan ke API backend untuk membuat invoice
//     fetch("http://localhost:8080/create-invoice", {
//         method: "POST",
//         headers: { "Content-Type": "application/json" },
//         body: JSON.stringify({
//             external_id: "order_" + new Date().getTime(),
//             amount: parseInt(amount),
//             payer_email: email,
//             description: "Pembayaran Produk"
//         })
//     })
//     .then(response => response.json())
//     .then(data => {
//         // Redirect ke halaman pembayaran Xendit
//         window.location.href = data.invoice_url;
//     })
//     .catch(error => console.error("Error:", error));
// });

document.getElementById("checkoutForm").addEventListener("submit", function (e) {
    e.preventDefault(); // Mencegah halaman reload

    // Mengambil data yang dimasukkan oleh pengguna
    let productElement = document.getElementById("product");
    let selectedProduct = productElement.options[productElement.selectedIndex];
    let amount = selectedProduct.dataset.price; // Ambil harga produk dari atribut data-price
    let email = document.getElementById("email").value;
    let name = document.getElementById("name").value;

    // Memastikan pengguna memilih produk dan metode pembayaran
    if (!amount || !email || !name) {
        alert("Silakan isi semua bidang sebelum checkout.");
        return;
    }

    // Membuat payload request untuk dikirim ke backend
    let requestBody = {
        external_id: "order_" + new Date().getTime(),
        amount: parseInt(amount),
        payer_email: email,
        description: "Pembayaran Produk",
    };

    // Kirim request ke backend
    fetchCheckout(requestBody);
});



function fetchCheckout(data) {
    fetch("http://localhost:8080/create-invoice", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        if (data.invoice_url) {
            // Redirect ke halaman pembayaran Xendit
            window.location.href = data.invoice_url;
        } else {
            alert("Gagal membuat invoice. Coba lagi.");
        }
    })
    .catch(error => {
        console.error("Error:", error);
        alert("Terjadi kesalahan saat membuat checkout.");
    });
}


function checkPaymentStatus(invoiceId) {
    fetch(`http://localhost:8080/transaction-status?invoice_id=${invoiceId}`)
    .then(response => response.json())
    .then(data => {
        if (data.status === "PAID") {
            document.getElementById("payment-status").innerHTML = "<div class='alert alert-success'>Pembayaran Berhasil!</div>";
        } else {
            document.getElementById("payment-status").innerHTML = "<div class='alert alert-warning'>Menunggu Pembayaran...</div>";
        }
    })
    .catch(error => {
        console.error("Error:", error);
        document.getElementById("payment-status").innerHTML = "<div class='alert alert-danger'>Gagal mengambil status pembayaran.</div>";
    });
}


