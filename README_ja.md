# Simple Gallery

軽量で堅牢な写真・動画管理アプリケーション。膨張したギャラリーアプリの代替として、本質的な機能に絞り込みました：写真・動画閲覧、アップロード、アカウント別アクセス制御。シングル Go バイナリで動作。Intel QSV ハードウェアエンコーディング。日本語全文検索対応。余ったインテル PC での家族サーバーに最適。

**[English version is here](README.md)**

## 主な機能

- **シングルバイナリ** — Go 1.26+ で単一の実行ファイル。Docker、Node.js、Python 不要
- **メディア管理** — 写真・動画のアップロード、タイムライン表示、アルバム管理
- **アカウント別アクセス制御** — argon2id 認証。ロールベース（admin/user）
- **ハードウェア高速化** — Intel QuickSync Video (QSV) による超高速エンコード
- **日本語全文検索** — SQLite FTS5 による日本語部分一致検索
- **PWA対応** — iPhone Safari、Chrome でホーム画面インストール可能
- **超軽量** — RAM 1-2GB で動作。余ったインテル PC に最適

## システム要件

### サーバー
- **CPU:** Intel Core i3 以上（QSV 対応には i5 以上推奨）
- **RAM:** 1-2 GB
- **ストレージ:** SSD + 外付けUSB ドライブ（RAID1 推奨）
- **OS:** Linux（Ubuntu 24.04 LTS、Debian 12）
- **外部依存:** FFmpeg のみ

### クライアント
- **ブラウザ:** Safari (iOS 16+)、Chrome、Firefox
- **PWA:** iOS 16.4+ でホーム画面追加対応

## クイックスタート

### 1. ビルド

```bash
make build
```

### 2. 設定

```bash
cp config.example.toml config.toml
# storage.media_root, storage.thumb_root, storage.db_path を編集
```

### 3. 起動

```bash
./bin/simple-gallery -config config.toml
```

ブラウザで http://localhost:8080 にアクセス。初回は管理者アカウント作成フローが表示されます。

### 4. 本番デプロイ

TLS 終端は Caddy（自動 HTTPS）または Tailscale を推奨：

```bash
caddy reverse-proxy --from https://photos.example.com --to http://localhost:8080
```

詳細は [docs/setup-guide.md](docs/setup-guide.md) を参照。

## システムアーキテクチャ

```
iPhone/Web → リバースプロキシ (Caddy/Tailscale) → Go ServeMux → SQLite + FTS5
                                                           ├→ FFmpeg (QSV)
                                                           └→ HTMX + Alpine.js + Tailwind
```

- **バックエンド:** Go 1.26+ `net/http.ServeMux`（外部ルータ不要）
- **データベース:** SQLite 3 with FTS5
- **フロントエンド:** HTMX 2.0.8 + Alpine.js + Tailwind CSS
- **メディア処理:** FFmpeg with Intel QSV ハードウェアエンコーディング

## コア機能一覧

| 機能 | 説明 |
|------|------|
| メディアアップロード | 写真・動画のアップロード。HEIC/HEIF ネイティブ対応。大容量ファイルのチャンク対応 |
| タイムライン表示 | EXIF 撮影日時でソート。月/年ごとのセクションヘッダー表示 |
| サムネイル自動生成 | AVIF + WebP の2種生成。Content Negotiation で最適配信 |
| BlurHash | サムネイル読み込み前のカラフルなプレースホルダー表示 |
| アルバム管理 | アルバムの作成・編集・削除。カバー画像設定 |
| 日本語検索 | SQLite FTS5 による部分一致検索 |
| 共有リンク | 有効期限・パスワード保護・アクセス回数制限付き共有 |
| 動画再生 | HLS ストリーミング。Safari / Chrome / Firefox 対応 |
| PWA | ホーム画面追加でスタンドアロンアプリ化 |

## 開発

### ビルドコマンド

```bash
make build      # ./bin/simple-gallery をビルド
make test       # go test ./... テスト実行
make clean      # ビルド成果物削除
make dev        # go run (開発モード)
make lint       # golangci-lint 実行
```

### プロジェクト構成

```
simple-gallery/
├── cmd/simple-gallery/          # エントリポイント
├── internal/
│   ├── config/                  # 設定管理
│   ├── server/                  # HTTP サーバー
│   ├── handler/                 # ハンドラ（認証、アップロード等）
│   ├── store/                   # SQLite CRUD
│   ├── auth/                    # 認証、CSRF、レート制限
│   ├── media/                   # FFmpeg、サムネイル、メタデータ
│   └── model/                   # データモデル
├── web/
│   ├── static/                  # CSS、JS、PWA 関連ファイル
│   └── templates/               # HTML テンプレート
├── deploy/                      # Systemd、Caddy 設定例
├── docs/                        # ドキュメント
└── Makefile, Dockerfile
```

### コーディング規約

- **Go 1.26+** の標準ライブラリを最大限活用
- **JSON:** `goccy/go-json` を使用（標準互換、高速）
- **FFmpeg:** `internal/media/ffmpeg.go` 経由で呼び出し（シェル実行禁止）
- **テスト:** `*_test.go` で単体・統合テストを実装
- **HTTPS:** バイナリは HTTP のみ提供。TLS はリバースプロキシに委譲

## コントリビューション

コントリビューションを歓迎します。

1. フォーク & 機能ブランチを作成
2. 機能・バグフィックスを実装
3. テスト通過確認（`make test`）
4. Lint 通過確認（`make lint`）
5. Pull Request を提出

詳細は [docs/contributing.md](docs/contributing.md) を参照。

## セキュリティ設計

- **認証:** argon2id（メモリハード、GPU耐性）
- **セッション:** UUID v4 + Secure Cookie（HttpOnly, SameSite=Strict）
- **CSRF対策:** トークンベース検証
- **レート制限:** ログイン試行の IP ベース制限
- **コマンドインジェクション対策:** シェル実行禁止
- **パストラバーサル対策:** `os.Root` + パス正規化
- **XSS対策:** `html/template` 自動エスケープ

詳細は [docs/security.md](docs/security.md) を参照。

## パフォーマンス最適化

- **SQLite チューニング:** WAL モード、メモリマップI/O
- **非同期サムネイル生成:** バックグラウンド処理で即座のレスポンス実現
- **ハードウェアエンコード:** Intel QSV で CPU 負荷を削減
- **HTTPキャッシュ:** Immutable キャッシュ戦略でブラウザキャッシュを最大活用
- **tmpfs for HLS:** 一時セグメントを RAM ディスク（`/tmp`）に出力し SSD 寿命を保護

## ライセンス

MIT or Apache-2.0

## ドキュメント

- [セットアップガイド](docs/setup-guide.md) — インストール・初期設定手順
- [USB ブートガイド](docs/usb-boot-guide.md) — Linux USB ブート環境構築
- [RAID 構築](docs/raid-setup.md) — RAID1 構成手順
- [セキュリティ詳細](docs/security.md) — セキュリティ設計詳細
- [コントリビューションガイド](docs/contributing.md) — 開発参加ガイド
- [完全仕様書](docs/SimpleGallery-Specifications_2.md) — 技術仕様（AI エージェント向け）
