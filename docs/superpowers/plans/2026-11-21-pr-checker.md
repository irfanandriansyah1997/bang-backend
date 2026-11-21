# Plan: PR checker otomatis (lint + unit test → review kelayakan merge oleh Claude)
> Status: **#21–#25 dikerjakan di PR #27; #26 menunggu token dan PR #2** · Last updated: 2026-11-21

## Tujuan

Kamu menulis comment **`@claude-code-review`** di sebuah PR, lalu pengecekan berjalan otomatis di GitHub:

1. **Gate: lint + unit test** dari CI biasa, yang jalan otomatis di setiap push. Saat kamu menulis comment, workflow mengecek hasil Lint dan Test untuk commit terakhir PR. Kalau belum hijau, Claude **tidak dipanggil**, dan bot cukup membalas *"perbaiki Lint/Test dulu"*.
2. **Review oleh Claude.** Hanya jalan kalau gate lolos. Isinya: cek arsitektur, cek kesalahan Go, kualitas test (fokusnya: **apakah kasus utama sudah dites, dan kasus mana yang masih terlewat**), dan konvensi repo. Hasilnya **otomatis di-posting sebagai comment di PR**, lengkap dengan verdict **Ready / Needs changes / Blocked**.

Kamu masih baru di Go, jadi setiap temuan wajib menjelaskan *apa masalahnya*, *kenapa itu masalah*, dan *bagaimana memperbaikinya*, lengkap dengan contoh kode.

## Keputusan yang sudah diambil

| # | Keputusan | Pilihan |
|---|---|---|
| 1 | Seberapa ketat lint | **Strict dari awal**: golangci-lint dengan banyak linter yang dinyalakan |
| 2 | Coverage | **Tanpa batas persentase.** Review mengecek apakah kasus utama (happy path, input tidak valid, error dari dependency, edge case) sudah tercakup, lalu menyebut kasus yang masih terlewat sebagai bahan belajar |
| 3 | Comment ke PR | **Otomatis** |
| 4 | Di mana pengecekan dijalankan | **GitHub Actions**: Claude dipanggil dari workflow, bukan dari laptop |
| 5 | Trigger review | **Comment `@claude-code-review`** di PR (manual), bukan otomatis di setiap push |

## Bisa pakai langganan Claude yang sekarang?

