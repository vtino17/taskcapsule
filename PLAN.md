# DevHibernate

## Product Requirements, Architecture, Implementation Plan, and Test Specification

**Status:** Draft siap implementasi
**Target release pertama:** `v0.1.0`
**Bahasa utama:** Go
**Tipe produk:** Local-first developer CLI
**Ketergantungan AI:** Tidak ada
**API key:** Tidak diperlukan
**Target awal:** Linux dan macOS
**Lisensi yang disarankan:** Apache-2.0 atau MIT

---

# 1. Ringkasan Produk

DevHibernate adalah CLI yang memungkinkan developer menyimpan satu pekerjaan coding sebagai sebuah **capsule**.

Satu capsule mewakili:

* Git branch
* Git worktree
* service development
* port service
* log service
* catatan pekerjaan
* hasil test terakhir
* status pekerjaan
* langkah untuk melanjutkan pekerjaan

Developer dapat membuat, menghentikan sementara, melanjutkan, menyerahkan, dan menghapus capsule tanpa kehilangan konteks.

Contoh penggunaan:

```bash
devhibernate start payment-timeout
devhibernate note payment-timeout "Investigasi retry gateway"
devhibernate pause payment-timeout

devhibernate start urgent-hotfix
devhibernate pause urgent-hotfix

devhibernate resume payment-timeout
devhibernate where payment-timeout
```

Tujuan utamanya bukan menyimpan isi RAM atau snapshot virtual machine.

DevHibernate menyimpan **cara merekonstruksi environment**, bukan isi memory process.

---

# 2. Masalah yang Diselesaikan

Developer sering mengerjakan beberapa task secara paralel:

* feature baru
* bug production
* code review
* eksperimen
* maintenance
* support issue

Setiap perpindahan task menimbulkan pekerjaan manual:

1. Menghentikan development server.
2. Menyimpan perubahan.
3. Berpindah branch.
4. Mengingat file yang sedang dikerjakan.
5. Mengingat command yang sedang dijalankan.
6. Mencari kembali issue terkait.
7. Menjalankan ulang frontend dan backend.
8. Menangani port yang bentrok.
9. Mengingat test terakhir yang gagal.

DevHibernate mengubah proses tersebut menjadi:

```bash
devhibernate pause task-a
devhibernate resume task-b
```

---

# 3. Product Positioning

Gunakan positioning berikut:

> Pause one coding task and resume another without losing your place.

Alternatif:

> A lightweight hibernation layer for developer tasks.

DevHibernate bukan:

* IDE
* Git client penuh
* container orchestrator
* cloud workspace
* task manager
* package manager
* terminal multiplexer
* AI coding assistant
* virtual machine snapshot tool

DevHibernate menghubungkan Git worktree, process, log, note, dan test state menjadi satu lifecycle task.

---

# 4. Target Pengguna

## Primary users

* software developer individual
* maintainer open-source
* freelancer
* developer yang menangani beberapa repository
* developer yang sering berpindah feature dan hotfix
* developer dengan RAM terbatas
* developer yang menggunakan coding agent secara paralel

## Secondary users

* technical lead
* reviewer
* pair programmer
* developer support
* contributor open-source

---

# 5. Sasaran Produk

## Sasaran MVP

DevHibernate versi pertama harus dapat:

1. Membuat Git worktree untuk sebuah task.
2. Membuat branch khusus task.
3. Menjalankan satu atau beberapa service.
4. Menyimpan PID dan log service.
5. Menghentikan seluruh service dengan aman.
6. Menjalankan ulang seluruh service.
7. Menyimpan catatan terakhir developer.
8. Menjalankan dan menyimpan hasil check/test.
9. Membuat handoff report.
10. Menolak penghapusan jika worktree masih memiliki perubahan.
11. Mendeteksi capsule dengan state rusak atau process yatim.
12. Berjalan tanpa server cloud.

## Sasaran nonfungsional

* startup CLI kurang dari 200 ms untuk command ringan
* tidak membutuhkan daemon untuk MVP
* tidak menyimpan secret
* tidak melakukan Git commit otomatis
* tidak melakukan Git push otomatis
* tidak melakukan stash otomatis
* state ditulis secara atomik
* command harus idempotent bila memungkinkan
* error harus memberikan tindakan perbaikan

---

# 6. Non-Goals MVP

Fitur berikut tidak boleh dikerjakan pada `v0.1.0`:

* menyimpan editor tabs
* menyimpan browser tabs
* menyimpan terminal scrollback
* cloud synchronization
* sharing capsule melalui internet
* Kubernetes support
* remote development
* Windows process lifecycle penuh
* Docker snapshot
* database snapshot otomatis
* local reverse proxy
* `.localhost` domain
* plugin marketplace
* GUI desktop
* VS Code extension
* JetBrains extension
* AI-generated handoff
* automatic commit
* automatic merge
* automatic GitHub issue creation

Fitur-fitur tersebut masuk roadmap setelah lifecycle utama stabil.

AI agent dilarang menambahkan fitur roadmap sebelum seluruh acceptance criteria MVP selesai.

---

# 7. Terminologi

## Capsule

Satu unit pekerjaan developer.

Contoh:

```text
payment-timeout
checkout-redesign
hotfix-login
GH-482
```

## Source repository

Repository Git utama tempat developer menjalankan DevHibernate.

## Worktree

Folder Git terpisah yang dibuat untuk capsule.

## Service

Process development yang berjalan selama capsule aktif.

Contoh:

* frontend
* backend
* worker
* documentation server

## Check

Command berjangka pendek yang dijalankan untuk memvalidasi pekerjaan.

Contoh:

* unit test
* lint
* type check
* build

## Handoff

Laporan Markdown yang memungkinkan developer lain melanjutkan pekerjaan.

---

# 8. Scope Release

## Release `v0.1.0`

Wajib tersedia:

```text
init
start
pause
resume
list
status
note
where
check
logs
handoff
delete
doctor
version
```

