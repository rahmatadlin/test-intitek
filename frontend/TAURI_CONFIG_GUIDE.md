# Tauri Configuration Guide

## Struktur Konfigurasi Tauri

Tauri menggunakan **dua file konfigurasi utama** yang memiliki fungsi berbeda:

### 1. `tauri.conf.json` - Konfigurasi Tauri Framework

**Lokasi**: `frontend/src-tauri/tauri.conf.json`

**Fungsi**:
- Konfigurasi utama untuk Tauri application
- Mengatur build process, bundling, dan app metadata
- **Version di sini adalah source of truth** untuk aplikasi
- Mengatur window properties, security, plugins, dll

**Contoh**:
```json
{
  "version": "1.0.0",           // ← Version aplikasi (source of truth)
  "productName": "My App",
  "identifier": "com.myapp",
  "build": {
    "beforeBuildCommand": "npm run build"
  },
  "app": {
    "windows": [...]
  },
  "bundle": {
    "targets": "all",
    "icon": [...]
  }
}
```

### 2. `Cargo.toml` - Konfigurasi Rust Package

**Lokasi**: `frontend/src-tauri/Cargo.toml`

**Fungsi**:
- Konfigurasi Rust package (Cargo)
- Mengatur dependencies Rust, build features
- **Version di sini harus sama dengan tauri.conf.json**
- Mengatur Rust edition, authors, license, dll

**Contoh**:
```toml
[package]
name = "my-app"
version = "1.0.0"              # ← Harus sama dengan tauri.conf.json
description = "A Tauri App"
edition = "2021"

[dependencies]
tauri = { version = "2.0", features = [] }
```

## Perbedaan Utama

| Aspek | `tauri.conf.json` | `Cargo.toml` |
|-------|-------------------|--------------|
| **Bahasa** | JSON | TOML (Rust) |
| **Fungsi** | Konfigurasi Tauri app | Konfigurasi Rust package |
| **Version** | Source of truth | Harus sync dengan tauri.conf.json |
| **Digunakan oleh** | Tauri CLI | Cargo (Rust build tool) |
| **Isi** | App metadata, windows, bundle | Rust dependencies, build config |

## Version Management

### Manual Update

Jika ingin update version secara manual, **update kedua file**:

1. **Update `tauri.conf.json`**:
```json
{
  "version": "1.1.1"  // ← Update di sini
}
```

2. **Update `Cargo.toml`**:
```toml
[package]
version = "1.1.1"  // ← Update juga di sini (harus sama)
```

### Auto-Increment dengan Script

Gunakan script yang sudah dibuat untuk auto-increment:

```bash
# Patch version (1.1.0 -> 1.1.1)
npm run version:patch

# Minor version (1.1.0 -> 1.2.0)
npm run version:minor

# Major version (1.1.0 -> 2.0.0)
npm run version:major
```

Script akan otomatis:
- ✅ Update `tauri.conf.json`
- ✅ Update `Cargo.toml`
- ✅ Update `package.json` (optional)
- ✅ Create git tag

### Build dengan Auto-Increment

```bash
# Build dengan auto-increment patch version
npm run build:patch

# Build dengan auto-increment minor version
npm run build:minor

# Build dengan auto-increment major version
npm run build:major
```

## Best Practices

1. **Selalu sync version** antara `tauri.conf.json` dan `Cargo.toml`
2. **Gunakan script** untuk menghindari human error
3. **Commit version changes** sebelum build
4. **Push git tags** untuk release tracking

## Troubleshooting

### Version Mismatch Error

Jika ada error version mismatch:
1. Pastikan `tauri.conf.json` dan `Cargo.toml` memiliki version yang sama
2. Run `npm run version:patch` untuk sync otomatis

### Git Tag Already Exists

Jika tag sudah ada:
- Script akan skip tag creation
- Buat tag manual dengan: `git tag -a v1.1.1 -m "Release version 1.1.1"`

