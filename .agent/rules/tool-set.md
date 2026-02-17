---
trigger: always_on
---

# プロジェクト共通ルール

## 1. 開発コンテキストの特定

- `Orders/` に含まれるファイルはユーザーである私の領分なので、変更したり勝手にコミットしたりしないでください。
- `Orders/` 内のファイルを実行する際は必ずはじめに「ガイドライン」に類する項目の有無をチェックし、存在する場合は読み込んで理解を示してください。
- 変更対象のファイルパスに基づき、以下のスタックを前提として動作すること。
- **backend/**
  - Language: Go
  - Framework: Echo, Google Wire
  - ORM: Ent
  - Validation: go-playground/validator
  - Testing: testing, stretchr/testify
  - Container: Podman
- **frontend/**
  - Language: TypeScript
  - Library: Next.js
  - UI Framework: shadcn ui, Tailwind CSS
  - HTTP Client: Axios

## 2. ツール選定と提案

### backend

- **ディレクトリ固有ルールの遵守:** `go.mod` および `mise.toml` を参照し、定義されているツール・設定を最優先すること。
- **一貫性の維持:** 既存のスタック（例： Ent 利用）がある場合、それと重複する機能を持つ別ライブラリを導入しないこと。
- **改善提案:** パフォーマンスや型安全性の観点で、現在のスタックよりも適したライブラリがある場合は、実装前に提案すること。

### frontend

- **ディレクトリ固有ルールの遵守:** `package.json` および `tsconfig.json` を参照し、定義されているツール・設定を最優先すること。
- **一貫性の維持:** 既存のスタック（例： shadcn 利用）がある場合、それと重複する機能を持つ別ライブラリを導入しないこと。
- **改善提案:** パフォーマンスや型安全性の観点で、現在のスタックよりも適したライブラリがある場合は、実装前に提案すること。

## 3. パッケージマネージャー

### Go

- tool は `go install` ではなく `mise use` を使って `mise.toml` に記録されるようにインストールする

### Node

- **pnpm の徹底:** パッケージ操作には必ず `pnpm` を使用すること。
- `npm` や `yarn`、`npx` は使用禁止。
  - `npm install` -> `pnpm add`
  - `npx <command>` -> `pnpm dlx <command>`
  - `npm run <script>` -> `pnpm <script>`

## 4. コード品質と一貫性

- **周辺参照:** 新規ファイル作成や編集の際は、必ず同階層や類似機能を持つ既存ファイルを参照し、プロジェクトの設計パターン（命名規則、ディレクトリ構成、エラーハンドリング）を模倣すること。
- **型定義:** `any` などの不明確な型を避け、可能な限り厳密な型定義を行うこと。
- **DI:** Google Wire を使う前提で実装を行うこと。
- **パスエイリアス:** `tsconfig.json` で設定されているパスエイリアスを優先して使用すること。

## 5. 言語

- 英語でのコミュニケーションを求められない限り日本語を使用すること。
