# Rencana Belajar Fullstack ala FDE (Forward Deployed Engineer) — v8

> v8 (8 Nov 2026): **onboarding carousel 3 halaman** (A0), **password bisa ditampilkan** (ikon mata) di semua field, autocomplete **tanpa** "tambah baru", logo baru **b Palu**, komponen Typography / Password Field / Page Indicator, layar loading disesuaikan; Figma dianggap selesai (Bagian 4 & 7). v7: **foto bukti wajib + persetujuan pelanggan** saat menyelesaikan tugas, **nonaktifkan mitra dengan aturan tugas aktif**, **riwayat mitra & pelanggan**, semua tambah/edit memakai **side sheet** (lihat Bagian 4). v6: brand **Benah** (Benah Admin & Benah Mitra), design system custom krem-cokelat via plugin Figma lokal, maskot **Mang Benah**, fitur baru **Call Center** (Bagian 4), total ±10 sesi. v5: tambah **observability (Elasticsearch + Kibana + Filebeat)** di Bagian 6. v4: **fitur MVP sudah ditetapkan** (lihat Bagian 4), Postgres via **OrbStack (Docker)**, sesi Minggu **5–6 jam**. Fokus kamu ada di cara coding; scope sudah dikunci supaya muat dalam ±9 sesi.

---

## 0. Keputusan yang Sudah Disepakati

| Topik | Keputusan |
|---|---|
| Project | **A. Field Service / Work Order**: Web untuk **admin**, Android untuk **teknisi** (client) |
| Database runtime | **PostgreSQL via OrbStack (Docker Compose)** |
| Observability | **Elasticsearch + Kibana + Filebeat** (log terstruktur, lokal) |
| Sesi | **5–6 jam** per Minggu |
| Backend | **Golang** + PostgreSQL (level: baru belajar) |
| Web | **SPA sederhana**: React + Vite + TypeScript (level: ~7 tahun FE) |
| Mobile | **Native Android**: Kotlin + Jetpack Compose (level: baru belajar) |
| Design | **Figma** — dibuat lewat plugin lokal *Benah Design Generator* (`~/workspace/benah-figma-plugin`) |
| Brand | **Benah** ("membenahi"): Benah Admin (web) & Benah Mitra (Android). Maskot: **Mang Benah** |
| Gaya UI | Ala UI Claude: krem hangat + cokelat kopi. Font **Plus Jakarta Sans** (judul), **DM Sans** (isi), **JetBrains Mono** (kode) |
| Picker tanggal/jam | Web: komponen custom. Android/iOS: picker bawaan OS |
| Login gagal | Satu pesan umum (anti user-enumeration); alasan spesifik hanya di log/Kibana |
| Komunikasi | REST API (JSON), kontrak dalam OpenAPI |
| Lingkungan | **Lokal saja**, Android pakai **emulator** |
| Mesin | MacBook Pro 2023, **M2 Pro, 16GB RAM** (Apple Silicon / ARM64) |
| Waktu | **1 hari/minggu (Minggu)** |
| Mode belajar | **Kode dibantu AI, lalu dijelaskan**. Tidak belajar dasar dari nol berbulan-bulan |
| Tujuan | **Belajar cara kerja FDE** |

---

## 1. Kenapa "FDE"

FDE bekerja dekat dengan klien: memahami masalah, memilih scope, membangun cepat dengan alat apa pun (termasuk AI), lalu iterasi. **Belajar dengan dibantu AI justru cocok dengan cara kerja FDE sekarang.** Yang dilatih adalah *judgment*: apa yang dibangun, apakah hasilnya benar, dan bagaimana menjelaskannya.

| Kebiasaan FDE | Cara melatihnya |
|---|---|
| Memahami masalah klien | Problem brief + persona di Sesi 1 |
| Scoping | Pilih fitur MVP sendiri, potong yang tidak penting |
| Kontrak yang jelas | `openapi.yaml` ditulis sebelum web/Android |
| Kirim value cepat | Per fitur *end-to-end* (DB → API → Web → Android) |
| Memvalidasi hasil AI | Jalankan, tes, dan pahami kode sebelum lanjut |
| Komunikasi | Catatan singkat tiap sesi + demo di akhir |

---

## 2. Mode Belajar: "AI Build → Explain-Back" (ringan, tidak makan waktu)

