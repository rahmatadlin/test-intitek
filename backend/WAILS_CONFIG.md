# Konfigurasi Wails

Dokumentasi ini menjelaskan struktur konfigurasi Wails dan file `wails.json`.

## File wails.json

File `wails.json` adalah file konfigurasi utama untuk aplikasi Wails. File ini berada di folder `backend/` dan digunakan oleh Wails CLI untuk build dan development.

### Struktur wails.json

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "warehouse-management",
  "outputfilename": "warehouse-management",
  "frontend": {
    "dir": "../frontend/dist",
    "install": "",
    "build": "",
    "bridge": "",
    "serve": ""
  },
  "author": {
    "name": "Warehouse Management",
    "email": ""
  },
  "info": {
    "productName": "Warehouse Management System",
    "productVersion": "1.0.0",
    "copyright": "",
    "comments": "Warehouse Management Desktop Application"
  }
}
```

### Penjelasan Field

#### 1. `$schema`
- **Tipe**: String
- **Deskripsi**: Schema JSON untuk validasi dan autocomplete di editor
- **Nilai**: URL ke schema Wails v2

#### 2. `name`
- **Tipe**: String
- **Deskripsi**: Nama aplikasi (untuk internal Wails)
- **Contoh**: `"warehouse-management"`

#### 3. `outputfilename`
- **Tipe**: String
- **Deskripsi**: Nama file executable yang akan dihasilkan saat build
- **Contoh**: `"warehouse-management"` → menghasilkan `warehouse-management.exe` (Windows)

#### 4. `frontend`
- **Tipe**: Object
- **Deskripsi**: Konfigurasi untuk frontend
- **Field**:
  - `dir`: Path ke folder frontend dist (relatif dari `wails.json`)
  - `install`: Command untuk install dependencies (contoh: `"npm install"`)
  - `build`: Command untuk build frontend (contoh: `"npm run build"`)
  - `bridge`: Path ke bridge directory (untuk Wails bindings)
  - `serve`: Command untuk serve frontend di development

#### 5. `author`
- **Tipe**: Object
- **Deskripsi**: Informasi author aplikasi
- **Field**:
  - `name`: Nama author
  - `email`: Email author

#### 6. `info`
- **Tipe**: Object
- **Deskripsi**: Informasi aplikasi (digunakan untuk metadata executable)
- **Field**:
  - `productName`: Nama produk aplikasi
  - `productVersion`: Versi aplikasi
  - `copyright`: Copyright notice
  - `comments`: Komentar tambahan

## Konfigurasi Tambahan (di app.go)

Selain `wails.json`, konfigurasi aplikasi juga dilakukan di `app.go` melalui `options.App`:

```go
appOptions := &options.App{
    Title:            "Warehouse Management System",  // Judul window
    Width:            1280,                           // Lebar window (px)
    Height:           800,                            // Tinggi window (px)
    AssetServer:      assetServerOptions,             // Konfigurasi asset server
    BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1}, // Warna background
    OnStartup:        app.OnStartup,                  // Callback saat startup
    OnDomReady:       app.OnDomReady,                 // Callback saat DOM ready
    OnBeforeClose:    app.OnBeforeClose,              // Callback sebelum close
    OnShutdown:       app.OnShutdown,                 // Callback saat shutdown
    Logger:           fileLogger,                     // Custom logger
    LogLevel:         2,                              // Log level (DEBUG = 2)
}
```

### Options yang Tersedia

- **Title**: Judul window aplikasi
- **Width/Height**: Ukuran window
- **MinWidth/MinHeight**: Ukuran minimum window
- **MaxWidth/MaxHeight**: Ukuran maksimum window
- **Resizable**: Apakah window bisa di-resize
- **Frameless**: Window tanpa frame (borderless)
- **StartHidden**: Start aplikasi dalam keadaan hidden
- **AlwaysOnTop**: Window selalu di atas
- **Fullscreen**: Start dalam mode fullscreen
- **AssetServer**: Konfigurasi untuk serve frontend assets
- **BackgroundColour**: Warna background window
- **Logger**: Custom logger instance
- **LogLevel**: Level logging (TRACE=1, DEBUG=2, INFO=3, WARNING=4, ERROR=5)
- **OnStartup**: Callback saat aplikasi start
- **OnDomReady**: Callback saat frontend DOM ready
- **OnBeforeClose**: Callback sebelum aplikasi close
- **OnShutdown**: Callback saat aplikasi shutdown

## Lokasi File

- **Development**: `backend/wails.json`
- **Production**: File ini tetap di `backend/` dan digunakan saat build

## Catatan

- File `wails.json` harus ada di folder yang sama dengan `main.go` atau di parent directory
- Path di `frontend.dir` relatif terhadap lokasi `wails.json`
- Untuk production build, pastikan frontend sudah di-build dan di-copy ke `frontend-dist/`

