# Plan: Issue #2: Infrastructure (config, logger, Postgres, docker-compose, Makefile)
> Issue: https://github.com/irfanandriansyah1997/bang-backend/issues/2 · Branch: `feat/2-infrastructure` · Last updated: 2026-11-21

**Goal:** Siapkan "kotak perkakas" yang dipakai semua fitur: membaca konfigurasi dari env, logger JSON yang bisa dicari di Kibana, koneksi ke Postgres, dan aplikasi yang bisa dimatikan dengan rapi (Ctrl+C).

**Depends on:** #1 ✅ (sudah di-merge)

**Already done in the repo:** `logs/` dan `.env` sudah ada di `.gitignore`. Selain itu belum ada apa-apa.

**Before you start:**
- Nyalakan Docker. Docker Desktop terinstal tapi daemon-nya mati. Pilih salah satu: buka Docker Desktop, atau install OrbStack (sesuai spec, lebih ringan) dengan `brew install orbstack` lalu buka aplikasinya.
- Buat branch:
  ```bash
  git switch main && git pull
  git switch -c feat/2-infrastructure
  ```

Semua kode di file ini sudah di-compile dan di-test di project percobaan: `go vet` bersih, `gofmt` bersih, dan semua test lolos dengan `-race`.

## Progress
- [ ] Step 1: Jalankan Postgres dengan Docker
- [ ] Step 2: Package `config` (+ `.env.example`)
- [ ] Step 3: Package `logger` (level, payload, SDC, x_request_id, pretty print)
- [ ] Step 4: Package `database` (pgxpool + ping)
- [ ] Step 4b: Postgres service container di CI, supaya test database ikut jalan
- [ ] Step 5: Sambungkan semuanya di `cmd/api/main.go` + graceful shutdown *(batch berikutnya)*
- [ ] Step 6: `Makefile` + `.air.toml` *(batch berikutnya)*

## What you'll build

```
docker-compose.yml                                   NEW   (Step 1)
.env.example                                         NEW   (Step 2)
internal/infrastructure/config/config.go             NEW   (Step 2)
internal/infrastructure/config/config_test.go        NEW   (Step 2)
internal/infrastructure/logger/logger.go             NEW   (Step 3)
internal/infrastructure/logger/pretty.go             NEW   (Step 3)
internal/infrastructure/logger/logger_test.go        NEW   (Step 3)
internal/infrastructure/database/database.go         NEW   (Step 4)
internal/infrastructure/database/database_test.go    NEW   (Step 4)
go.mod, go.sum                                       EDIT  (Step 4, otomatis lewat go get)
.github/workflows/ci.yml                             EDIT  (Step 4b)
cmd/api/main.go, cmd/api/main_test.go                EDIT  (Step 5)
Makefile, .air.toml                                  EDIT / NEW (Step 6)
```

---

## Step 1: Jalankan Postgres dengan Docker

**File:** `docker-compose.yml` (NEW, di root repo; bukan bagian dari layer Go)

**Why:** Kamu tidak perlu menginstal Postgres langsung di Mac. Docker menjalankannya di dalam container yang bisa dinyalakan dan dimatikan kapan saja. Kira-kira mirip `npx`, tapi untuk database: tinggal pakai, tanpa instalasi permanen.

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

**What this does:**

| Baris | Artinya |
|---|---|
| `image: postgres:16` | Pakai Postgres versi 16 |
| `environment` | Username, password, dan nama database yang dibuat saat container pertama kali jalan |
| `ports: "5432:5432"` | Port 5432 di Mac diteruskan ke port 5432 di dalam container |
| `volumes: pgdata` | Data disimpan di volume `pgdata`, jadi tidak hilang saat container dimatikan |

**Check:**
```bash
docker compose up -d
docker compose ps        # service "db" berstatus running / Up
```
Setelah itu, coba connect dari TablePlus/DBeaver dengan: host `localhost`, port `5432`, user `app`, password `app`, database `fieldservice`.

**Common mistakes:**
- **Docker belum menyala.** Muncul error `Cannot connect to the Docker daemon`. Buka Docker Desktop atau OrbStack dulu.
- **Port 5432 sudah dipakai.** Muncul error `port is already allocated`. Matikan Postgres lain yang sedang jalan, atau ganti mapping port jadi `"5433:5432"`.

---

## Step 2: Package `config`

### 2a. `.env.example` (NEW, di root)

**Why:** File ini berisi daftar env yang dibutuhkan aplikasi, lengkap dengan contoh nilainya, dan **ikut di-commit**. File `.env` (yang berisi nilai sungguhan) **tidak** di-commit karena sudah ada di `.gitignore`.

```bash
# Copy to .env and fill in real values: cp .env.example .env
PORT=8080
DATABASE_URL=postgres://app:app@localhost:5432/fieldservice?sslmode=disable
JWT_SECRET=change-me-to-a-long-random-string
LOG_FILE=logs/app.log
LOG_LEVEL=info        # debug, info, warn, error
APP_ENV=local         # local, staging, production
APP_VERSION=dev       # git tag or commit
```

Buat salinan `.env` untuk dirimu sendiri:
```bash
cp .env.example .env
```

### 2b. `config.go`

**File:** `internal/infrastructure/config/config.go` (NEW, layer: `infrastructure`)

**Why:** Semua setting dibaca **di satu tempat** dan divalidasi saat aplikasi start. Package lain tidak boleh memanggil `os.Getenv`; mereka cukup menerima `Config`. Padanannya di frontend: satu file `env.ts` yang mem-parse `import.meta.env`, lalu di-import ke mana-mana.