Siklus per fitur/tugas, ±15 menit overhead saja:

1. **Tentukan target**: "endpoint X yang mengembalikan Y".
2. **AI generate**: minta kode sesuai struktur project dan kontrak OpenAPI.
3. **Jalankan & tes**: curl/Bruno untuk API, klik langsung untuk web/Android.
4. **Explain-back (5 menit)**: minta AI jelaskan *hanya bagian yang belum kamu pahami*, lalu tulis 2–3 poin di `docs/session-notes.md`.
5. **Ubah 1 hal kecil sendiri**, misal tambah field atau ubah validasi. Ini cara cepat membuktikan kamu paham tanpa belajar teori panjang.

Prompt yang berguna:
- "Jelaskan kode ini seolah saya FE engineer React yang baru belajar Go/Kotlin, bandingkan dengan konsep React."
- "Apa yang tidak idiomatic di kode Go ini?"
- "Apa yang bisa salah di endpoint ini kalau dipakai di produksi?" (melatih pikiran FDE)

---

## 3. Skenario Utama

> Admin menerima permintaan servis dari pelanggan → membuat **Work Order** di web → menugaskan ke teknisi → teknisi melihat tugas di **Android** → mulai kerja → selesaikan dengan catatan → admin melihat status berubah di web.

Kalau skenario ini berjalan *end-to-end* **dan setiap langkahnya bisa dilacak di Kibana**, MVP dianggap selesai.

---

## 4. Fitur MVP (dikunci)

### Role
| Role | Platform | Bisa apa |
|---|---|---|
| `admin` | Web SPA | Kelola teknisi, pelanggan, work order |
| `technician` | Android | Lihat & kerjakan work order miliknya |

### Status Work Order
```
open ──assign──▶ assigned ──start──▶ in_progress ──complete──▶ completed
  │                  │
  └──────cancel──────┴──▶ cancelled        (cancel hanya oleh admin)
```
Aturan bisnis ini dicek di **backend (service layer)**, bukan di client. Ini latihan penting: client boleh salah, server tidak.

Aturan tambahan (v7):
- **Selesaikan** butuh ≥1 foto dan `customer_approved=true`; kalau tidak → `422 validation_failed`. Server menyimpan `customer_approved_at`.
- **Nonaktifkan mitra**: ada tugas `in_progress` → `409 technician_has_active_work` (UI: dialog diblokir). Ada tugas `assigned` → wajib pilih `reassign`: `{"mode":"open"}` (lepas ke Open) atau `{"mode":"technician","to":"<id>"}`. Nonaktif + pindah tugas + catat `work_order_events` dalam **1 transaksi DB**.

### Web (Admin)
| # | Fitur | Detail |
|---|---|---|
| W1 | Login / logout | Hanya role `admin`; password bisa ditampilkan (ikon mata) |
| W2 | Dashboard | Jumlah work order per status |
| W3 | Kelola teknisi | Tabel + aksi ikon (tooltip). Tambah/edit di **side sheet**. Nonaktifkan via Switch + konfirmasi (aturan di bawah) |
| W4 | Kelola pelanggan | Tabel + search, tambah/edit di **side sheet** (nama, telepon, alamat, catatan) |
| W5 | Work order: list | Tabel + filter status/teknisi + pagination |
| W6 | Work order: buat/edit | Judul, deskripsi, pelanggan, prioritas (low/medium/high), jadwal |
| W7 | Work order: assign & cancel | Pilih teknisi; batalkan WO |
| W8 | Work order: detail | Info lengkap + **riwayat status** + **bukti pekerjaan** (foto, catatan, mulai/selesai/durasi, waktu persetujuan pelanggan) + lightbox foto. Edit WO di side sheet |
| W12 | Riwayat mitra | Halaman detail mitra: profil, statistik (selesai bulan ini, rata-rata durasi), tabel riwayat WO; klik baris → side sheet ringkasan WO |
| W13 | Riwayat pelanggan | Halaman detail pelanggan: info, statistik, riwayat order, tombol "Buat WO untuk pelanggan ini" |

