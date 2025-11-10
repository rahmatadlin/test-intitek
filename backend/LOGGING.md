# Logging Documentation

Aplikasi ini menggunakan Wails runtime logger yang menulis log ke file di folder `logs/`.

## Struktur Logging

### Folder Logs
- **Development**: `backend/logs/`
- **Production**: `backend/build/logs/` (relatif ke executable)

### Format File Log
- Nama file: `warehouse-management-YYYY-MM-DD.log`
- Format: `[YYYY-MM-DD HH:MM:SS] [LEVEL] message`

### Log Levels

Wails logger mendukung level berikut (dari terendah ke tertinggi):

1. **TRACE** - Informasi detail untuk debugging
2. **DEBUG** - Informasi debugging
3. **INFO** - Informasi umum
4. **WARNING** - Peringatan
5. **ERROR** - Error yang terjadi
6. **FATAL** - Error fatal yang menyebabkan aplikasi exit

## Penggunaan di Go

### Menggunakan File Logger

```go
import "warehouse-management/logger"

// Di dalam App struct
if a.fileLogger != nil {
    a.fileLogger.Info("Pesan informasi")
    a.fileLogger.Debug("Pesan debug")
    a.fileLogger.Warning("Pesan peringatan")
    a.fileLogger.Error("Pesan error")
}
```

### Contoh di Controller

```go
// Di auth_controller.go atau controller lainnya
// Dapatkan logger dari App (jika tersedia)
if appLogger := getAppLogger(); appLogger != nil {
    appLogger.Info("User login successful")
    appLogger.Error("Login failed: " + err.Error())
}
```

## Penggunaan di JavaScript/Frontend

### Import Runtime Logger

```javascript
import { LogInfo, LogError, LogWarning, LogDebug } from '@wailsio/runtime';

// Gunakan di frontend
LogInfo("Frontend message");
LogError("Error occurred");
LogWarning("Warning message");
LogDebug("Debug information");
```

## Konfigurasi

Log level dapat diatur di `app.go`:

```go
appOptions := &options.App{
    // ...
    Logger:   fileLogger,
    LogLevel: logger.DEBUG, // Set level: TRACE, DEBUG, INFO, WARNING, ERROR
}
```

## Catatan

- Log file dibuat otomatis saat aplikasi start
- Log ditulis ke file dan console (stdout)
- Log file di-append (tidak di-overwrite)
- Setiap hari membuat file log baru dengan timestamp

