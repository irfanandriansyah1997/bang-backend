# Plan: Issue #3: HTTP base: router, middleware, error format, cursor pagination
> Issue: https://github.com/irfanandriansyah1997/bang-backend/issues/3 · Plan branch: `docs/3-http-plan` (lapisan paling bawah dari stack; setiap step punya branch sendiri di atasnya) · Last updated: 2026-11-22

**Goal:** Siapkan fondasi layer HTTP supaya semua endpoint nanti berperilaku sama: satu format error, satu cara menulis JSON, satu format list (cursor), dan satu `request_id` di setiap baris log. Setelah task ini, endpoint baru tinggal menulis handler tipis.

**Depends on:** #2 ✅ (sudah di-merge: logger, config, dan database sudah ada).
**Already done in the repo:** `logger.WithRequestID` dan `logger.WithSDC` (dari #2) yang akan dipakai middleware. Selebihnya belum ada: `internal/delivery/http/` baru berisi `README.md`.

**Before you start:**
- File plan ini ada di branch `docs/3-http-plan` (PR #36). Semua branch step dibuat **di atas** branch ini, bukan dari `main`, sebagai stacked PR (lihat "Cara kerja").
- Pasang ekstensi `gh stack` sekali saja (ekstensi resmi GitHub, butuh `gh` 2.0+ dan Git 2.36+):
  ```bash
  gh extension install github/gh-stack
  ```
- Docker dan database hanya dibutuhkan untuk tes terakhir di Step 8. Library baru hanya satu: `validator` di Step 7.

Semua kode di langkah 1–8 sudah di-compile dan di-test di salinan project: `go vet` bersih, `golangci-lint` `0 issues`, test lolos dengan `-race` (termasuk test yang menyalakan server dan memanggilnya lewat HTTP), coverage `internal/delivery/http` 96% dan `internal/domain` 100%. Server sungguhan juga sudah dicoba dengan `curl` dan Ctrl+C.

## Cara kerja: 1 step = 1 PR, ditumpuk (stacked PR)

Setiap step adalah satu PR kecil yang bisa di-review sendiri, dan semuanya **ditumpuk**: branch Step 1 dibuat di atas branch plan, branch Step 2 di atas Step 1, dan seterusnya. Setiap PR memakai branch di bawahnya sebagai *base*, jadi reviewer (kamu maupun Claude) hanya melihat perubahan step itu, bukan semua step sebelumnya.

```
feat/3-8-router-server        → PR (base: feat/3-7-decode-validate)   ← atas
...
feat/3-2-responder            → PR (base: feat/3-1-domain-errors)
feat/3-1-domain-errors        → PR (base: docs/3-http-plan)
docs/3-http-plan              → PR #36 (base: main)                   ← bawah
main
```

Keuntungannya dibanding menunggu merge satu per satu: kamu bisa lanjut ke step berikutnya tanpa menunggu review, dan tiap PR tetap kecil. `gh stack` mengurus bagian yang membosankan: membuat branch di atas branch sebelumnya, merebase semuanya kalau ada yang berubah, dan mengatur base tiap PR.

**Setup (sekali, setelah ekstensi terpasang):**
```bash
git switch docs/3-http-plan && git pull
gh stack init docs/3-http-plan feat/3-1-domain-errors   # menyusun stack; branch yang sudah ada langsung dipakai
gh stack view                                            # lihat susunannya
```
Kalau `feat/3-1-domain-errors` belum ada, perintah di atas membuatnya di atas branch plan. Branch plan menjadi lapisan paling bawah, dan PR #36 yang sudah ada jadi PR pertama di stack.

**Ulangi untuk setiap step** (kamu berada di branch teratas stack):
```bash
# ketik kode step itu, lalu jalankan bagian "Check"
make check                                   # harus lolos sebelum commit
/atomic-conventional-commits                 # commit (rencana commit ada di bagian "PR" tiap step)
gh stack add feat/3-<N+1>-<slug>             # branch untuk step berikutnya, di atas yang sekarang
```
Perintah `gh stack add` hanya dijalankan **setelah** step itu di-commit, dan harus dari branch teratas. Pindah antar lapisan dengan `gh stack up`, `gh stack down`, `gh stack top`, atau `gh stack switch`.

**Mengirim ke GitHub** (setelah beberapa step selesai, atau semuanya):
```bash
gh stack push                                # push semua branch sekaligus
gh stack submit                              # satu PR per branch, base = branch di bawahnya, ditautkan sebagai Stack
```
`gh stack submit` membuka editor di terminal untuk judul dan deskripsi tiap PR; pakai judul dari bagian "PR" tiap step. PR dibuat sebagai draft kecuali kamu memilih ready di editor atau memakai `--open`. Untuk menulis ulang body sesuai template, jalankan `/create-pr` di tiap lapisan (ia mengedit PR yang sudah ada).

**Review dan merge:**
1. Setelah CI hijau, comment `@claude-code-review` di tiap PR. Review hanya mencakup perubahan lapisan itu.
2. Kalau ada temuan di Step 3, pindah ke branchnya (`gh stack switch`), commit perbaikannya, lalu `gh stack rebase --upstack` supaya perbaikan itu turun ke step 4 sampai 8, dan `gh stack push`.
3. Merge dari bawah: PR #36 dulu, lalu satu per satu ke atas, atau sekaligus dengan `gh stack merge` (memilih sampai PR mana). Setelah ada yang di-merge, jalankan `gh stack sync`: ia merebase sisa stack ke `main` dan mengatur ulang base PR.

Aturan:
- Step berikutnya memakai kode step sebelumnya (Step 2 memakai error dari Step 1, dan seterusnya), jadi urutan lapisan tidak boleh dibalik.
- Semua PR memakai `Refs #3`. Hanya PR Step 8 memakai `Closes #3`.
- Issue #3 tetap `In progress` sampai Step 8 di-merge.
- Kalau `gh stack` bermasalah, tumpukan ini tetap bisa dikerjakan tanpa ekstensi: buat branch dengan `git switch -c` dari branch sebelumnya, dan atur base PR dengan `gh pr create --base <branch-di-bawahnya>`.

## Progress
Centang kotak setelah PR step itu **di-merge** ke `main`. Branch di bawah ini membentuk stack dari bawah ke atas, di atas `docs/3-http-plan`.
- [ ] Step 1: Error bisnis di `domain` · branch `feat/3-1-domain-errors`
- [ ] Step 2: `Responder`: satu tempat menulis JSON dan error, plus mapper domain error → HTTP · `feat/3-2-responder`
- [ ] Step 3: Helper cursor pagination · `feat/3-3-cursor-pagination`
- [ ] Step 4: Middleware `X-Request-ID` dan `X-Client` · `feat/3-4-request-id`
- [ ] Step 5: Middleware log per request dan panic recovery · `feat/3-5-request-log-recover`
- [ ] Step 6: Middleware CORS · `feat/3-6-cors`
- [ ] Step 7: Decode body dan validasi (library `validator`) · `feat/3-7-decode-validate`
- [ ] Step 8: Router `/api/v1`, rantai middleware, dan server di `main.go` · `feat/3-8-router-server`

## What you'll build

```
internal/domain/errors.go                       NEW   (Step 1)
internal/domain/errors_test.go                  NEW   (Step 1)
internal/delivery/http/response.go              NEW   (Step 2)
internal/delivery/http/errors.go                NEW   (Step 2)
internal/delivery/http/helpers_test.go          NEW   (Step 2, dipakai semua test di folder ini)
internal/delivery/http/response_test.go         NEW   (Step 2)
internal/delivery/http/errors_test.go           NEW   (Step 2)
internal/delivery/http/cursor.go                NEW   (Step 3)
internal/delivery/http/cursor_test.go           NEW   (Step 3)
internal/delivery/http/request_id.go            NEW   (Step 4)
internal/delivery/http/client.go                NEW   (Step 4)
internal/delivery/http/request_id_test.go       NEW   (Step 4)
internal/delivery/http/logging.go               NEW   (Step 5)
internal/delivery/http/recover.go               NEW   (Step 5)
internal/delivery/http/logging_test.go          NEW   (Step 5)
internal/delivery/http/recover_test.go          NEW   (Step 5)
internal/delivery/http/cors.go                  NEW   (Step 6)
internal/delivery/http/cors_test.go             NEW   (Step 6)
internal/delivery/http/decode.go                NEW   (Step 7)
internal/delivery/http/decode_test.go           NEW   (Step 7)
go.mod, go.sum                                  EDIT  (Step 7, otomatis lewat go get)
internal/delivery/http/router.go                NEW   (Step 8)
internal/delivery/http/router_test.go           NEW   (Step 8)
internal/delivery/http/README.md                EDIT  (Step 8)
cmd/api/main.go                                 EDIT  (Step 8)
cmd/api/main_test.go                            EDIT  (Step 8)
```

## Keputusan desain (kenapa begini)

- **Satu paket `http`, file datar, tanpa subfolder `middleware/`.** README folder menyarankan `middleware/`, tapi middleware `Recover` perlu menulis error lewat `Responder`. Kalau `middleware/` jadi paket terpisah dan `http` mengimpornya (untuk router), lalu `middleware/` mengimpor `http` (untuk `Responder`), Go menolak karena **import cycle** (A mengimpor B, B mengimpor A). Satu paket menghilangkan masalah itu.
- **Nama paket `http` sama dengan `net/http`.** Itu legal: di dalam paket kita, `http.StatusOK` tetap berarti `net/http`. Nama paket kita sendiri tidak masuk scope file.
- **`domain.ValidationError` ada di domain.** Aturan seperti "menyelesaikan work order butuh minimal 1 foto" adalah aturan bisnis dan harus berujung `422 validation_failed` (lihat spec). Domain tidak boleh tahu HTTP, jadi ia hanya mengembalikan error bertipe; yang memetakan ke 422 adalah mapper di layer HTTP.
- **Urutan middleware (dari luar ke dalam):** CORS → RequestID → Client → RequestLog → Recover → router. `RequestLog` harus membungkus langsung `Recover` dan router, dan meneruskan `*http.Request` yang **sama**, karena router mengisi `r.Pattern` (nama route) pada request yang ia terima. Middleware yang membuat salinan request (seperti `RequestID`) membuat Pattern tidak terlihat oleh yang di luarnya. Saya sudah mengujinya: Pattern hanya terbaca di request yang sama.
- **Jawaban `404` untuk route yang tidak ada dibuat sendiri** (Step 8), karena `ServeMux` bawaan menjawab dengan teks biasa. Konsekuensinya method yang salah juga menjadi `404`, bukan `405`.
- **CORS paling luar** (Step 6), supaya preflight `OPTIONS` dijawab sebelum menyentuh router (tidak perlu route `OPTIONS` di tiap endpoint).
- **Satu step = satu PR.** Setiap step bisa di-compile dan lolos `make check` sendiri, dan tidak ada kode yang menggantung tanpa pemakai yang dipaksakan: tipe dan fungsi yang di-export tidak dianggap "tidak dipakai" oleh linter.
- **`Recover` di bagian dalam**, supaya panic berubah jadi respons 500 *sebelum* `RequestLog` mencatat statusnya (yang tercatat jadi `500`, bukan `200`).

## Steps

### Step 1: Error bisnis di `domain`

**File:** `internal/domain/errors.go` (NEW, layer: `domain`)
**Why:** Error seperti "tidak ditemukan" atau "status tidak boleh berubah begitu" adalah bahasa bisnis, jadi tempatnya di `domain`. Layer HTTP nanti memetakannya ke status code. Domain hanya memakai library standar (`errors`), tanpa tag `json`.

```go
package domain

import "errors"

// Business errors. A function returns one of these when a rule is broken.
// Callers check them with errors.Is, even when they are wrapped with fmt.Errorf("...: %w", err).
var (
	// ErrNotFound means the thing that was asked for does not exist.
	ErrNotFound = errors.New("not found")
	// ErrForbidden means the user is known but is not allowed to do this.
	ErrForbidden = errors.New("forbidden")
	// ErrInvalidTransition means a work order status cannot change that way.
	ErrInvalidTransition = errors.New("invalid status transition")
)

// ValidationError means the input broke a rule, for example
// "completing a work order needs at least one photo".
// Message is written for the person who sent the request.
type ValidationError struct {
	Message string
}

// Error makes *ValidationError an error.
func (e *ValidationError) Error() string {
	return "validation failed: " + e.Message
}
```

**What this does:** Tiga error "sentinel" (nilai error yang dibandingkan lewat `errors.Is`) dan satu tipe error `ValidationError` yang membawa pesan untuk pengirim request. Pemanggil tidak pernah membandingkan teks error, hanya nilainya.

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `var ErrNotFound = errors.New("not found")` | Nilai error yang dibuat sekali, lalu dibandingkan | Konstanta seperti `const NOT_FOUND = Symbol(...)` |
| `errors.Is(err, ErrNotFound)` | Apakah `err` adalah `ErrNotFound`, juga kalau sudah dibungkus `fmt.Errorf("...: %w", err)` | `err.cause` ditelusuri sampai ketemu |
| `type ValidationError struct{ Message string }` | Error buatan sendiri yang membawa data | `class ValidationError extends Error` |
| `func (e *ValidationError) Error() string` | Method ini membuat tipe tersebut memenuhi interface `error` | Mengimplementasi `toString()` |
| `errors.As(err, &ve)` | "Apakah ada `*ValidationError` di dalam rantai `err`? Kalau ya, isi `ve`" | `err instanceof ValidationError` |

**File:** `internal/domain/errors_test.go` (NEW)

```go
package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestBusinessErrors_SurviveWrapping(t *testing.T) {
	tests := []struct {
		name   string
		target error
	}{
		{name: "not found", target: ErrNotFound},
		{name: "forbidden", target: ErrForbidden},
		{name: "invalid transition", target: ErrInvalidTransition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := fmt.Errorf("find work order 12: %w", tt.target)

			if !errors.Is(wrapped, tt.target) {
				t.Errorf("errors.Is(%q, %q) = false, want true", wrapped, tt.target)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := fmt.Errorf("complete work order: %w", &ValidationError{Message: "at least one photo is required"})

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatal("errors.As did not find a *ValidationError")
	}

	if ve.Message != "at least one photo is required" {
		t.Errorf("Message = %q", ve.Message)
	}

	if got, want := ve.Error(), "validation failed: at least one photo is required"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
```

**Check:**
```bash
go test -v ./internal/domain
```
Hasilnya `--- PASS: TestBusinessErrors_SurviveWrapping`, `--- PASS: TestValidationError`, dan baris terakhir `ok  	.../internal/domain`.

**Common mistakes:**
- **Membandingkan dengan `==` setelah error dibungkus.** `err == ErrNotFound` gagal kalau `err` hasil `fmt.Errorf("...: %w", ErrNotFound)`. Pakai `errors.Is`. Linter `errorlint` akan menegur.
- **Lupa `&` saat memakai `errors.As`.** `errors.As(err, ve)` memberi `second argument to errors.As must be a non-nil pointer`. Tulis `errors.As(err, &ve)`.

**PR (step ini):**
- Branch: `feat/3-1-domain-errors` (dibuat di atas `docs/3-http-plan`)
- Base PR: `docs/3-http-plan` (PR plan; stack mulai dari sini)
- Commit:
  - `feat(http): add domain business errors` (`errors.go`, `errors_test.go`)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add domain business errors` · Body: `Refs #3`

---

### Step 2: `Responder`: satu tempat menulis JSON dan error

**Why:** Semua respons API (sukses maupun error) harus berbentuk sama. Kalau tiap handler menulis JSON sendiri, formatnya cepat berbeda-beda. `Responder` adalah satu-satunya tempat yang memanggil `WriteHeader` dan menulis body. Mapper error juga hanya ada di sini (sesuai `.ai/rules/http.md`): handler tidak pernah memilih status code untuk error domain.

**File:** `internal/delivery/http/response.go` (NEW, layer: `delivery/http`)

```go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// Responder writes every JSON response of the API: normal bodies and errors.
// Handlers get one Responder and never write JSON themselves.
type Responder struct {
	log *logger.Logger
}

// NewResponder creates a Responder that logs through log.
func NewResponder(log *logger.Logger) *Responder {
	return &Responder{log: log}
}

// JSON writes v as a JSON body with the given status code.
func (rs *Responder) JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	// Encode first. Once WriteHeader is called, the status can no longer change,
	// so a value that cannot be encoded must be found before that.
	body, err := json.Marshal(v)
	if err != nil {
		rs.log.Error(r.Context(), "encode response", logger.Payload{"error": err.Error()})
		status = http.StatusInternalServerError
		body = []byte(`{"error":{"code":"internal_error","message":"something went wrong"}}`)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(append(body, '\n')); err != nil {
		rs.log.Warn(r.Context(), "write response", logger.Payload{"error": err.Error()})
	}
}
```

**File:** `internal/delivery/http/errors.go` (NEW)

```go
package http

import (
	"errors"
	"net/http"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// errorBody is the one error shape of the API: { "error": { "code": "...", "message": "..." } }.
type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error turns err into an HTTP response. It is the only place that decides
// which status code a domain error gets, so handlers never choose one themselves.
// A 5xx error is logged here, once, and its details are never sent to the client.
func (rs *Responder) Error(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := mapError(err)

	if status >= http.StatusInternalServerError {
		rs.log.Error(r.Context(), "request failed",
			logger.Payload{"error": err.Error()},
			logger.SDC{"error_code": code},
		)
	}

	rs.JSON(w, r, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}

// mapError is the table from domain errors to HTTP (see .ai/rules/http.md).
func mapError(err error) (status int, code, message string) {
	var validation *domain.ValidationError

	switch {
	case errors.As(err, &validation):
		return http.StatusUnprocessableEntity, "validation_failed", validation.Message
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "not_found", "the resource was not found"
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "forbidden", "you are not allowed to do this"
	case errors.Is(err, domain.ErrInvalidTransition):
		return http.StatusConflict, "INVALID_STATUS_TRANSITION", "the status cannot change this way"
	default:
		return http.StatusInternalServerError, "internal_error", "something went wrong"
	}
}
```

**What this does:**
1. `Responder.JSON` meng-encode nilai dulu ke `[]byte`, baru menulis header. Setelah `WriteHeader` dipanggil, status tidak bisa diubah lagi, jadi kalau encode gagal kita masih bisa mengirim `500`.
2. `Responder.Error` memanggil `mapError`, yang berisi tabel dari aturan HTTP: `ErrNotFound` → 404 `not_found`, `ErrForbidden` → 403 `forbidden`, `ErrInvalidTransition` → 409 `INVALID_STATUS_TRANSITION`, `ValidationError` → 422 `validation_failed`, selain itu → 500 `internal_error`.
3. Error 5xx dicatat **sekali** di sini dengan `error_code` di SDC (field yang bisa dicari di Kibana). Isi error aslinya (misalnya `pq: password authentication failed`) hanya masuk log, **tidak pernah** masuk respons.

| Konsep Go | Artinya | Di Node/React |
|---|---|---|
| `type Responder struct{ log *logger.Logger }` + `NewResponder(log)` | Struct yang menyimpan dependency; dibuat sekali di `main.go` | Class dengan constructor injection |
| `func (rs *Responder) JSON(...)` | Method milik `Responder` (`rs` mirip `this`) | Method class |
| `http.ResponseWriter` | Tempat menulis respons; `w` | `res` di Express |
| `*http.Request` | Request yang masuk; `r` | `req` di Express |
| `v any` | Nilai bertipe apa saja | `unknown` di TypeScript |
| `json.Marshal(v)` | Mengubah nilai jadi JSON (`[]byte`) | `JSON.stringify` |
| `switch { case ...: }` tanpa nilai | Rangkaian `if / else if` yang rapi | `if/else` berantai |
| Return bernama `(status int, code, message string)` | Nilai kembalian diberi nama supaya terbaca | Destructuring `{ status, code, message }` |

**File:** `internal/delivery/http/helpers_test.go` (NEW; dipakai oleh semua test di folder ini)

```go
package http

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// newTestLogger returns a logger that writes into the returned buffer, so a test can read the log lines.
func newTestLogger() (*logger.Logger, *bytes.Buffer) {
	var logs bytes.Buffer

	log := logger.New(&logs, logger.Options{Level: slog.LevelDebug, Service: "test", Env: "test", Version: "test"})

	return log, &logs
}

func newTestResponder() (*Responder, *bytes.Buffer) {
	log, logs := newTestLogger()

	return NewResponder(log), logs
}

// get builds a GET request for a test. http.NoBody means "no body" (better than nil), and the
// context is passed explicitly because the linter asks for it.
func get(target string) *http.Request {
	return httptest.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
}
```

**File:** `internal/delivery/http/response_test.go` (NEW)

```go
package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponder_JSON(t *testing.T) {
	rs, _ := newTestResponder()
	rec := httptest.NewRecorder()

	rs.JSON(rec, get("/"), http.StatusCreated, map[string]string{"name": "Budi"})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}

	if got := strings.TrimSpace(rec.Body.String()); got != `{"name":"Budi"}` {
		t.Errorf("body = %q", got)
	}
}

func TestResponder_JSON_ValueThatCannotBeEncoded(t *testing.T) {
	rs, logs := newTestResponder()
	rec := httptest.NewRecorder()

	// A channel has no JSON form, so Marshal fails.
	rs.JSON(rec, get("/"), http.StatusOK, make(chan int))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"code":"internal_error"`) {
		t.Errorf("body is not the error shape: %s", rec.Body.String())
	}

	if !strings.Contains(logs.String(), "encode response") {
		t.Errorf("the failure was not logged: %s", logs.String())
	}
}
```

**File:** `internal/delivery/http/errors_test.go` (NEW)

```go
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
)

