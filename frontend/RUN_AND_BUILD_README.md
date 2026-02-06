# 🚀 Echo Chat App - Development & Build Guide

Dokumen ini menjelaskan cara menjalankan aplikasi dalam berbagai lingkungan (Development & Production) serta cara melakukan build APK.

## 🛠 Konfigurasi Backend

Aplikasi ini sudah diprogram untuk mendeteksi lingkungan secara otomatis menggunakan class `PlatformConfig`.

| Mode            | Platform            | Backend URL (Default)                                        |
| :-------------- | :------------------ | :----------------------------------------------------------- |
| **Development** | Android Emulator    | `http://10.0.2.2:8080`                                       |
| **Development** | iOS / Web / Desktop | `http://localhost:8080`                                      |
| **Production**  | APK / Release Build | `https://api.echo-chat.com` (Ubah di `platform_config.dart`) |

---

## 🏃 Menjalankan Aplikasi (Development)

### 1. Standar Run

Cukup jalankan tanpa argumen tambahan:

```bash
flutter run
```

### 2. Custom Backend URL

Jika kamu ingin menghubungkan aplikasi ke IP tertentu (misalnya saat menggunakan HP Fisik), gunakan flag `--dart-define`:

```bash
flutter run --dart-define=BASE_URL=http://192.168.1.XX:8080
```

---

## 📦 Build APK (Production)

### 1. Persiapan

Pastikan `_prodBaseUrl` di `lib/core/config/platform_config.dart` sudah sesuai dengan backend production kamu.

### 2. Command Build

Jalankan perintah berikut untuk menghasilkan APK yang optimal:

```bash
flutter build apk --release
```

Atau jika ingin menentukan URL backend secara dinamis saat build:

```bash
flutter build apk --release --dart-define=BASE_URL=https://api.domain-kamu.com
```

### 3. Lokasi File

Setelah selesai, file APK bisa ditemukan di:
`build/app/outputs/flutter-apk/app-release.apk`

---

## ⚠️ Catatan Penting Android (AndroidManifest)

Untuk rilis produksi, file `android/app/src/main/AndroidManifest.xml` **wajib** memiliki permission internet (sudah ditambahkan):

```xml
<uses-permission android:name="android.permission.INTERNET"/>
```

---

## 💡 Troubleshooting

- **Connection Refused (Android Emulator):** Pastikan backend kamu berjalan di port 8080 dan gunakan IP `10.0.2.2` (sudah dihandle otomatis oleh `PlatformConfig`).
- **Issues with HTTP (bukan HTTPS):** Android secara default memblokir trafik `http` di mode release. Sangat disarankan menggunakan `https` untuk production.
