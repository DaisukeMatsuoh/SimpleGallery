# Simple Gallery

Lightweight alternative to bloated gallery apps. Stripped down to essentials: photo/video browsing, upload, per-account access control. Single Go binary. Hardware encoding (Intel QSV). Japanese full-text search. Perfect for family servers on spare Intel machines.

**[日本語版はこちら](README_ja.md)**

## Features

- **Single Binary** — Go 1.26+ single binary. No Docker, Node.js, or Python required
- **Media Management** — Upload, organize, and browse photos and videos with timeline view and albums
- **Per-Account Access Control** — argon2id authentication with role-based access (admin/user)
- **Hardware Acceleration** — Intel QuickSync Video (QSV) for ultra-fast encoding
- **Full-Text Search** — SQLite FTS5 with support for Japanese character search
- **PWA Ready** — Install as standalone app on iPhone Safari and Chrome
- **Lightweight** — Runs on 1-2 GB RAM on spare Intel machines with USB storage

## System Requirements

### Server
- **CPU:** Intel Core i3 or newer (i5+ recommended for QSV support)
- **RAM:** 1-2 GB
- **Storage:** SSD + External USB drives (RAID1 recommended)
- **OS:** Linux (Ubuntu 24.04 LTS, Debian 12)
- **External Dependency:** FFmpeg only

### Client
- **Browser:** Safari (iOS 16+), Chrome, Firefox
- **PWA:** iOS 16.4+ for home screen installation

## Quick Start

### 1. Build

```bash
make build
```

### 2. Configure

```bash
cp config.example.toml config.toml
# Edit storage.media_root, storage.thumb_root, storage.db_path
```

### 3. Run

```bash
./bin/simple-gallery -config config.toml
```

Visit http://localhost:8080 in your browser. On first run, you'll see the admin account setup flow.

### 4. Production Deployment

Use Caddy (with auto HTTPS) or Tailscale for TLS termination:

```bash
caddy reverse-proxy --from https://photos.example.com --to http://localhost:8080
```

See [docs/setup-guide.md](docs/setup-guide.md) for details.

## Architecture

```
iPhone/Web → Reverse Proxy (Caddy/Tailscale) → Go ServeMux → SQLite + FTS5
                                                           ├→ FFmpeg (QSV)
                                                           └→ HTMX + Alpine.js + Tailwind
```

- **Backend:** Go 1.26+ `net/http.ServeMux` (no external router needed)
- **Database:** SQLite 3 with FTS5
- **Frontend:** HTMX 2.0.8 + Alpine.js + Tailwind CSS
- **Media Processing:** FFmpeg with Intel QSV hardware encoding

## Core Features

| Feature | Description |
|---------|-------------|
| Media Upload | Photos and videos. Native HEIC/HEIF support. Chunked uploads for large files |
| Timeline View | Sorted by EXIF capture date with month/year section headers |
| Auto Thumbnails | AVIF + WebP generation with Content Negotiation for optimal delivery |
| BlurHash | Colorful placeholder while thumbnails load |
| Album Management | Create, edit, and organize albums with custom cover images |
| Full-Text Search | SQLite FTS5 with partial-match support |
| Share Links | Configurable expiration, password protection, and access limits |
| Video Playback | HLS streaming compatible with Safari, Chrome, and Firefox |
| PWA | Install to home screen for standalone app experience |

## Development

### Build Commands

```bash
make build      # ./bin/simple-gallery
make test       # go test ./...
make clean      # Clean build artifacts
make dev        # go run (development mode)
make lint       # golangci-lint
```

### Project Structure

```
simple-gallery/
├── cmd/simple-gallery/          # Entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── server/                  # HTTP server
│   ├── handler/                 # Handlers
│   ├── store/                   # SQLite CRUD
│   ├── auth/                    # Authentication, CSRF, rate limiting
│   ├── media/                   # FFmpeg, thumbnails, metadata
│   └── model/                   # Data models
├── web/
│   ├── static/                  # CSS, JS, PWA assets
│   └── templates/               # HTML templates
├── deploy/                      # Systemd, Caddy examples
├── docs/                        # Documentation
└── Makefile, Dockerfile
```

### Code Standards

- **Go 1.26+** leveraging standard library
- **JSON:** Use `goccy/go-json` (standard-compatible, faster)
- **FFmpeg:** Go through `internal/media/ffmpeg.go` (shell execution prohibited)
- **Testing:** Unit and integration tests in `*_test.go` files
- **HTTPS:** Binary serves HTTP only; delegate TLS to reverse proxy

## Contributing

Contributions are welcome!

1. Fork & create a feature branch
2. Implement changes
3. Verify tests pass (`make test`)
4. Verify linting passes (`make lint`)
5. Submit a Pull Request

See [docs/contributing.md](docs/contributing.md) for details.

## Security

- **Authentication:** argon2id (memory-hard, GPU-resistant)
- **Sessions:** UUID v4 with Secure Cookies (HttpOnly, SameSite=Strict)
- **CSRF Protection:** Token-based validation
- **Rate Limiting:** IP-based login attempt throttling
- **Command Injection:** Shell execution prohibited
- **Path Traversal:** `os.Root` with path normalization
- **XSS:** `html/template` auto-escaping

See [docs/security.md](docs/security.md) for details.

## Performance

- **SQLite Tuning:** WAL mode, memory-mapped I/O
- **Async Thumbnail Generation:** Background processing for instant upload response
- **Hardware Encoding:** Intel QSV reduces CPU load
- **HTTP Caching:** Immutable cache strategy with content hashing
- **tmpfs for HLS:** Temporary segments on RAM disk to protect SSD lifespan

## License

MIT or Apache-2.0

## Documentation

- [Setup Guide](docs/setup-guide.md) — Installation and configuration
- [USB Boot Guide](docs/usb-boot-guide.md) — Linux USB boot setup
- [RAID Setup](docs/raid-setup.md) — RAID1 configuration
- [Security Details](docs/security.md) — Security architecture
- [Contributing](docs/contributing.md) — Contribution guidelines
- [Full Specifications](docs/SimpleGallery-Specifications_2.md) — Complete technical spec