### Android (Teknisi)
| # | Fitur | Detail |
|---|---|---|
| A0 | Getting Started | Carousel 3 halaman (Tugas masuk → Mulai kerja → Kirim bukti), tombol Lanjut/Lewati; tampil sekali (flag `onboarding_done` di DataStore) |
| A1 | Login / logout | Hanya role `technician`; token disimpan di DataStore; password bisa ditampilkan (ikon mata) |
| A2 | Daftar tugas saya | Tab: Aktif (assigned + in_progress) / Selesai; pull-to-refresh |
| A3 | Detail tugas | Pelanggan, alamat, telepon, deskripsi, jadwal, prioritas |
| A4 | Mulai kerja | Konfirmasi → `in_progress`; detail menampilkan **timer** dari `started_at` (jam server) + snackbar |
| A5 | Selesaikan | Bottom sheet: catatan + **foto bukti wajib (1–5)** + checkbox *"sudah dicek & disetujui pelanggan di lokasi"* → `completed`. Tombol aktif hanya jika valid |
| A6 | Profil | Nama, email, tombol logout |
| A7 | Pekerjaan selesai | Layar terima kasih berilustrasi + ringkasan (durasi, jam selesai) + tugas berikutnya |
| A8 | Hubungi admin | Dari login gagal: form nama, email, keluhan, **upload foto** (maks 3) → tiket Call Center |

### Call Center (Web, fitur baru)
| # | Fitur | Detail |
|---|---|---|
| W9 | Form "Hubungi admin" (publik) | Tanpa login, dari halaman login web/Android: nama, email, keluhan, upload gambar (JPG/PNG, 5 MB, maks 3) |
| W10 | Call Center: list tiket | Tabel: no. tiket, nama, keluhan, jumlah lampiran, status (Baru/Diproses/Selesai), waktu masuk + filter status |
| W11 | Call Center: detail tiket | Klik baris → **side sheet**: data pengirim, keluhan, galeri lampiran, ubah status, catatan admin |

### Pola UI tambahan
- **Autocomplete** Pelanggan & Mitra di form WO: debounce 300 ms, min. 2 huruf, query `ILIKE '%' || $1 || '%'` (+ index `pg_trgm` kalau mau lebih cepat). **Tidak ada** tombol "tambah baru" di panel: data baru didaftarkan di menu Pelanggan / Mitra (satu tempat untuk validasi).
- **Password field** sama di semua tempat: ikon gembok + tombol mata. Web: toggle `type="password" ↔ "text"` + `aria-label`. Android: `PasswordVisualTransformation` ↔ `VisualTransformation.None`.
- **Side sheet** untuk semua tambah/edit dan ringkasan (WO dari riwayat mitra **dan** pelanggan); header & footer sticky.
- **Shimmer / skeleton** di semua halaman saat request pertama > 300 ms.
- **Pagination cursor** di semua tabel (web) & infinite scroll (Android).

### Pesan & ilustrasi (UX)
- **Toast** (web) / **Snackbar** (Android): aksi kecil.
- **Feedback Dialog** (satu komponen reusable: ilustrasi + judul + teks + CTA): Modal (web: WO dibuat, keluhan terkirim), Fullscreen (Android: pekerjaan selesai), Inline (empty state, 404), Compact (kartu "Tidak bisa masuk?").
- **Login gagal**: error di field; setelah 3× gagal muncul kartu bantuan → Hubungi admin.
- Halaman **404** bermerek di web.

### Di luar MVP (iterasi berikutnya, kalau mau lanjut)
GPS check-in, tanda tangan pelanggan, push notification, mode offline (Room), export laporan, multi-admin/permission detail.

### Model Data
```
users
  id, name, email (unique), password_hash, phone,
  role (admin|technician), is_active, created_at, updated_at

customers
  id, name, phone, address, notes, created_at, updated_at

work_orders
  id, code (mis. WO-0001), title, description,
  customer_id → customers, technician_id → users (nullable),
  status, priority, scheduled_at,
  started_at, completed_at, completion_note,
  customer_approved_at,                ← v7: diisi saat complete
  created_by → users, created_at, updated_at

work_order_photos          ← v7: bukti pekerjaan (min 1, maks 5)
  id, work_order_id → work_orders, uploaded_by → users,
  file_path, content_type, size_bytes, created_at

work_order_events          ← audit log / riwayat status
  id, work_order_id, actor_id → users,
  from_status, to_status, note, created_at

support_tickets            ← Call Center
  id, code (mis. TK-0012), name, email, message,
  status (new|in_progress|resolved), admin_note,
  handled_by → users (nullable), created_at, updated_at

ticket_attachments
  id, ticket_id → support_tickets, file_path,
  content_type, size_bytes, created_at
```