## Release `v0.2.0`

Direncanakan:

* internal reverse proxy
* stable `.localhost` URL
* automatic browser opening
* Docker Compose integration
* editor integration

## Release `v0.3.0`

Direncanakan:

* capsule export/import
* GitHub issue integration
* pull-request integration
* team handoff package

---

# 9. Kontrak CLI

Binary:

```bash
devhibernate
```

Alias opsional setelah MVP:

```bash
devh
```

## 9.1 `devhibernate init`

Membuat file konfigurasi awal.

```bash
devhibernate init
```

Output:

```text
Created .devhibernate.json

Next steps:
1. Review the generated configuration
2. Run: devhibernate start my-task
```

Error jika file sudah ada:

```text
Configuration already exists: .devhibernate.json
Use --force to replace it.
```

Flag:

```bash
devhibernate init --force
```

`--force` hanya mengganti file konfigurasi, bukan menghapus state capsule.

---

## 9.2 `devhibernate start`

Membuat capsule baru.

```bash
devhibernate start payment-timeout
```

Optional flags:

```bash
devhibernate start payment-timeout --base main
devhibernate start payment-timeout --branch fix/payment-timeout
devhibernate start payment-timeout --no-services
```

Urutan proses:

1. Validasi repository.
2. Validasi konfigurasi.
3. Validasi nama capsule.
4. Periksa capsule dengan nama sama.
5. Tentukan base branch.
6. Tentukan branch capsule.
7. Buat worktree.
8. Jalankan setup commands.
9. Alokasikan port.
10. Jalankan services.
11. Jalankan health checks.
12. Simpan state.
13. Tampilkan hasil.

Output:

```text
Capsule started: payment-timeout

Branch:    fix/payment-timeout
Worktree:  ~/.devhibernate/worktrees/shop/payment-timeout
Status:    running

Services:
  api       running  pid=24122  port=43102
  frontend  running  pid=24129  port=43103

Use:
  devhibernate pause payment-timeout
```

Jika setup gagal:

* service tidak boleh dijalankan
* worktree boleh dipertahankan
* capsule menjadi `error`
* error log harus disimpan
* pengguna diberi instruksi perbaikan

---

## 9.3 `devhibernate pause`

Menghentikan seluruh service capsule.

```bash
devhibernate pause payment-timeout
```

Urutan:

1. Baca capsule state.
2. Periksa process yang masih hidup.
3. Kirim graceful termination.
4. Tunggu grace period.
5. Kirim force termination jika masih hidup.
6. Simpan status `paused`.
7. Simpan waktu pause.
8. Simpan summary terakhir.

Output:

```text
Capsule paused: payment-timeout

Stopped:
  api       pid=24122
  frontend  pid=24129

Resources released.
```

Jika capsule sudah paused:

```text
Capsule is already paused: payment-timeout
```

Command harus tetap exit code `0`.

---

## 9.4 `devhibernate resume`

Menjalankan kembali service dari capsule paused.

```bash
devhibernate resume payment-timeout
```

Urutan:

1. Validasi worktree masih ada.
2. Validasi branch masih benar.
3. Periksa capsule tidak sedang running.
4. Alokasikan port baru jika diperlukan.
5. Jalankan services.
6. Jalankan health checks.
7. Simpan PID baru.
8. Ubah status menjadi `running`.

Output:

```text
Capsule resumed: payment-timeout

Services:
  api       running  pid=25101  port=43210
  frontend  running  pid=25107  port=43211

Last note:
  Investigate retry after gateway timeout.
```

Resume tidak menjalankan setup ulang secara default.

Flag opsional:

```bash
devhibernate resume payment-timeout --setup
```

---

## 9.5 `devhibernate list`

Menampilkan seluruh capsule dalam repository aktif.

```bash
devhibernate list
```

Output:

```text
NAME                STATUS    BRANCH                     UPDATED
payment-timeout     paused    fix/payment-timeout        2h ago
checkout-redesign   running   feat/checkout-redesign     14m ago
hotfix-login        error     hotfix/login                1d ago
```

Flag global:

```bash
devhibernate list --all
```

`--all` menampilkan capsule dari semua repository.

---

## 9.6 `devhibernate status`

Menampilkan detail capsule.

```bash
devhibernate status payment-timeout
```

Output:

```text
Capsule: payment-timeout
Status: paused

Repository:
  /home/user/projects/shop

Worktree:
  /home/user/.devhibernate/worktrees/shop/payment-timeout

Git:
  Branch: fix/payment-timeout
  Base: main
  Dirty: yes
  Changed files: 3

Services:
  api       stopped
  frontend  stopped

Last check:
  pnpm test payment-retry
  Failed
  Exit code: 1
  Finished: 2026-07-13 14:41

Last note:
  Integration test expects 2 retries but receives 3.
```

---

## 9.7 `devhibernate note`

Menyimpan catatan konteks.

```bash
devhibernate note payment-timeout "Investigate duplicate retry"
```

Perilaku:

* satu capsule memiliki satu `current note`
* note lama masuk history
* note disimpan dengan timestamp
* note tidak boleh menjalankan command
* note tidak boleh diinterpretasikan sebagai konfigurasi

Output:

```text
Note saved for payment-timeout.
```

---

## 9.8 `devhibernate where`

Menampilkan ringkasan untuk melanjutkan pekerjaan.

```bash
devhibernate where payment-timeout
```

Output:

```text
You were working on: payment-timeout

Last note:
  Investigate duplicate retry after gateway timeout.

Last modified files:
  src/payment/retry-policy.ts
  tests/payment-retry.test.ts

Last check:
  pnpm test payment-retry
  Result: failed

Suggested next action:
  Open the worktree and review the last failing check.
```

Tidak boleh menggunakan AI.

Data berasal dari:

* note terakhir
* Git status
* file modification time
* check terakhir
* state service

---

## 9.9 `devhibernate check`