func TestResponder_Error(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{name: "not found", err: domain.ErrNotFound, wantStatus: 404, wantCode: "not_found", wantMessage: "the resource was not found"},
		{name: "forbidden", err: domain.ErrForbidden, wantStatus: 403, wantCode: "forbidden", wantMessage: "you are not allowed to do this"},
		{name: "invalid transition", err: domain.ErrInvalidTransition, wantStatus: 409, wantCode: "INVALID_STATUS_TRANSITION", wantMessage: "the status cannot change this way"},
		{name: "validation", err: &domain.ValidationError{Message: "at least one photo is required"}, wantStatus: 422, wantCode: "validation_failed", wantMessage: "at least one photo is required"},
		{name: "wrapped not found", err: fmt.Errorf("find work order 12: %w", domain.ErrNotFound), wantStatus: 404, wantCode: "not_found", wantMessage: "the resource was not found"},
		{name: "wrapped validation", err: fmt.Errorf("complete: %w", &domain.ValidationError{Message: "bad"}), wantStatus: 422, wantCode: "validation_failed", wantMessage: "bad"},
		{name: "unknown error", err: errors.New("boom"), wantStatus: 500, wantCode: "internal_error", wantMessage: "something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, _ := newTestResponder()
			rec := httptest.NewRecorder()

			rs.Error(rec, get("/"), tt.err)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var body errorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body is not the error shape: %v: %s", err, rec.Body.String())
			}

			if body.Error.Code != tt.wantCode || body.Error.Message != tt.wantMessage {
				t.Errorf("error = %+v, want code %q message %q", body.Error, tt.wantCode, tt.wantMessage)
			}
		})
	}
}

func TestResponder_Error_HidesDetailsOf5xx(t *testing.T) {
	rs, logs := newTestResponder()
	rec := httptest.NewRecorder()

	rs.Error(rec, get("/"), errors.New(`pq: password authentication failed for user "app"`))

	if strings.Contains(rec.Body.String(), "password") {
		t.Errorf("response leaks the error details: %s", rec.Body.String())
	}

	// The details are for us: logged once, with the error code in the indexed SDC.
	if got := strings.Count(logs.String(), "request failed"); got != 1 {
		t.Errorf("5xx was logged %d times, want 1: %s", got, logs.String())
	}

	if !strings.Contains(logs.String(), `"error_code":"internal_error"`) {
		t.Errorf("log line has no error_code in sdc: %s", logs.String())
	}
}