### Draft Endpoint REST (`/api/v1`)
| Method | Path | Role | Fitur |
|---|---|---|---|
| POST | `/auth/login` | publik | W1, A1 |
| POST | `/auth/refresh` | login | W1, A1 |
| GET | `/me` | login | A6 |
| GET | `/dashboard/summary` | admin | W2 |
| GET, POST | `/technicians` (`?q=&active=true&limit=` untuk autocomplete) | admin | W3, W6 |
| GET, PATCH | `/technicians/{id}` (`{is_active:false, reassign:{mode, to}}` saat nonaktif; 409 jika ada tugas Dikerjakan) | admin | W3 |
| GET | `/technicians/{id}/work-orders?limit=15&cursor=` | admin | W12 |
| GET, POST | `/customers` (`?q=&limit=8` untuk autocomplete) | admin | W4, W6 |
| GET, PATCH | `/customers/{id}` | admin | W4 |
| GET | `/customers/{id}/work-orders?limit=15&cursor=` | admin | W13 |
| GET, POST | `/work-orders` | admin | W5, W6 |
| GET, PATCH | `/work-orders/{id}` | admin | W6, W8 |
| POST | `/work-orders/{id}/assign` | admin | W7 |
| POST | `/work-orders/{id}/cancel` | admin | W7 |
| GET | `/work-orders/{id}/events` | admin | W8 |
| GET | `/me/work-orders?status=` | technician | A2 |
| GET | `/me/work-orders/{id}` | technician | A3 |
| POST | `/me/work-orders/{id}/start` | technician | A4 |
| POST | `/me/work-orders/{id}/complete` (**multipart**: `note`, `customer_approved`, `photos[]` 1–5) | technician | A5, A7 (respons berisi `duration`, `completed_at`, `next_task`) |
| GET | `/work-orders/{id}/photos/{photoId}` | admin | W8 (lightbox) |
| POST | `/support-tickets` (multipart) | **publik** + rate limit | W9, A8 |
| GET | `/support-tickets?status=` | admin | W10 |
| GET, PATCH | `/support-tickets/{id}` | admin | W11 |
| GET | `/support-tickets/{id}/attachments/{attachmentId}` | admin | W11 |

Konvensi: response error konsisten `{ "error": { "code": "...", "message": "..." } }`, list pakai **cursor pagination** `?limit=15&cursor=<opaque>` (limit 15–20, default 15) → respons `{ "data": [...], "next_cursor": "…", "prev_cursor": "…" }`; cursor = base64 dari `(created_at, id)` item terakhir, query `WHERE (created_at, id) < ($1, $2) ORDER BY created_at DESC, id DESC LIMIT 16` (ambil 1 lebih untuk tahu ada halaman berikutnya). Web: tombol Sebelumnya/Berikutnya; Android: infinite scroll, waktu dalam ISO 8601 (UTC).

### Layar Figma
- **Web (±30 + 8 loading):** Login, Dashboard, List WO, Form WO (+ autocomplete, popover tanggal), Detail WO, Detail WO selesai (bukti foto), Lightbox, Edit WO (side sheet), Mitra, Tambah/Edit mitra (side sheet), Nonaktifkan ×3 (konfirmasi / pilih nasib tugas / diblokir), Detail mitra + riwayat (+ side sheet WO), Pelanggan, Edit pelanggan, Detail pelanggan + riwayat (+ side sheet WO), Call Center + Detail tiket (side sheet), Login gagal, Hubungi admin, Keluhan terkirim, WO dibuat, 404, WO kosong. Loading: dashboard, WO, detail WO, mitra, pelanggan, call center, detail mitra, detail pelanggan
- **Android (±17 + 3 loading):** Getting Started (carousel 3 halaman), Login, Daftar Tugas, Detail, Konfirmasi mulai, Sedang dikerjakan (timer), Selesaikan (sheet + foto + persetujuan), Selesaikan belum lengkap, Profil, Error, Pekerjaan selesai, Belum ada tugas, Login gagal, Hubungi admin, Keluhan terkirim
- **Journey** di Figma memakai **cabang** (baris baru + panah turun): mis. Mitra → Edit / Nonaktifkan (3 varian) / Riwayat
- **Brand:** logo **b Palu** (siluet rumah, huruf "b" berkepala palu, kilau "beres"); lockup **Benah Admin** / **Benah Mitra**; frame "Logo Options" menyimpan 6 alternatif.
- **Design system:** token, Typography (komponen Style × Tone, text style tetap jadi acuan kode), 40 ikon, Password Field, Page Indicator, Icon Button, Tooltip, Switch, Radio, Confirm Dialog, Photo Lightbox, Side/Bottom Sheet, Button, Input/Select/Textarea, Date & Time Picker (web), Segmented, Badge/Tag/Avatar, Progress, Timeline, Toast, Snackbar, Alert, Error State, Feedback Dialog, Image Upload, Ticket Status, Sidebar, Bottom Nav, App Bar, kartu, 8 ilustrasi Mang Benah

