# Tauri 2.0 Setup Guide

## Prerequisites

Sebelum menjalankan aplikasi Tauri, pastikan Anda sudah menginstall:

1. **Rust** - Diperlukan untuk build Tauri
   - Download dari: https://www.rust-lang.org/tools/install
   - Setelah install, jalankan: `rustup update`

2. **Node.js** (sudah terinstall)
   - Versi 18.x atau lebih tinggi

3. **Backend Server** - Pastikan backend Go berjalan di `http://localhost:8080`

## Setup Icons

Untuk menggunakan custom icons, buat file icons di folder `src-tauri/icons/`:

- `32x32.png` - 32x32 pixels
- `128x128.png` - 128x128 pixels  
- `128x128@2x.png` - 256x256 pixels (untuk retina display)
- `icon.icns` - Icon untuk macOS
- `icon.ico` - Icon untuk Windows

Jika tidak ada icons, Tauri akan menggunakan default icon.

## Development

Untuk menjalankan aplikasi dalam mode development:

```bash
cd frontend
npm run tauri:dev
```

Pastikan backend server sudah berjalan di `http://localhost:8080`.

## Build untuk Production

Untuk membuat build aplikasi desktop:

```bash
cd frontend
npm run tauri:build
```

File hasil build akan berada di `frontend/src-tauri/target/release/`.

## Troubleshooting

### Error: "rustc not found"
- Install Rust: https://www.rust-lang.org/tools/install
- Setelah install, restart terminal

### Error: "Failed to connect to backend"
- Pastikan backend Go berjalan di `http://localhost:8080`
- Atau ubah `VITE_API_BASE_URL` di `.env` file

### Error: "Port 5173 already in use"
- Tutup aplikasi lain yang menggunakan port 5173
- Atau ubah port di `vite.config.js`

