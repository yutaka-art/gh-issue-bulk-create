# gh-issue-bulk-create

![](https://github.com/ntsk/gh-issue-bulk-create/actions/workflows/ci.yml/badge.svg)

複数のGitHub Issueを一括作成するGitHub CLI拡張機能です。

## 前提条件

この拡張機能を使用するには、gh CLIがインストールされ、`PATH`に設定されている必要があります。また、`gh auth`でユーザー認証が完了している必要があります。

## インストール

このプロジェクトはGitHub CLI拡張機能です。`gh` CLIをインストール後、コマンドラインから以下を実行してください：

```bash
gh extension install ntsk/gh-issue-bulk-create
```

## 使用方法

この拡張機能は、テンプレートマークダウンファイルとデータを含むCSVファイルを使用して、複数のGitHub Issueを一括作成します。

```bash
gh issue-bulk-create --template <template_file> --csv <csv_file> [--repo <owner/repo>] [--dry-run]
```

### オプション

- `--template`: テンプレートマークダウンファイルのパス（必須）
- `--csv`: データを含むCSVファイルのパス（必須）
- `--repo`: 対象リポジトリ（owner/repo形式）（デフォルト: 現在のリポジトリ）
- `--dry-run`: Issueを実際に作成せずに内容のみを表示

### テンプレートファイル

テンプレートファイルは、ファイルの先頭にフロントマターメタデータを含むGitHub Issueテンプレート形式に従います：

```markdown
---
title: "{{title}}"
labels: "{{label1}}, {{label2}}"
assignees: "{{assignee}}"
---

## 概要
{{description}}

## 再現手順
{{steps}}
```

テンプレート内でCSVファイルのデータを埋め込むために、Mustache記法（`{{variable_name}}`）を使用できます。

### CSVファイル

CSVファイルには**ヘッダー行が必須**で、テンプレートで使用する変数名と一致する列名を含んでいる必要があります。
最初の行がヘッダー行で、それ以降の各行が個別のIssue作成に使用されます。

```csv
title,label1,label2,assignee,description,steps
ログインページエラー,bug,frontend,username,ログインボタンクリック時にエラーが発生,ログインボタンをクリック
検索が機能しない,bug,backend,username,"検索時に結果が表示されない","検索ボックスに""test""と入力して検索ボタンをクリック"
```

#### CSV要件

- ファイルはRFC 4180仕様に従った標準的なカンマ区切り値（CSV）形式である必要があります
- ヘッダーは必須で、CSVファイルの最初の行に配置する必要があります
- 各ヘッダー（列名）は空であってはならず、テンプレートで使用される変数と一致している必要があります
- カンマ、改行、ダブルクォートを含むフィールドは、ダブルクォートで囲む必要があります
- クォートされたフィールド内のダブルクォートは、二重にしてエスケープする必要があります（例：`"`は`""`になります）
- 特殊文字を含む適切にフォーマットされたCSVの例：
  ```csv
  title,description
  "カンマ, を含むタイトル","ダブル""クォート""を含む説明"
  "改行
  を含むテキスト","別のフィールド"
  ```

#### 警告動作
- テンプレートで使用されていないCSVヘッダーがある場合：警告が表示されますが、処理は続行されます
- 対応するCSVヘッダーがないテンプレート変数がある場合：警告が表示され、続行するかどうかの確認が求められます。続行する場合、それらの不足している変数は生成されるIssueで空のままになります

## 例

リポジトリに含まれているサンプルファイルで試すことができます：

```bash
# gitリポジトリディレクトリ内から
gh issue-bulk-create --template sample-template.md --csv sample-data.csv --dry-run

# または異なるリポジトリを指定
gh issue-bulk-create --template sample-template.md --csv sample-data.csv --repo owner/repo-name
```