```go
package config

import (
	"errors"
	"log/slog"
	"os"
)

// Config holds every setting the app needs, read once at startup.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	LogFile     string
	LogLevel    string // debug, info, warn, error
	AppEnv      string // local, staging, production
	AppVersion  string // git tag or commit
}

// Load reads the config from real environment variables.
func Load() (Config, error) {
	return LoadFrom(os.Getenv)
}

// LoadFrom reads the config using getenv, so tests can pass a fake one.
func LoadFrom(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:        valueOr(getenv("PORT"), "8080"),
		DatabaseURL: valueOr(getenv("DATABASE_URL"), "postgres://app:app@localhost:5432/fieldservice?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET"),
		LogFile:     valueOr(getenv("LOG_FILE"), "logs/app.log"),
		LogLevel:    valueOr(getenv("LOG_LEVEL"), "info"),
		AppEnv:      valueOr(getenv("APP_ENV"), "local"),
		AppVersion:  valueOr(getenv("APP_VERSION"), "dev"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("config: JWT_SECRET is required")
	}

	return cfg, nil
}

// LogValue controls what slog prints for a Config: secrets are never logged.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("port", c.Port),
		slog.String("log_file", c.LogFile),
		slog.String("log_level", c.LogLevel),
	)
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
```

**What this does:**
- `Load()` membaca env dari sistem. Semua env punya nilai default, kecuali `JWT_SECRET`. Kalau `JWT_SECRET` kosong, aplikasi menolak start dan menampilkan pesan error yang jelas.
- `LogLevel` sengaja disimpan sebagai string. Tugas config cuma membaca env. Mengubah string itu jadi level `slog` adalah tugas package `logger` (lewat `ParseLevel` di Step 3).
- `LogValue` memastikan bahwa saat `Config` di-log, yang tercetak hanya port, lokasi file log, dan level. Password database dan JWT secret **tidak pernah** ikut tercetak.

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `type Config struct { ... }` | Kumpulan field dengan tipe yang sudah ditentukan | `type Config = { port: string; ... }` |
| `func Load() (Config, error)` | Fungsi yang mengembalikan **dua nilai**: hasil dan error | Mirip `return [config, error]` |
| `errors.New("...")` | Membuat sebuah error. Error di Go di-`return`, bukan di-`throw` | `new Error("...")` |
| `getenv func(string) string` | Parameter yang isinya sebuah fungsi | `(getenv: (key: string) => string)` |
| `func (c Config) LogValue()` | **Method**: fungsi yang menempel ke tipe `Config`. `c` disebut *receiver* | Method di class |
| `Load` vs `valueOr` | Huruf **besar** = public (bisa dipakai package lain), huruf **kecil** = private | `export function` vs fungsi tanpa `export` |

**Kenapa ada `LoadFrom(getenv)`?** Supaya bisa dites. Di test kita memberi "env palsu" berupa sebuah map. Polanya sama seperti `run(io.Writer)` di issue #1: dependency dimasukkan lewat parameter.

### 2c. `config_test.go`

**File:** `internal/infrastructure/config/config_test.go` (NEW)

```go
package config

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

// fakeEnv turns a map into a getenv function.
func fakeEnv(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoadFrom(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		wantPort     string
		wantLogLevel string
		wantErr      bool
	}{
		{name: "defaults", env: map[string]string{"JWT_SECRET": "s3cret"}, wantPort: "8080", wantLogLevel: "info"},
		{name: "custom values", env: map[string]string{"JWT_SECRET": "s3cret", "PORT": "9090", "LOG_LEVEL": "debug"}, wantPort: "9090", wantLogLevel: "debug"},
		{name: "missing secret", env: map[string]string{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFrom(fakeEnv(tt.env))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Port != tt.wantPort {
				t.Errorf("Port = %q, want %q", cfg.Port, tt.wantPort)
			}
			if cfg.LogLevel != tt.wantLogLevel {
				t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, tt.wantLogLevel)
			}
		})
	}
}

func TestConfig_LogValue_HidesSecrets(t *testing.T) {
	cfg := Config{
		Port:        "8080",
		DatabaseURL: "postgres://app:supersecret@localhost:5432/db",
		JWTSecret:   "jwt-secret-value",
		LogFile:     "logs/app.log",
		LogLevel:    "info",
	}

	var buf bytes.Buffer
	slog.New(slog.NewJSONHandler(&buf, nil)).Info("config loaded", "config", cfg)
	out := buf.String()

	for _, secret := range []string{"supersecret", "jwt-secret-value"} {
		if strings.Contains(out, secret) {
			t.Errorf("log output leaks %q: %s", secret, out)
		}
	}
	if !strings.Contains(out, `"port":"8080"`) {
		t.Errorf("log output should include the port: %s", out)
	}
}
```

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `map[string]string{...}` | Dictionary string ke string | `Record<string, string>` |
| `values[key]` | Kalau key tidak ada, hasilnya `""`, **bukan error** | `obj[key] ?? ""` |
| `return func(...) {...}` | Fungsi yang mengembalikan fungsi (closure) | `return (key) => values[key]` |
| `nil` | Artinya "tidak ada", mirip `null`. `err == nil` berarti tidak ada error | `err === null` |

**Check:**
```bash
go test -v ./internal/infrastructure/config
```
Hasilnya: `TestLoadFrom` beserta 3 sub-test-nya `PASS`, `TestConfig_LogValue_HidesSecrets` `PASS`, dan baris terakhir `ok ...`.

**Common mistakes:**
- **Nama package di satu folder berbeda-beda.** Ada file yang baris pertamanya `package main` di folder yang sama dengan `package config`. Muncul error `found packages main (x.go) and config (y.go)`. Semua file dalam satu folder harus memakai nama package yang sama.
- **Fungsi ditulis huruf kecil, misalnya `loadFrom`.** Fungsi itu tidak bisa dipanggil dari `main.go`, dan muncul error `undefined: config.loadFrom`.

---

## Step 3: Package `logger`

Cara pakainya (bentuk yang kamu minta):

