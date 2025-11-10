# Version Bump Scripts

Script untuk auto-increment version dan create git tag secara otomatis.

## Usage

### Auto-Increment Version

```bash
# Patch version (1.1.0 -> 1.1.1)
npm run version:patch

# Minor version (1.1.0 -> 1.2.0)
npm run version:minor

# Major version (1.1.0 -> 2.0.0)
npm run version:major
```

### Build dengan Auto-Increment

```bash
# Build dengan auto-increment patch version
npm run build:patch

# Build dengan auto-increment minor version
npm run build:minor

# Build dengan auto-increment major version
npm run build:major
```

## Files yang Diupdate

Script akan otomatis update version di:

1. ✅ `src-tauri/tauri.conf.json` (source of truth)
2. ✅ `src-tauri/Cargo.toml` (harus sync)
3. ✅ `package.json` (optional, untuk consistency)

## Git Tag

Script akan otomatis create git tag dengan format:
- Tag name: `v1.1.1`
- Tag message: `Release version 1.1.1`

## Manual Usage

Jika ingin run script langsung:

```bash
# Windows
node scripts/version-bump.js patch

# Linux/Mac
node scripts/version-bump.js patch
```

## Next Steps Setelah Version Bump

1. **Review changes**:
   ```bash
   git diff
   ```

2. **Commit changes**:
   ```bash
   git add .
   git commit -m "Bump version to 1.1.1"
   ```

3. **Push dengan tags**:
   ```bash
   git push
   git push --tags
   ```

## Troubleshooting

### Error: Version mismatch

Pastikan `tauri.conf.json` dan `Cargo.toml` memiliki version yang sama. Script akan otomatis sync keduanya.

### Error: Git tag already exists

Script akan skip tag creation jika tag sudah ada. Anda bisa:
- Delete tag: `git tag -d v1.1.1`
- Atau gunakan tag yang sudah ada

