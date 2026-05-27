---
name: design-to-shadcn-css
description: Reads a specified DESIGN.md and syncs its color tokens into shadcn/ui CSS variables (:root and .dark) in globals.css or index.css. Use when the user asks to apply DESIGN.md to shadcn theme variables, sync design tokens to CSS, or map semantic colors (background, foreground, primary, destructive, border, input).
disable-model-invocation: true
---

# DESIGN.md → shadcn/ui CSS 変数

指定されたDESIGN.md の内容を読み取り、指定されたglobals.css（または index.css）にある shadcn/ui のCSS変数定義（:root と .dark）に正確に反映してください。shadcn/uiが要求するセマンティックカラー（background, foreground, primary, destructive, border, input等）に適切にマッピングして上書きしてください。

## 入力の特定

ユーザーがパスを省略した場合は確認する。両方指定されている場合はそのパスを使う。

| 入力 | 例 |
|------|-----|
| デザイントークン | リポジトリルートの `DESIGN.md`、またはユーザー指定パス |
| CSS ターゲット | `src/globals.css`、`src/index.css` などユーザー指定パス |

## ワークフロー

1. **DESIGN.md を読む**
   - 先頭 YAML の `colors:`（および必要なら `rounded:`）を抽出する
   - 本文の `## Colors` などでセマンティック用途（CTA・hairline・dark-only 等）を補足として参照する

2. **ターゲット CSS を読む**
   - `:root` と `.dark` ブロック内の `--*` 変数だけを更新対象とする
   - `@theme inline`・`@layer base`・`@import` は触らない（既存の shadcn/Tailwind 配線を維持）

3. **マッピングして上書き**
   - ターゲットファイルが `oklch(...)` 形式なら、hex を **oklch に変換**して既存形式に揃える（変換不能な場合のみ hex を使う）

4. **検証**
   - `:root` / `.dark` の両方に、shadcn が参照する主要変数が欠けていないか確認する
   - DESIGN が dark-only と明記している場合は「ダークモード方針」に従う（下記）

## ダークモード方針

DESIGN.md に light 用パレットが無く「near-black canvas」「no light-mode」とある場合:

- **推奨**: `:root` と `.dark` の両方に同じダークパレットを適用する（アプリが常にダーク UI のとき）
- ユーザーがライト/ダーク切替を求める場合のみ、`:root` をライト用に分離する（そのときは DESIGN から推論できる範囲で対比を作る）

## 編集ルール

- **上書きのみ**: 既存の変数名・ブロック構造・宣言順は可能な限り維持する
- **触らない**: `--font-*`、`@custom-variant`、chart/sidebar 以外の `@theme` キー（ユーザーが chart/sidebar 同期を明示した場合を除く）
- **radius**: DESIGN の `rounded:` があるときのみ `--radius` を更新（例: `md: 6px` → `0.375rem`）
- **M トライカラー**（`m-blue-light` 等）は CTA/背景に使わない — `chart-*` や装飾用 CSS 変数（プロジェクトで定義済みの場合）に限定

## 完了チェックリスト

- [ ] DESIGN.md の `colors:` トークンを漏れなく参照した
- [ ] `:root` と `.dark` の `--background` 〜 `--sidebar-ring`（存在するもの）を更新した
- [ ] `primary` / `destructive` / `border` / `input` / `muted-foreground` が用途と矛盾しない
- [ ] 色形式がターゲット CSS の既存表記（oklch 等）と一致している

## 追加リソース

- トークン対応表・oklch 変換・chart/sidebar の割当: [reference.md](reference.md)