func TestResponder_Error_DoesNotLog4xx(t *testing.T) {
	rs, logs := newTestResponder()

	rs.Error(httptest.NewRecorder(), get("/"), domain.ErrNotFound)

	if logs.Len() != 0 {
		t.Errorf("a 404 must not be logged by the error mapper, got: %s", logs.String())
	}
}
```

Hal baru di test:
- `httptest.NewRecorder()` adalah "respons palsu" yang menyimpan status dan body; mirip `supertest` di Node, tapi tanpa membuka port.
- `newTestLogger()` mengembalikan logger yang menulis ke `bytes.Buffer`, jadi test bisa membaca baris log.
- Test `HidesDetailsOf5xx` membuktikan error teknis tidak bocor ke client dan dicatat tepat satu kali.

**Check:**
```bash
go test -v ./internal/delivery/http
```
Hasilnya `--- PASS: TestResponder_JSON`, `TestResponder_JSON_ValueThatCannotBeEncoded`, `TestResponder_Error` (dengan 7 sub-test), `TestResponder_Error_HidesDetailsOf5xx`, `TestResponder_Error_DoesNotLog4xx`, lalu `ok`.

**Common mistakes:**
- **Error `"net/http" imported and not used`** di `errors_test.go`. Itu terjadi kalau kamu memakai helper `get(...)` dan tidak lagi memakai `http.` di file itu. Hapus baris import yang tidak dipakai.
- **Lint `httpNoBody` atau `noctx`.** `httptest.NewRequest(..., nil)` ditegur: pakai `http.NoBody` dan `httptest.NewRequestWithContext(context.Background(), ...)`. Itu sebabnya ada helper `get` di `helpers_test.go`.
- **Menulis `w.WriteHeader(...)` lalu baru `json.Marshal`.** Kalau Marshal gagal, status sudah terkirim dan tidak bisa diubah. Itu sebabnya `JSON` meng-encode lebih dulu.

**PR (step ini):**
- Branch: `feat/3-2-responder` (dibuat di atas `feat/3-1-domain-errors`)
- Base PR: `feat/3-1-domain-errors` (bukan main: ini stack)
- Commit:
  - `feat(http): add responder and domain error mapper` (semua file di Step 2)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add responder and domain error mapper` · Body: `Refs #3`

---

### Step 3: Helper cursor pagination

**Why:** Semua list (teknisi, customer, work order) memakai pola yang sama: `?limit=15&cursor=<teks acak>` dan respons `{ "data": [...], "next_cursor": "...", "prev_cursor": "..." }`. Cursor menunjuk posisi `(created_at, id)` item terakhir, jadi list tetap benar walau ada data baru (beda dengan `?page=2`). Client tidak boleh membacanya, ia hanya mengirim balik teks dari respons sebelumnya.

**File:** `internal/delivery/http/cursor.go` (NEW, layer: `delivery/http`)

```go
package http

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
)

const (
	defaultLimit = 15
	maxLimit     = 20
)

// Cursor is a position in a list that is sorted by (created_at, id).
// The client never reads it: it only sends back the opaque text from the last response.
type Cursor struct {
	CreatedAt time.Time
	ID        int64
}

// EncodeCursor turns a Cursor into the opaque text that goes into ?cursor=.
// The time is stored as microseconds, which is what PostgreSQL keeps for timestamptz.
func EncodeCursor(c Cursor) string {
	raw := fmt.Sprintf("%d|%d", c.CreatedAt.UnixMicro(), c.ID)

	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor is the reverse of EncodeCursor. Text that was not made by EncodeCursor
// gives a validation error, which the error mapper turns into a 422.
func DecodeCursor(text string) (Cursor, error) {
	invalid := &domain.ValidationError{Message: "cursor is not valid"}

	raw, err := base64.RawURLEncoding.DecodeString(text)
	if err != nil {
		return Cursor{}, invalid
	}

	micros, id, found := strings.Cut(string(raw), "|")
	if !found {
		return Cursor{}, invalid
	}

	createdAt, err := strconv.ParseInt(micros, 10, 64)
	if err != nil {
		return Cursor{}, invalid
	}

	itemID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return Cursor{}, invalid
	}

	return Cursor{CreatedAt: time.UnixMicro(createdAt).UTC(), ID: itemID}, nil
}

// PageParams is what a list endpoint reads from ?limit=&cursor=.
type PageParams struct {
	Limit int
	After *Cursor // nil means "from the start"
}

// ParsePageParams reads ?limit= (default 15, from 1 to 20) and ?cursor= from a query string.
func ParsePageParams(query url.Values) (PageParams, error) {
	params := PageParams{Limit: defaultLimit}

	if text := query.Get("limit"); text != "" {
		limit, err := strconv.Atoi(text)
		if err != nil || limit < 1 || limit > maxLimit {
			return PageParams{}, &domain.ValidationError{Message: fmt.Sprintf("limit must be a number from 1 to %d", maxLimit)}
		}

		params.Limit = limit
	}

	if text := query.Get("cursor"); text != "" {
		cursor, err := DecodeCursor(text)
		if err != nil {
			return PageParams{}, err
		}

		params.After = &cursor
	}

	return params, nil
}

// Page is the response of every list endpoint: { "data": [...], "next_cursor": "...", "prev_cursor": "..." }.
// A cursor that does not exist is null, for example next_cursor on the last page.
type Page[T any] struct {
	Data       []T     `json:"data"`
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
}

// NewPage builds a Page. next and prev may be nil.
func NewPage[T any](data []T, next, prev *Cursor) Page[T] {
	if data == nil {
		data = []T{} // "data": [] and not "data": null, so clients can always loop over it
	}

	return Page[T]{Data: data, NextCursor: encodeOptional(next), PrevCursor: encodeOptional(prev)}
}

func encodeOptional(c *Cursor) *string {
	if c == nil {
		return nil
	}

	text := EncodeCursor(*c)

	return &text
}
```

**What this does:**
1. `EncodeCursor` menyimpan `created_at` (dalam mikrodetik, sama dengan ketelitian PostgreSQL) dan `id` sebagai teks `"1763807415000000|42"`, lalu membungkusnya dengan base64 URL-safe (aman ditaruh di URL tanpa di-escape).
2. `DecodeCursor` membalikkan itu. Teks yang bukan buatan `EncodeCursor` menghasilkan `ValidationError`, yang oleh mapper di Step 2 jadi **422**.
3. `ParsePageParams` membaca `?limit=` (default 15, harus 1 sampai 20, di luar itu 422) dan `?cursor=`.
4. `Page[T]` adalah bentuk respons semua list. `NewPage` memastikan `data` selalu `[]` (bukan `null`) dan cursor yang tidak ada jadi `null`.

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `type Page[T any] struct` | **Generics**: satu tipe untuk banyak jenis isi | `interface Page<T>` |
| `*string` | Pointer ke string: bisa `nil`, jadi muncul sebagai `null` di JSON | `string \| null` |
| `strings.Cut(s, "\|")` | Memotong di pemisah pertama; mengembalikan sebelum, sesudah, dan apakah ketemu | `s.split("\|", 2)` |
| `strconv.ParseInt` / `Atoi` | Teks jadi angka, dan bisa gagal | `Number(s)` yang mengembalikan error |
| `url.Values` | Query string yang sudah di-parse (`r.URL.Query()`) | `URLSearchParams` |
| `json:"next_cursor"` | Tag yang menamai field di JSON | Nama key di objek |

**File:** `internal/delivery/http/cursor_test.go` (NEW)

```go
package http

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
)

func TestCursor_RoundTrip(t *testing.T) {
	createdAt := time.Date(2026, 11, 22, 10, 30, 15, 123456000, time.UTC) // microseconds survive

	got, err := DecodeCursor(EncodeCursor(Cursor{CreatedAt: createdAt, ID: 42}))
	if err != nil {
		t.Fatalf("DecodeCursor() error: %v", err)
	}

	if !got.CreatedAt.Equal(createdAt) || got.ID != 42 {
		t.Errorf("got %+v, want created_at %v and id 42", got, createdAt)
	}
}

func TestCursor_IsOpaque(t *testing.T) {
	text := EncodeCursor(Cursor{CreatedAt: time.Unix(1763807415, 0), ID: 42})

	// Safe to put in a URL without escaping, and nothing a client could read or edit by hand.
	if strings.ContainsAny(text, "+/=|") {
		t.Errorf("cursor %q has characters that need URL escaping", text)
	}
}

func TestDecodeCursor_Invalid(t *testing.T) {
	notBase64 := "!!!"
	noSeparator := base64.RawURLEncoding.EncodeToString([]byte("12345"))
	badTime := base64.RawURLEncoding.EncodeToString([]byte("abc|42"))
	badID := base64.RawURLEncoding.EncodeToString([]byte("1763807415000000|x"))

	for name, text := range map[string]string{
		"not base64":     notBase64,
		"no separator":   noSeparator,
		"time is no int": badTime,
		"id is no int":   badID,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := DecodeCursor(text)

			var validation *domain.ValidationError
			if !errors.As(err, &validation) {
				t.Errorf("error = %v, want a *domain.ValidationError", err)
			}
		})
	}
}

func TestParsePageParams(t *testing.T) {
	valid := EncodeCursor(Cursor{CreatedAt: time.UnixMicro(1763807415000000), ID: 7})

	tests := []struct {
		name      string
		query     string
		wantLimit int
		wantAfter bool
		wantErr   bool
	}{
		{name: "defaults", query: "", wantLimit: 15},
		{name: "custom limit", query: "limit=8", wantLimit: 8},
		{name: "max limit", query: "limit=20", wantLimit: 20},
		{name: "with cursor", query: "limit=10&cursor=" + valid, wantLimit: 10, wantAfter: true},
		{name: "limit too big", query: "limit=21", wantErr: true},
		{name: "limit zero", query: "limit=0", wantErr: true},
		{name: "limit negative", query: "limit=-3", wantErr: true},
		{name: "limit not a number", query: "limit=many", wantErr: true},
		{name: "bad cursor", query: "cursor=%21%21%21", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatalf("test query: %v", err)
			}

			got, err := ParsePageParams(values)

			if tt.wantErr {
				var validation *domain.ValidationError
				if !errors.As(err, &validation) {
					t.Fatalf("error = %v, want a *domain.ValidationError", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Limit != tt.wantLimit || (got.After != nil) != tt.wantAfter {
				t.Errorf("got %+v, want limit %d and after set = %v", got, tt.wantLimit, tt.wantAfter)
			}
		})
	}
}

func TestNewPage_JSONShape(t *testing.T) {
	next := &Cursor{CreatedAt: time.UnixMicro(1763807415000000), ID: 7}

	tests := []struct {
		name string
		page Page[string]
		want string
	}{
		{
			name: "no items and no cursors",
			page: NewPage[string](nil, nil, nil),
			want: `{"data":[],"next_cursor":null,"prev_cursor":null}`,
		},
		{
			name: "items with a next cursor",
			page: NewPage([]string{"a", "b"}, next, nil),
			want: `{"data":["a","b"],"next_cursor":"` + EncodeCursor(*next) + `","prev_cursor":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.page)
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("json = %s, want %s", got, tt.want)
			}
		})
	}
}
```

**Check:**
```bash
go test -v -run 'Cursor|PageParams|NewPage' ./internal/delivery/http
```
Hasilnya `TestCursor_RoundTrip`, `TestCursor_IsOpaque`, `TestDecodeCursor_Invalid` (4 sub-test), `TestParsePageParams` (9 sub-test), `TestNewPage_JSONShape` semuanya `PASS`.

**Common mistakes:**
- **`"data": null` di JSON.** Slice kosong di Go bisa `nil`, dan `nil` di-encode jadi `null`. Client yang melakukan `data.map(...)` akan error. `NewPage` mengubahnya jadi `[]`.
- **Memakai `base64.URLEncoding` (dengan `=`)** membuat cursor punya `=` di belakang. Pakai `RawURLEncoding` supaya tidak ada padding.

**PR (step ini):**
- Branch: `feat/3-3-cursor-pagination` (dibuat di atas `feat/3-2-responder`)
- Base PR: `feat/3-2-responder` (bukan main: ini stack)
- Commit:
  - `feat(http): add cursor pagination helper` (`cursor.go`, `cursor_test.go`)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add cursor pagination helper` · Body: `Refs #3`

---

### Step 4: Middleware `X-Request-ID` dan `X-Client`

**Why:** Saat client melapor "tombol Selesai error jam 10", kamu perlu mencari **satu** `request_id` di Kibana dan melihat semua baris log request itu. Web dan Android mengirim `X-Request-ID`; backend memakainya (atau membuat yang baru), menyimpannya di `context`, dan mengirimnya balik di header respons. `X-Client` (`web` atau `android`) menjadi field yang bisa difilter.

**File:** `internal/delivery/http/request_id.go` (NEW)

```go
package http

import (
	"crypto/rand"
	"net/http"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

const (
	headerRequestID = "X-Request-ID"
	maxRequestIDLen = 64
)

// RequestID gives every request an id. It reads X-Request-ID from the client
// (the web app and the Android app send one) or makes a new one, puts it in the
// context so every log line of this request carries it, and sends it back in the
// response header so someone reporting a problem can quote it.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(headerRequestID)
		if !validRequestID(id) {
			id = rand.Text() // 26 random letters and digits
		}

		w.Header().Set(headerRequestID, id)

		// A request is read-only. To add something, make a copy that carries the new context.
		next.ServeHTTP(w, r.WithContext(logger.WithRequestID(r.Context(), id)))
	})
}

// validRequestID accepts letters, digits, "-" and "_" up to 64 characters.
// The value comes from outside and ends up in our logs and in a response header,
// so anything else is replaced instead of trusted.
func validRequestID(id string) bool {
	if id == "" || len(id) > maxRequestIDLen {
		return false
	}

	for _, c := range id {
		isLetter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isDigit := c >= '0' && c <= '9'

		if !isLetter && !isDigit && c != '-' && c != '_' {
			return false
		}
	}

	return true
}
```