Menjalankan command validation dalam worktree capsule.

```bash
devhibernate check payment-timeout -- pnpm test
```

Aturan:

* semua argumen setelah `--` dianggap command
* command dijalankan langsung tanpa shell
* stdout dan stderr disimpan
* exit code disimpan
* check tidak mengubah status capsule
* check bisa dijalankan saat capsule paused

Output berhasil:

```text
Check passed: payment-timeout

Command: pnpm test
Duration: 12.4s
Exit code: 0
```

Output gagal:

```text
Check failed: payment-timeout

Command: pnpm test
Duration: 8.1s
Exit code: 1

Log:
  ~/.devhibernate/capsules/.../checks/20260713-144100.log
```

Exit code DevHibernate harus mengikuti exit code command jika memungkinkan.

---

## 9.10 `devhibernate logs`

Menampilkan log service.

```bash
devhibernate logs payment-timeout
devhibernate logs payment-timeout api
devhibernate logs payment-timeout api --follow
devhibernate logs payment-timeout api --lines 100
```

Default:

* semua service
* 50 baris terakhir
* tanpa follow

---

## 9.11 `devhibernate handoff`

Menghasilkan Markdown report.

```bash
devhibernate handoff payment-timeout
```

Output file:

```text
.devhibernate/handoff/payment-timeout.md
```

Isi minimal:

````markdown
# Handoff: payment-timeout

## Status

Paused

## Current objective

Fix duplicate retries after gateway timeout.

## Git

- Branch: fix/payment-timeout
- Base: main
- Dirty: yes

## Changed files

- src/payment/retry-policy.ts
- tests/payment-retry.test.ts

## Last check

Command: `pnpm test payment-retry`
Result: failed
Exit code: 1

## Services

- api
- frontend

## How to continue

```bash
devhibernate resume payment-timeout
````

## Security

No environment variable values are included.

````

Handoff tidak boleh mengandung:

- nilai environment variable
- API key
- authorization header
- cookie
- password
- isi `.env`
- private key
- database password

---

## 9.12 `devhibernate delete`

Menghapus capsule.

```bash
devhibernate delete payment-timeout
````

Default behavior:

* menolak jika capsule running
* menolak jika worktree dirty
* tidak menghapus branch
* menghapus worktree
* menghapus state
* mempertahankan handoff dan archived logs

Untuk menghapus capsule running:

```bash
devhibernate pause payment-timeout
devhibernate delete payment-timeout
```

Force:

```bash
devhibernate delete payment-timeout --force
```

`--force` boleh:

* menghentikan process
* menghapus dirty worktree
* menghapus state

`--force` tidak boleh otomatis menghapus branch.

Menghapus branch harus command terpisah pada versi masa depan.

---

## 9.13 `devhibernate doctor`

Memeriksa instalasi dan state.

```bash
devhibernate doctor
```

Pemeriksaan:

* Git tersedia
* repository valid
* config valid
* worktree root dapat ditulis
* state root dapat ditulis
* capsule state dapat dibaca
* PID masih valid
* worktree masih terdaftar
* branch masih ada
* log directory dapat ditulis
* lock stale
* orphan process
* port conflict

Output:

```text
DevHibernate Doctor

✓ Git available
✓ Configuration valid
✓ State directory writable
✓ Worktree directory writable
! Capsule checkout-redesign has a stale PID
! Capsule hotfix-login is missing its worktree

2 issues detected.
```

Doctor tidak boleh memperbaiki otomatis tanpa flag.

Future flag:

```bash
devhibernate doctor --repair
```

Tidak termasuk MVP.

---

## 9.14 `devhibernate version`

```bash
devhibernate version
```

Output:

```text
devhibernate 0.1.0
commit: abc1234
built: 2026-07-13
go: go1.xx
```

---

# 10. Exit Codes

Gunakan kontrak berikut:

| Exit code | Arti                                        |
| --------: | ------------------------------------------- |
|       `0` | sukses                                      |
|       `1` | command atau operation gagal                |
|       `2` | penggunaan CLI atau konfigurasi tidak valid |
|       `3` | capsule atau resource tidak ditemukan       |
|       `4` | operasi ditolak karena tidak aman           |
|       `5` | dependency sistem tidak tersedia            |
|      `10` | internal error                              |

Jangan menggunakan exit code acak.

---

# 11. State Machine

Status capsule:

```text
preparing
running
pausing
paused
resuming
error
deleting
```

Transisi utama:

```text
ABSENT
  ↓ start
PREPARING
  ↓ success
RUNNING
  ↓ pause
PAUSING
  ↓ success
PAUSED
  ↓ resume
RESUMING
  ↓ success
RUNNING
  ↓ pause
PAUSED
  ↓ delete
DELETING
  ↓ success
ABSENT
```

Transisi error:

```text
PREPARING → ERROR
PAUSING   → ERROR
RESUMING  → ERROR
DELETING  → ERROR
```

## Aturan state

* hanya satu command lifecycle boleh berjalan per capsule
* state transition harus ditulis sebelum operasi berisiko
* final state hanya ditulis setelah operasi selesai
* operasi gagal harus menyimpan `lastError`
* `pause` terhadap paused harus idempotent
* `resume` terhadap running harus idempotent
* `delete` terhadap capsule tidak ditemukan menghasilkan exit code `3`
* `start` dengan nama capsule yang sudah ada harus ditolak

---

# 12. Arsitektur

```text
CLI
 │
 ▼
Application Layer
 ├── Config Service
 ├── Repository Service
 ├── Capsule Service
 ├── Process Service
 ├── Health Service
 ├── State Store
 ├── Port Allocator
 ├── Check Runner
 ├── Handoff Generator
 └── Doctor Service
```

## Komponen

### CLI layer

Bertanggung jawab untuk:

* parsing argument
* validation awal
* output terminal
* exit code

CLI layer tidak boleh berisi business logic berat.

### Application layer

