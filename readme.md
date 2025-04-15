# Order Service

## 1. Persiapan

Pastikan kamu sudah memiliki PostgreSQL yang terinstal di sistem atau menggunakan Docker untuk menjalankan PostgreSQL.

Jika menggunakan Docker, kamu bisa menjalankan PostgreSQL menggunakan perintah berikut:

```bash
docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -v /var/lib/postgresql/data -d postgres
```

Untuk melakukan koneksi ke PostgreSQL, pastikan kamu menggunakan \`host=localhost\`, \`port=5432\`, \`user=postgres\`, dan \`password=postgres\`.

## 2. Database Name / Ekstensi UUID

Untuk mendukung penggunaan tipe data UUID di PostgreSQL, pastikan ekstensi \`uuid-ossp\` sudah diaktifkan. Jalankan perintah berikut pada database PostgreSQL untuk mengaktifkan ekstensi \`uuid-ossp\`:

```sql
CREATE DATABASE geb1_order_service;
```

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

Ekstensi ini digunakan untuk menghasilkan UUID secara otomatis.

## 3. Persiapan Tabel orders dan order_details

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    total_price NUMERIC(19, 2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    qty INT NOT NULL DEFAULT 0,
    price NUMERIC(19, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);
```