**File:** `internal/delivery/http/client.go` (NEW)

```go
package http

import (
	"net/http"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

const headerClient = "X-Client"

// Client reads X-Client ("web" or "android") and adds it to the searchable log fields
// of this request, so in Kibana you can filter every line by client.
// Any other value, or none, becomes "unknown".
func Client(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := r.Header.Get(headerClient)
		if client != "web" && client != "android" {
			client = "unknown"
		}

		ctx := logger.WithSDC(r.Context(), logger.SDC{"client": client})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

**What this does:**
1. Middleware adalah fungsi yang menerima handler dan mengembalikan handler baru yang membungkusnya: melakukan sesuatu sebelum atau sesudah `next.ServeHTTP`.
2. `RequestID` memakai header dari client **hanya kalau valid** (huruf, angka, `-`, `_`, maksimal 64). Nilai itu masuk ke log dan header respons, jadi nilai aneh (spasi, `<script>`, terlalu panjang) diganti dengan id baru dari `rand.Text()`.
3. `logger.WithRequestID` menaruh id di `context`. Setelah itu setiap `log.Info(ctx, ...)` otomatis membawa `x_request_id`, tanpa handler perlu melakukan apa-apa.
4. `Client` menaruh `client` di SDC (field log yang di-index). Nilai selain `web` atau `android` menjadi `unknown`.
5. `r.WithContext(ctx)` membuat **salinan** request. Request di Go bersifat read-only: untuk menambah sesuatu, kamu membuat salinan baru yang membawa context baru.

| Konsep Go | Artinya | Di Node/React |
|---|---|---|
| `func(next http.Handler) http.Handler` | Fungsi yang membungkus handler (middleware) | `app.use((req, res, next) => ...)` di Express |
| `http.HandlerFunc(func(w, r) {...})` | Mengubah fungsi biasa jadi `http.Handler` | Handler berbentuk fungsi |
| `next.ServeHTTP(w, r)` | Meneruskan ke langkah berikutnya | `next()` |
| `context.Context` | Wadah nilai dan sinyal batal yang dibawa sepanjang request | `req.locals` / AsyncLocalStorage |
| `r.WithContext(ctx)` | Salinan request dengan context baru | `{ ...req, ctx }` (immutable, bukan mutasi) |
| `rand.Text()` | String acak yang aman secara kriptografis (`crypto/rand`) | `crypto.randomUUID()` |
| `for _, c := range id` | Loop per karakter (`rune`) | `for (const c of id)` |

**File:** `internal/delivery/http/request_id_test.go` (NEW; ini test yang diminta issue)

```go
package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// serve runs the request through middleware and returns the id the handler saw in its context.
func serve(t *testing.T, mw func(http.Handler) http.Handler, req *http.Request) (seenID string, rec *httptest.ResponseRecorder) {
	t.Helper()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenID = logger.RequestID(r.Context())

		w.WriteHeader(http.StatusNoContent)
	}))

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return seenID, rec
}

func TestRequestID_KeepsTheIDFromTheClient(t *testing.T) {
	req := get("/")
	req.Header.Set("X-Request-ID", "b7f2-web_01")

	seen, rec := serve(t, RequestID, req)

	if seen != "b7f2-web_01" {
		t.Errorf("id in context = %q, want b7f2-web_01", seen)
	}

	if got := rec.Header().Get("X-Request-ID"); got != "b7f2-web_01" {
		t.Errorf("X-Request-ID in response = %q, want b7f2-web_01", got)
	}
}

func TestRequestID_CreatesOneWhenMissing(t *testing.T) {
	seen, rec := serve(t, RequestID, get("/"))

	if seen == "" {
		t.Fatal("no id was created")
	}

	if got := rec.Header().Get("X-Request-ID"); got != seen {
		t.Errorf("response header = %q, context = %q, want the same", got, seen)
	}

	again, _ := serve(t, RequestID, get("/"))
	if again == seen {
		t.Errorf("two requests got the same id %q", seen)
	}
}

func TestRequestID_ReplacesValuesWeDoNotTrust(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "too long", id: strings.Repeat("a", 65)},
		{name: "has a space", id: "abc def"},
		{name: "has markup", id: "<script>alert(1)</script>"},
		{name: "has a newline", id: "abc\ndef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := get("/")
			req.Header["X-Request-Id"] = []string{tt.id} // set directly: Header.Set would not accept a newline

			seen, rec := serve(t, RequestID, req)

			if seen == tt.id || seen == "" {
				t.Errorf("id = %q, want a new one instead of %q", seen, tt.id)
			}

			if got := rec.Header().Get("X-Request-ID"); got != seen {
				t.Errorf("response header = %q, want %q", got, seen)
			}
		})
	}
}

func TestClient_GoesToTheLogLines(t *testing.T) {
	tests := []struct {
		header string
		want   string
	}{
		{header: "web", want: `"sdc":{"client":"web"}`},
		{header: "android", want: `"sdc":{"client":"android"}`},
		{header: "ios", want: `"sdc":{"client":"unknown"}`},
		{header: "", want: `"sdc":{"client":"unknown"}`},
	}

	for _, tt := range tests {
		t.Run("X-Client="+tt.header, func(t *testing.T) {
			log, logs := newTestLogger()
			handler := Client(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				log.Info(r.Context(), "inside the handler")
			}))

			req := get("/")
			if tt.header != "" {
				req.Header.Set("X-Client", tt.header)
			}

			handler.ServeHTTP(httptest.NewRecorder(), req)

			if !strings.Contains(logs.String(), tt.want) {
				t.Errorf("log line = %s, want it to contain %s", logs.String(), tt.want)
			}
		})
	}
}
```

**Check:**
```bash
go test -v -run 'RequestID|Client' ./internal/delivery/http
```
Hasilnya `TestRequestID_KeepsTheIDFromTheClient`, `TestRequestID_CreatesOneWhenMissing`, `TestRequestID_ReplacesValuesWeDoNotTrust` (4 sub-test), dan `TestClient_GoesToTheLogLines` (4 sub-test) semuanya `PASS`.

**Common mistakes:**
- **Mengubah `r` langsung** (misalnya `r.Header.Set(...)` untuk menyimpan id). Itu mengubah request milik semua orang. Simpan di `context` lewat `WithContext`.
- **Memakai nilai header apa adanya.** Header datang dari luar. Tanpa validasi, nilai aneh bisa masuk ke log (log injection) dan ke header respons.
- **Test header dengan newline memakai `Header.Set`** akan menolak nilainya. Itu sebabnya test menulis `req.Header["X-Request-Id"] = ...` langsung (nama key di map header selalu berbentuk kanonik: `X-Request-Id`).

**PR (step ini):**
- Branch: `feat/3-4-request-id` (dibuat di atas `feat/3-3-cursor-pagination`)
- Base PR: `feat/3-3-cursor-pagination` (bukan main: ini stack)
- Commit:
  - `feat(http): add request id and client middleware` (`request_id.go`, `client.go`, `request_id_test.go`)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add request id and client middleware` · Body: `Refs #3`

---

### Step 5: Log per request dan panic recovery

**Why:** Dua hal yang wajib ada di setiap server: (1) satu baris log per request (method, route, status, lama) supaya bisa dilacak; (2) kalau ada handler yang `panic`, server tidak boleh membuang koneksi tanpa jawaban. Client harus tetap mendapat `500` dengan bentuk error yang standar.

**File:** `internal/delivery/http/logging.go` (NEW)

```go
package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

// statusRecorder remembers the status code that the handler wrote, so we can log it.
type statusRecorder struct {
	http.ResponseWriter

	status int
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

// Unwrap lets http.ResponseController reach the real writer underneath (for Flush and similar).
func (s *statusRecorder) Unwrap() http.ResponseWriter {
	return s.ResponseWriter
}

// RequestLog writes one log line per request, after it is done: method, path, route,
// status, and how long it took. The request id and the client come from the context.
//
// It must sit right around the router and pass the SAME *http.Request on: the router
// fills in r.Pattern (the route, like "GET /api/v1/things/{id}") on the request it
// receives, and this middleware reads it back after the handler is done.
func RequestLog(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK} // 200 unless the handler says otherwise

			next.ServeHTTP(rec, r)

			log.Info(r.Context(), "http request",
				logger.SDC{
					"method": r.Method,
					"path":   r.URL.Path,
					"route":  routeOf(r),
					"status": strconv.Itoa(rec.status),
				},
				logger.Payload{"duration_ms": time.Since(start).Milliseconds()},
			)
		})
	}
}

// routeOf returns the route without the method: "GET /things/{id}" becomes "/things/{id}".
// Many different paths share one route, which is what you want to group by in Kibana.
func routeOf(r *http.Request) string {
	if r.Pattern == "" {
		return "unmatched" // no route matched: a 404 or 405 made by the router
	}

	if _, route, found := strings.Cut(r.Pattern, " "); found {
		return route
	}

	return r.Pattern
}
```

**File:** `internal/delivery/http/recover.go` (NEW)

```go
package http

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Recover stops a panic in a handler from killing the request with no answer.
// It logs the panic with its stack (once, through the error mapper) and sends the
// client a normal 500 error body. Without it, the connection is just dropped.
func Recover(rs *Responder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				value := recover()
				if value == nil {
					return
				}

				// http.ErrAbortHandler is how a handler says "drop this connection on purpose".
				// Pass it on instead of answering.
				if err, ok := value.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(value)
				}

				rs.Error(w, r, fmt.Errorf("panic: %v\n%s", value, debug.Stack()))
			}()

			next.ServeHTTP(w, r)
		})
	}
}
```

**What this does:**
1. `statusRecorder` membungkus `ResponseWriter` supaya kita tahu status yang ditulis handler. Defaultnya `200`, karena handler yang hanya memanggil `Write` berarti sukses. Method `Unwrap` membiarkan `http.ResponseController` tetap menjangkau writer aslinya.
2. `RequestLog` mencatat satu baris `"http request"` setelah request selesai: `method`, `path`, `route`, `status` di SDC (bisa dicari) dan `duration_ms` di payload (tidak di-index). `route` diambil dari `r.Pattern` tanpa method (`/api/v1/things/{id}`), supaya banyak path yang berbeda terkelompok jadi satu route di Kibana. Kalau tidak ada route yang cocok, nilainya `unmatched`.
3. `Recover` memakai `defer` + `recover()`: kalau handler panic, ia mencatat panic beserta stack-nya **sekali** (lewat `Responder.Error`) dan mengirim `500` standar. Isi panic tidak pernah bocor ke client. `http.ErrAbortHandler` sengaja diteruskan, karena itu cara handler berkata "putuskan koneksi ini".

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `type statusRecorder struct { http.ResponseWriter; status int }` | **Struct embedding**: `statusRecorder` otomatis punya semua method `ResponseWriter` dan kita menimpa yang perlu | Mirip `extends`, tapi lewat komposisi |
| `defer func() { ... }()` | Jalankan saat fungsi selesai, termasuk saat panic | `finally { }` |
| `recover()` | Menangkap panic dan mengembalikan nilainya | `catch (e)` |
| `panic(value)` | Menghentikan alur seperti `throw`; untuk kondisi yang tidak seharusnya terjadi | `throw new Error()` |
| `debug.Stack()` | Teks stack trace saat ini | `new Error().stack` |
| `r.Pattern` | Pola route yang cocok (Go 1.23+), diisi oleh router | `req.route.path` di Express |
| `strconv.Itoa(n)` | Angka jadi teks (SDC hanya boleh berisi string) | `String(n)` |

**File:** `internal/delivery/http/logging_test.go` (NEW)