Mengorkestrasi operasi:

* start
* pause
* resume
* delete
* handoff
* check

### Git adapter

Bertanggung jawab untuk:

* menemukan repository root
* menemukan default branch
* membuat branch
* membuat worktree
* membaca Git status
* membaca changed files
* menghapus worktree
* validasi branch

### Process manager

Bertanggung jawab untuk:

* menjalankan service
* mengelola process group
* menyimpan PID
* graceful shutdown
* force shutdown
* redirect stdout/stderr
* memeriksa process masih hidup

### State store

Bertanggung jawab untuk:

* membaca state
* atomic write
* lock
* versioning state schema
* migrasi state di masa depan

### Health checker

Tipe health check MVP:

* `none`
* `process`
* `tcp`
* `http`

### Report generator

Bertanggung jawab untuk:

* handoff Markdown
* safe redaction
* changed file list
* last check summary
* continuation command

---

# 13. Struktur Repository

```text
devhibernate/
├── cmd/
│   └── devhibernate/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   ├── start.go
│   │   ├── pause.go
│   │   ├── resume.go
│   │   ├── delete.go
│   │   ├── check.go
│   │   └── handoff.go
│   │
│   ├── cli/
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── start.go
│   │   ├── pause.go
│   │   ├── resume.go
│   │   └── output.go
│   │
│   ├── config/
│   │   ├── config.go
│   │   ├── load.go
│   │   ├── validate.go
│   │   └── template.go
│   │
│   ├── capsule/
│   │   ├── model.go
│   │   ├── state_machine.go
│   │   └── paths.go
│   │
│   ├── git/
│   │   ├── git.go
│   │   ├── repository.go
│   │   ├── branch.go
│   │   └── worktree.go
│   │
│   ├── process/
│   │   ├── manager.go
│   │   ├── start_unix.go
│   │   ├── stop_unix.go
│   │   ├── start_windows.go
│   │   └── stop_windows.go
│   │
│   ├── health/
│   │   ├── checker.go
│   │   ├── http.go
│   │   └── tcp.go
│   │
│   ├── ports/
│   │   └── allocator.go
│   │
│   ├── state/
│   │   ├── store.go
│   │   ├── atomic.go
│   │   └── lock.go
│   │
│   ├── checks/
│   │   └── runner.go
│   │
│   ├── report/
│   │   ├── handoff.go
│   │   └── redact.go
│   │
│   ├── doctor/
│   │   └── doctor.go
│   │
│   └── version/
│       └── version.go
│
├── test/
│   ├── integration/
│   ├── fixtures/
│   └── helpers/
│
├── docs/
│   ├── architecture.md
│   ├── configuration.md
│   ├── security.md
│   ├── testing.md
│   └── roadmap.md
│
├── examples/
│   ├── node/
│   ├── go/
│   └── fullstack/
│
├── .github/
│   ├── workflows/
│   │   ├── ci.yml
│   │   └── release.yml
│   ├── ISSUE_TEMPLATE/
│   └── pull_request_template.md
│
├── .gitignore
├── AGENTS.md
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── LICENSE
├── README.md
├── SECURITY.md
├── go.mod
└── go.sum
```

---

# 14. Konfigurasi

Nama file:

```text
.devhibernate.json
```

Gunakan JSON pada MVP agar parser tersedia dari Go standard library.

Contoh:

```json
{
  "version": 1,
  "defaults": {
    "baseBranch": "main",
    "branchPrefix": "task/",
    "gracefulShutdownSeconds": 5,
    "healthTimeoutSeconds": 30
  },
  "setup": [
    {
      "command": ["go", "mod", "download"]
    }
  ],
  "services": {
    "api": {
      "command": ["go", "run", "./cmd/api"],
      "workingDirectory": ".",
      "environment": {
        "PORT": "${PORT:api}"
      },
      "inheritEnvironment": [
        "DATABASE_URL"
      ],
      "health": {
        "type": "http",
        "url": "http://127.0.0.1:${PORT:api}/health",
        "timeoutSeconds": 30
      }
    }
  },
  "checks": {
    "test": {
      "command": ["go", "test", "./..."]
    },
    "vet": {
      "command": ["go", "vet", "./..."]
    }
  }
}
```

## Aturan command

Command harus berupa array:

```json
["pnpm", "dev"]
```

Jangan menerima string shell:

```json
"pnpm dev && echo done"
```

Tujuannya:

* menghindari shell injection
* menjaga argument parsing
* lebih mudah diuji
* lebih mudah dijalankan lintas platform

Future option untuk shell harus eksplisit:

```json
{
  "shell": true,
  "command": "pnpm dev"
}
```

Tidak termasuk MVP.

---

# 15. State Data

Lokasi global:

```text
~/.devhibernate/
```

Struktur:

```text
~/.devhibernate/
├── capsules/
│   └── <repository-id>/
│       └── <capsule-name>/
│           ├── state.json
│           ├── capsule.lock
│           ├── logs/
│           ├── checks/
│           └── handoffs/
│
└── worktrees/
    └── <repository-name>/
        └── <capsule-name>/
```

## Contoh state

```json
{
  "schemaVersion": 1,
  "name": "payment-timeout",
  "status": "paused",
  "repositoryRoot": "/home/user/projects/shop",
  "repositoryID": "3fcb4c...",
  "worktreePath": "/home/user/.devhibernate/worktrees/shop/payment-timeout",
  "branch": "fix/payment-timeout",
  "baseBranch": "main",
  "createdAt": "2026-07-13T10:00:00Z",
  "updatedAt": "2026-07-13T12:00:00Z",
  "lastPausedAt": "2026-07-13T12:00:00Z",
  "currentNote": "Investigate duplicate retry",
  "services": {
    "api": {
      "status": "stopped",
      "command": ["go", "run", "./cmd/api"],
      "pid": 0,
      "processGroupID": 0,
      "port": 43102,
      "logPath": ".../logs/api.log",
      "lastStartedAt": "2026-07-13T10:01:00Z",
      "lastStoppedAt": "2026-07-13T12:00:00Z"
    }
  },
  "lastCheck": {
    "command": ["go", "test", "./..."],
    "exitCode": 1,
    "startedAt": "2026-07-13T11:30:00Z",
    "finishedAt": "2026-07-13T11:30:10Z",
    "logPath": ".../checks/20260713-113000.log"
  },
  "lastError": null
}
```

