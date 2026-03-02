# Frontend Survey Result

Frontendディレクトリ (`frontend/`) のリファクタリング箇所の調査結果を以下にまとめます。

## 1. コンポーネントの肥大化とロジックの混在

- **`frontend/src/app/page.tsx` (`HomeContent`)**
    - 状態管理、API通信 (CRUD)、トースト通知、ページネーション、モーダル制御など、約200行に渡るロジックが1つのコンポーネントに集中しています。
    - **改善案**: 
        - カスタムフック (`useTodos`) を作成し、状態管理とAPI通信ロジックを分離する。
        - フィルタリング (完了したものを非表示) やトースト通知のロジックを分離する。
- **`TodoModal.tsx` / `LoginPage.tsx`**
    - フォームの状態管理とバリデーションが手動で行われています。

## 2. API通信・通信基盤の課題

- **`frontend/src/lib/api.ts`**
    - `baseURL` が `http://localhost:8080` とハードコードされています。環境変数 (`NEXT_PUBLIC_API_URL`) を使用すべきです。
    - `X-CSRF-Token` の取得が `document.cookie` からの直接抽出であり、やや不安定な実装です。
- **コンポーネント内での直接通信**
    - 複数のコンポーネント (`page.tsx`, `login/page.tsx`) で `axios` (api インスタンス) を直接呼び出しています。
    - **改善案**: Service層を作成し、API通信を抽象化する。

## 3. バリデーションと型定義

- **Zodの未使用**
    - `package.json` に `zod` が含まれていますが、現状のコードベースでは使用されていません。
    - **改善案**: `TodoModal` や `LoginPage` のフォームバリデーションに `zod` を導入し、スキーマに基づいたバリデーションを行う。
- **型安全性の向上**
    - 一部で `as Error` や型アサーション (`error as { ... }`) が見られます。`axios` のエラーハンドリングや、独自の型ガードを導入することで安全性を高められます。

## 4. ディレクトリ構成とコンポーネントの分類

- **`src/app/components`**
    - shadcn/ui由来の純粋なUIコンポーネント (`ui/`) と、ドメインに依存するコンポーネント (`TodoItem`, `TodoModal`) が混在しています。
    - **改善案**: `src/components/ui` (共通UI) と `src/components/todo` (ドメイン固有) のように整理する。
- **UI/UXの一貫性**
    - `page.tsx` 内に `toastHandler` という独自のラッパー関数が定義されています。トースト通知のデフォルト設定 (duration等) は、共通化または `Toaster` の設定で一括管理すべきです。

## 5. その他

- **`suppressHydrationWarning`**
    - `RootLayout` で使用されていますが、不要な警告を抑制している可能性があります。原因となっているコンポーネントを特定し、適切に修正できるか検討が必要です。
- **Dead Dependencies**
    - `zod` が現状未使用であるため、導入するか削除するかの判断が必要です (前述の通りバリデーションへの導入を推奨)。

## 6. ディレクトリ固有ルールへの準拠状況

- `pnpm` は `pnpm-lock.yaml` および `pnpm-workspace.yaml` が存在し、適切に利用されています。
- `eslint`, `prettier` の設定ファイルが存在し、基本的なコード品質は保たれています。
