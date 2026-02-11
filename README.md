# 💰 Go Finance Wallet API

> REST API untuk sistem e-wallet digital yang aman — dibangun dengan Go, Gin, PostgreSQL, dan JWT Authentication.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)

---

## 📋 Daftar Isi

- [Tentang Project](#-tentang-project)
- [Fitur Utama](#-fitur-utama)
- [Tech Stack](#-tech-stack)
- [Arsitektur Project](#-arsitektur-project)
- [Struktur Folder](#-struktur-folder)
- [Menjalankan Secara Lokal (VSCode)](#-menjalankan-secara-lokal-vscode)
- [Menjalankan via Docker Compose](#-menjalankan-via-docker-compose)
- [Environment Variables](#-environment-variables)
- [API Endpoints](#-api-endpoints)
  - [Register](#1--register)
  - [Login](#2--login)
  - [Cek Saldo (Get Balance)](#3--cek-saldo-get-balance)
  - [Top Up](#4--top-up)
  - [Withdraw (Tarik Saldo)](#5--withdraw-tarik-saldo)
- [Keamanan](#-keamanan)
- [Troubleshooting](#-troubleshooting)

---

## 🏦 Tentang Project

**Go Finance Wallet API** adalah backend REST API untuk sistem dompet digital (e-wallet). Aplikasi ini memungkinkan pengguna untuk:

- Mendaftar akun baru dengan PIN keamanan 6 digit
- Login dan mendapatkan token JWT untuk autentikasi
- Melihat saldo wallet secara real-time
- Melakukan top up saldo
- Melakukan withdraw (penarikan) saldo dengan verifikasi PIN

Setiap transaksi dicatat dan saldo dijamin integritasnya menggunakan **HMAC-SHA256 signature** — jika ada manipulasi langsung di database, sistem akan mendeteksinya.

---

## ✨ Fitur Utama

| Fitur | Deskripsi |
|---|---|
| 🔐 **JWT Authentication** | Token JWT (expired 24 jam) untuk mengamankan endpoint |
| 🔑 **PIN 6 Digit** | Withdraw memerlukan PIN yang di-hash dengan bcrypt |
| 🛡️ **HMAC Signature** | Setiap saldo ditandatangani HMAC-SHA256, anti-manipulasi database |
| 🔒 **Row-Level Locking** | `SELECT ... FOR UPDATE` mencegah race condition saat transaksi bersamaan |
| 📒 **Transaction Log** | Semua top up & withdraw tercatat sebagai CREDIT/DEBIT |
| 🐳 **Docker Ready** | Satu perintah `docker compose up` langsung jalan |
| 🗄️ **Auto Migration** | Tabel otomatis dibuat saat aplikasi pertama kali dijalankan |

---

## 🛠 Tech Stack

| Teknologi | Kegunaan |
|---|---|
| [Go](https://golang.org/) | Bahasa pemrograman utama |
| [Gin](https://github.com/gin-gonic/gin) | HTTP web framework |
| [GORM](https://gorm.io/) | ORM untuk PostgreSQL |
| [PostgreSQL](https://www.postgresql.org/) | Database relasional |
| [JWT (golang-jwt)](https://github.com/golang-jwt/jwt) | Token autentikasi |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Hashing password & PIN |
| [HMAC-SHA256](https://pkg.go.dev/crypto/hmac) | Signature integritas saldo |
| [Docker](https://www.docker.com/) | Containerization |

---

## 🏗 Arsitektur Project

Aplikasi ini menggunakan **Clean Architecture** pattern:

```
Request → Handler → Service → Repository → Database
                       ↓
                  pkg/crypto (hash, jwt, signature)
```

```
┌─────────────────────────────────────────────────────┐
│                    CLIENT (Postman / cURL)           │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP Request
                       ▼
┌──────────────────────────────────────────────────────┐
│                  GIN ROUTER (:5000)                  │
│  POST /api/v1/register    POST /api/v1/login         │
│  GET  /api/v1/balance     POST /api/v1/topup         │
│  POST /api/v1/withdraw                               │
└──────────────────────┬───────────────────────────────┘
                       │
            ┌──────────▼──────────┐
            │  AUTH MIDDLEWARE    │  ← JWT Validation
            │  (protected routes)│
            └──────────┬─────────┘
                       │
         ┌─────────────▼─────────────┐
         │        HANDLER LAYER      │
         │  auth_handler.go          │
         │  wallet_handler.go        │
         └─────────────┬─────────────┘
                       │
         ┌─────────────▼─────────────┐
         │       SERVICE LAYER       │
         │  auth_service.go          │
         │  wallet_service.go        │
         └─────────────┬─────────────┘
                       │
         ┌─────────────▼─────────────┐
         │     REPOSITORY LAYER      │
         │  user_repo.go             │
         │  wallet_repo.go           │
         │  transaction_repo.go      │
         └─────────────┬─────────────┘
                       │
              ┌────────▼────────┐
              │   PostgreSQL    │
              │   (wallet_db)   │
              └─────────────────┘
```

---

## 📂 Struktur Folder

```
go-finance-wallet/
├── main.go                          # Entry point aplikasi
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── .env                             # Environment variables (jangan commit!)
├── .gitignore                       # Ignore .env
├── dockerfile                       # Docker image build
├── docker-compose.yaml              # Multi-container setup
│
├── internal/                        # Private application code
│   ├── handler/
│   │   ├── auth_handler.go          # Handler Register & Login
│   │   └── wallet_handler.go        # Handler Balance, TopUp, Withdraw
│   ├── middleware/
│   │   └── auth_middleware.go       # JWT token validation middleware
│   ├── model/
│   │   └── entity.go               # Struct: User, Wallet, Transaction
│   ├── repository/
│   │   ├── user_repo.go            # CRUD User + create wallet
│   │   ├── wallet_repo.go          # Get & update wallet (with lock)
│   │   └── transaction_repo.go     # Create & list transactions
│   └── service/
│       ├── auth_service.go          # Logika register & login
│       └── wallet_service.go        # Logika balance, topup, withdraw
│
└── pkg/                             # Reusable packages
    ├── crypto/
    │   ├── hash.go                  # bcrypt hash & verify
    │   ├── jwt.go                   # Generate & validate JWT token
    │   └── signature.go             # HMAC-SHA256 signature untuk saldo
    └── database/
        └── postgres.go              # Koneksi PostgreSQL via GORM
```

---

## 🚀 Menjalankan Secara Lokal (VSCode)

### Prerequisites

Pastikan sudah terinstall di komputer kamu:

- **[Go](https://golang.org/dl/)** (versi 1.25 atau lebih baru)
- **[PostgreSQL](https://www.postgresql.org/download/)** (versi 15 direkomendasikan)
- **[Git](https://git-scm.com/downloads)**
- **[VSCode](https://code.visualstudio.com/)** + Extension: [Go for VSCode](https://marketplace.visualstudio.com/items?itemName=golang.Go)

### Langkah-langkah

**1. Clone repository**

```bash
git clone https://github.com/hanif411/mampu-wallet
cd mampu-wallet
```

**2. Buat database PostgreSQL**

Buka terminal / pgAdmin, lalu buat database:

```sql
CREATE DATABASE wallet_db;
```

**3. Buat file `.env`**

Buat file `.env` di root project (sejajar dengan `main.go`):

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password_kamu
DB_NAME=wallet_db
DB_PORT=5432
DB_SSLMODE=disable

SECRET_KEY=ganti-dengan-secret-key-minimal-32-karakter
```

> ⚠️ **Penting:** Ganti `DB_PASSWORD` dengan password PostgreSQL kamu dan `SECRET_KEY` dengan string acak yang panjang.

**4. Install dependencies**

```bash
go mod download
```

**5. Jalankan aplikasi**

```bash
go run main.go
```

Jika berhasil, kamu akan melihat output:

```
Database Connected
migration succes
Server running on :5000
```

**6. Test API**

Buka Postman atau gunakan `curl` untuk mencoba endpoint (lihat bagian [API Endpoints](#-api-endpoints)).

---

## 🐳 Menjalankan via Docker Compose

Cara paling mudah — **tidak perlu install Go atau PostgreSQL** di komputer kamu.

### Prerequisites

- **[Docker Desktop](https://www.docker.com/products/docker-desktop/)** (sudah termasuk Docker Compose)

### Langkah-langkah

**1. Clone repository**

```bash
git clone https://github.com/hanif411/mampu-wallet
cd mampu-wallet
```

**2. Jalankan dengan Docker Compose**

```bash
docker compose up --build
```

> Perintah ini akan:
> - Pull image PostgreSQL 15 Alpine
> - Build aplikasi Go dari Dockerfile
> - Membuat database `wallet_db` secara otomatis
> - Menjalankan API di port `5000`

**3. Tunggu sampai muncul**

```
wallet_api_container  | Database Connected
wallet_api_container  | migration succes
wallet_api_container  | Server running on :5000
```

**4. Test API**

API sudah siap di `http://localhost:5000`. Gunakan Postman atau `curl`.

**5. Menghentikan**

```bash
# Ctrl+C di terminal, lalu:
docker compose down

# Untuk menghapus data database juga:
docker compose down -v
```

---

## ⚙️ Environment Variables

| Variable | Deskripsi | Default |
|---|---|---|
| `DB_HOST` | Hostname database PostgreSQL | `localhost` |
| `DB_USER` | Username database | `postgres` |
| `DB_PASSWORD` | Password database | — |
| `DB_NAME` | Nama database | `wallet_db` |
| `DB_PORT` | Port database | `5432` |
| `DB_SSLMODE` | SSL mode koneksi | `disable` |
| `SECRET_KEY` | Secret key untuk JWT & HMAC signature | — |

> 💡 Pada Docker Compose, environment variables sudah di-set di `docker-compose.yaml`. Untuk lokal, gunakan file `.env`.

---

## 📡 API Endpoints

**Base URL:** `http://localhost:5000/api/v1`

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| `POST` | `/register` | ❌ | Daftar akun baru |
| `POST` | `/login` | ❌ | Login & dapatkan JWT token |
| `GET` | `/balance` | ✅ Bearer Token | Cek saldo wallet |
| `POST` | `/topup` | ✅ Bearer Token | Top up saldo |
| `POST` | `/withdraw` | ✅ Bearer Token | Tarik saldo (perlu PIN) |

---

### 1. 📝 Register

Daftar akun baru. Setiap akun akan otomatis dibuatkan wallet dengan saldo 0.

**Request:**

```bash
curl -X POST http://localhost:5000/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "hanif",
    "password": "password123",
    "pin": "123456"
  }'
```

| Field | Type | Rules | Keterangan |
|---|---|---|---|
| `username` | string | required, unique | Username untuk login |
| `password` | string | required | Password akun |
| `pin` | string | required, numeric, 6 digit | PIN untuk withdraw |

**Response Sukses (201):**

```json
{
  "message": "user berhasil didaftarkan"
}
```

**Response Error (400):**

```json
{
  "error": "Key: 'Pin' Error:Field validation for 'Pin' failed on the 'len' tag"
}
```

---

### 2. 🔑 Login

Login untuk mendapatkan JWT token. Token berlaku selama **24 jam**.

**Request:**

```bash
curl -X POST http://localhost:5000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "hanif",
    "password": "password123"
  }'
```

**Response Sukses (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Error (401):**

```json
{
  "error": "password salah"
}
```

> 💡 **Simpan token ini!** Kamu akan membutuhkannya untuk semua endpoint yang memerlukan autentikasi.

---

### 3. 💵 Cek Saldo (Get Balance)

Lihat saldo wallet kamu. Sistem juga memverifikasi integritas saldo via HMAC signature.

**Request:**

```bash
curl -X GET http://localhost:5000/api/v1/balance \
  -H "Authorization: Bearer <TOKEN_KAMU>"
```

**Response Sukses (200):**

```json
{
  "balance": 500000
}
```

**Response Error (500) — Jika saldo dimanipulasi:**

```json
{
  "error": "data saldo tidak valid (manipulasi terdeteksi!)"
}
```

---

### 4. 💳 Top Up

Tambah saldo ke wallet kamu.

**Request:**

```bash
curl -X POST http://localhost:5000/api/v1/topup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_KAMU>" \
  -d '{
    "amount": 500000
  }'
```

| Field | Type | Rules | Keterangan |
|---|---|---|---|
| `amount` | integer | required, > 0 | Jumlah top up (dalam satuan terkecil) |

**Response Sukses (200):**

```json
{
  "message": "top up berhasil"
}
```

**Response Error (400):**

```json
{
  "error": "jumlah harus lebih dari 0"
}
```

---

### 5. 💸 Withdraw (Tarik Saldo)

Tarik saldo dari wallet. Memerlukan **PIN 6 digit** yang dibuat saat register.

**Request:**

```bash
curl -X POST http://localhost:5000/api/v1/withdraw \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_KAMU>" \
  -d '{
    "amount": 100000,
    "pin": "123456"
  }'
```

| Field | Type | Rules | Keterangan |
|---|---|---|---|
| `amount` | integer | required, > 0 | Jumlah withdraw |
| `pin` | string | required | PIN 6 digit kamu |

**Response Sukses (200):**

```json
{
  "message": "withdraw berhasil"
}
```

**Response Error (400) — PIN salah:**

```json
{
  "error": "PIN salah"
}
```

**Response Error (400) — Saldo tidak cukup:**

```json
{
  "error": "saldo tidak mencukupi"
}
```

---

## 🛡 Keamanan

Fitur keamanan yang diimplementasikan:

### 1. Password Hashing (bcrypt)
Password dan PIN **tidak pernah disimpan dalam bentuk plain text**. Semua di-hash menggunakan bcrypt dengan cost factor default (10).

### 2. JWT Token Authentication
Endpoint yang memerlukan autentikasi dilindungi JWT token. Token dikirim via header `Authorization: Bearer <token>` dan expired otomatis setelah **24 jam**.

### 3. HMAC-SHA256 Balance Signature
Setiap saldo wallet ditandatangani menggunakan HMAC-SHA256. Saat membaca saldo, sistem memverifikasi signature — jika ada yang mengubah saldo langsung di database (bypass API), sistem akan menolak dengan error **"manipulasi terdeteksi"**.

### 4. Row-Level Locking
Operasi top up dan withdraw menggunakan `SELECT ... FOR UPDATE` di PostgreSQL untuk mencegah **race condition** saat ada transaksi bersamaan pada wallet yang sama.

### 5. PIN Verification
Withdraw memerlukan PIN 6 digit yang telah di-hash. Ini memberikan lapisan keamanan tambahan di atas autentikasi JWT.

---

## 🔥 Troubleshooting

### ❌ `error load .env`
- Pastikan file `.env` ada di root project (sejajar dengan `main.go`)
- Jika running via Docker Compose, environment sudah di-set di `docker-compose.yaml`

### ❌ `gagal koneksi ke database`
- Pastikan PostgreSQL sudah berjalan
- Cek kembali `DB_HOST`, `DB_USER`, `DB_PASSWORD` di `.env`
- Jika via Docker: tunggu beberapa detik, PostgreSQL butuh waktu untuk startup

### ❌ `butuh token`
- Endpoint `/balance`, `/topup`, `/withdraw` memerlukan JWT token
- Login dulu untuk mendapatkan token, lalu tambahkan header: `Authorization: Bearer <token>`

### ❌ `token expired atau salah`
- Token JWT expired setelah 24 jam — login kembali
- Pastikan format header benar: `Bearer <spasi> <token>`

### ❌ `data saldo tidak valid (manipulasi terdeteksi!)`
- Saldo di database telah diubah secara langsung (bypass API)
- Ini fitur keamanan — saldo hanya boleh berubah melalui API

---

## 📄 Database Schema

Aplikasi menggunakan 3 tabel yang otomatis dibuat saat pertama kali dijalankan:

```
┌──────────────┐       ┌──────────────────┐       ┌──────────────────┐
│    users     │       │     wallets      │       │  transactions    │
├──────────────┤       ├──────────────────┤       ├──────────────────┤
│ id (PK)      │──┐    │ id (PK)          │──┐    │ id (PK)          │
│ username     │  └───▶│ user_id (FK)     │  └───▶│ wallet_id (FK)   │
│ password     │       │ balance          │       │ amount           │
└──────────────┘       │ pin              │       │ type (CREDIT/    │
                       │ signature        │       │       DEBIT)     │
                       └──────────────────┘       └──────────────────┘
```

---

## 🧪 Contoh Alur Penggunaan Lengkap

Berikut contoh alur dari awal sampai akhir menggunakan `curl`:

```bash
# 1. Register akun baru
curl -X POST http://localhost:5000/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "hanif", "password": "rahasia123", "pin": "123456"}'

# 2. Login untuk mendapatkan token
curl -X POST http://localhost:5000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "hanif", "password": "rahasia123"}'
# Simpan token dari response!

# 3. Cek saldo awal (harusnya 0)
curl -X GET http://localhost:5000/api/v1/balance \
  -H "Authorization: Bearer eyJhbGci..."

# 4. Top up Rp 1.000.000
curl -X POST http://localhost:5000/api/v1/topup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGci..." \
  -d '{"amount": 1000000}'

# 5. Cek saldo setelah top up (harusnya 1000000)
curl -X GET http://localhost:5000/api/v1/balance \
  -H "Authorization: Bearer eyJhbGci..."

# 6. Withdraw Rp 250.000
curl -X POST http://localhost:5000/api/v1/withdraw \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGci..." \
  -d '{"amount": 250000, "pin": "123456"}'

# 7. Cek saldo akhir (harusnya 750000)
curl -X GET http://localhost:5000/api/v1/balance \
  -H "Authorization: Bearer eyJhbGci..."
```

---