| Gaya JS | Di Go |
|---|---|
| `logger.Log("message")` | `log.Info(ctx, "message")` |
| `logger.Log("message", { user_id: 1, sample_message: 'asd asd' })` | `log.Info(ctx, "message", logger.Payload{"user_id": 1, "sample_message": "asd asd"})` |
| `logger.Log("message", {}, { correlation_id: '12121211', path: '/usr/asd/asd' })` | `log.Info(ctx, "message", logger.Payload{}, logger.SDC{"correlation_id": "12121211", "path": "/usr/asd/asd"})` |

- `Log` dipecah jadi empat method sesuai level: `Debug`, `Info`, `Warn`, `Error`. Cara pakai keempatnya sama.
- `ctx` selalu di depan. Ini aturan repo, dan lewat `ctx` inilah `x_request_id` serta SDC dari middleware ikut terbawa.
- `Payload` dan `SDC` boleh ditulis dengan urutan bebas. `Payload{}` yang kosong boleh dihilangkan saja.

**Output sungguhan** (hasil menjalankan kode di bawah):
```json
{"timestamp":"2026-11-08T15:49:22.380Z","level":"INFO","message":"message","service":"bang-backend","env":"local","version":"dev","x_request_id":"b7f2"}
{"timestamp":"2026-11-08T15:49:22.381Z","level":"INFO","message":"message","service":"bang-backend","env":"local","version":"dev","x_request_id":"b7f2","payload":{"sample_message":"asd asd","user_id":1}}
{"timestamp":"2026-11-08T15:49:22.381Z","level":"INFO","message":"message","service":"bang-backend","env":"local","version":"dev","x_request_id":"b7f2","sdc":{"correlation_id":"12121211","path":"/usr/asd/asd"}}
```

**Field di setiap baris log:**