## State write requirements

* gunakan temporary file
* lakukan `fsync` jika tersedia
* rename temporary file ke target
* permission file `0600`
* jangan menulis langsung ke file target
* jangan menyimpan secret
* gunakan schema version

---

# 16. Locking dan Concurrency

Setiap capsule memiliki lock:

```text
capsule.lock
```

Saat command lifecycle dimulai:

1. Coba membuat lock dengan exclusive create.
2. Jika lock sudah ada, baca PID pemilik lock.
3. Jika PID hidup, tolak operation.
4. Jika PID mati, anggap stale lock.
5. MVP hanya melaporkan stale lock.
6. Jangan menghapus stale lock otomatis kecuali operasi aman.

Output:

```text
Capsule is busy: payment-timeout
Another DevHibernate process is operating on it.
```

---

# 17. Git Worktree Behavior

## Start

Default branch:

```text
task/<capsule-name>
```

Contoh:

```text
task/payment-timeout
```

Jika branch belum ada:

```bash
git worktree add -b task/payment-timeout <path> main
```

Jika branch sudah ada dan belum digunakan worktree lain:

```bash
git worktree add <path> task/payment-timeout
```

Jika branch sedang digunakan:

* operasi harus ditolak
* jangan menggunakan `--force`
* tampilkan worktree yang menggunakan branch tersebut

## Delete

Sebelum delete:

```bash
git status --porcelain
```

Jika output tidak kosong:

```text
Cannot delete capsule: worktree has uncommitted changes.

Use:
  devhibernate status payment-timeout
  devhibernate delete payment-timeout --force
```

DevHibernate tidak boleh:

* commit otomatis
* stash otomatis
* reset otomatis
* checkout otomatis pada repository utama
* rebase otomatis
* merge otomatis
* push otomatis
* menghapus branch otomatis

---

# 18. Process Management

## Unix

Setiap service harus dijalankan dalam process group baru.

Tujuannya agar child process ikut dihentikan.

Start:

```go
SysProcAttr.Setpgid = true
```

Pause:

1. Ambil process group ID.
2. Kirim `SIGTERM` ke group.
3. Tunggu grace period.
4. Periksa process masih hidup.
5. Kirim `SIGKILL` jika diperlukan.

Jangan hanya menghentikan parent process.

## Windows

Untuk MVP:

* source harus dapat dikompilasi
* lifecycle process penuh boleh ditandai experimental
* implementasi penuh menggunakan Job Objects masuk post-MVP

README harus transparan:

```text
Linux and macOS are fully supported in v0.1.
Windows support is experimental.
```

## Process output

* stdout dan stderr ditulis ke file yang sama
* setiap startup membuat session marker
* log tidak boleh ditimpa tanpa marker
* log rotation belum diperlukan pada MVP
* maksimal log tidak dibatasi pada MVP, tetapi roadmap harus mencatatnya

Contoh marker:

```text
--- DevHibernate service start: 2026-07-13T10:01:00Z ---
```

---

# 19. Port Allocation

Service dapat meminta dynamic port menggunakan:

```text
${PORT:service-name}
```

Port allocator:

1. Bind ke `127.0.0.1:0`.
2. Ambil port yang diberikan OS.
3. Tutup temporary listener.
4. Gunakan port tersebut untuk service.
5. Simpan port di state.

Risiko race antara port release dan service bind harus didokumentasikan.

Future improvement:

* pass inherited listener
* local proxy
* reserved port daemon

Tidak perlu untuk MVP.

---

# 20. Health Checks

## `none`

Service dianggap berjalan setelah process berhasil dimulai.

## `process`

Service dianggap sehat jika process tetap hidup setelah periode minimum, misalnya 500 ms.

## `tcp`

Contoh:

```json
{
  "type": "tcp",
  "host": "127.0.0.1",
  "port": "${PORT:api}",
  "timeoutSeconds": 30
}
```

## `http`

Contoh:

```json
{
  "type": "http",
  "url": "http://127.0.0.1:${PORT:api}/health",
  "expectedStatus": 200,
  "timeoutSeconds": 30
}
```

Retry:

* interval 500 ms
* berhenti saat timeout
* tidak menggunakan exponential backoff pada MVP

Jika satu service gagal:

1. tandai service gagal
2. hentikan seluruh service yang sebelumnya berhasil dimulai
3. ubah capsule menjadi `error`
4. simpan log
5. jangan meninggalkan partial-running capsule

---

# 21. Security Requirements

## Secret handling

DevHibernate boleh menyimpan nama environment variable:

```json
["DATABASE_URL", "JWT_SECRET"]
```

DevHibernate tidak boleh menyimpan nilainya.

Saat service dijalankan:

1. Ambil nilai dari environment process DevHibernate.
2. Pass ke child process.
3. Jangan tulis ke state.
4. Jangan tampilkan ke terminal.
5. Jangan masukkan ke handoff.

## Redaction

Sebelum membuat handoff, redaksi pola:

* `Bearer ...`
* `Authorization: ...`
* `password=...`
* `token=...`
* `api_key=...`
* `secret=...`
* private key block
* database URL dengan credential

Handoff harus memiliki bagian:

```text
Security:
No environment variable values are included.
```

## Path security

Nama capsule harus diubah menjadi slug aman.

Valid:

```text
payment-timeout
GH-482
checkout_redesign
```

Invalid:

```text
../../etc
/task
task\name
.
..
```

Aturan slug:

