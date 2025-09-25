# SnapShare - 写真共有アプリケーション

結婚式や特別なイベントで、ゲストが写真を簡単に共有できるWebアプリケーションです。

## 🏗️ アーキテクチャ

- **フロントエンド**: Next.js 15 (TypeScript, Tailwind CSS)
- **バックエンド**: Go (Echo, GORM)
- **データベース**: PostgreSQL
- **ストレージ**: MinIO (S3互換)
- **開発環境**: Docker Compose

## 🚀 クイックスタート

### Docker Compose での起動

```bash
# 全サービス起動
docker compose up --build

# バックグラウンド起動
docker compose up -d --build
```

### アクセス方法

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **MinIO管理画面**: http://localhost:9001 (minioadmin/minioadmin)

### 開発用サンプルデータ

以下のイベントコードでテストできます：
- `WEDDING1` - 山田太郎・花子 結婚式
- `TRAVEL02` - 田中家 家族旅行
- `REUNION2` - 高校同窓会 2025

## 📋 利用フロー

1. **イベント参加**: QRコード読み取り → 名前入力
2. **写真アップロード**: ドラッグ&ドロップで複数ファイル対応
3. **写真共有**: リアルタイムで他のゲストと共有

## 🛠️ 技術スタック

### フロントエンド
- Next.js 15, TypeScript, Tailwind CSS
- Biome, Vitest, Zustand, React Hook Form

### バックエンド
- Go 1.23, Echo, GORM, JWT認証

### インフラ
- Docker Compose, PostgreSQL, MinIO