### Data Awal (seed)
1 admin, 3 teknisi, 5 pelanggan, ±10 work order dengan status campuran (2 yang Selesai punya foto bukti), 5 tiket Call Center (2 dengan lampiran). Ini supaya web dan Android langsung punya data untuk dites.

---

## 5. Arsitektur (lokal)

```
                    Figma (design system + layar)
                 ┌────────────────┴────────────────┐
                 ▼                                 ▼
      ┌────────────────────┐            ┌─────────────────────┐
      │ Web SPA            │            │ Android Emulator    │
      │ React + Vite :5173 │            │ Kotlin + Compose    │
      │ → localhost:8080   │            │ → 10.0.2.2:8080     │
      └─────────┬──────────┘            └──────────┬──────────┘
                └──────────────┬────────────────────┘
                               ▼   REST API (JSON, OpenAPI)
                 ┌───────────────────────────┐
                 │ Backend Go  :8080         │── JSON log ──▶ Filebeat
                 └─────────────┬─────────────┘                  │
                               ▼                                ▼
                 ┌───────────────────────────┐      Elasticsearch :9200
                 │ PostgreSQL :5432          │                  │
                 └───────────────────────────┘           Kibana :5601
```

---

## 6. Observability (Elasticsearch + Kibana)

### Tujuan (dari sudut pandang FDE)
Di lapangan, laporan klien biasanya seperti ini: *"Teknisi Budi bilang tombol Selesai error tadi jam 10."* FDE harus bisa menjawab **apa yang terjadi** tanpa menebak. Jadi target observability di project ini sederhana: **setiap request bisa dilacak dari Android/Web sampai backend lewat satu `request_id`, dan semuanya bisa dicari di Kibana.**

### Alur
```
Web / Android ──(header X-Request-ID, X-Client)──▶ Backend Go
                                                     │ log/slog JSON
                                                     ▼
                                      backend/logs/app.log  (+ stdout)
                                                     │
                                                     ▼
                                    Filebeat ──▶ Elasticsearch :9200
                                                     │
                                                     ▼
                                                Kibana :5601
```

### Lisensi
Elasticsearch & Kibana sejak 2024 kembali tersedia dengan lisensi open source **AGPLv3** (selain lisensi Elastic/SSPL). Untuk belajar di lokal, ini aman dipakai. Alternatif yang 100% Apache 2.0 adalah **OpenSearch + OpenSearch Dashboards**, dengan konsep yang hampir sama.

### Apa yang di-log (field standar)
| Field | Contoh | Sumber |
|---|---|---|
| `time`, `level`, `msg` | `2026-11-15T03:10:00Z`, `INFO`, `request completed` | slog |
| `request_id` | `b7f2…` | Header `X-Request-ID` dari client, atau dibuat backend kalau kosong |
| `client` | `web` / `android` | Header `X-Client` |
| `method`, `path`, `route` | `POST`, `/api/v1/me/work-orders/12/complete`, `/me/work-orders/{id}/complete` | middleware |
| `status`, `duration_ms` | `409`, `23` | middleware |
| `user_id`, `role` | `7`, `technician` | dari JWT |
| `error_code` | `INVALID_STATUS_TRANSITION` | error handler |
| `event` (bisnis) | `work_order.status_changed` + `work_order_id`, `from`, `to` | service layer |

Aturan: **jangan log password, token, atau data pribadi lengkap**.

