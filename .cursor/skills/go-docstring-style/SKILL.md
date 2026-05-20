---
name: go-docstring-style
description: Enforces docstring writing in this Go project for functions, methods, and struct fields, with approved exceptions and multiline parameter supplements. Use when editing Go code, adding APIs, or refactoring functions, methods, and structs.
---

# Go Docstring Style

このgoのプロジェクトにおいて、関数、メソッド、または構造体（Struct）のフィールドにおいてはdocstringを書くようにしてください。

## 適用対象

- すべての関数
- すべてのメソッド
- 構造体のフィールド

`exported` な要素は常に docstring 必須です。

## 例外（承認済み）

以下のみ例外として docstring を省略できます。

1. 明らかな一時変数用の匿名構造体フィールド
2. enum のように `mode` が複数分かれる説明箇所

## 引数説明の特別ルール（承認済み）

例外的に引数の説明を文章の中で改行を入れて補足してください。

運用ルール:

- 引数説明が長くなる場合は、docstring 本文で改行して補足を入れる
- 特に `mode` 引数は、取りうる値と挙動差分を改行付きで明示する

## 書き方ルール

1. 先頭で要約を1文で書く
2. 必要なら詳細を複数行で補足する
3. 引数の意味、単位、前提条件、失敗条件を必要に応じて書く
4. 挙動が分岐する `mode` は値ごとの差を明示する

## 最低テンプレート

```go
// DoThing は X を実行する。
//
// mode は実行モードを表す。
// "fast": 検証を簡略化して高速実行する。
// "safe": 追加検証を行って安全に実行する。
func DoThing(mode string) error {
    // ...
    return nil
}
```

## 実装時チェック

- docstring が不足している要素を追加したか
- 例外が承認済み2条件に該当するかを確認したか
- `mode` など分岐引数の説明を改行補足で記述したか