```go
package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestServer wires the middleware in the order the router will use:
// RequestID → Client → RequestLog → Recover → router.
func newTestServer(t *testing.T) (handler http.Handler, readLines func() []map[string]any) {
	t.Helper()

	log, logs := newTestLogger()
	rs := NewResponder(log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/things/{id}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /api/v1/boom", func(http.ResponseWriter, *http.Request) {
		panic("something broke")
	})
	mux.HandleFunc("GET /api/v1/silent", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("no WriteHeader call: the status is 200"))
	})

	handler = RequestID(Client(RequestLog(log)(Recover(rs)(mux))))

	readLines = func() []map[string]any {
		var out []map[string]any

		for _, line := range strings.Split(strings.TrimSpace(logs.String()), "\n") {
			var entry map[string]any
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				t.Fatalf("log line is not JSON: %v: %s", err, line)
			}

			out = append(out, entry)
		}

		return out
	}

	return handler, readLines
}

func TestRequestLog_OneLinePerRequest(t *testing.T) {
	handler, readLines := newTestServer(t)

	req := get("/api/v1/things/7")
	req.Header.Set("X-Request-ID", "b7f2")
	req.Header.Set("X-Client", "android")
	handler.ServeHTTP(httptest.NewRecorder(), req)

	got := readLines()
	if len(got) != 1 {
		t.Fatalf("got %d log lines, want 1: %v", len(got), got)
	}

	line := got[0]
	sdc, _ := line["sdc"].(map[string]any)

	want := map[string]string{"method": "GET", "path": "/api/v1/things/7", "route": "/api/v1/things/{id}", "status": "200", "client": "android"}
	for key, value := range want {
		if sdc[key] != value {
			t.Errorf("sdc.%s = %v, want %s", key, sdc[key], value)
		}
	}

	if line["x_request_id"] != "b7f2" || line["message"] != "http request" {
		t.Errorf("line = %v, want x_request_id b7f2 and message \"http request\"", line)
	}

	payload, _ := line["payload"].(map[string]any)
	if _, ok := payload["duration_ms"]; !ok {
		t.Errorf("payload has no duration_ms: %v", payload)
	}
}

func TestRequestLog_StatusAndRoute(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus string
		wantRoute  string
	}{
		{name: "handler never calls WriteHeader", path: "/api/v1/silent", wantStatus: "200", wantRoute: "/api/v1/silent"},
		{name: "no route matches", path: "/nope", wantStatus: "404", wantRoute: "unmatched"},
		{name: "panic becomes 500", path: "/api/v1/boom", wantStatus: "500", wantRoute: "/api/v1/boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, readLines := newTestServer(t)

			handler.ServeHTTP(httptest.NewRecorder(), get(tt.path))

			all := readLines()
			last := all[len(all)-1] // the request line is written last
			sdc, _ := last["sdc"].(map[string]any)

			if sdc["status"] != tt.wantStatus || sdc["route"] != tt.wantRoute {
				t.Errorf("status = %v, route = %v, want %s and %s", sdc["status"], sdc["route"], tt.wantStatus, tt.wantRoute)
			}
		})
	}
}
```

**File:** `internal/delivery/http/recover_test.go` (NEW)

```go
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecover_PanicBecomes500WithTheErrorShape(t *testing.T) {
	rs, logs := newTestResponder()
	handler := Recover(rs)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("index out of range")
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, get("/"))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}

	var body errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Error.Code != "internal_error" {
		t.Errorf("body = %s, want the internal_error shape (err: %v)", rec.Body.String(), err)
	}

	if strings.Contains(rec.Body.String(), "index out of range") {
		t.Errorf("response leaks the panic message: %s", rec.Body.String())
	}

	// The panic and its stack go to the log, once.
	if got := strings.Count(logs.String(), "request failed"); got != 1 || !strings.Contains(logs.String(), "index out of range") {
		t.Errorf("want one log line with the panic, got %d: %s", got, logs.String())
	}
}

func TestRecover_NoPanicLeavesTheResponseAlone(t *testing.T) {
	rs, logs := newTestResponder()
	handler := Recover(rs)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, get("/"))

	if rec.Code != http.StatusTeapot || logs.Len() != 0 {
		t.Errorf("status = %d, log = %q, want 418 and no log", rec.Code, logs.String())
	}
}

func TestRecover_LetsErrAbortHandlerThrough(t *testing.T) {
	rs, _ := newTestResponder()
	handler := Recover(rs)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	defer func() {
		got := recover()

		if err, ok := got.(error); !ok || !errors.Is(err, http.ErrAbortHandler) {
			t.Errorf("recovered %v, want http.ErrAbortHandler to pass through", got)
		}
	}()

	handler.ServeHTTP(httptest.NewRecorder(), get("/"))
}
```

`newTestServer` di `logging_test.go` merakit middleware dalam urutan yang akan dipakai router sungguhan: `RequestID → Client → RequestLog → Recover → router`. Test `StatusAndRoute` membuktikan tiga hal: handler yang tidak memanggil `WriteHeader` tercatat `200`, path yang tidak cocok tercatat `404` dengan route `unmatched`, dan panic tercatat `500`.

**Check:**
```bash
go test -v -race ./internal/delivery/http
make check
```
`go test` menampilkan semua test di folder ini `PASS` dan `ok`. `make check` harus `0 issues` dan semua package `ok` (test database-nya `SKIP` kalau `TEST_DATABASE_URL` tidak di-set).

**Common mistakes:**
- **Route selalu `unmatched` di log.** Penyebabnya `RequestLog` tidak membungkus router secara langsung, atau ada middleware di antaranya yang membuat salinan request. Router hanya mengisi `r.Pattern` pada request yang ia terima. Ikuti urutan di bagian "Keputusan desain".
- **Status selalu `200` walau panic.** Berarti `Recover` ada di luar `RequestLog`. Letakkan `Recover` di dalam.
- **Linter `errorlint` pada `value == http.ErrAbortHandler`.** Hasil `recover()` bertipe `any`; ubah jadi `error` dulu lalu pakai `errors.Is`, seperti di `recover.go`.

**PR (step ini):**
- Branch: `feat/3-5-request-log-recover` (dibuat di atas `feat/3-4-request-id`)
- Base PR: `feat/3-4-request-id` (bukan main: ini stack)
- Commit:
  - `feat(http): add request log and panic recovery middleware` (`logging.go`, `recover.go`, dan kedua test)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add request log and panic recovery middleware` · Body: `Refs #3`

---

### Step 6: Middleware CORS

**Why:** Web admin (React) berjalan di `http://localhost:5173`, API di `http://localhost:8080`. Beda port berarti beda *origin*, dan browser **memblokir** halaman membaca jawaban dari origin lain kecuali server berkata "origin itu boleh". Itu aturan browser (CORS), bukan aturan server. Android dan `curl` tidak terpengaruh.

Sebelum request sungguhan yang memakai `Authorization` atau JSON, browser mengirim **preflight**: request `OPTIONS` yang bertanya "boleh aku mengirim ini?". Server harus menjawab `204` dengan daftar yang diizinkan.

**File:** `internal/delivery/http/cors.go` (NEW)

```go
package http

import (
	"net/http"
	"slices"
)

const corsMaxAge = "600" // seconds a browser may remember the answer to a preflight

// CORS lets a browser app that runs on one of the allowed origins call the API.
// Without it, the browser blocks the React app (http://localhost:5173) from reading
// our answers, because it is a different origin than the API (port 8080).
//
// Requests from other origins, and requests without an Origin header (curl, the
// Android app), pass through untouched: CORS is a rule that browsers enforce, not us.
func CORS(allowedOrigins ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The answer depends on the Origin header, so a cache must keep one copy per origin.
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")
			if origin == "" || !slices.Contains(allowedOrigins, origin) {
				next.ServeHTTP(w, r)

				return
			}

			header := w.Header()
			header.Set("Access-Control-Allow-Origin", origin)
			// By default a browser app cannot read response headers. The web app shows the
			// request id in its error messages, so it must be allowed to read this one.
			header.Set("Access-Control-Expose-Headers", headerRequestID)

			// A preflight is the browser asking "may I send this request?" before it sends the real one.
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				header.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID, X-Client")
				header.Set("Access-Control-Max-Age", corsMaxAge)
				w.WriteHeader(http.StatusNoContent)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
```

**What this does:**
1. Hanya origin yang ada di daftar (`http://localhost:5173`) yang mendapat header `Access-Control-Allow-Origin`. Origin lain dan request tanpa header `Origin` (curl, Android) diteruskan apa adanya.
2. Preflight (`OPTIONS` dengan `Access-Control-Request-Method`) dijawab langsung dengan `204` dan **tidak** diteruskan ke router, jadi tidak perlu route `OPTIONS` di tiap endpoint.
3. `Access-Control-Expose-Headers: X-Request-ID` membuat web admin boleh membaca header `X-Request-ID` (spec: web menampilkan "Kode error: b7f2…"). Tanpa ini browser menyembunyikannya.
4. `Vary: Origin` ditambahkan di semua jawaban, supaya cache tidak memberikan jawaban untuk satu origin ke origin lain.

| Konsep Go | Artinya | Di React/Node |
|---|---|---|
| `allowedOrigins ...string` | **Variadic**: parameter yang bisa diisi nol sampai banyak nilai | `...origins: string[]` |
| `slices.Contains(list, x)` | Apakah `x` ada di `list` | `list.includes(x)` |
| `w.Header().Set(k, v)` / `.Add` | Menulis header respons; `Add` menambah tanpa menimpa | `res.setHeader` / `res.append` |
| `w.WriteHeader(204)` lalu `return` | Menjawab sendiri dan tidak memanggil `next` | `res.status(204).end()` tanpa `next()` |
| Alternatif | Vite punya `server.proxy` yang membuat browser mengira API satu origin; kita pakai CORS karena itu cara kerja di produksi | |

**File:** `internal/delivery/http/cors_test.go` (NEW)

```go
package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const webOrigin = "http://localhost:5173"

// newCORSRequest builds a request with the given method, an Origin header, and optional extra headers.
func newCORSRequest(method, origin string, extra map[string]string) *http.Request {
	req := httptest.NewRequestWithContext(context.Background(), method, "/api/v1/things", http.NoBody)

	if origin != "" {
		req.Header.Set("Origin", origin)
	}

	for key, value := range extra {
		req.Header.Set(key, value)
	}

	return req
}

// serveCORS runs the request through CORS and reports whether the next handler was called.
func serveCORS(req *http.Request) (rec *httptest.ResponseRecorder, nextCalled bool) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalled = true

		w.WriteHeader(http.StatusOK)
	})

	rec = httptest.NewRecorder()
	CORS(webOrigin)(next).ServeHTTP(rec, req)

	return rec, nextCalled
}

func TestCORS_AllowedOrigin(t *testing.T) {
	rec, nextCalled := serveCORS(newCORSRequest(http.MethodGet, webOrigin, nil))

	if !nextCalled || rec.Code != http.StatusOK {
		t.Errorf("next called = %v, status = %d, want the request to reach the handler", nextCalled, rec.Code)
	}

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != webOrigin {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, webOrigin)
	}

	if got := rec.Header().Get("Access-Control-Expose-Headers"); got != "X-Request-ID" {
		t.Errorf("Access-Control-Expose-Headers = %q, want X-Request-ID", got)
	}
}

func TestCORS_Preflight(t *testing.T) {
	rec, nextCalled := serveCORS(newCORSRequest(http.MethodOptions, webOrigin, map[string]string{
		"Access-Control-Request-Method":  "POST",
		"Access-Control-Request-Headers": "authorization, content-type",
	}))

	if nextCalled {
		t.Error("a preflight must be answered by CORS, not passed to the handler")
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}

	for header, want := range map[string]string{
		"Access-Control-Allow-Origin":  webOrigin,
		"Access-Control-Allow-Methods": "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Authorization, Content-Type, X-Request-ID, X-Client",
		"Access-Control-Max-Age":       "600",
	} {
		if got := rec.Header().Get(header); got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}
}

func TestCORS_OptionsWithoutPreflightHeaderGoesToTheHandler(t *testing.T) {
	_, nextCalled := serveCORS(newCORSRequest(http.MethodOptions, webOrigin, nil))

	if !nextCalled {
		t.Error("an OPTIONS request that is not a preflight must reach the handler")
	}
}

func TestCORS_OtherOriginsGetNoCORSHeaders(t *testing.T) {
	tests := []struct {
		name   string
		origin string
	}{
		{name: "another website", origin: "https://evil.example.com"},
		{name: "same host, other port", origin: "http://localhost:3000"},
		{name: "no Origin header (curl, Android)", origin: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, nextCalled := serveCORS(newCORSRequest(http.MethodGet, tt.origin, nil))

			if !nextCalled {
				t.Error("the request must still reach the handler")
			}

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
				t.Errorf("Access-Control-Allow-Origin = %q, want none", got)
			}

			if got := rec.Header().Get("Vary"); got != "Origin" {
				t.Errorf("Vary = %q, want Origin on every answer", got)
			}
		})
	}
}
```

**Check:**
```bash
go test -v -run CORS ./internal/delivery/http
```
Hasilnya `TestCORS_AllowedOrigin`, `TestCORS_Preflight`, `TestCORS_OptionsWithoutPreflightHeaderGoesToTheHandler`, dan `TestCORS_OtherOriginsGetNoCORSHeaders` (3 sub-test) semuanya `PASS`.