### Implementasi singkat
- **Backend Go**
  - `slog.NewJSONHandler` ke `io.MultiWriter(os.Stdout, file)` → `backend/logs/app.log`
  - Middleware: buat/ambil `request_id`, simpan di `context`, kirim balik di response header `X-Request-ID`, log 1 baris per request (status + durasi)
  - Log event bisnis di service saat status work order berubah
- **Web**: wrapper fetch menambahkan `X-Request-ID: crypto.randomUUID()` dan `X-Client: web`; tampilkan `request_id` di pesan error ("Kode error: b7f2…")
- **Android**: OkHttp `Interceptor` menambahkan `X-Request-ID` (UUID) dan `X-Client: android`; `request_id` ditampilkan di Snackbar error
- **Kibana**
  - Data view `fieldservice-logs-*`
  - **Discover**: cari `request_id : "b7f2…"` atau `level : "ERROR"`
  - **Dashboard sederhana**: request per endpoint, error rate (status ≥ 400), latency p95, jumlah perubahan status work order per hari, perbandingan request web vs android

### Docker Compose (profil terpisah supaya hemat RAM)
Tambahkan ke `docker-compose.yml`. Service ini **hanya jalan kalau dipanggil** dengan `--profile observability`.

`.env`:
```
ELASTIC_VERSION=9.x.x   # isi dengan versi terbaru dari elastic.co/downloads
```

```yaml
services:
  # ... db (sudah ada)

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_VERSION}
    profiles: ["observability"]
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false      # HANYA untuk lokal
      - ES_JAVA_OPTS=-Xms512m -Xmx512m    # batasi RAM
    ports:
      - "9200:9200"
    volumes:
      - esdata:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:${ELASTIC_VERSION}
    profiles: ["observability"]
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on: [elasticsearch]

  filebeat:
    image: docker.elastic.co/beats/filebeat:${ELASTIC_VERSION}
    profiles: ["observability"]
    user: root
    command: ["-e", "--strict.perms=false"]
    volumes:
      - ./infra/filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - ./backend/logs:/logs:ro
    depends_on: [elasticsearch]

volumes:
  pgdata:
  esdata:
```

`infra/filebeat.yml`:
```yaml
filebeat.inputs:
  - type: filestream
    id: fieldservice-api
    paths: ["/logs/*.log"]
    parsers:
      - ndjson:
          target: ""
          overwrite_keys: true
          add_error_key: true

output.elasticsearch:
  hosts: ["http://elasticsearch:9200"]
  index: "fieldservice-logs-%{+yyyy.MM.dd}"

setup.template.name: "fieldservice-logs"
setup.template.pattern: "fieldservice-logs-*"
setup.ilm.enabled: false
```

Perintah:
```
docker compose up -d                              # Postgres saja
docker compose --profile observability up -d      # + Elasticsearch, Kibana, Filebeat
docker compose --profile observability stop       # matikan saat sesi Android
```

### Latihan FDE di Sesi 9
Sengaja buat skenario masalah, lalu cari penyebabnya **hanya dari Kibana**:
1. Teknisi mencoba "Selesaikan" work order yang sudah di-cancel admin → 409
2. Token kedaluwarsa di Android → 401 lalu refresh
3. Endpoint list work order lambat (tambah `time.Sleep` sementara) → terlihat di latency p95

### Opsional nanti
Traces dengan **OpenTelemetry → Elastic APM** untuk melihat waktu per query DB dalam satu request.

---

## 7. Tech Stack

### Backend: Go
- `net/http` (Go 1.22+ routing) + **chi** kalau perlu middleware lebih rapi
- **sqlc** + **pgx** (SQL → kode Go type-safe), migrasi dengan **goose**
- JWT auth, `log/slog` (JSON → Elasticsearch), go-playground/validator, **air** (hot reload)

```
backend/
  cmd/api/main.go
  internal/{handler,service,repository,middleware}/
  db/{migrations,queries}/
```

### Web: SPA
- React + Vite + TypeScript, React Router, **TanStack Query**
- Styling bebas (Tailwind direkomendasikan supaya mudah dari Figma)
- Opsional: `openapi-typescript` untuk generate type dari kontrak