* huruf
* angka
* dash
* underscore
* maksimum 64 karakter
* tidak boleh dimulai dengan titik
* tidak boleh mengandung separator path

## Command security

* tidak memakai `eval`
* tidak memakai `sh -c`
* tidak menggabungkan argument menjadi satu shell string
* gunakan `exec.Command`
* working directory harus berada dalam worktree
* tolak working directory yang keluar dari worktree

---

# 22. Error Handling

Setiap error harus memiliki:

1. Ringkasan masalah.
2. Penyebab teknis singkat.
3. Lokasi resource.
4. Tindakan yang dapat dilakukan.

Contoh buruk:

```text
operation failed
```

Contoh benar:

```text
Failed to start service: api

Command:
  go run ./cmd/api

Reason:
  Process exited with code 1 before passing its health check.

Log:
  ~/.devhibernate/capsules/shop/payment-timeout/logs/api.log

Next:
  Review the log, fix the problem, then run:
  devhibernate resume payment-timeout
```

---

# 23. CLI Output Guidelines

Gunakan output plain text yang tetap bagus tanpa warna.

Warna bersifat opsional jika terminal mendukung.

Simbol:

```text
✓ success
! warning
× failure
● running
○ stopped
```

Dukung environment:

```text
NO_COLOR=1
```

Output harus dapat dibaca oleh script.

Future flag:

```bash
--json
```

Tidak wajib di MVP kecuali waktu mencukupi setelah seluruh acceptance criteria selesai.

---

# 24. Implementation Phases

## Phase 0 — Repository Foundation

Task:

* buat module Go
* buat struktur folder
* tambahkan license
* tambahkan README awal
* tambahkan AGENTS.md
* tambahkan CI minimal
* implement version package

Acceptance criteria:

* `go test ./...` dapat dijalankan
* `go vet ./...` berhasil
* binary dapat dibuild
* `devhibernate version` bekerja

---

## Phase 1 — Configuration

Task:

* definisikan struct configuration
* implement JSON loader
* implement default values
* implement config validation
* implement `init`
* implement variable substitution untuk `${PORT:name}`

Acceptance criteria:

* config valid dapat dibaca
* unknown schema version ditolak
* duplicate service name tidak mungkin
* command kosong ditolak
* unsafe working directory ditolak
* test unit mencakup config invalid

---

## Phase 2 — Repository Discovery

Task:

* temukan Git root
* baca remote origin
* buat repository ID stabil
* baca branch aktif
* baca default branch
* implement Git command runner
* implement Git status

Acceptance criteria:

* bekerja dari subdirectory repository
* non-Git folder menghasilkan exit code `5`
* repository ID konsisten
* Git error tidak ditelan
* test menggunakan temporary Git repository

---

## Phase 3 — State Store

Task:

* definisikan capsule state
* implement path resolver
* implement atomic write
* implement load/save/delete
* implement lock
* implement schema version

Acceptance criteria:

* write interruption tidak merusak state lama
* state permission aman
* concurrent operation ditolak
* malformed state menghasilkan error yang jelas
* seluruh state test lulus dengan race detector

---

## Phase 4 — Git Worktree Lifecycle

Task:

* validasi capsule name
* generate branch name
* create branch
* create worktree
* detect existing worktree
* detect dirty worktree
* remove worktree
* preserve branch

Acceptance criteria:

* start membuat branch baru
* start membuat worktree
* duplicate capsule ditolak
* delete dirty worktree ditolak
* force-delete dapat menghapus worktree
* branch tetap ada setelah capsule dihapus
* repository utama tidak berpindah branch

---

## Phase 5 — Process Lifecycle

Task:

* service process model
* dynamic environment
* stdout/stderr log
* process group
* graceful termination
* force termination
* process liveness check

Acceptance criteria:

* process benar-benar hidup setelah start
* child process ikut dihentikan
* pause menghentikan seluruh service
* pause dua kali aman
* resume menghasilkan PID baru
* partial startup failure membersihkan service lain
* log tersedia

---

## Phase 6 — Health Checks dan Ports

Task:

* port allocator
* process health
* TCP health
* HTTP health
* timeout
* cleanup on failure

Acceptance criteria:

* service sehat dikenali
* service timeout menghasilkan error
* port tersedia diberikan ke environment
* semua service dihentikan jika satu gagal
* tidak ada capsule berstatus running dengan service gagal

---

## Phase 7 — User Context

Task:

* `note`
* note history
* `where`
* Git changed files
* recent modified files
* last check
* `status`
* `list`

Acceptance criteria:

* note tersimpan
* note lama masuk history
* where menampilkan note terbaru
* status menampilkan dirty state
* list mengurutkan berdasarkan updated time
* tidak ada secret di output

---

## Phase 8 — Check Runner

Task:

* parse argument setelah `--`
* menjalankan command dalam worktree
* menyimpan stdout/stderr
* menyimpan duration
* menyimpan exit code
* menampilkan summary

Acceptance criteria:

* successful check exit `0`
* failed check mengikuti exit code
* check dapat dijalankan ketika paused
* check log tersimpan
* last check muncul di status dan handoff

---

## Phase 9 — Handoff dan Doctor

Task:

* handoff Markdown
* redaction
* changed file list
* continuation instructions
* doctor checks
* stale PID detection
* missing worktree detection

Acceptance criteria:

* handoff dapat dibaca manusia
* handoff tidak berisi secret value
* doctor menemukan missing worktree
* doctor menemukan stale PID
* doctor tidak mengubah state

---

## Phase 10 — Documentation dan Release

Task:

* README final
* installation guide
* configuration guide
* architecture guide
* security guide
* examples
* contribution guide
* changelog
* release workflow

Acceptance criteria:

* pengguna baru dapat menjalankan demo dari README
* seluruh command terdokumentasi
* keterbatasan Windows disebutkan
* roadmap dipisahkan dari fitur selesai
* tag `v0.1.0` dapat menghasilkan binary release

---