**Common mistakes:**
- **`Access-Control-Allow-Origin: *`.** Terlalu longgar dan tidak boleh dipakai bersama `Authorization`. Kita menyebut origin yang persis.
- **Lupa `Expose-Headers`.** Web admin tidak bisa membaca `X-Request-ID`, dan tidak ada error yang menjelaskan kenapa; header-nya "hilang" di browser.
- **Preflight sampai ke router.** Hasilnya `404` atau `401` pada `OPTIONS`, dan browser menolak request aslinya dengan pesan CORS yang membingungkan. Itu sebabnya CORS ada paling luar.

**PR (step ini):**
- Branch: `feat/3-6-cors` (dibuat di atas `feat/3-5-request-log-recover`)
- Base PR: `feat/3-5-request-log-recover` (bukan main: ini stack)
- Commit:
  - `feat(http): add cors middleware` (`cors.go`, `cors_test.go`)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add cors middleware` · Body: `Refs #3`

---

### Step 7: Decode body dan validasi

**Why:** Hampir tiap endpoint `POST`/`PATCH` melakukan hal yang sama: baca JSON dari body, tolak yang rusak, cek field wajib. Satu helper membuat semua handler konsisten, dan semua kesalahan input berujung `422 validation_failed` dengan pesan yang jelas untuk client. Aturan **format** (wajib, email valid, panjang) dicek di sini; aturan **bisnis** tetap di domain.

**Perintah (library baru):**
```bash
go get github.com/go-playground/validator/v10
go mod tidy
```
`go get` mengunduh library dan mencatatnya di `go.mod` dan `go.sum` (seperti `npm install` yang mengubah `package.json` dan `package-lock.json`). `go mod tidy` merapikannya dan menandai `validator` sebagai dependency langsung. Hasilnya: `go.mod` bertambah `validator/v10` ditambah beberapa library bawaannya yang bertanda `// indirect`, dan `golang.org/x/sync` serta `golang.org/x/text` naik versi. Itu normal.

**File:** `internal/delivery/http/decode.go` (NEW)

```go
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
)

const maxBodyBytes = 1 << 20 // 1 MiB, far more than any JSON body of this API needs

// RequestDecoder reads a JSON request body into a struct and checks the
// `validate:"..."` tags on that struct. It only checks the FORMAT of the input
// (required, a valid email, a length). Business rules stay in the domain.
type RequestDecoder struct {
	validate *validator.Validate
}

// NewRequestDecoder creates a RequestDecoder. Make one in main.go and share it.
func NewRequestDecoder() *RequestDecoder {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Name a field by its JSON name ("email"), which is what the client knows,
	// and not by its Go name ("Email").
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")

		return name
	})

	return &RequestDecoder{validate: validate}
}

// Decode reads the body of r into dst (a pointer to a struct) and validates it.
// Every problem with the input comes back as a *domain.ValidationError, which the
// Responder turns into a 422, so a handler only has to pass the error on.
func (d *RequestDecoder) Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // a misspelled field is a mistake, not something to ignore silently

	if err := decoder.Decode(dst); err != nil {
		return bodyError(err)
	}

	// After the first JSON value there must be nothing left.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return &domain.ValidationError{Message: "request body must contain a single JSON object"}
	}

	if err := d.validate.Struct(dst); err != nil {
		return fieldErrors(err)
	}

	return nil
}

// bodyError explains, in words for the client, why the body could not be read.
func bodyError(err error) error {
	var (
		syntaxErr   *json.SyntaxError
		typeErr     *json.UnmarshalTypeError
		tooLargeErr *http.MaxBytesError
		invalidDst  *json.InvalidUnmarshalError
	)

	message := "request body is not valid JSON"

	switch {
	case errors.As(err, &invalidDst):
		// dst was not a pointer: a bug in our code, not a problem with the input, so a 500.
		return fmt.Errorf("decode request: %w", err)
	case errors.Is(err, io.EOF):
		message = "request body is required"
	case errors.As(err, &tooLargeErr):
		message = "request body is too large"
	case errors.As(err, &typeErr):
		message = fmt.Sprintf("field %q has the wrong type", typeErr.Field)
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
		message = "request body is not valid JSON"
	default:
		// encoding/json has no error type for an unknown field, only this text.
		if name, found := strings.CutPrefix(err.Error(), "json: unknown field "); found {
			message = "unknown field " + name
		}
	}

	return &domain.ValidationError{Message: message}
}

// fieldErrors turns the validator's errors into one message: "name is required; email must be a valid email".
func fieldErrors(err error) error {
	var invalid validator.ValidationErrors
	if !errors.As(err, &invalid) {
		// Not a problem with the input but with our struct tags: a bug, so a 500.
		return fmt.Errorf("validate request: %w", err)
	}

	messages := make([]string, 0, len(invalid))
	for _, fieldErr := range invalid {
		messages = append(messages, fieldErr.Field()+" "+describe(fieldErr))
	}

	return &domain.ValidationError{Message: strings.Join(messages, "; ")}
}

// describe says what is wrong with a field, by the validator tag that failed.
func describe(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + fieldErr.Param()
	case "max":
		return "must be at most " + fieldErr.Param()
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(fieldErr.Param(), " ", ", ")
	default:
		return "is not valid"
	}
}
```

**What this does:**
1. `NewRequestDecoder` membuat validator dan mengajarinya memakai nama JSON (`email`), bukan nama Go (`Email`), supaya pesan error cocok dengan yang dilihat client.
2. `Decode` membatasi body 1 MiB (`MaxBytesReader`), menolak field yang tidak dikenal (`DisallowUnknownFields`, karena salah ketik nama field adalah kesalahan, bukan sesuatu yang diam-diam diabaikan), memastikan hanya ada satu objek JSON, lalu memvalidasi tag `validate:"..."`.
3. Semua masalah input menjadi `*domain.ValidationError` dengan pesan untuk client, misalnya `name is required; email must be a valid email`. Mapper di Step 2 mengubahnya jadi `422`.
4. Kalau `dst` bukan pointer, itu bug di kode kita, bukan salah client, jadi hasilnya error biasa (`500`).

| Konsep Go | Artinya | Di React/TS |
|---|---|---|
| `` `validate:"required,email"` `` | **Struct tag**: catatan di field yang dibaca library | Skema Zod: `z.string().email()` |
| `json.NewDecoder(r.Body).Decode(&dst)` | Membaca JSON dari stream langsung ke struct | `await req.json()` + parse ke tipe |
| `errors.As(err, &typeErr)` dengan tipe error bawaan | Mencari jenis error tertentu (`*json.SyntaxError`, dll.) | `instanceof SyntaxError` |
| `fmt.Sprintf("%q", x)` | Teks dengan tanda kutip: `"name"` | `JSON.stringify(x)` |
| `strings.CutPrefix(s, p)` | Menghapus awalan kalau ada, dan memberi tahu apakah ada | `s.startsWith(p) ? s.slice(p.length)` |
| `reflect.StructField` | Informasi tentang sebuah field struct (nama, tag) | Metadata decorator |

**File:** `internal/delivery/http/decode_test.go` (NEW)

```go
package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
)

type createTechnicianRequest struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"omitempty,oneof=admin technician"`
}

func post(body string) *http.Request {
	return httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(body))
}

func TestRequestDecoder_Valid(t *testing.T) {
	var got createTechnicianRequest

	err := NewRequestDecoder().Decode(httptest.NewRecorder(), post(`{"name":"Budi","email":"budi@example.com","role":"technician"}`), &got)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	if got.Name != "Budi" || got.Email != "budi@example.com" || got.Role != "technician" {
		t.Errorf("got %+v", got)
	}
}

func TestRequestDecoder_Problems(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantMessage string
	}{
		{name: "empty body", body: ``, wantMessage: "request body is required"},
		{name: "not JSON", body: `{name:`, wantMessage: "request body is not valid JSON"},
		{name: "cut off", body: `{"name":"Bu`, wantMessage: "request body is not valid JSON"},
		{name: "unknown field", body: `{"name":"Budi","email":"b@example.com","age":30}`, wantMessage: `unknown field "age"`},
		{name: "wrong type", body: `{"name":5,"email":"b@example.com"}`, wantMessage: `field "name" has the wrong type`},
		{name: "two objects", body: `{"name":"Budi","email":"b@example.com"}{"name":"x"}`, wantMessage: "request body must contain a single JSON object"},
		{name: "too large", body: `{"name":"` + strings.Repeat("a", maxBodyBytes) + `"}`, wantMessage: "request body is too large"},
		{name: "missing fields", body: `{}`, wantMessage: "name is required; email is required"},
		{name: "name too short", body: `{"name":"B","email":"b@example.com"}`, wantMessage: "name must be at least 2"},
		{name: "bad email", body: `{"name":"Budi","email":"not-an-email"}`, wantMessage: "email must be a valid email"},
		{name: "bad role", body: `{"name":"Budi","email":"b@example.com","role":"boss"}`, wantMessage: "role must be one of: admin, technician"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst createTechnicianRequest

			err := NewRequestDecoder().Decode(httptest.NewRecorder(), post(tt.body), &dst)

			var validation *domain.ValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("error = %v, want a *domain.ValidationError", err)
			}

			if validation.Message != tt.wantMessage {
				t.Errorf("message = %q, want %q", validation.Message, tt.wantMessage)
			}
		})
	}
}