### Android: Native
- Kotlin, **Jetpack Compose** + Material 3, MVVM (`ViewModel` + `StateFlow`)
- **Retrofit** + OkHttp + kotlinx.serialization, Coroutines, Navigation Compose, DataStore
- DI manual dulu (Hilt nanti kalau perlu)

| React | Android / Compose |
|---|---|
| Component | `@Composable` |
| `useState` | `remember { mutableStateOf() }` |
| TanStack Query / store | `ViewModel` + `StateFlow` |
| `useEffect` | `LaunchedEffect` |
| React Router | Navigation Compose |
| fetch/axios | Retrofit |
| localStorage | DataStore |
| Theme / design token | `MaterialTheme` |

### Design: Figma
- Design system **custom Benah** (krem-cokelat, gaya Claude) dibuat lewat plugin lokal "Benah Design Generator" (`workspace/benah-figma-plugin`) di akun Figma pribadi
- **Variables** "Benah Tokens" dengan code syntax WEB `var(--x)` dan ANDROID `BenahColors.x` → dipetakan ke Tailwind dan `MaterialTheme`/`BenahTheme`
- 3 halaman: Brand & Components, Web Admin, Android Mitra; layar dikelompokkan per **journey** (bercabang) dengan label aksi + endpoint API
- Web = admin/dispatcher (desktop), Android = teknisi lapangan (mobile)

---

## 8. Setup untuk MacBook M2 Pro 16GB

### Tools
| Tool | Install | Catatan |
|---|---|---|
| Homebrew | brew.sh | Package manager |
| Go | `brew install go` | |
| sqlc, goose, air | `brew install sqlc goose` / `go install github.com/air-verse/air@latest` | |
| OrbStack | `brew install orbstack` | Pengganti Docker Desktop yang lebih ringan; Postgres jalan di container |
| Node.js | `brew install node` (atau fnm/nvm) | Kemungkinan sudah ada |
| Android Studio | developer.android.com/studio (**versi Apple Silicon**) | Termasuk SDK & emulator |
| API client | Bruno / Postman / curl | Bruno: koleksi API bisa disimpan di repo |
| DB viewer | TablePlus / DBeaver | |
| Figma | Desktop app | |

### Catatan RAM 16GB
Android Studio + emulator + Docker + browser + editor bisa terasa berat di 16GB. Saran:
- Pakai **OrbStack**, bukan Docker Desktop (lebih hemat RAM di Apple Silicon). Kalau masih berat, `docker compose stop` saat tidak dipakai.
- **Elasticsearch + Kibana** butuh ±1,5–2GB RAM. Jalankan hanya lewat `--profile observability` dan **jangan dinyalakan bersamaan dengan emulator** kecuali saat Sesi 9 (integrasi). Heap ES dibatasi 512MB.
- Emulator: pakai system image **ARM64** (`arm64-v8a`), device Pixel, **RAM emulator 2GB** cukup.
- Saat sesi Android, tutup yang tidak perlu (tab browser, Figma).
- Saat sesi backend/web, Android Studio tidak perlu dibuka.

### PostgreSQL via OrbStack
`docker-compose.yml` di root repo:
```yaml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: fieldservice
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
volumes:
  pgdata:
```
Jalankan `docker compose up -d`, lalu connect: `postgres://app:app@localhost:5432/fieldservice?sslmode=disable`.
Reset total (hapus data): `docker compose down -v`.

### Koneksi Emulator → API lokal
1. Backend listen di `:8080`.
2. Di Android, base URL = **`http://10.0.2.2:8080`** (alias emulator untuk localhost Mac). Tidak perlu static IP.
3. Izinkan HTTP cleartext khusus `10.0.2.2` lewat `res/xml/network_security_config.xml` (Android 9+ memblokir HTTP secara default).
4. Simpan base URL di `BuildConfig` (lewat `build.gradle.kts`), jangan hardcode.
5. Web: aktifkan CORS untuk `http://localhost:5173` di backend, atau pakai proxy Vite.

### Struktur repo
```
fde-field-service/
  api/openapi.yaml
  backend/
  web/
  android/
  design/            ← link Figma
  docs/              ← problem-brief, ERD, session-notes, retro
  infra/filebeat.yml
  docker-compose.yml ← Postgres + profil observability
  .env               ← ELASTIC_VERSION
```

---

## 9. Roadmap Padat (±10 Minggu, ~2,5 bulan)