# 25. Test Plan

## Unit tests

Wajib mencakup:

* capsule slug validation
* repository ID
* config parsing
* config validation
* environment substitution
* path containment
* state transitions
* atomic state write
* lock acquisition
* stale lock detection
* port allocation
* health timeout
* secret redaction
* handoff generation
* exit-code mapping

---

## Integration tests

Gunakan temporary directory dan temporary Git repository.

### Scenario 1 — Happy path

1. Buat Git repository.
2. Commit file awal.
3. Buat config.
4. Jalankan start.
5. Pastikan branch dibuat.
6. Pastikan worktree dibuat.
7. Pastikan process hidup.
8. Jalankan note.
9. Jalankan check.
10. Jalankan pause.
11. Pastikan process mati.
12. Jalankan resume.
13. Pastikan PID baru.
14. Jalankan handoff.
15. Jalankan delete.
16. Pastikan state hilang.
17. Pastikan branch tetap ada.

### Scenario 2 — Dirty deletion

1. Start capsule.
2. Pause capsule.
3. Ubah file.
4. Delete tanpa force.
5. Pastikan ditolak.
6. Delete dengan force.
7. Pastikan worktree dihapus.
8. Pastikan branch tetap ada.

### Scenario 3 — Setup failure

1. Config setup command gagal.
2. Start capsule.
3. Pastikan state `error`.
4. Pastikan service tidak berjalan.
5. Pastikan log tersedia.

### Scenario 4 — Service health failure

1. Service process hidup tetapi health endpoint gagal.
2. Start capsule.
3. Pastikan timeout.
4. Pastikan process dihentikan.
5. Pastikan capsule `error`.

### Scenario 5 — Partial startup

1. Service A berhasil.
2. Service B gagal.
3. Pastikan A dihentikan.
4. Pastikan tidak ada service tersisa.

### Scenario 6 — Idempotency

* pause paused capsule
* resume running capsule
* list saat state kosong
* logs untuk service yang belum pernah dimulai

### Scenario 7 — Secret safety

1. Set fake API key di environment.
2. Start service.
3. Buat handoff.
4. Periksa seluruh state dan handoff.
5. Pastikan nilai API key tidak ditemukan.

### Scenario 8 — Orphan state

1. Simpan PID yang tidak hidup.
2. Jalankan doctor.
3. Pastikan stale PID dilaporkan.

---

# 26. Race dan Reliability Testing

Jalankan:

```bash
go test -race ./...
```

Race test wajib untuk:

* state store
* lock
* process state
* parallel list/status
* concurrent read selama check berjalan

Jalankan test berulang:

```bash
go test ./... -count=20
```

Target:

* tidak ada flaky test
* tidak ada leaked process
* tidak ada leaked temporary worktree
* tidak ada file lock tertinggal

---

# 27. Manual Smoke Test

Sebelum release, lakukan:

```bash
mkdir demo-project
cd demo-project
git init
```

Buat aplikasi sederhana dengan HTTP server.

Kemudian:

```bash
devhibernate init
devhibernate start demo-task
devhibernate status demo-task
devhibernate note demo-task "Testing lifecycle"
devhibernate check demo-task -- go test ./...
devhibernate pause demo-task
devhibernate resume demo-task
devhibernate where demo-task
devhibernate handoff demo-task
devhibernate pause demo-task
devhibernate delete demo-task
```

Pastikan:

* semua output dapat dipahami
* service benar-benar mati setelah pause
* worktree tidak hilang saat pause
* state tidak mengandung secret
* handoff valid Markdown
* branch tidak dihapus

---

# 28. CI Pipeline

File:

```text
.github/workflows/ci.yml
```

Jobs:

## Format

```bash
gofmt -w
git diff --exit-code
```

CI tidak boleh memodifikasi file, sehingga gunakan pemeriksaan setara:

```bash
test -z "$(gofmt -l .)"
```

## Vet

```bash
go vet ./...
```

## Test

```bash
go test ./...
```

## Race

```bash
go test -race ./...
```

## Build

```bash
go build ./cmd/devhibernate
```

## Platform matrix

* Ubuntu latest
* macOS latest
* Windows latest untuk compile dan unit test non-process

Windows integration test process group dapat dilewati dengan alasan eksplisit.

---

# 29. Release Pipeline

Saat tag dibuat:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Release workflow membangun:

* Linux amd64
* Linux arm64
* macOS amd64
* macOS arm64
* Windows amd64 experimental

Artifact:

```text
devhibernate_0.1.0_linux_amd64.tar.gz
devhibernate_0.1.0_linux_arm64.tar.gz
devhibernate_0.1.0_darwin_amd64.tar.gz
devhibernate_0.1.0_darwin_arm64.tar.gz
devhibernate_0.1.0_windows_amd64.zip
checksums.txt
```

Gunakan version injection:

```bash
go build \
  -ldflags "-X internal/version.Version=0.1.0 \
            -X internal/version.Commit=$GITHUB_SHA \
            -X internal/version.BuildDate=$BUILD_DATE"
```

---

# 30. GitHub Repository Setup

Nama repository:

```text
devhibernate
```

Description:

```text
Pause and resume coding tasks with isolated Git worktrees and managed local services.
```

Topics:

```text
developer-tools
git-worktree
cli
golang
local-first
productivity
dev-environment
process-manager
```

Repository harus memiliki:

* README
* license
* contributing guide
* code of conduct
* security policy
* issue templates
* pull-request template
* CI badge
* release badge
* example configuration

---

# 31. Commit Strategy

AI agent harus membuat commit berdasarkan fase.

Contoh:

```text
chore: initialize Go CLI project
feat(config): add project configuration loader
feat(git): add repository and worktree management
feat(state): add atomic capsule state store
feat(process): add Unix process group management
feat(cli): implement start and pause commands
feat(cli): implement resume and delete commands
feat(context): add notes and where command
feat(checks): add command check runner
feat(report): add handoff generation
feat(doctor): add capsule diagnostics
test: add lifecycle integration suite
docs: add usage and architecture guides
ci: add multi-platform test workflow
```