func TestRequestDecoder_ErrorBecomesA422(t *testing.T) {
	rs, _ := newTestResponder()

	var dst createTechnicianRequest

	err := NewRequestDecoder().Decode(httptest.NewRecorder(), post(`{}`), &dst)

	rec := httptest.NewRecorder()
	rs.Error(rec, get("/"), err)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"code":"validation_failed"`) || !strings.Contains(rec.Body.String(), "name is required") {
		t.Errorf("body = %s", rec.Body.String())
	}
}

func TestRequestDecoder_DestinationThatIsNotAPointerIsABug(t *testing.T) {
	var dst createTechnicianRequest // should be &dst

	err := NewRequestDecoder().Decode(httptest.NewRecorder(), post(`{"name":"Budi"}`), dst)

	var validation *domain.ValidationError
	if err == nil || errors.As(err, &validation) {
		t.Errorf("error = %v, want a plain error (it becomes a 500), not a validation error", err)
	}
}
```

**Check:**
```bash
go test -v -run RequestDecoder ./internal/delivery/http
make check
```
Hasilnya `TestRequestDecoder_Valid`, `TestRequestDecoder_Problems` (11 sub-test), `TestRequestDecoder_ErrorBecomesA422`, dan `TestRequestDecoder_DestinationThatIsNotAPointerIsABug` semuanya `PASS`.

**Common mistakes:**
- **`go: updates to go.mod needed`** atau build gagal dengan `no required module provides package`. Jalankan `go get ...` lalu `go mod tidy`.
- **`validator` bertanda `// indirect` di `go.mod`.** Berarti kamu belum menjalankan `go mod tidy` setelah kode mengimpornya.
- **Pesan memakai nama Go (`Email is required`).** Itu terjadi kalau `RegisterTagNameFunc` tidak dipasang, atau field tidak punya tag `json`.

**PR (step ini):**
- Branch: `feat/3-7-decode-validate` (dibuat di atas `feat/3-6-cors`)
- Base PR: `feat/3-6-cors` (bukan main: ini stack)
- Commit:
  - `feat(http): add json decode and validation helper` (`decode.go`, `decode_test.go`, `go.mod`, `go.sum`)
  - Footer commit: `Refs: #3`
- Judul PR: `feat(http): add json decode and validation helper` · Body: `Refs #3`

---

### Step 8: Router `/api/v1`, rantai middleware, dan server

**Why:** Sekarang semua bagian disatukan: satu tempat mendaftarkan route (semuanya di bawah `/api/v1`), middleware dipasang dalam urutan yang benar, dan `main.go` menyalakan server HTTP yang berhenti dengan rapi saat Ctrl+C (menunggu request yang sedang berjalan selesai). Setelah ini endpoint pertama (`GET /health`, issue #4) tinggal satu handler dan satu baris `router.HandleFunc`.

**File:** `internal/delivery/http/router.go` (NEW)

```go
package http

import (
	"net/http"
	"strings"

	"github.com/irfanandriansyah1997/bang-backend/internal/domain"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

const apiPrefix = "/api/v1"

// Middleware wraps a handler with something every request needs (an id, a log line, CORS).
type Middleware func(http.Handler) http.Handler

// Chain wraps h with the middlewares. The first one is the outermost:
// it runs first on the way in and last on the way out.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

// Router collects the routes of the API and builds the handler that serves them.
type Router struct {
	mux         *http.ServeMux
	responder   *Responder
	middlewares []Middleware
}

// NewRouter creates a Router with all the middlewares of the API, from the outside in:
// CORS → RequestID → Client → RequestLog → Recover → the routes.
// RequestLog must stay directly around Recover and the routes, see logging.go.
func NewRouter(log *logger.Logger, responder *Responder, corsOrigins ...string) *Router {
	return &Router{
		mux:       http.NewServeMux(),
		responder: responder,
		middlewares: []Middleware{
			CORS(corsOrigins...),
			RequestID,
			Client,
			RequestLog(log),
			Recover(responder),
		},
	}
}

// Handle registers a route. The pattern is "METHOD /path", and /api/v1 is put in
// front of the path: Handle("GET /technicians/{id}", h) serves GET /api/v1/technicians/{id}.
func (rt *Router) Handle(pattern string, h http.Handler) {
	rt.mux.Handle(withPrefix(pattern), h)
}

// HandleFunc is Handle for a plain function.
func (rt *Router) HandleFunc(pattern string, fn http.HandlerFunc) {
	rt.Handle(pattern, fn)
}

// Handler returns the whole API: the routes with the middlewares around them.
func (rt *Router) Handler() http.Handler {
	return Chain(http.HandlerFunc(rt.serve), rt.middlewares...)
}

func (rt *Router) serve(w http.ResponseWriter, r *http.Request) {
	// When no route matches, the ServeMux answers with plain text ("404 page not found").
	// Every answer of this API has the same JSON error shape, so answer it ourselves.
	if _, pattern := rt.mux.Handler(r); pattern == "" {
		rt.responder.Error(w, r, domain.ErrNotFound)

		return
	}

	rt.mux.ServeHTTP(w, r)
}

// withPrefix turns "GET /things" into "GET /api/v1/things".
func withPrefix(pattern string) string {
	if method, path, found := strings.Cut(pattern, " "); found {
		return method + " " + apiPrefix + path
	}

	return apiPrefix + pattern
}
```

**What this does:**
1. `Chain` membungkus handler dengan middleware; yang pertama di daftar paling luar. `NewRouter` memasang urutan `CORS → RequestID → Client → RequestLog → Recover → route`.
2. `Handle("GET /technicians/{id}", h)` menambahkan `/api/v1` di depan path, jadi nanti tidak ada yang lupa prefix.
3. `serve` mengganti jawaban `404` bawaan `ServeMux` (teks biasa `404 page not found`) dengan error JSON standar, supaya semua jawaban API berbentuk sama. Middleware tetap berjalan untuk `404`, jadi `X-Request-ID` ada dan baris log ditulis dengan route `unmatched`.
4. Batasan yang diketahui: path yang ada tapi method-nya salah (misalnya `POST` ke route `GET`) juga menjadi `404`, bukan `405`. Cukup untuk sekarang.

| Konsep Go | Artinya | Di React/Node |
|---|---|---|
| `type Middleware func(http.Handler) http.Handler` | Memberi nama pada tipe fungsi | `type Middleware = (next: Handler) => Handler` |
| `http.NewServeMux()` + `"GET /things/{id}"` | Router bawaan Go (1.22+): method dan parameter path ada di pola | `app.get("/things/:id", ...)` |
| `r.PathValue("id")` | Membaca bagian `{id}` dari path | `req.params.id` |
| `for i := len(x) - 1; i >= 0; i--` | Loop dari belakang ke depan | `[...x].reverse().reduce(...)` |

**File:** `internal/delivery/http/router_test.go` (NEW)

```go
package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRouter(t *testing.T) (router *Router, logs func() string) {
	t.Helper()

	log, buf := newTestLogger()
	router = NewRouter(log, NewResponder(log), "http://localhost:5173")

	router.HandleFunc("GET /things/{id}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("thing " + r.PathValue("id")))
	})
	router.HandleFunc("GET /boom", func(http.ResponseWriter, *http.Request) {
		panic("something broke")
	})

	return router, buf.String
}

func TestChain_FirstMiddlewareIsOutermost(t *testing.T) {
	var order []string

	mark := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+" in")
				next.ServeHTTP(w, r)
				order = append(order, name+" out")
			})
		}
	}

	handler := Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		order = append(order, "handler")
	}), mark("a"), mark("b"))

	handler.ServeHTTP(httptest.NewRecorder(), get("/"))

	if got, want := strings.Join(order, ", "), "a in, b in, handler, b out, a out"; got != want {
		t.Errorf("order = %q, want %q", got, want)
	}
}

func TestRouter_RoutesLiveUnderAPIV1(t *testing.T) {
	router, _ := newTestRouter(t)

	rec := httptest.NewRecorder()
	router.Handler().ServeHTTP(rec, get("/api/v1/things/7"))

	if rec.Code != http.StatusOK || rec.Body.String() != "thing 7" {
		t.Errorf("status = %d, body = %q, want 200 and \"thing 7\"", rec.Code, rec.Body.String())
	}
}

func TestRouter_UnknownRouteIsAJSON404(t *testing.T) {
	router, _ := newTestRouter(t)

	for _, path := range []string{"/things/7", "/api/v1/nope", "/"} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			router.Handler().ServeHTTP(rec, get(path))

			var body errorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body is not the error shape: %v: %s", err, rec.Body.String())
			}

			if rec.Code != http.StatusNotFound || body.Error.Code != "not_found" {
				t.Errorf("status = %d, code = %q, want 404 and not_found", rec.Code, body.Error.Code)
			}

			if rec.Header().Get("X-Request-ID") == "" {
				t.Error("the middlewares must also run for a 404: X-Request-ID is missing")
			}
		})
	}
}

func TestRouter_AllMiddlewaresAreWiredIn(t *testing.T) {
	router, logs := newTestRouter(t)

	req := get("/api/v1/things/7")
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("X-Request-ID", "b7f2")
	req.Header.Set("X-Client", "web")

	rec := httptest.NewRecorder()
	router.Handler().ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("CORS header = %q", got)
	}

	if got := rec.Header().Get("X-Request-ID"); got != "b7f2" {
		t.Errorf("X-Request-ID = %q, want b7f2", got)
	}

	for _, want := range []string{
		`"x_request_id":"b7f2"`,
		`"client":"web"`,
		`"route":"/api/v1/things/{id}"`,
		`"status":"200"`,
	} {
		if !strings.Contains(logs(), want) {
			t.Errorf("log line is missing %s: %s", want, logs())
		}
	}
}

func TestRouter_PreflightNeverReachesTheRoutes(t *testing.T) {
	router, _ := newTestRouter(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/api/v1/things/7", http.NoBody)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "GET")

	rec := httptest.NewRecorder()
	router.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204 (no route is registered for OPTIONS)", rec.Code)
	}
}

func TestRouter_PanicIsA500WithTheErrorShape(t *testing.T) {
	router, logs := newTestRouter(t)

	rec := httptest.NewRecorder()
	router.Handler().ServeHTTP(rec, get("/api/v1/boom"))

	if rec.Code != http.StatusInternalServerError || !strings.Contains(rec.Body.String(), `"code":"internal_error"`) {
		t.Errorf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	if !strings.Contains(logs(), `"status":"500"`) {
		t.Errorf("the request line must say 500: %s", logs())
	}
}
```

**File:** `cmd/api/main.go` (EDIT: ganti seluruh isi file). Yang berubah: import (`net`, `net/http`, `time`, dan paket HTTP kita dengan alias `apihttp`), tiga konstanta, `run` membuat `Responder` dan `Router` lalu memanggil `serve`, dan fungsi baru `serve`.

```go
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apihttp "github.com/irfanandriansyah1997/bang-backend/internal/delivery/http"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/config"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/database"
	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/logger"
)

const (
	serviceName = "bang-backend"

	// webDevOrigin is where the React admin runs during development (`npm run dev`).
	webDevOrigin = "http://localhost:5173"

	// shutdownTimeout is how long requests that are still running may take to finish after Ctrl+C.
	shutdownTimeout = 10 * time.Second
)

func main() {
	os.Exit(realMain())
}

// realMain returns the exit code instead of calling os.Exit itself,
// so the deferred stop() below always runs.
func realMain() int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// ctx is canceled when the user presses Ctrl+C (SIGINT) or the
	// platform sends SIGTERM (e.g. Docker stop).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// run wires the app together and blocks until ctx is canceled.
func run(ctx context.Context, cfg config.Config, stdout io.Writer) error {
	level, err := logger.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}

	logFile, err := logger.OpenFile(cfg.LogFile)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	defer logFile.Close()

	// The terminal is indented only in local development.
	// The log file always gets the raw JSON line, one per entry.
	// The file comes first: MultiWriter stops at the first error, and the file is the full record.
	terminal := stdout
	if cfg.AppEnv == "local" {
		terminal = logger.NewPrettyWriter(stdout)
	}

	log := logger.New(io.MultiWriter(logFile, terminal), logger.Options{
		Level:   level,
		Service: serviceName,
		Env:     cfg.AppEnv,
		Version: cfg.AppVersion,
	})

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	log.Info(ctx, "database connected")

	responder := apihttp.NewResponder(log)
	router := apihttp.NewRouter(log, responder, webDevOrigin)
	// Routes are registered here, one feature at a time. The first one is GET /health (issue #4).

	return serve(ctx, log, cfg, router.Handler())
}

// serve runs the HTTP server until ctx is canceled (Ctrl+C), then shuts it down:
// it stops taking new requests and waits for the running ones to finish.
func serve(ctx context.Context, log *logger.Logger, cfg config.Config, handler http.Handler) error {
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("listen on port %s: %w", cfg.Port, err)
	}

	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second, // without it, a client that sends its headers very slowly can hold a connection forever
	}

	failed := make(chan error, 1)
	go func() { failed <- server.Serve(listener) }() // Serve blocks, so it runs in its own goroutine

	log.Info(ctx, "app started", logger.Payload{
		"addr":      listener.Addr().String(),
		"port":      cfg.Port,
		"log_level": cfg.LogLevel,
	})

	select {
	case err := <-failed:
		return fmt.Errorf("serve http: %w", err)
	case <-ctx.Done():
	}

	log.Info(ctx, "shutting down")

	// ctx is already canceled, so the shutdown needs its own time limit (it keeps the values of ctx).
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down http server: %w", err)
	}

	return nil
}
```

**What this does:**
1. Paket kita bernama `http`, sama dengan `net/http`, jadi di `main.go` ia diimpor dengan alias `apihttp` supaya tidak bentrok.
2. `serve` membuka port (`net.ListenConfig`), menjalankan `server.Serve(listener)` di goroutine karena fungsi itu memblokir, lalu menunggu salah satu: server gagal, atau `ctx` dibatalkan (Ctrl+C).
3. Saat Ctrl+C: log `"shutting down"`, lalu `server.Shutdown` berhenti menerima request baru dan menunggu yang sedang berjalan sampai 10 detik. `ctx` sudah dibatalkan, jadi shutdown memakai context baru dari `context.WithoutCancel(ctx)` (nilainya terbawa, batalnya tidak) dengan batas waktu sendiri.
4. Pool database ditutup oleh `defer pool.Close()` di `run`, **sesudah** server berhenti, jadi request yang sedang berjalan masih bisa memakai database sampai selesai.
5. `ReadHeaderTimeout` wajib: tanpa itu client yang mengirim header sangat pelan bisa menahan koneksi selamanya (linter `gosec` menegur).

| Konsep Go | Artinya | Di React/Node |
|---|---|---|
| `go func() { ... }()` | **Goroutine**: menjalankan sesuatu secara paralel | Memulai promise tanpa `await` |
| `chan error` + `select` | Channel mengirim hasil antar goroutine; `select` menunggu yang lebih dulu siap | `Promise.race([...])` |
| `server.Shutdown(ctx)` | Menghentikan server dengan rapi | `server.close()` yang menunggu koneksi selesai |
| `context.WithoutCancel(ctx)` | Salinan `ctx` yang tidak ikut dibatalkan | Tidak ada padanan langsung |
| `apihttp "…/delivery/http"` | Alias import untuk menghindari bentrok nama | `import { x as y }` |
| Port `"0"` | Sistem memilih port kosong; dipakai test supaya tidak bentrok | `server.listen(0)` |

**File:** `cmd/api/main_test.go` (EDIT: ganti seluruh isi file). Yang berubah: port test jadi `"0"`, pengecekan `TEST_DATABASE_URL` dipindah ke helper `testDatabaseURL`, dan ada test baru `TestRun_ServesTheAPI` yang benar-benar menyalakan server, memanggilnya lewat HTTP, lalu memastikan jawaban dan baris lognya.

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/irfanandriansyah1997/bang-backend/internal/infrastructure/config"
)

const testSecret = "test-secret-do-not-log"

func testConfig(t *testing.T) config.Config {
	t.Helper()

	return config.Config{
		Port:        "0", // 0 lets the system pick a free port, so tests never clash with a running server
		DatabaseURL: "postgres://app:app@127.0.0.1:1/fieldservice?sslmode=disable&connect_timeout=1",
		JWTSecret:   testSecret,
		LogFile:     filepath.Join(t.TempDir(), "logs", "app.log"),
		LogLevel:    "info",
		AppEnv:      "test",
		AppVersion:  "test",
	}
}

func TestRun_InvalidLogLevel(t *testing.T) {
	cfg := testConfig(t)
	cfg.LogLevel = "random level"

	err := run(context.Background(), cfg, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "log level") {
		t.Fatalf("run() error = %v, want an error about the log level", err)
	}
}

func TestRun_DatabaseUnreachable(t *testing.T) {
	err := run(context.Background(), testConfig(t), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "connect database") {
		t.Fatalf("run() error = %v, want an error about the database connection", err)
	}
}

// testDatabaseURL returns TEST_DATABASE_URL. Without it the test is skipped on a laptop,
// but fails in CI, so a broken CI setup cannot hide the tests that need a real database.
func testDatabaseURL(t *testing.T) string {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		if os.Getenv("CI") == "true" {
			t.Fatalf("TEST_DATABASE_URL must be set in CI")
		}

		t.Skip("TEST_DATABASE_URL not set")
	}

	return dbURL
}

