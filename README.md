# ⚡ Bimskrip API

API untuk menunjang tujuan dan kebutuhan **bimskrip-app** dengan menggunakan **Go Fiber**.  
Dibuat untuk mengelola proses bimbingan skripsi antara **mahasiswa** dan **dosen pembimbing**.

---

## ✨ Fitur API
- **Autentikasi**
  - Register & Login (Mahasiswa & Dosen)

- **Mahasiswa**
  - CRUD progress penelitian
  - Upload hasil penelitian
  - Melihat status bimbingan 

- **Dosen**
  - Melihat daftar mahasiswa bimbingan
  - Mengubah status progress mahasiswa
  - Membuat dan mengatur jadwal bimbingan

- **Umum**
  - Manajemen jadwal
  - Upload dokumen penelitian (`/storage/upload/photo`)

---

## 🛠️ Tech Stack
- **Framework:** Go Fiber
- **Database:** MySQL
- **ORM:** GORM
- **Storage:** Local file storage (`/storage/upload/photo`)

---

- link repo aplikasi : https://github.com/azkaainurridho514/bimskrip-app