**Bisa, tanpa biaya API tambahan**, asalkan kamu punya langganan **Pro, Max, Ultra, Team, atau Enterprise**. Sumber: [Claude Code GitHub Actions docs](https://code.claude.com/docs/en/github-actions).

- Autentikasinya memakai **`CLAUDE_CODE_OAUTH_TOKEN`**. Token ini dibuat sekali di lokal dengan `claude setup-token`, lalu disimpan sebagai secret repo. Dokumentasinya menyebut: *"If you authenticate with an OAuth token, runs use your Claude subscription instead of API billing."*
- Alternatifnya `ANTHROPIC_API_KEY` dari Claude Console. Opsi ini **dibayar per token** (pay-as-you-go) dan terpisah dari langganan.

**Biaya yang tetap ada atau perlu diperhatikan:**

| Hal | Untuk repo ini |
|---|---|
| Kuota langganan | Setiap review **memakai kuota langganan yang sama** dengan pemakaian Claude Code kamu sehari-hari. Kalau kuota habis, review di CI gagal sampai kuota reset |
| GitHub Actions minutes | Repo `bang-backend` **public**, dan runner standar GitHub untuk repo public **gratis** |
| Token pribadi | OAuth token terikat ke akunmu. Jangan dipakai bersama untuk banyak repo atau orang lain. Untuk tim/organisasi, dokumentasi menyarankan API key |
| PR dari fork | GitHub tidak memberikan secret ke PR dari fork, jadi review tidak jalan di sana (aman) |

**Cara menghemat kuota:**
- Claude **hanya jalan kalau kamu minta** lewat comment, jadi push yang masih setengah jadi tidak memakan kuota.
- Claude hanya jalan kalau Lint dan Test untuk commit terakhir sudah lolos.
- Batasi dengan `--max-turns` dan `timeout-minutes`.
- Pakai `concurrency` per PR: comment baru membatalkan review lama yang masih berjalan.

## Alur

```
push ke PR ──▶ CI biasa (otomatis, gratis): Lint + Test   ← dijaga ruleset main
                         │
kamu comment di PR:  "@claude-code-review"   (opsional: "@claude-code-review fokus ke test")
   │
   └─ workflow claude-review.yml (on: issue_comment)
        1. Comment ini di PR, bukan di issue? Pengirimnya punya akses write?   → kalau tidak, diam
        2. Beri reaksi 👀 di comment                                          → tanda sedang dikerjakan
        3. Cek status Lint + Test untuk commit terakhir PR
              ├─ gagal / belum selesai → balas "Lint/Test belum hijau, perbaiki dulu" → STOP (kuota aman)
              └─ hijau ↓
        4. Checkout kode PR (bukan main)
        5. Claude: /check-pr <PR> --ci  [+ catatan tambahan dari comment]
             a. aturan arsitektur
             b. checklist kesalahan Go
             c. kasus test utama: sudah dites vs masih terlewat
             d. konvensi repo (commit, template PR, OpenAPI, migration)
        6. Comment di PR: verdict + temuan + inline comment di baris yang bermasalah
```

Lint dan unit test tetap dijalankan oleh **CI biasa** di setiap push, bukan oleh Claude: cepat, gratis, dan hasilnya pasti. Workflow review cuma **membaca hasilnya**, tidak menjalankan ulang. Kalau mau review ulang setelah ada perbaikan, tulis lagi comment `@claude-code-review`.

## Yang perlu ditambahkan

### 1. Tooling: golangci-lint (strict)

| File | Isi |
|---|---|
| `.golangci.yml` | Daftar linter strict (tabel di bawah), plus pengecualian yang wajar untuk file test |
| `Makefile` | `make lint` (golangci-lint), `make test` (`go test -race -cover`), dan `make check` menjalankan keduanya |
| `.github/workflows/ci.yml` | Job **Lint** memakai `golangci/golangci-lint-action` |
| `lefthook.yml` | Hook pre-commit memanggil golangci-lint juga, jadi kesalahan sudah tertangkap sebelum push |
| `README.md` | golangci-lint ditambahkan ke tabel requirements |

Linter yang dinyalakan (mode strict):

| Kelompok | Linter | Menangkap |
|---|---|---|
| Bug | `govet`, `staticcheck`, `errcheck`, `ineffassign`, `unused` | Error yang diabaikan, kode yang tidak pernah dipakai, bug umum |
| Resource | `bodyclose`, `sqlclosecheck`, `rowserrcheck` | `Body` atau `rows` yang lupa di-`Close()` |
| Context | `contextcheck`, `noctx` | `ctx` tidak diteruskan ke pemanggilan berikutnya |
| Error | `errorlint`, `wrapcheck` | Error dibandingkan dengan `==` padahal seharusnya `errors.Is`, error tidak dibungkus |
| Keamanan | `gosec` | SQL dirangkai dengan string, izin file terlalu longgar |
| Gaya | `gofmt`, `goimports`, `revive`, `gocritic`, `misspell` | Penulisan Go yang tidak idiomatis |
| Kompleksitas | `gocyclo` / `funlen` (batas longgar) | Fungsi yang terlalu panjang atau terlalu bercabang |

Opsional: job `govulncheck` untuk mengecek celah keamanan di dependency.

### 2. Rules baru (`.ai/rules/`)

Rules ini dipakai **AI saat menulis kode** dan juga **Claude saat me-review di CI**.

| Rule | `paths` | Isi |
|---|---|---|
| `testing.md` | `**/*_test.go` | Kasus utama yang wajib ada (happy path, input tidak valid, error dari dependency, edge case), table-driven, `t.Helper()`, tanpa sleep, test DB di-skip kalau env tidak ada |
| `http.md` | `internal/delivery/http/**` | Handler tipis, DTO terpisah dari domain, error dipetakan di satu tempat, format error standar |
| `domain-usecase.md` | `internal/domain/**`, `internal/usecase/**` | Domain tidak meng-import package luar, `now` dikirim lewat parameter, interface didefinisikan di sisi usecase, `ctx` di parameter pertama |

### 3. Skill baru: `check-pr`

```
.ai/skills/check-pr/
├── SKILL.md
└── references/
    ├── go-review-checklist.md     cara mendeteksi + penjelasan untuk pemula, per item
    ├── test-cases-checklist.md    kasus utama per layer (domain / usecase / repository / http)
    └── report-template.md         format comment PR
```

Skill ini punya dua mode:
- **`--ci`** (dipanggil dari workflow `@claude-code-review`): workflow sudah memastikan Lint dan Test hijau, jadi skill langsung me-review lalu **mem-posting comment otomatis**. Comment berupa ringkasan, ditambah inline comment di baris yang bermasalah.
- **Lokal** (`/check-pr 23` dari laptop): skill menjalankan gate sendiri di worktree sementara. Hasilnya ditampilkan di chat, dan comment baru di-posting kalau kamu minta.

**Bagian "kasus test"** di laporan berbentuk tabel per fungsi baru:

| Fungsi | Kasus | Status |
|---|---|---|
| `WorkOrder.Start` | dari `assigned` → sukses | ✅ dites |
| | dari `open` → `ErrInvalidTransition` | ✅ dites |
| | dari `completed` → `ErrInvalidTransition` | ❌ **terlewat**: kenapa penting + contoh test |

**Aturan verdict:**
- Ada **blocker** → **Blocked**
- Ada **major** → **Needs changes**
- Hanya minor/nit → **Ready**

| Tingkat | Contoh |
|---|---|
| **blocker** | Arah import dilanggar, secret masuk log, SQL injection, logic baru tanpa test sama sekali |
| **major** | Kasus gagal yang utama tidak dites, `ctx` tidak diteruskan, OpenAPI tidak di-update, migration tanpa `Down` |
| **minor** | Nama kurang jelas, fungsi terlalu panjang |
| **nit** | Selera gaya penulisan |

Skill ini **tidak pernah approve atau merge**. Kalau PR yang sama sudah punya comment review dari Claude, comment itu di-update, bukan ditambah comment baru.

### 4. Workflow: `.github/workflows/claude-review.yml` (file baru)

Workflow ini dipisah dari `ci.yml` karena trigger-nya berbeda: comment, bukan push.

```yaml
name: Claude Review
on:
  issue_comment:
    types: [created]

concurrency:
  group: claude-review-${{ github.event.issue.number }}
  cancel-in-progress: true

jobs:
  review:
    # hanya comment di PR (bukan issue) yang berisi trigger
    if: github.event.issue.pull_request && contains(github.event.comment.body, '@claude-code-review')
    runs-on: ubuntu-latest
    timeout-minutes: 15
    permissions:
      contents: read
      pull-requests: write
      issues: write        # untuk reaksi 👀 dan comment balasan
      checks: read         # untuk membaca hasil Lint/Test
      id-token: write
    steps:
      - name: React to the comment
        run: gh api repos/${{ github.repository }}/issues/comments/${{ github.event.comment.id }}/reactions -f content=eyes
        env: { GH_TOKEN: "${{ github.token }}" }

      - name: Require green Lint and Test on the PR head
        run: |
          # gagal → posting comment "perbaiki Lint/Test dulu" lalu exit 1
          # (detail ditulis saat implementasi: gh pr checks / API check-runs untuk head SHA)

      - uses: actions/checkout@v6
        with:
          ref: refs/pull/${{ github.event.issue.number }}/head   # kode PR, bukan main
          fetch-depth: 0

      - uses: anthropics/claude-code-action@v1
        with:
          claude_code_oauth_token: ${{ secrets.CLAUDE_CODE_OAUTH_TOKEN }}
          prompt: "/check-pr ${{ github.event.issue.number }} --ci"
          claude_args: |
            --model claude-haiku-5-5
            --max-turns 30
            --allowedTools "mcp__github_inline_comment__create_inline_comment,Bash(gh pr view:*),Bash(gh pr diff:*),Bash(gh pr comment:*),Read,Grep,Glob"
```

Ini masih sketsa. Daftar tool yang diizinkan, cara meneruskan catatan tambahan dari comment, dan versi action dicek ulang saat implementasi.

Hal-hal penting:
- **Checkout `refs/pull/N/head`.** Event `issue_comment` secara default men-checkout `main`. Tanpa `ref` ini, yang di-review justru kode yang salah.
- **Siapa yang bisa memicu.** `claude-code-action` hanya mau jalan kalau yang menulis comment punya akses *write* ke repo, dan bot ditolak. Jadi orang asing yang menulis `@claude-code-review` di repo public kamu **tidak** akan memakai kuotamu.
- **`@claude-code-review` tidak sama dengan `@claude`.** Workflow ini memfilter sendiri dengan `if: contains(...)`, jadi tidak bentrok dengan workflow `@claude` kalau suatu saat kamu memasangnya.
- **Tidak dijadikan required check** di ruleset, karena review ini manual dan kuotanya bisa habis. Lint dan Test tetap wajib lolos.

### 5. Hubungan dengan skill yang sudah ada

```
/start-task #N → ngoding → /atomic-conventional-commits → /create-pr
   → CI: Lint + Test (otomatis) → kamu comment @claude-code-review → Claude Review (comment di PR)
   → /review-task <PR>   (apakah task-nya selesai? opsional, lokal)
   → merge (kamu)
```

## Persiapan dari sisimu (sekali saja)

1. **Install Claude Code CLI.** Saat ini kamu pakai desktop app, dan `claude` belum ada di terminal. Ikuti [cara install di docs](https://code.claude.com/docs/en/setup).
2. Jalankan `claude setup-token` di terminal, lalu login dengan akun langgananmu. Hasilnya sebuah token panjang.
3. Simpan token itu sebagai secret repo bernama `CLAUDE_CODE_OAUTH_TOKEN`: **Settings → Secrets and variables → Actions → New repository secret**. Jangan paste token ini ke chat.
4. Install **Claude GitHub App** ke repo `bang-backend`: https://github.com/apps/claude

Langkah 2–4 juga bisa dilakukan sekaligus dengan `/install-github-app` dari Claude Code CLI.

## Urutan pengerjaan

| # | Pekerjaan | Hasil | Perkiraan |
|---|---|---|---|
| 1 | golangci-lint strict: `.golangci.yml`, `make lint`/`make test`, CI, lefthook, README | Lint ketat di lokal dan CI | S |
| 2 | Rules: `testing.md`, `http.md`, `domain-usecase.md` + `make ai-link` | Panduan menulis dan me-review | S |
| 3 | Skill `check-pr` (mode lokal + `--ci`) + `references/` | `/check-pr` bisa dipakai | M |
| 4 | Workflow `claude-review.yml` (trigger `@claude-code-review`) | Review saat diminta lewat comment | S |
| 5 | Update `create-pr` / `review-task` supaya saling merujuk | Alurnya nyambung | XS |
| 6 | Uji coba di PR #2 (infrastructure) | Bukti bahwa alurnya jalan dari awal sampai akhir | S |

Langkah 4 baru bisa dites setelah kamu menyelesaikan bagian **Persiapan** (token dan GitHub App).

## Langganan: Claude Pro

- Pro **termasuk** Claude Code dan `claude setup-token`, jadi review di CI **tidak memakai biaya API**.
- Batas pemakaian Pro **lebih kecil** daripada Max, dan **dipakai bersama** oleh chat, Claude Code di laptop, dan review di CI. Angka pastinya bisa dilihat di pengaturan akun Claude kamu.
- Karena itu desain ini sengaja hemat:
  - review hanya jalan kalau kamu menulis `@claude-code-review`
  - Claude dilewati kalau Lint atau Test masih merah
  - `--max-turns` dibatasi (mulai dari 20–30, disesuaikan setelah uji coba)
  - satu review per comment, dan comment baru membatalkan review yang masih berjalan
- **Model untuk review di CI** (diatur lewat `--model` di `claude_args`):
  - **Default: Claude Haiku 5.5** (`claude-haiku-5-5`), model paling murah. Cukup untuk menjalankan checklist dan untuk belajar cara kerja agentic review.
  - **Hanya Haiku untuk sekarang.** Opsi `deep` → Sonnet dihapus (2026-11-21). Kalau nanti review Haiku terasa terlalu dangkal, ganti `MODEL` di `claude-review.yml`.
  - Opus tidak dipakai di CI, karena paling berat untuk kuota dan tidak diperlukan untuk tujuan belajar.

  | Model | Harga API per 1 juta token (input / output) | Kira-kira |
  |---|---|---|
  | Claude Haiku 5.5 | $0.10 / $0.50 | paling hemat |
  | Claude Sonnet 5.5 | $2 / $10 | ±20× Haiku |
  | Claude Opus 5.5 | $4 / $20 | ±40× Haiku |

  Kamu pakai langganan Pro, jadi angka di tabel ini **tidak ditagihkan**. Tabel ini hanya gambaran seberapa "berat" tiap model. Umumnya, model yang lebih berat menghabiskan kuota langganan lebih cepat.
- Kalau suatu saat kuota sering habis, ada dua opsi: upgrade ke Max, atau ganti ke `ANTHROPIC_API_KEY` (bayar per token). Workflow-nya cukup diubah di satu baris.

## Task di board

Label `tooling`, dan setiap issue punya acceptance criteria + Definition of Done.

| Issue | Pekerjaan | Status | Priority | Size |
|---|---|---|---|---|
| #21 | Strict golangci-lint (Makefile, CI, hooks) | Ready | P1 | S |
| #22 | Review rules (tests, HTTP, domain/usecase) | Backlog | P1 | S |
| #23 | Skill `check-pr` (lokal + `--ci`), depends on #21, #22 | Backlog | P1 | M |
| #24 | Workflow `@claude-code-review`, depends on #23 + setup token | Backlog | P1 | S |
| #25 | Hubungkan `create-pr` / `review-task` ke `check-pr`, depends on #23 | Backlog | P2 | XS |
| #26 | Uji coba di PR #2, depends on #24 + PR #2 | Backlog | P2 | S |

Review comment di PR ditulis dalam bahasa Inggris, sama seperti issue dan PR lain di repo ini.

## Catatan

- Kamu sedang berada di branch `feat/init`, padahal branch itu sudah di-merge. Pindahkan pekerjaan #2 ke branch baru dengan `git switch -c feat/2-infrastructure`. File yang belum di-commit ikut terbawa.
- Anthropic juga punya produk **[Code Review](https://code.claude.com/docs/en/code-review.md)**: review otomatis tanpa perlu menulis workflow. Tapi produk itu tidak menjalankan checklist dan aturan khusus repo ini seperti `check-pr`, jadi untuk tujuan belajar, workflow sendiri lebih cocok.

## Changelog
- 2026-11-21: plan dibuat; keputusan 1–5 diambil; issue #21–#26 dibuat
- 2026-11-21: #21–#25 dikerjakan di branch `chore/pr-checker-plan` (PR #27):
  - #21: strict golangci-lint v2.15.0 (`.golangci.yml`, `make lint`/`fmt`/`test`/`check`, CI job Lint, lefthook, README). Lint menemukan error yang tidak dibungkus di `cmd/api/main.go`, sudah diperbaiki.
  - #22: rules `testing.md`, `http.md`, `domain-usecase.md`
  - #23: skill `check-pr` + `references/`. Gate diuji pada branch `feat/2-infrastructure`: lint gagal (26 temuan), jadi review berhenti di gate sesuai rancangan.
  - #24: workflow `claude-review.yml` lolos actionlint, dan logika gate disimulasikan. **Belum dites end-to-end**: perlu secret `CLAUDE_CODE_OAUTH_TOKEN` dan workflow harus sudah ada di `main` (GitHub menjalankan workflow `issue_comment` dari branch default).
  - #25: `create-pr` dan `review-task` sekarang merujuk ke `check-pr`, dan README diperbarui
- Catatan untuk #2: dengan lint baru, branch `feat/2-infrastructure` punya 26 temuan (20 revive: komentar doc pada fungsi exported; 4 wrapcheck; 2 gosec: permission file/folder). Ada juga bug di `logger.go`: `os.Mkdir` seharusnya `os.MkdirAll`.
- 2026-11-21: review comments switched to **one inline comment per finding** on its line, with GitHub `suggestion` blocks (one-click *Commit suggestion*). The summary is shorter and only lists the findings. Workflow limits raised to 60 turns / 20 minutes.
- 2026-11-21: review sekarang **hanya memakai Claude Haiku 5.5**. Opsi `deep` → Sonnet dihapus dari workflow, skill, README, dan issue #24/#26.