| Field | Asal | Di-index di Kibana? |
|---|---|---|
| `timestamp`, `level`, `message` | Otomatis. `timestamp` dalam UTC dengan presisi milidetik | ✅ |
| `service`, `env`, `version` | `Options` saat `New`, nilainya dari config | ✅ |
| `x_request_id` | `WithRequestID(ctx, id)`. Nanti dipasang di middleware (#3) | ✅ |
| `sdc` | `SDC{...}` saat logging, ditambah `WithSDC(ctx, ...)`. **Nilainya harus string** | ✅ |
| `payload` | `Payload{...}` saat logging. Nilainya bebas | ❌ (diatur lewat mapping Elasticsearch nanti, lihat Notes) |

**Prioritas SDC:** SDC dari context (`WithSDC`) dipakai lebih dulu. Kalau key yang sama juga dikirim langsung saat logging, nilai yang dikirim saat logging yang menang.

### 3a. `logger.go`

**File:** `internal/infrastructure/logger/logger.go` (NEW, layer: `infrastructure`)

```go
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

// Payload is extra detail for one log line. It is stored under "payload"
// and is NOT indexed in Kibana: use it for data you only read, never search.
type Payload map[string]any

// SDC (Standardized Diagnostic Context) holds fields you want to SEARCH in
// Kibana, e.g. correlation_id or path. It is stored under "sdc" and indexed.
// Values are strings only, so Elasticsearch always sees the same field type.
type SDC map[string]string

// Extra is anything that can be passed after the message: Payload or SDC.
type Extra interface {
	apply(e *entry)
}

type entry struct {
	payload Payload
	sdc     SDC
}

func (p Payload) apply(e *entry) { maps.Copy(e.payload, p) }
func (s SDC) apply(e *entry)     { maps.Copy(e.sdc, s) }

// Options are the fixed fields added to every log line.
type Options struct {
	Level   slog.Level
	Service string // e.g. "bang-backend"
	Env     string // e.g. "local", "staging", "production"
	Version string // e.g. a git tag or commit
}

// Logger writes JSON log lines with a message, an optional Payload, and an optional SDC.
type Logger struct {
	base *slog.Logger
}

const timestampFormat = "2006-01-02T15:04:05.000Z07:00" // UTC, milliseconds

// New creates a JSON logger that writes to w.
func New(w io.Writer, opts Options) *Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       opts.Level,
		ReplaceAttr: renameKeys,
	})
	base := slog.New(handler).With(
		slog.String("service", opts.Service),
		slog.String("env", opts.Env),
		slog.String("version", opts.Version),
	)
	return &Logger{base: base}
}

// renameKeys turns slog's "time" and "msg" into "timestamp" (UTC, ms) and "message".
func renameKeys(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	switch a.Key {
	case slog.TimeKey:
		return slog.String("timestamp", a.Value.Time().UTC().Format(timestampFormat))
	case slog.MessageKey:
		return slog.String("message", a.Value.String())
	}
	return a
}

// Debug, Info, Warn, and Error log at their level:
//
//	log.Info(ctx, "message")
//	log.Info(ctx, "message", logger.Payload{"user_id": 1})
//	log.Info(ctx, "message", logger.Payload{}, logger.SDC{"correlation_id": "abc"})
func (l *Logger) Debug(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelDebug, msg, extras)
}

func (l *Logger) Info(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelInfo, msg, extras)
}

func (l *Logger) Warn(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelWarn, msg, extras)
}

func (l *Logger) Error(ctx context.Context, msg string, extras ...Extra) {
	l.log(ctx, slog.LevelError, msg, extras)
}

func (l *Logger) log(ctx context.Context, level slog.Level, msg string, extras []Extra) {
	if !l.base.Enabled(ctx, level) {
		return
	}

	// SDC priority: context (WithSDC) first, then values passed to this call win.
	e := entry{payload: Payload{}, sdc: SDC{}}
	maps.Copy(e.sdc, sdcFrom(ctx))
	for _, x := range extras {
		x.apply(&e)
	}

	var attrs []slog.Attr
	if id := RequestID(ctx); id != "" {
		attrs = append(attrs, slog.String("x_request_id", id))
	}
	if len(e.payload) > 0 {
		attrs = append(attrs, slog.Any("payload", map[string]any(e.payload)))
	}
	if len(e.sdc) > 0 {
		attrs = append(attrs, slog.Any("sdc", map[string]string(e.sdc)))
	}
	l.base.LogAttrs(ctx, level, msg, attrs...)
}

// --- values carried in context.Context ---

type sdcKey struct{}
type requestIDKey struct{}

// WithSDC returns a copy of ctx that carries SDC fields (e.g. user_id).
// Every log call made with that ctx includes them under "sdc".
func WithSDC(ctx context.Context, sdc SDC) context.Context {
	merged := SDC{}
	maps.Copy(merged, sdcFrom(ctx))
	maps.Copy(merged, sdc)
	return context.WithValue(ctx, sdcKey{}, merged)
}

func sdcFrom(ctx context.Context) SDC {
	sdc, _ := ctx.Value(sdcKey{}).(SDC)
	return sdc
}

// WithRequestID returns a copy of ctx that carries the request ID.
// Every log call made with that ctx gets a top-level "x_request_id".
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID returns the request ID stored in ctx, or "" if there is none.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// ParseLevel turns "debug", "info", "warn", or "error" into a slog.Level.
func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(s))); err != nil {
		return slog.LevelInfo, fmt.Errorf("logger: unknown level %q (use debug, info, warn, error)", s)
	}
	return level, nil
}

// OpenFile opens (or creates) the log file, creating its folder if needed.
// The caller must Close the returned file.
func OpenFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
}
```

**What this does:**
1. **Payload dan SDC:** keduanya tipe map. Keduanya punya method `apply`, sehingga keduanya otomatis memenuhi interface `Extra`. Karena itu `Info(ctx, msg, extras ...Extra)` bisa menerima keduanya.
2. **Kenapa SDC hanya string:** di Elasticsearch, tipe sebuah field ditentukan oleh nilai pertama yang masuk. Misalnya `user_id` pertama kali masuk sebagai angka; begitu kemudian ada `"abc"`, log itu **ditolak** karena *mapping conflict*. Kalau SDC dipaksa string sejak di Go, masalah ini tidak mungkin terjadi.
3. **Rename key:** `renameKeys` mengganti key bawaan slog `time` dan `msg` menjadi `timestamp` (UTC, presisi milidetik) dan `message`.
4. **Level:** `log()` langsung berhenti kalau level-nya berada di bawah level minimum. Contoh: saat `LOG_LEVEL=info`, panggilan `log.Debug(...)` tidak mencetak apa-apa.
5. **Context:** `WithRequestID` dan `WithSDC` menyimpan nilai di dalam `context`. Setiap log yang memakai `ctx` itu otomatis membawa `x_request_id` dan isi SDC-nya.

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `type SDC map[string]string` | Tipe baru berbasis map, bisa diberi method sendiri | `type SDC = Record<string, string>`, tapi bisa punya method |
| `type Extra interface { apply(e *entry) }` | **Interface**: "tipe apa pun yang punya method `apply`". Tidak perlu menulis `implements` | `interface Extra { apply(e: Entry): void }`, terpenuhi otomatis |
| `extras ...Extra` | *Variadic*: boleh mengirim 0, 1, atau banyak argumen | `(...extras: Extra[])` |
| `func (l *Logger) Info(...)` | Method dengan *pointer receiver*: `l` menunjuk ke Logger yang sama, bukan salinannya | Method di class |
| `*slog.Logger`, `&Logger{...}` | `*` = pointer (alamat object). `&` = "ambil alamatnya" | Mengoper object by reference |
| `maps.Copy(dst, src)` | Salin semua key dari `src` ke `dst` | `Object.assign(dst, src)` |
| `context.WithValue(ctx, key, val)` | Membuat context **baru** yang membawa nilai. Context lama tidak berubah | Mirip React Context Provider |
| `type sdcKey struct{}` | Tipe kosong yang dipakai sebagai key unik di context | `Symbol()` |
| `x.(SDC)` | *Type assertion*: "anggap nilai ini bertipe SDC". Kalau gagal, hasilnya nilai kosong | `x as SDC`, tapi dicek saat runtime |
| `ReplaceAttr` | Hook yang mengubah setiap key dan value sebelum ditulis | Mirip `JSON.stringify(obj, replacer)` |

### 3b. `pretty.go`

**File:** `internal/infrastructure/logger/pretty.go` (NEW)

**Why:** Saat development lokal, JSON satu baris susah dibaca. `NewPrettyWriter` mencetaknya dengan indentasi. Pakai **hanya untuk output terminal**. File log tetap harus satu baris per entri, karena Filebeat membaca file itu baris per baris.

```go
package logger

import (
	"bytes"
	"encoding/json"
	"io"
)

// prettyWriter re-indents each JSON log line so it is easy to read in a
// terminal. Use it for local development only: it is slower, and the log
// file must stay one line per entry for Filebeat.
type prettyWriter struct {
	w io.Writer
}

// NewPrettyWriter wraps w so JSON log lines are printed indented.
func NewPrettyWriter(w io.Writer) io.Writer {
	return prettyWriter{w: w}
}

func (p prettyWriter) Write(line []byte) (int, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, bytes.TrimSpace(line), "", "  "); err != nil {
		return p.w.Write(line) // not JSON: print it unchanged
	}
	out.WriteByte('\n')
	if _, err := p.w.Write(out.Bytes()); err != nil {
		return 0, err
	}
	return len(line), nil
}
```

**Output sungguhan** dengan pretty writer:
```json
{
  "timestamp": "2026-11-08T15:49:22.698Z",
  "level": "ERROR",
  "message": "database ping failed",
  "service": "bang-backend",
  "env": "local",
  "version": "dev",
  "x_request_id": "b7f2",
  "payload": {
    "error": "connection refused"
  },
  "sdc": {
    "correlation_id": "12121211"
  }
}
```

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `Write(line []byte) (int, error)` | Method yang membuat `prettyWriter` jadi sebuah `io.Writer` | Mengimplementasikan interface `Writable` |
| `json.Indent(&out, data, "", "  ")` | Merapikan JSON tanpa mengubah urutan key | `JSON.stringify(JSON.parse(s), null, 2)` |
| `bytes.Buffer` | Penampung byte di memori | Array string yang di-`join` |

### 3c. `logger_test.go`

**File:** `internal/infrastructure/logger/logger_test.go` (NEW)

```go
package logger

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var testOpts = Options{Level: slog.LevelInfo, Service: "bang-backend", Env: "test", Version: "dev"}

// readLines parses every JSON line written to buf.
func readLines(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var lines []map[string]any
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		var line map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			t.Fatalf("log line is not JSON: %v\n%s", err, scanner.Text())
		}
		lines = append(lines, line)
	}
	return lines
}

func TestLog_Shapes(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, testOpts)
	ctx := context.Background()

	// 1. message only
	log.Info(ctx, "message")
	// 2. message + payload
	log.Info(ctx, "message", Payload{"user_id": 1, "sample_message": "asd asd"})
	// 3. message + empty payload + sdc
	log.Info(ctx, "message", Payload{}, SDC{"correlation_id": "12121211", "path": "/usr/asd/asd"})

	lines := readLines(t, &buf)
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}

	// 1. no payload, no sdc
	if _, ok := lines[0]["payload"]; ok {
		t.Errorf("line 1 must not have payload: %v", lines[0])
	}
	if _, ok := lines[0]["sdc"]; ok {
		t.Errorf("line 1 must not have sdc: %v", lines[0])
	}

	// 2. payload nested under "payload"
	payload, _ := lines[1]["payload"].(map[string]any)
	if payload["user_id"] != float64(1) || payload["sample_message"] != "asd asd" {
		t.Errorf("line 2 payload = %v", lines[1]["payload"])
	}

	// 3. empty payload is omitted, sdc nested under "sdc"
	if _, ok := lines[2]["payload"]; ok {
		t.Errorf("line 3 empty payload must be omitted: %v", lines[2])
	}
	sdc, _ := lines[2]["sdc"].(map[string]any)
	if sdc["correlation_id"] != "12121211" || sdc["path"] != "/usr/asd/asd" {
		t.Errorf("line 3 sdc = %v", lines[2]["sdc"])
	}

	// every line has the fixed fields
	for i, line := range lines {
		if line["service"] != "bang-backend" || line["env"] != "test" || line["version"] != "dev" {
			t.Errorf("line %d missing service/env/version: %v", i+1, line)
		}
	}
}

func TestLog_KeyNamesAndTimestamp(t *testing.T) {
	var buf bytes.Buffer
	New(&buf, testOpts).Info(context.Background(), "hello")

	line := readLines(t, &buf)[0]
	if line["message"] != "hello" {
		t.Errorf("message = %v, want hello", line["message"])
	}
	for _, old := range []string{"msg", "time"} {
		if _, ok := line[old]; ok {
			t.Errorf("old key %q must be renamed: %v", old, line)
		}
	}

	ts, _ := line["timestamp"].(string)
	if !strings.HasSuffix(ts, "Z") {
		t.Errorf("timestamp %q must be UTC (end with Z)", ts)
	}
	if _, err := time.Parse(timestampFormat, ts); err != nil {
		t.Errorf("timestamp %q has the wrong format: %v", ts, err)
	}
}

func TestLog_Levels(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, testOpts) // minimum level: info
	ctx := context.Background()

	log.Debug(ctx, "hidden")
	log.Info(ctx, "info")
	log.Warn(ctx, "warn")
	log.Error(ctx, "error")

	lines := readLines(t, &buf)
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3 (debug must be dropped)", len(lines))
	}
	for i, want := range []string{"INFO", "WARN", "ERROR"} {
		if lines[i]["level"] != want {
			t.Errorf("line %d level = %v, want %s", i+1, lines[i]["level"], want)
		}
	}
}

func TestContext_RequestIDAndSDC(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, testOpts)

	ctx := WithRequestID(context.Background(), "b7f2")
	ctx = WithSDC(ctx, SDC{"user_id": "7", "role": "technician"})

	// a value passed to the call overrides the same key from context
	log.Info(ctx, "work order started", SDC{"work_order_id": "12", "role": "admin"})
	log.Info(context.Background(), "no context")

	lines := readLines(t, &buf)
	if lines[0]["x_request_id"] != "b7f2" {
		t.Errorf("x_request_id = %v, want b7f2", lines[0]["x_request_id"])
	}
	sdc, _ := lines[0]["sdc"].(map[string]any)
	if sdc["user_id"] != "7" || sdc["work_order_id"] != "12" || sdc["role"] != "admin" {
		t.Errorf("sdc = %v, want user_id=7, work_order_id=12, role=admin", sdc)
	}
	if _, ok := lines[1]["x_request_id"]; ok {
		t.Errorf("line without context must not have x_request_id: %v", lines[1])
	}
}

func TestPrettyWriter(t *testing.T) {
	var buf bytes.Buffer
	New(NewPrettyWriter(&buf), testOpts).Info(context.Background(), "pretty")

	out := buf.String()
	if !strings.Contains(out, "\n  \"level\": \"INFO\"") {
		t.Errorf("output is not indented:\n%s", out)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in      string
		want    slog.Level
		wantErr bool
	}{
		{in: "debug", want: slog.LevelDebug},
		{in: "info", want: slog.LevelInfo},
		{in: "WARN", want: slog.LevelWarn},
		{in: "error", want: slog.LevelError},
		{in: "verbose", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseLevel(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestOpenFile_CreatesFolder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "app.log")

	f, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile() error: %v", err)
	}
	defer f.Close()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("log file was not created: %v", err)
	}
}
```

Beberapa hal di test ini yang mungkin baru:
- `t.Helper()` membuat error di dalam `readLines` dilaporkan di baris yang memanggilnya.
- Angka hasil parse JSON di Go selalu bertipe `float64`. Karena itu ada perbandingan `float64(1)`.
- `_, ok := m["key"]` adalah cara mengecek apakah sebuah key ada di map.

**Check:**
```bash
go test -v ./internal/infrastructure/logger
```
```
--- PASS: TestLog_Shapes (0.00s)
--- PASS: TestLog_KeyNamesAndTimestamp (0.00s)
--- PASS: TestLog_Levels (0.00s)
--- PASS: TestContext_RequestIDAndSDC (0.00s)
--- PASS: TestPrettyWriter (0.00s)
--- PASS: TestParseLevel (0.00s)
--- PASS: TestOpenFile_CreatesFolder (0.00s)
PASS
```

Untuk menjalankan satu test saja: `go test -run TestLog_Shapes -v ./internal/infrastructure/logger`

**Common mistakes:**
- **Mengisi SDC dengan angka.** Contoh `logger.SDC{"correlation_id": 12121211}` memunculkan error `cannot use 12121211 (untyped int constant) as string value in map literal`. Bungkus angkanya jadi string: `"12121211"`. Untuk variabel angka, pakai `strconv.Itoa(id)`.
- **Mengirim map biasa.** Contoh `log.Info(ctx, "m", map[string]any{"a": 1})` memunculkan error `map[string]any does not implement logger.Extra (missing method apply)`. Pakai `logger.Payload{"a": 1}` atau `logger.SDC{...}`.
- **Menaruh data sensitif di `payload`.** Walaupun tidak di-index, isinya tetap tersimpan di log. Jangan masukkan password, token, atau nomor HP lengkap ke sana.

---

## Step 4: Package `database`

Download dulu driver Postgres-nya:
```bash
go get github.com/jackc/pgx/v5
```
```
go: added github.com/jackc/pgx/v5 v5.12.0
```
Perintah ini mengubah `go.mod` dan membuat `go.sum`. Mirip `npm install` yang mengubah `package.json` dan `package-lock.json`.

**File:** `internal/infrastructure/database/database.go` (NEW, layer: `infrastructure`)

**Why:** Satu tempat untuk membuka *connection pool*, yaitu sekumpulan koneksi DB yang dipakai bergantian supaya tidak perlu membuka koneksi baru di setiap request. Saat start, kode langsung melakukan **ping** untuk memastikan database benar-benar bisa dihubungi.

```go
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a Postgres connection pool and checks it with a ping.
// The caller must call pool.Close() when the app shuts down.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("database: parse url: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("database: create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	return pool, nil
}
```

**What this does:**
1. Mem-parse URL database lalu mengatur pool (maksimal 10 koneksi).
2. Membuat pool.
3. Melakukan ping dengan batas waktu 5 detik. Kalau gagal, pool ditutup lagi dan fungsi mengembalikan error yang menjelaskan di tahap mana kegagalannya terjadi.

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `context.Context` | Membawa **deadline dan sinyal batal** ke operasi I/O. Selalu jadi parameter pertama | `AbortController.signal` |
| `context.WithTimeout(ctx, 5*time.Second)` | Context yang otomatis batal setelah 5 detik | `AbortSignal.timeout(5000)` |
| `defer cancel()` | Lepaskan timer saat fungsi selesai. Ini wajib | Membersihkan timeout |
| `fmt.Errorf("...: %w", err)` | Membungkus error dengan konteks tambahan. `%w` tetap menyimpan error aslinya | `new Error("ping failed", { cause: err })` |
| Import `github.com/...` | Package pihak ketiga, di-download lewat `go get` | `import x from "pkg"` setelah `npm install` |

### `database_test.go`

**File:** `internal/infrastructure/database/database_test.go` (NEW)

```go
package database

import (
	"context"
	"os"
	"testing"
)

func TestConnect_InvalidURL(t *testing.T) {
	_, err := Connect(context.Background(), "this is not a url")
	if err == nil {
		t.Fatal("expected an error for an invalid URL, got nil")
	}
}

func TestConnect_Unreachable(t *testing.T) {
	// Port 1 has nothing listening, so the ping must fail quickly.
	url := "postgres://app:app@127.0.0.1:1/fieldservice?sslmode=disable&connect_timeout=1"

	_, err := Connect(context.Background(), url)
	if err == nil {
		t.Fatal("expected an error for an unreachable database, got nil")
	}
}

// Runs only when TEST_DATABASE_URL is set: locally against docker compose,
// and in CI against the Postgres service container (see Step 4b).
func TestConnect_RealDatabase(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		// GitHub Actions always sets CI=true. Fail instead of skipping there,
		// so a broken CI setup can't hide this test.
		if os.Getenv("CI") == "true" {
			t.Fatal("TEST_DATABASE_URL must be set in CI")
		}
		t.Skip("TEST_DATABASE_URL not set")
	}

	pool, err := Connect(context.Background(), url)
	if err != nil {
		t.Fatalf("Connect() error: %v", err)
	}
	pool.Close()
}
```

- Dua test pertama **tidak butuh database**, jadi selalu jalan, di laptop maupun di CI.
- Test ketiga butuh Postgres sungguhan, jadi hanya jalan kalau `TEST_DATABASE_URL` di-set:
  - **Di laptop:** kalau tidak di-set, test-nya dilewati (SKIP). `make check` tetap bisa jalan tanpa Docker.
  - **Di CI:** GitHub Actions menyalakan Postgres dan men-set `TEST_DATABASE_URL` (lihat Step 4b). Kalau variabel itu lupa di-set, test-nya **gagal**, bukan SKIP. Tanpa pengecekan `CI`, salah setup di workflow akan membuat test ini diam-diam tidak pernah jalan, padahal CI tetap hijau.

**Check (tanpa DB):**
```bash
go test -v ./internal/infrastructure/database
```
Hasilnya: `TestConnect_InvalidURL` dan `TestConnect_Unreachable` `PASS`, sedangkan `TestConnect_RealDatabase` `SKIP` dengan pesan `TEST_DATABASE_URL not set`.

**Check (dengan DB dari Step 1):**
```bash
TEST_DATABASE_URL="postgres://app:app@localhost:5432/fieldservice?sslmode=disable" \
  go test -v -run TestConnect_RealDatabase ./internal/infrastructure/database
```
Hasilnya harus `--- PASS`. Check ini belum diuji di project percobaan, karena Docker belum menyala saat plan ini ditulis.

### 4b. Jalankan test database di CI

**File:** `.github/workflows/ci.yml` (EDIT, hanya job `Test`)

**Why:** Mesin CI (runner GitHub) tidak punya Postgres, dan `docker compose` di laptop kamu tidak ikut jalan di sana. GitHub Actions punya fitur **service container**: sebelum step pertama dijalankan, GitHub menyalakan container lain di samping job, dan container itu dimatikan lagi setelah job selesai. Kita pakai image yang sama dengan `docker-compose.yml` (`postgres:16`) supaya hasil test di CI sama dengan di laptop.

Analogi frontend: mirip menjalankan mock server (misalnya `json-server`) di `beforeAll` lalu mematikannya di `afterAll`, tapi yang mengatur hidup-matinya GitHub, bukan kode test.

Ganti job `test` di `ci.yml` menjadi seperti ini. Job `lint` tidak berubah, dan nama job tetap `Test` karena ruleset di `main` mewajibkan nama itu.

```yaml
  test:
    name: Test # required by the main ruleset; keep this name
    runs-on: ubuntu-latest

    # A throwaway Postgres for tests that need a real database.
    # Same image as docker-compose.yml, so CI matches local.
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: app
          POSTGRES_PASSWORD: app
          POSTGRES_DB: fieldservice_test
        ports:
          - 5432:5432
        # Wait until Postgres accepts connections before the steps start
        options: >-
          --health-cmd "pg_isready -U app -d fieldservice_test"
          --health-interval 5s
          --health-timeout 5s
          --health-retries 10

    env:
      # Read by TestConnect_RealDatabase (and later by repository tests)
      TEST_DATABASE_URL: postgres://app:app@localhost:5432/fieldservice_test?sslmode=disable

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: go build
        run: go build -o /dev/null ./...

      - name: go test
        run: go test -race -count=1 -cover ./...
```

**What this does:**
1. `services.postgres` menyalakan Postgres 16 dan membuka port 5432 di runner, jadi test bisa connect ke `localhost:5432`.
2. `options: --health-cmd ...` membuat GitHub **menunggu** sampai Postgres siap menerima koneksi. Tanpa ini, step `go test` bisa mulai lebih dulu dan gagal dengan `connection refused`.
3. `env.TEST_DATABASE_URL` di level job berlaku untuk semua step, jadi `go test` ikut membacanya.
4. Nama database-nya `fieldservice_test`, bukan `fieldservice`. Kebiasaan ini berguna nanti: database untuk test terpisah dari database untuk development, supaya test yang menghapus data tidak menyentuh data lokal kamu.

| Konsep | Artinya | Di frontend |
|---|---|---|
| `services:` | Container pendamping yang hidup selama job berjalan | Mock server di `beforeAll` / `afterAll` |
| `--health-cmd` | Perintah untuk mengecek apakah container sudah siap | `await waitOn("tcp:5432")` |
| `env:` di level job | Environment variable untuk semua step | `process.env` di `.env.test` |
| `CI=true` | Di-set otomatis oleh GitHub Actions | `process.env.CI` yang dibaca Jest |

**Password `app:app` aman ditulis di sini?** Ya. Database ini hanya hidup beberapa menit di mesin CI, tidak bisa diakses dari luar, dan isinya data test. Secret sungguhan, misalnya `JWT_SECRET` production, tetap harus disimpan di GitHub Secrets dan tidak boleh ditulis di file workflow.

**Check (di CI):** push branch lalu buka tab Actions di PR. Di log step `go test`, cari baris `ok  .../internal/infrastructure/database`. Untuk memastikan test DB benar-benar jalan (bukan SKIP), jalankan `go test -v` sekali di CI atau lihat log-nya: harus ada `--- PASS: TestConnect_RealDatabase`. Check ini belum diuji, karena workflow-nya baru akan jalan setelah kamu push.

**Check (di laptop, meniru CI):**
```bash
docker compose up -d
CI=true TEST_DATABASE_URL="postgres://app:app@localhost:5432/fieldservice?sslmode=disable" \
  go test -v ./internal/infrastructure/database
```
Ketiga test harus `PASS`. Lalu coba tanpa `TEST_DATABASE_URL`:
```bash
CI=true go test -v -run TestConnect_RealDatabase ./internal/infrastructure/database
```
Hasilnya harus `FAIL` dengan pesan `TEST_DATABASE_URL must be set in CI`. Ini membuktikan pengamannya bekerja.

**Kalau branch kamu sudah di-rebase ke `main` setelah PR #27 di-merge:** job `lint` di `ci.yml` sudah memakai `golangci-lint-action`. Biarkan job itu apa adanya, dan cukup ganti job `test`.

**Alternatif yang tidak dipilih:** library [testcontainers-go](https://golang.testcontainers.org/) bisa menyalakan Postgres dari dalam kode test, jadi tidak perlu setup apa pun di workflow. Kekurangannya, Docker wajib menyala setiap kali `go test` dijalankan di laptop, dan konsepnya lebih banyak untuk dipelajari sekarang. Service container sudah cukup untuk project ini.

**Common mistakes:**
- **Lupa `go get`.** Muncul error `no required module provides package github.com/jackc/pgx/v5/pgxpool; to add it: go get github.com/jackc/pgx/v5/pgxpool`. Jalankan saja perintah yang disarankan di pesan error itu.
- **Lupa `defer cancel()`.** `go vet` akan memperingatkan `the cancel function is not used on all paths (possible context leak)`.

---

## Setelah Step 1–4

```bash
make check     # gofmt, vet, build, test → semua lolos
```

Ini titik yang pas untuk commit. Usulan rencana commit (jalankan `/atomic-conventional-commits`):
- `build(setup): add docker compose for postgres` → `docker-compose.yml`
- `feat(infra): add config package` → `.env.example`, `config/`
- `feat(infra): add json logger with payload and sdc` → `logger/`
- `feat(infra): add postgres connection pool` → `database/`, `go.mod`, `go.sum`
- `ci(github): run database tests against postgres service` → `.github/workflows/ci.yml`
- `docs(infra): add plan for #2` → file plan ini

## Next batch (Step 5–6)

- **Step 5: `cmd/api/main.go` + graceful shutdown.** Urutannya: load config → `logger.ParseLevel` → buka file log → buat logger. Output terminal memakai `NewPrettyWriter` saat `APP_ENV=local`, dan file log tetap satu baris per entri, digabung lewat `io.MultiWriter`. Setelah itu connect DB, lalu tunggu Ctrl+C dengan `signal.NotifyContext`. Saat Ctrl+C ditekan, log "shutting down", tutup DB, dan exit 0. `main_test.go` juga diperbarui.
- **Step 6: `Makefile` + `.air.toml`.** Target `run` (otomatis membaca `.env`), `dev` (hot reload dengan air), `test`, dan `lint`. Target `migrate-up/down`, `sqlc`, dan `seed` dibuat sebagai placeholder yang menunjuk ke #5 dan #7.

Kalau Step 1–4 sudah selesai, bilang **"lanjut step 5"**. Bagian ini akan diisi di file yang sama.

## Definition of Done (issue #2)
- [ ] `docker-compose.yml`: Postgres 16 dengan volume `pgdata`; `docker compose up -d` berjalan
- [ ] `.env.example` berisi `PORT`, `DATABASE_URL`, `JWT_SECRET`, `LOG_FILE`, `LOG_LEVEL`, `APP_ENV`, `APP_VERSION`
- [ ] `config.Load()` membaca env beserta default-nya, dan gagal dengan error yang jelas kalau `JWT_SECRET` kosong
- [ ] `os.Getenv` hanya dipanggil di `internal/infrastructure/config`
- [ ] Logger JSON menulis ke stdout **dan** `logs/app.log`; file log satu baris per entri
- [ ] Level `debug` / `info` / `warn` / `error`; level minimum diatur lewat `LOG_LEVEL` (nilai yang tidak dikenal → error yang jelas)
- [ ] Setiap baris log punya `timestamp` (UTC, milidetik), `level`, `message`, `service`, `env`, `version`
- [ ] `log.Info(ctx, msg, Payload{...}, SDC{...})`: `payload` tidak di-index, `sdc` di-index, nilai SDC berupa string
- [ ] `WithRequestID` menghasilkan `x_request_id` di level teratas; `WithSDC` menambahkan SDC dari context
- [ ] Pretty print hanya dipakai untuk terminal saat `APP_ENV=local`
- [ ] Secret (password DB, JWT secret) tidak pernah muncul di log, dan ada test-nya
- [ ] `database.Connect` membuka `pgxpool` dan melakukan ping saat start; error-nya jelas kalau DB mati
- [ ] Graceful shutdown: Ctrl+C → log "shutting down" → pool DB ditutup → exit 0
- [ ] `Makefile`: `run`, `dev`, `test`, `lint` berfungsi; `migrate-up/down`, `sqlc`, `seed` berupa placeholder yang menunjuk ke #5 dan #7
- [ ] Ada test untuk config, logger, dan database; test database yang butuh Postgres ikut jalan di CI lewat service container (`TestConnect_RealDatabase` `PASS`, bukan SKIP)
- [ ] `make check` lolos
- [ ] Commit memakai format `<type>(<feature>): <message>`; body PR berisi `Closes #2`
- [ ] PR dibuat dengan template; CI Lint dan Test hijau

## Notes

- **Inspirasi desain logger:** sebagian ide diambil dari library logging internal tim lama (dipakai sebagai referensi, kodenya tidak disalin):
  - SDC hanya berisi string
  - `x_request_id` di level teratas
  - key `timestamp` dan `message`
  - pretty writer untuk lokal
  - error code masuk ke SDC (rencana untuk nanti)
- **Rencana di #3:** saat format error `{ "error": { "code", "message" } }` dibuat, error handler akan menambahkan `SDC{"error_code": code}` ke log error, supaya error bisa dicari berdasarkan kodenya di Kibana.
- **Payload tidak di-index:** keputusan sebuah field di-index atau tidak ada di Elasticsearch, bukan di Go. Saat setup observability nanti (Sesi 6 di spec), tambahkan index template berikut. Template ini belum diuji karena Elasticsearch belum ada:
  ```json
  "mappings": {
    "properties": {
      "payload": { "type": "object", "enabled": false },
      "sdc":     { "type": "object" }
    }
  }
  ```

## When you're stuck
- "Aku dapat error `<paste error>` di Step 3. Jelaskan penyebabnya dulu, jangan langsung kasih kode."
- "Jelaskan bedanya `Logger` dan `*Logger` dengan analogi React."
- "Kenapa `context` harus jadi parameter pertama? Kasih contoh kalau tidak dipakai."

## Changelog
- 2026-11-08: plan dibuat (Step 1–4 lengkap, Step 5–6 masih berupa judul)
- 2026-11-08: logger ditambah level (`LOG_LEVEL`), field `service`/`env`/`version`, dan API `Payload` / `SDC`
- 2026-11-08: logger mengambil 5 ide: SDC hanya string, `x_request_id` di level teratas, key `timestamp`/`message`, pretty writer, error code ke SDC (rencana di #3)
- 2026-11-21: Step 4b ditambahkan: Postgres service container di CI dan `TEST_DATABASE_URL`, plus pengaman `CI=true` supaya test DB gagal (bukan SKIP) kalau setup CI salah