Setiap Minggu (5–6 jam) mengikuti siklus di Bagian 2. Tidak ada fase "belajar dasar" terpisah; dasar Go/Kotlin dipelajari sambil membangun.

| Sesi | Fokus | Output |
|---|---|---|
| **1** | Setup semua tools + Postgres jalan + `/health` dari emulator; tulis problem brief singkat dari Bagian 3–4 | Lingkungan siap, `docs/problem-brief.md` |
| **2** | **Figma** ✅ (selesai lewat plugin): baca journey per role, catat token → `tailwind.config` & `Color.kt`/`Type.kt`, daftar komponen yang akan dibuat ulang di React & Compose | `docs/ui-inventory.md` |
| **3** | **Data & kontrak**: migrasi tabel inti (`users`, `customers`, `work_orders`, `work_order_events`, `work_order_photos`), seed data, `openapi.yaml`, struktur backend Go + **slog JSON & middleware request_id** | DB + kontrak + log terstruktur |
| **4** | **Backend**: auth (login/refresh/me, JWT, middleware role) + teknisi & pelanggan CRUD | API W1, W3, W4 (+ aturan nonaktif & reassign dalam 1 transaksi) |
| **5** | **Backend**: work order + aturan status + events + dashboard + endpoint teknisi (complete **multipart** + foto, `customer_approved_at`) + endpoint riwayat mitra/pelanggan; test service status; **log event bisnis** | API MVP lengkap |
| **6** | **Web SPA**: W1–W8, W12–W13 (pola side sheet sekali, dipakai ulang) + header `X-Request-ID`; sisa waktu: **nyalakan Elasticsearch/Kibana/Filebeat, buat data view & dashboard** | Web MVP + Kibana jalan |
| **7** | **Android**: struktur MVVM, theme dari Figma, Retrofit + **OkHttp interceptor request_id**, A0 onboarding (HorizontalPager), A1 login + A2 daftar tugas | App jalan di emulator |
| **8** | **Android**: A3 detail, A4 mulai + timer, A5 selesaikan (pilih/ambil foto, upload multipart, validasi), A6 profil | Android MVP |
| **9** | **Call Center**: tabel `support_tickets` + `ticket_attachments`, endpoint publik multipart + rate limit, simpan file di disk lokal, layar web W9–W11 & Android A8 | Fitur Call Center jalan |
| **10** | **Integrasi & demo**: skenario end-to-end (admin di web → teknisi di Android), **latihan debugging via Kibana** (Bagian 6), perbaikan bug, retro | `docs/retro.md`, demo |

Kalau ternyata lebih cepat, Sesi 4–8 bisa diganti dengan pola **vertical slice**: satu fitur per Minggu langsung di backend, web, dan Android sekaligus.

---

## 10. Checklist Sesi 1

**Pagi: Setup (±2–3 jam)**
- [ ] Install Homebrew, Go, Node, sqlc, goose, air
- [ ] Install OrbStack, `docker compose up -d`, cek connect via TablePlus/DBeaver
- [ ] Install Android Studio, buat emulator Pixel ARM64, jalankan template "Empty Compose Activity"
- [ ] Buat repo `fde-field-service` dengan struktur folder di atas
- [ ] Go: jalankan server `GET /health` di `:8080`
- [ ] Tes dari emulator: buka `http://10.0.2.2:8080/health` di Chrome emulator

**Siang: Discovery ringan (±1–2 jam)**
- [ ] Tulis `docs/problem-brief.md` (1 halaman) berdasarkan Bagian 3–4: klien, masalah, persona, skenario utama
- [ ] Gambar user flow admin & teknisi (boleh di FigJam)
- [ ] Sisa waktu: buka file Figma Benah, jalankan plugin sekali, lihat journey Web J1–J7 & Android J1–J4 sebagai persiapan Sesi 2

---

## 11. Pertanyaan Terbuka

1. Ada yang mau diubah dari fitur MVP di Bagian 4?
2. Mau saya siapkan isi **Sesi 1** (repo awal, `docker-compose.yml`, server `/health`, `network_security_config.xml`) supaya tinggal jalan?
3. Sesi 2 (Figma) sudah banyak terbantu plugin: mau sisa waktunya dipakai untuk mulai Sesi 3 (migrasi DB) lebih awal?
