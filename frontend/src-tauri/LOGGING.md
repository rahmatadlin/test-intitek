# 📝 Logging Configuration untuk Tauri App

Dokumen ini menjelaskan konfigurasi logging di aplikasi Tauri untuk troubleshooting.

## 📍 Lokasi Log Files

Log files disimpan di folder project untuk memudahkan troubleshooting:

### Lokasi
```
frontend/src-tauri/logs/
```

File log akan bernama: `warehouse-management.log` atau `warehouse-management-{number}.log` (jika ada rotation)

## 📋 Konfigurasi Logging

### Log Targets
Logging dikonfigurasi untuk menulis ke 3 target:

1. **Stdout** - Terminal/console output
2. **Webview** - Browser console di DevTools
3. **Folder** - File di folder `logs/` di project untuk troubleshooting

### Log Level
- **Level**: `Info` (Info, Warn, Error)
- Debug dan Trace logs tidak disimpan (untuk mengurangi ukuran file)

### Log Rotation
- **Strategy**: `KeepAll` - Semua log file dipertahankan
- **Max File Size**: 10MB per file
- File baru dibuat otomatis saat mencapai limit

### Timezone
- Menggunakan **Local Timezone** (bukan UTC)

## 🔧 Penggunaan di Code

### Import Logger
```javascript
import logger from '@/utils/logger';
```

### Basic Logging
```javascript
// Info log
logger.info('User logged in successfully');

// Warning log
logger.warn('Low stock detected', { productId: 123 });

// Error log
logger.error('Failed to load products', error);

// Debug log (hanya di console, tidak ke file)
logger.debug('Debug information', { data });
```

### API Logging
```javascript
// Otomatis di-log oleh API interceptor
// Tidak perlu manual logging untuk API calls
```

### Manual API Logging
```javascript
import logger from '@/utils/logger';

logger.apiRequest('GET', '/api/products');
logger.apiResponse('GET', '/api/products', 200, data);
logger.apiError('GET', '/api/products', error);
```

## 📊 Format Log

Format log di file:
```
[YYYY-MM-DD HH:MM:SS][TARGET][LEVEL] MESSAGE
```

Contoh:
```
[2025-11-10 10:30:45][warehouse-management][INFO] User logged in successfully
[2025-11-10 10:30:46][warehouse-management][ERROR] API Error: GET /api/products - Error: Network error
```

## 🐛 Troubleshooting

### Cara Membaca Log Files

1. **Buka folder logs** di `frontend/src-tauri/logs/`
2. **Cari file** dengan nama `warehouse-management.log` atau `warehouse-management-*.log`
3. **Buka dengan text editor** (Notepad, VS Code, dll)
4. **Cari berdasarkan timestamp** atau keyword error

### Tips untuk Troubleshooting

1. **Cari Error Level**: Cari `[ERROR]` untuk menemukan semua error
2. **Cari API Calls**: Cari `API Request` atau `API Response` untuk tracking API calls
3. **Cari Timestamp**: Gunakan timestamp untuk melacak urutan kejadian
4. **Check Log Rotation**: Jika file terlalu besar, cek file dengan nomor urut

### Contoh Log Analysis

```bash
# Windows PowerShell - Cari semua error
Select-String -Path "frontend\src-tauri\logs\*.log" -Pattern "\[ERROR\]"

# Linux/Mac - Cari semua error
grep "\[ERROR\]" frontend/src-tauri/logs/*.log
```

## 🔒 Security Notes

- Log files mungkin berisi sensitive data (tokens, user info)
- Jangan share log files tanpa review terlebih dahulu
- Log files di folder user local, tidak di-share otomatis

## 📝 Best Practices

1. **Gunakan log level yang sesuai**:
   - `info`: Informasi umum
   - `warn`: Warning yang perlu diperhatikan
   - `error`: Error yang perlu diinvestigasi

2. **Jangan log sensitive data**:
   - Jangan log password, tokens, atau data pribadi
   - Gunakan masking jika perlu

3. **Gunakan context**:
   - Sertakan informasi yang relevan (user ID, action, dll)
   - Gunakan structured logging jika memungkinkan

4. **Monitor log size**:
   - Log rotation otomatis, tapi monitor ukuran folder
   - Hapus log lama jika tidak diperlukan

## 🔧 Configuration Files

- **Rust Config**: `src-tauri/src/main.rs`
- **Permissions**: `src-tauri/capabilities/default.json`
- **Logger Utility**: `src/utils/logger.js`

## 📚 References

- [Tauri Logging Plugin Documentation](https://v2.tauri.app/plugin/logging/)
- [Log Crate Documentation](https://docs.rs/log/)