Jangan membuat satu commit berisi seluruh project.

---

# 32. Pull Request Requirements

Sebelum membuka PR:

```bash
gofmt
go vet ./...
go test ./...
go test -race ./...
go build ./cmd/devhibernate
```

Pull request harus menjelaskan:

* fitur yang ditambahkan
* command yang berubah
* state schema yang berubah
* test yang ditambahkan
* risiko
* platform yang diuji
* bagian roadmap yang belum dikerjakan

Checklist PR:

```markdown
- [ ] Code formatted
- [ ] Vet passed
- [ ] Unit tests passed
- [ ] Race tests passed
- [ ] Integration tests passed
- [ ] No secrets stored
- [ ] No automatic Git commit/stash/reset
- [ ] Documentation updated
- [ ] Windows limitation documented
```

---

# 33. README Structure

README final harus memiliki urutan:

1. Logo atau nama project.
2. One-line explanation.
3. Demo terminal GIF.
4. Masalah yang diselesaikan.
5. Installation.
6. Quick start.
7. Example configuration.
8. Command reference.
9. Lifecycle explanation.
10. Security behavior.
11. Platform support.
12. Architecture summary.
13. Roadmap.
14. Contributing.
15. License.

Quick start:

```bash
devhibernate init
devhibernate start feature-a
devhibernate note feature-a "Working on checkout validation"
devhibernate pause feature-a
devhibernate resume feature-a
```

---

# 34. AGENTS.md

File ini wajib ada agar AI agent tidak keluar scope.

Isi minimal:

```markdown
# AI Agent Instructions

## Source of truth

Read these files before making changes:

1. AGENTS.md
2. docs/product-requirements.md
3. docs/architecture.md
4. docs/testing.md

## Hard constraints

- Do not add AI features.
- Do not add cloud dependencies.
- Do not automatically commit, stash, reset, rebase, merge, or push.
- Do not store environment variable values.
- Do not delete dirty worktrees without explicit force.
- Do not delete Git branches.
- Do not use shell evaluation for configured commands.
- Every lifecycle change requires an integration test.
- Every state schema change requires schema-version handling.
- Linux and macOS are the primary platforms for v0.1.

## Work sequence

1. Read the current milestone.
2. Implement only tasks in that milestone.
3. Add tests.
4. Run format, vet, test, race, and build.
5. Update documentation.
6. Stop after the milestone acceptance criteria pass.
```

---

# 35. Definition of Done MVP

Project dianggap selesai hanya jika:

## Functional

* semua command MVP tersedia
* worktree dapat dibuat
* process dapat dijalankan
* process dapat dihentikan
* process dapat dilanjutkan
* note tersimpan
* check tersimpan
* handoff dihasilkan
* dirty deletion ditolak
* doctor menemukan state bermasalah

## Safety

* tidak ada secret value di state
* tidak ada shell evaluation
* tidak ada automatic commit
* tidak ada automatic stash
* tidak ada automatic reset
* branch tidak dihapus
* partial startup dibersihkan
* dirty worktree tidak dihapus tanpa force

## Quality

* seluruh unit test lulus
* seluruh integration test lulus
* race detector lulus
* vet lulus
* build lulus
* tidak ada leaked process
* tidak ada flaky test setelah 20 pengulangan

## Documentation

* README lengkap
* configuration guide lengkap
* security guide lengkap
* architecture guide lengkap
* AGENTS.md lengkap
* roadmap jelas
* keterbatasan Windows transparan

## GitHub

* CI aktif
* pull request lulus
* repository publik
* tag `v0.1.0`
* binary release tersedia
* checksums tersedia

---

# 36. Roadmap Setelah MVP

## v0.2 — Stable Local URLs

Tambahkan internal reverse proxy:

```text
payment-timeout-api.localhost
payment-timeout-web.localhost
```

Port service boleh berubah, URL tetap sama.

## v0.3 — Editor Context

Integrasi awal:

```bash
code <worktree>
```

Simpan daftar file yang dibuka melalui workspace file, bukan membaca internal state editor.

## v0.4 — Capsule Export

```bash
devhibernate export payment-timeout
devhibernate import payment-timeout.capsule
```

Export tidak boleh membawa secret.

## v0.5 — GitHub Integration

```bash
devhibernate start GH-482
```

Membaca metadata issue melalui GitHub CLI jika tersedia.

## v1.0

Target stabil:

* Linux/macOS production-ready
* Windows process Job Objects
* proxy stabil
* import/export
* plugin hooks
* shell completion
* package manager distribution
* Homebrew
* Scoop
* Nix package

---

# 37. Instruksi Akhir untuk AI Agent

AI agent harus mengikuti urutan:

```text
Foundation
→ Configuration
→ Git discovery
→ State store
→ Worktree lifecycle
→ Process lifecycle
→ Health checks
→ Context commands
→ Check runner
→ Handoff
→ Doctor
→ Integration tests
→ Documentation
→ CI
→ Release
```

AI agent tidak boleh melompat ke:

* reverse proxy
* editor integration
* GitHub integration
* export/import
* GUI
* AI features

sebelum MVP mencapai Definition of Done.

Jika sebuah detail belum ditentukan:

1. Pilih solusi paling sederhana.
2. Jangan menambah dependency eksternal tanpa alasan kuat.
3. Prioritaskan Go standard library.
4. Prioritaskan keselamatan data.
5. Tambahkan test sebelum menganggap task selesai.
6. Dokumentasikan asumsi pada pull request.

Tujuan pertama bukan membuat project terlihat besar.

Tujuan pertama adalah menghasilkan lifecycle berikut secara reliabel:

```text
start
→ work
→ pause
→ switch task
→ resume
→ handoff
→ delete safely
```

Itulah fondasi utama DevHibernate.