func TestRun_ShutDownOnCancel(t *testing.T) {
	dbURL := testDatabaseURL(t)

	tests := []struct {
		name       string
		appEnv     string
		wantPretty bool
	}{
		{name: "local env indents the terminal", appEnv: "local", wantPretty: true},
		{name: "other env keeps one line per entry", appEnv: "production", wantPretty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testConfig(t)
			cfg.DatabaseURL = dbURL
			cfg.AppEnv = tt.appEnv

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var out bytes.Buffer
			done := make(chan error, 1)
			go func() { done <- run(ctx, cfg, &out) }()

			waitForLog(t, cfg.LogFile, "app started")
			cancel() // same as pressing Ctrl+C

			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("run() error: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("run() did not return after the context was canceled")
			}

			file, err := os.ReadFile(cfg.LogFile)
			if err != nil {
				t.Fatalf("read log file: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(string(file)), "\n")
			if len(lines) != 3 {
				t.Errorf("log file has %d lines, want 3: %q", len(lines), file)
			}

			if !strings.Contains(string(file), `"message":"shutting down"`) {
				t.Errorf("log file is missing the shutdown message: %q", file)
			}

			if strings.Contains(string(file), testSecret) || strings.Contains(out.String(), testSecret) {
				t.Errorf("JWT secret was written to the logs")
			}

			if pretty := strings.Contains(out.String(), "\n  \"message\""); pretty != tt.wantPretty {
				t.Errorf("terminal output indented = %v, want %v:\n%s", pretty, tt.wantPretty, out.String())
			}
		})
	}
}

func TestRun_ServesTheAPI(t *testing.T) {
	cfg := testConfig(t)
	cfg.DatabaseURL = testDatabaseURL(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- run(ctx, cfg, &bytes.Buffer{}) }()

	waitForLog(t, cfg.LogFile, "app started")

	// No route is registered yet (the first one is GET /health in #4),
	// so every path must give a 404 in the standard JSON error shape.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+serverAddr(t, cfg.LogFile)+"/api/v1/nope", http.NoBody)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	req.Header.Set("X-Request-ID", "test-request-1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /api/v1/nope: %v", err)
	}
	defer resp.Body.Close()

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&body); decodeErr != nil {
		t.Fatalf("body is not the error shape: %v", decodeErr)
	}

	if resp.StatusCode != http.StatusNotFound || body.Error.Code != "not_found" {
		t.Errorf("status = %d, code = %q, want 404 and not_found", resp.StatusCode, body.Error.Code)
	}

	if got := resp.Header.Get("X-Request-ID"); got != "test-request-1" {
		t.Errorf("X-Request-ID = %q, want test-request-1", got)
	}

	cancel() // same as pressing Ctrl+C

	select {
	case runErr := <-done:
		if runErr != nil {
			t.Fatalf("run() error: %v", runErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run() did not return after the context was canceled")
	}

	// The request left one line in the log file, with the id the client sent.
	file, err := os.ReadFile(cfg.LogFile)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	for _, want := range []string{`"message":"http request"`, `"x_request_id":"test-request-1"`, `"route":"unmatched"`, `"status":"404"`} {
		if !strings.Contains(string(file), want) {
			t.Errorf("log file is missing %s:\n%s", want, file)
		}
	}
}

// serverAddr reads the address of the running server (like 127.0.0.1:54321) from the "app started" log line.
func serverAddr(t *testing.T, logPath string) string {
	t.Helper()

	file, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	for _, line := range strings.Split(string(file), "\n") {
		var entry struct {
			Message string `json:"message"`
			Payload struct {
				Addr string `json:"addr"`
			} `json:"payload"`
		}

		if json.Unmarshal([]byte(line), &entry) != nil || entry.Message != "app started" {
			continue
		}

		_, port, err := net.SplitHostPort(entry.Payload.Addr)
		if err != nil {
			t.Fatalf("addr %q: %v", entry.Payload.Addr, err)
		}

		return "127.0.0.1:" + port
	}

	t.Fatal("no \"app started\" line in the log file")

	return ""
}

// waitForLog polls the log file until it contains want, or fails after 5 seconds.
func waitForLog(t *testing.T, path, want string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(path); err == nil && strings.Contains(string(data), want) {
			return
		}

		// Sleep is fine here: we poll a file written by another goroutine, and there is no signal to wait on.
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("log %q did not appear in %s within 5s", want, path)
}
```

**File:** `internal/delivery/http/README.md` (EDIT: ganti bagian "Suggested files (later)" dengan daftar file yang sebenarnya)

````markdown
## Files

```
router.go            Router: routes under /api/v1, and the middleware chain (Chain)
response.go          Responder.JSON: the one place that writes a JSON response
errors.go            Responder.Error: domain error → HTTP status + error code
cursor.go            cursor encode/decode, ParsePageParams, Page[T]
decode.go            RequestDecoder: read a JSON body and validate it
request_id.go        X-Request-ID middleware
client.go            X-Client middleware
logging.go           one log line per request
recover.go           panic → 500
cors.go              CORS for the web admin
*_handler.go         one file per resource (work_order, technician, customer, ...)
```

All of it lives in one package, so the middleware can use the `Responder` without an import cycle.
````

**Check:**
```bash
make check
```
Hasilnya `0 issues` dan semua package `ok`. Di laptop, tes yang butuh database di-`SKIP` kalau `TEST_DATABASE_URL` tidak di-set. Untuk menjalankannya (setelah `docker compose up -d`):
```bash
TEST_DATABASE_URL="postgres://app:app@localhost:5432/fieldservice?sslmode=disable" \
  go test -race -v ./cmd/api
```
`TestRun_ServesTheAPI` harus `PASS`.

Lalu coba server sungguhan (butuh `.env`, lihat README):
```bash
make run
```
Di terminal lain:
```bash
curl -i -H 'X-Request-ID: demo-1' -H 'X-Client: web' http://localhost:8080/api/v1/nope
```
Hasilnya:
```
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
X-Request-Id: demo-1

{"error":{"code":"not_found","message":"the resource was not found"}}
```
Preflight dari origin web admin:
```bash
curl -i -X OPTIONS -H 'Origin: http://localhost:5173' -H 'Access-Control-Request-Method: POST' \
  http://localhost:8080/api/v1/technicians
```
Hasilnya `204 No Content` dengan header `Access-Control-Allow-Origin: http://localhost:5173`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`, `Access-Control-Expose-Headers: X-Request-ID`, dan `Vary: Origin`.

Tekan Ctrl+C di terminal `make run`: muncul log `shutting down` dan proses berhenti dengan kode 0. Di `logs/app.log` ada baris `"message":"http request"` dengan `x_request_id` `demo-1`, `sdc.client` `web`, dan `sdc.route` `unmatched`.

**Common mistakes:**
- **`listen tcp :8080: bind: address already in use`.** Port 8080 dipakai proses lain (mungkin `make dev` yang masih berjalan). Cari dengan `lsof -i :8080`, atau jalankan dengan `PORT=8081`.
- **Lint `G112` (gosec): Potential Slowloris Attack.** `http.Server` tanpa `ReadHeaderTimeout`. Isi seperti di `serve`.
- **`"net/http" redeclared` atau nama `http` bentrok di `main.go`.** Impor paket kita dengan alias `apihttp`.
- **Test bentrok dengan server yang sedang berjalan.** Test memakai port `"0"` supaya sistem memilih port kosong.
- **Route selalu `unmatched` untuk route yang ada.** Ada middleware di antara `RequestLog` dan router yang membuat salinan request. Periksa urutan di `NewRouter`.

**PR (step ini):**
- Branch: `feat/3-8-router-server` (dibuat di atas `feat/3-7-decode-validate`)
- Base PR: `feat/3-7-decode-validate` (bukan main: ini stack)
- Commit (3 commit, urut):
  - `feat(http): add router with api prefix and middleware chain` (`router.go`, `router_test.go`)
  - `feat(http): serve the api from main with graceful shutdown` (`cmd/api/main.go`, `cmd/api/main_test.go`)
  - `docs(http): list the real files in the http readme` (`internal/delivery/http/README.md`)
  - Footer commit terakhir: `Closes: #3`
- Judul PR: `feat(http): add router and serve the api` · Body: `Closes #3` (ini PR terakhir, jadi issue #3 tertutup saat di-merge)

---

## Definition of Done (issue #3)
Per step: centang saat PR-nya di-merge. Issue #3 tertutup oleh PR Step 8.
- [ ] Router `net/http` Go 1.22+ dengan prefix `/api/v1` *(Step 8)*
- [ ] Middleware `X-Request-ID`: dibaca atau dibuat, disimpan di context, dikirim balik di header *(Step 4)*
- [ ] Middleware `X-Client` di log *(Step 4)*
- [ ] Satu baris log per request: status, `duration_ms`, route *(Step 5)*
- [ ] Panic recovery: panic jadi `500` dengan bentuk error standar *(Step 5)*
- [ ] CORS untuk `http://localhost:5173` *(Step 6)*
- [ ] Mapper error domain → HTTP; respons error `{ "error": { "code", "message" } }` *(Step 1 dan 2)*
- [ ] Helper respons JSON *(Step 2)* dan decode + validasi request dengan `go-playground/validator` *(Step 7)*
- [ ] Helper cursor pagination: encode/decode `(created_at, id)`, respons `{ data, next_cursor, prev_cursor }`, default limit 15, maksimal 20 *(Step 3)*
- [ ] Unit test untuk middleware `request_id` *(Step 4)* dan helper cursor *(Step 3)*
- [ ] Server HTTP menyala dari `main.go` dan berhenti dengan rapi saat Ctrl+C *(Step 8)*
- [ ] `make check` lolos di setiap PR
- [ ] Mengikuti `docs/architecture.md`: import ke arah dalam, handler tipis, tanpa tipe sqlc di luar `repository/postgres`
- [ ] `api/openapi.yaml`: tidak berubah karena belum ada endpoint (bagian "Error" dan "Pagination" di OpenAPI ditulis bersama endpoint pertama di #4)
- [ ] Commit memakai `<type>(<feature>): <message>` (feature `http`); PR Step 1–7 memakai `Refs #3`, PR Step 8 memakai `Closes #3`
- [ ] Setiap PR dibuka dengan template; CI Lint + Test hijau

## When you're stuck
- "Aku dapat error `<paste error>` di Step 4. Jelaskan penyebabnya dulu, jangan langsung kasih kode."
- "Jelaskan bedanya `r` dan `r.WithContext(ctx)` dengan analogi Express: kenapa request tidak boleh diubah langsung?"
- "Kenapa `Recover` harus di dalam `RequestLog`? Gambarkan urutan middleware dengan contoh satu request yang panic."

## Changelog
- 2026-11-22: plan dibuat (Step 1–5 lengkap dengan kode)
- 2026-11-22: Step 6 (CORS), Step 7 (decode + validasi), dan Step 8 (router, server, `main.go`) ditambahkan lengkap dengan kode. Plan disusun ulang menjadi **1 step = 1 PR**: tiap step punya branch, rencana commit, dan judul PR; hanya PR Step 8 yang memakai `Closes #3`
- 2026-11-22: alur kerja diganti dari "tunggu merge sebelum step berikutnya" menjadi **stacked PR dengan `gh stack`**: branch step dibuat di atas branch plan dan saling menumpuk, base tiap PR adalah branch di bawahnya, dan merge dari bawah ke atas
