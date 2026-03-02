---
description: バックエンド インフラ・ミドルウェア リファクタリング
---

# ガイドライン: バックエンド インフラ・ミドルウェア リファクタリング

対象の担当者: Backend-Infra-Agent

## 目的

バックエンドの環境変数設定、ミドルウェア、バリデーションロジックの改善を行います。
他のエージェントと並行稼働するため、コンフリクトを避けるべく、指定された【対象ファイル】以外は極力変更しないでください。

## 対象ファイル

- `backend/providers/ent_provider.go`
- `backend/server.go`
- `backend/routes/router.go`
- `backend/routes/todo.go`
- `backend/middleware/auth.go`
- `backend/validators/*.go`
- `backend/app_errors/*.go`
  （分離のために新規作成するファイルは自由に作成して問題ありません）

## やるべきこと

1. **設定値の環境変数化**:
   - 環境変数の設定は `backend/envs/` 内のファイルで行ってください。
   - `ent_provider.go` のDB接続文字列。
   - `router.go` のCORS許可オリジン。
   - `auth.go` のJWTシークレットのエッジケース（フォールバック）値。
   - `server.go` のポート番号 (`:8080`)。
2. **CSP Nonce の動的生成**:
   - `router.go` 等で固定値 (`'nonce-random123'`) となっている箇所を、リクエストごとに動的に生成するよう修正してください。
3. **バリデーションロジックの改善 (go-playground/validator)**:
   - `validators/` 配下に存在する、構造体フィールドと日本語エラーメッセージを手動で紐付ける巨大な switch 文を廃止してください。
   - `go-playground/validator` の Universal Translator 機能などを使い、タグによるエコシステムを用いたエラーメッセージのマッピングに切り替えてください。
4. **テストの通過（必須要件）**:
   - 実装作業後、必ず `backend/` ディレクトリに移動して `go test ./...` を実行し、**すべてのテストがパスすることを確認**してから作業を完了してください。
   - **注意: 絶対に `git commit` や `git add` 等のgit操作は実行しないでください。** 修正のコミットや `pre-commit` フックの実行はすべてのエージェントの作業完了後にメインエージェントがまとめて行います。コードの修正とテスト確認のみを行って終了してください。
