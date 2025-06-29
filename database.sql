-- Membuat database
CREATE DATABASE golang_api_assignment;

-- Membuat users tabel
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Membuat products tabel
CREATE TABLE products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    harga DECIMAL(15,2) NOT NULL,
    kategori VARCHAR(255) NOT NULL,
    foto_produk VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Membuat inventories tabel
CREATE TABLE inventories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT UNSIGNED NOT NULL,
    jumlah INT DEFAULT 0,
    lokasi VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Membuat orders tabel
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    total_harga DECIMAL(15,2) DEFAULT 0,
    status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    tanggal_order TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Membuat order_items tabel
CREATE TABLE order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    jumlah INT NOT NULL,
    harga DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Menambahkan 5 dummy users
INSERT INTO users (name, email, password) VALUES
('Ahmad Rizki', 'ahmad.rizki@gmail.com', '$2a$10$hashedpassword1'),
('Siti Nurhaliza', 'siti.nurhaliza@yahoo.com', '$2a$10$hashedpassword2'),
('Budi Santoso', 'budi.santoso@hotmail.com', '$2a$10$hashedpassword3'),
('Dewi Kartika', 'dewi.kartika@gmail.com', '$2a$10$hashedpassword4'),
('Eko Prasetyo', 'eko.prasetyo@outlook.com', '$2a$10$hashedpassword5');

-- Menambahkan 5 dummy products
INSERT INTO products (nama, deskripsi, harga, kategori, foto_produk) VALUES
('Laptop ASUS ROG Strix', 'Gaming laptop dengan prosesor Intel Core i7 dan NVIDIA RTX 3060', 15000000.00, 'Elektronik', '6887-13525-83491944.jpg'),
('Smartphone Samsung Galaxy S23', 'Smartphone flagship dengan kamera 108MP dan 5G', 12000000.00, 'Elektronik', '18941-18249-2255.jpg'),
('Sepatu Nike Air Force 1', 'Sepatu casual dengan desain klasik dan nyaman untuk sehari-hari', 1500000.00, 'Fashion', '66556-1574-2363.jpg'),
('Kemeja Formal Pria', 'Kemeja formal berkualitas tinggi untuk acara resmi', 350000.00, 'Fashion', '1244-18249-1245.jpg'),
('Meja Kerja Kayu Jati', 'Meja kerja dari kayu jati solid dengan laci penyimpanan', 2500000.00, 'Furniture', '23433-18249-1612.jpg');

-- Menambahkan 5 dummy inventories
INSERT INTO inventories (product_id, jumlah, lokasi) VALUES
(1, 25, 'Gudang Jakarta Pusat'),
(2, 50, 'Gudang Surabaya'),
(3, 100, 'Gudang Bandung'),
(4, 75, 'Gudang Jakarta Pusat'),
(5, 15, 'Gudang Yogyakarta');

-- Menambahkan 5 dummy orders
INSERT INTO orders (user_id, total_harga, status, tanggal_order) VALUES
(1, 15000000.00, 'delivered', '2025-06-15 10:30:00'),
(2, 12000000.00, 'shipped', '2025-06-20 14:15:00'),
(3, 1500000.00, 'confirmed', '2025-06-22 09:45:00'),
(4, 350000.00, 'pending', '2025-06-25 16:20:00'),
(5, 2500000.00, 'delivered', '2025-06-18 11:10:00');


-- Menambahkan 5 dummy order items
INSERT INTO order_items (order_id, product_id, jumlah, harga, subtotal) VALUES
(1, 1, 1, 15000000.00, 15000000.00),
(2, 2, 1, 12000000.00, 12000000.00),
(3, 3, 1, 1500000.00, 1500000.00),
(4, 4, 1, 350000.00, 350000.00),
(5, 5, 1, 2500000.00, 2500000.00);

-- Menampilkan semua produk beserta jumlah stok di tiap lokasi
SELECT 
    p.id AS product_id,
    p.nama AS nama_produk,
    p.kategori,
    i.lokasi,
    i.jumlah AS stok
FROM products p
JOIN inventories i ON p.id = i.product_id
WHERE p.deleted_at IS NULL AND i.deleted_at IS NULL;

-- Menampilkan detail pesanan dengan pengguna dan produk
SELECT 
    o.id AS order_id,
    u.name AS nama_pelanggan,
    p.nama AS nama_produk,
    oi.jumlah,
    oi.harga,
    oi.subtotal,
    o.status,
    o.tanggal_order
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
WHERE o.deleted_at IS NULL AND oi.deleted_at IS NULL AND p.deleted_at IS NULL;

-- Total jumlah unit produk yang terjual dari order_items
SELECT 
    p.nama AS nama_produk,
    SUM(oi.jumlah) AS total_terjual
FROM order_items oi
JOIN products p ON oi.product_id = p.id
JOIN orders o ON oi.order_id = o.id
WHERE oi.deleted_at IS NULL AND o.status != 'cancelled'
GROUP BY p.id
ORDER BY total_terjual DESC;

-- Total stok produk di semua lokasi
SELECT 
    p.nama AS nama_produk,
    SUM(i.jumlah) AS total_stok
FROM products p
JOIN inventories i ON p.id = i.product_id
WHERE i.deleted_at IS NULL
GROUP BY p.id
ORDER BY total_stok DESC;

