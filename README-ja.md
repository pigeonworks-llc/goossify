# 🚀 goossify

**Go言語のOSSプロジェクトを瞬時に立ち上げる次世代ボイラープレートジェネレーター**

[![Go Report Card](https://goreportcard.com/badge/github.com/pigeonworks-llc/goossify)](https://goreportcard.com/report/github.com/pigeonworks-llc/goossify)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/release/pigeonworks-llc/goossify.svg)](https://github.com/pigeonworks-llc/goossify/releases)
[![CI](https://github.com/pigeonworks-llc/goossify/workflows/CI/badge.svg)](https://github.com/pigeonworks-llc/goossify/actions)

> ⚠️ **開発中のプロジェクトです** - 現在アクティブに開発中で、一部機能は実装予定です。[実装状況](#-実装状況)をご確認ください。

## 🎯 コンセプト

goossifyは、Go言語のOSSプロジェクトを**完全自動化**で立ち上げるツールです。
従来のボイラープレート生成を超えて、プロジェクトのライフサイクル全体を管理します。

## ✨ 特徴

### 🏗️ 完全自動化されたプロジェクト初期化
- 📁 最適なディレクトリ構造の生成
- 📄 必須ファイル群の自動作成（README、LICENSE、.gitignore等）
- 🔧 開発ツール設定（golangci-lint、GoReleaser、GitHub Actions等）
- 📊 品質管理ツールの統合（カバレッジ、ベンチマーク等）

### 🤖 継続的メンテナンス自動化
- 🔄 依存関係の自動更新と脆弱性監視
- 📈 リリース管理の完全自動化（セマンティックバージョニング）
- 📝 チェンジログとドキュメントの自動生成
- 👥 コミュニティファイルの管理

### 🌟 Go言語エコシステム最適化
- 🐹 Go Modules完全対応
- 📦 pkg.go.devへの自動インデックス化
- 🏆 Go言語のベストプラクティス適用
- 🇯🇵 日本語環境での開発を考慮

## 🚀 クイックスタート

### インストール

```bash
# 開発版をビルド (推奨)
git clone https://github.com/pigeonworks-llc/goossify.git
cd goossify
go build -o goossify

# または将来的にリリース版から
go install github.com/pigeonworks-llc/goossify@latest
```

### 新しいOSSプロジェクトの作成

```bash
# 🚧 実装予定: 対話モードで新規プロジェクト作成
goossify init my-awesome-project

# 🚧 実装予定: テンプレートから作成
goossify create --template cli-tool my-cli-app
goossify create --template library my-go-lib
```

### 既存プロジェクトのOSS化 ✅

```bash
# 既存プロジェクトをOSS対応に変換
cd existing-project
goossify ossify .

# カレントディレクトリをOSS化
goossify ossify

# このgoossifyプロジェクト自体もdog foodingで実行済み
# 生成されたファイル: LICENSE, CONTRIBUTING.md, SECURITY.md, .github/workflows/ci.yml
```

## 📈 推奨ワークフロー: プライベート → パブリック

### 1. **プライベートリポジトリで開発開始** 🔒

```bash
# GitHubでプライベートリポジトリを作成
git clone git@github.com:yourusername/your-project.git
cd your-project

# OSS準備ファイルを生成
goossify ossify .
```

### 2. **GitHub Personal Access Token 発行** 🔑

GitHubで Personal Access Token を発行してください：
- GitHub Settings → Developer settings → Personal access tokens → Generate new token
- 必要な権限: `repo`, `read:org` (ブランチ保護設定チェック用)

### 3. **段階的品質向上** 📊

```bash
# 基本的な健全性チェック
goossify status .

# GitHub設定も含む完全チェック
goossify status --github --github-token ghp_xxxxxxxxxxxx .

# スコア100/100を目指して改善
# 不足項目があれば再度 ossify で修正
goossify ossify .
```

### 4. **パブリック公開準備完了チェック** 🚀

```bash
# 公開準備完了かの最終確認
goossify ready --github-token ghp_xxxxxxxxxxxx .

# トークンなしでも基本チェックは可能
goossify ready .
```

### 5. **パブリック化** 🌍

- GitHubでリポジトリを **Public** に変更
- 初回リリースタグを作成
- pkg.go.dev での自動インデックス化を確認

### 💡 このワークフローの利点

- **🔒 セキュリティ**: 機密情報の誤公開を防止
- **📈 品質保証**: 100%準備完了してから公開
- **🔄 反復改善**: プライベート環境で安全に試行錯誤
- **🤖 CI/CD確認**: GitHub Actionsが正常動作することを事前確認

## 📁 生成されるプロジェクト構造

```
my-project/
├── 📄 README.md                 # 包括的なプロジェクト説明
├── 📜 LICENSE                   # ライセンスファイル
├── 🐹 go.mod                    # Go Modules設定
├── 🔧 .golangci.yml             # リンター設定
├── 🚀 .goreleaser.yml           # リリース自動化設定
├── 📋 .goossify.yml             # goossify管理設定
├── 🔄 .github/
│   ├── workflows/               # GitHub Actions CI/CD
│   ├── ISSUE_TEMPLATE/          # Issueテンプレート
│   ├── PULL_REQUEST_TEMPLATE.md # PRテンプレート
│   ├── CONTRIBUTING.md          # コントリビューションガイド
│   ├── CODE_OF_CONDUCT.md       # 行動規範
│   └── SECURITY.md              # セキュリティポリシー
├── 📚 docs/                     # ドキュメント
├── 🧪 cmd/                      # エントリーポイント
├── 📦 internal/                 # 内部パッケージ
├── 🔧 pkg/                      # 公開パッケージ
├── ✅ tests/                    # テスト
└── 📊 examples/                 # 使用例
```

## 🎛️ 設定ファイル例

`.goossify.yml`：

```yaml
project:
  name: "my-awesome-project"
  description: "An awesome Go project"
  type: "cli-tool"  # cli-tool, library, web-api, service
  license: "MIT"
  author: "Your Name"
  email: "your.email@example.com"

automation:
  release:
    strategy: "auto"              # auto, manual, scheduled
    versioning: "semantic"       # semantic, calendar
    changelog: true
    goreleaser: true

  dependencies:
    auto_update: true
    security_check: true
    schedule: "weekly"

  quality:
    coverage_threshold: 80
    benchmarks: true
    performance_regression: true

  community:
    issue_templates: true
    pr_template: true
    contributing_guide: true
    code_of_conduct: true
    security_policy: true

integrations:
  github:
    branch_protection: true
    required_reviews: 2
    status_checks: ["test", "lint", "security"]

  ci_cd:
    provider: "github-actions"
    go_versions: ["1.21", "1.22", "1.23"]
    platforms: ["linux", "darwin", "windows"]

monitoring:
  health_checks: true
  dependency_updates: "weekly"
  security_scans: "daily"
```

## 📋 利用可能なテンプレート

### 🔧 CLI Tool
コマンドラインツール用の完全なボイラープレート
- Cobra CLI framework
- サブコマンド対応
- 設定ファイル管理
- バイナリ配布設定

### 📚 Library
Go言語ライブラリ用のボイラープレート
- 公開API設計
- サンプルコード
- ベンチマーク
- パッケージドキュメント

## 🤖 自動化機能

### リリース管理
- コミット履歴からの自動バージョン判定
- チェンジログ自動生成
- GitHub/GitLab Release作成
- バイナリのクロスコンパイル・配布

### 品質管理
- コードカバレッジ監視
- パフォーマンスベンチマーク
- セキュリティスキャン
- 依存関係脆弱性チェック

### コミュニティ管理
- Issue/PRテンプレート自動更新
- コントリビューター認識
- ドキュメント自動生成
- コミュニティ健全性レポート

## 🔧 コマンドリファレンス

### ✅ 実装済み機能
```bash
# OSS化機能
goossify ossify [path]                 # 既存プロジェクトのOSS化
```

### 🚧 実装予定機能
```bash
# プロジェクト初期化
goossify init [project-name]           # 対話的プロジェクト作成
goossify create --template cli-tool NAME    # CLI ツールテンプレート
goossify create --template library NAME     # ライブラリテンプレート

# GitHub 統合
goossify github setup                  # GitHub API による設定自動化

# 管理・メンテナンス
goossify status                        # プロジェクト健全性確認
goossify deps update                   # 依存関係更新
```

## 🏆 なぜgoossify？

### 従来の課題
- ❌ プロジェクト立ち上げに時間がかかる
- ❌ ベストプラクティスの適用が困難
- ❌ 継続的なメンテナンスが負担
- ❌ Go言語特有の設定が複雑
- ❌ コミュニティ管理が手作業

### goossifyの解決策
- ✅ 数分でプロダクション品質のプロジェクト作成
- ✅ Go言語のベストプラクティス自動適用
- ✅ 完全自動化されたメンテナンス
- ✅ Go言語エコシステム最適化
- ✅ コミュニティ管理の自動化

## 📊 実装状況

### ✅ 完了済み機能
- **OSS化コマンド** (`goossify ossify`)
  - LICENSE ファイル自動生成 (MIT License)
  - GitHub CI/CD workflow 生成
  - コミュニティファイル生成 (CONTRIBUTING.md, SECURITY.md)
  - Issue/PR テンプレート生成 (ファイルベース)
  - Git リポジトリ初期化
  - 依存関係整理 (`go mod tidy`)
  - Renovate設定ファイル生成
  - **Dog Fooding**: このgoossifyプロジェクト自体が `goossify ossify .` で自動OSS化済み

### 🚧 開発中機能
- **基本テンプレート生成**
  - CLI Tool テンプレート
  - Library テンプレート
- **GitHub API統合**
  - ブランチ保護ルール設定
  - ラベル自動設定
  - リポジトリ設定自動化

### 📋 実装予定機能 (優先度順)

#### v0.2.0 - Core Features
- [ ] `goossify create` - CLI/Library テンプレート生成
- [ ] `goossify init` - 対話的プロジェクト作成
- [ ] `goossify status` - プロジェクト健全性チェック
- [ ] 基本設定ファイル生成 (.golangci.yml, .goreleaser.yml)

#### v0.3.0 - GitHub Integration
- [ ] `goossify github setup` - GitHub API統合完成
- [ ] ブランチ保護ルール自動設定
- [ ] ラベル・リポジトリ設定自動化
- [ ] GitHub Releases自動化

#### v0.4.0 - Quality & Automation
- [ ] `goossify deps update` - 依存関係更新
- [ ] 静的解析ツール統合 (golangci-lint, govulncheck)
- [ ] テストカバレッジレポート
- [ ] セキュリティスキャン自動化

#### v1.0.0 - Production Ready
- [ ] 包括的テストスイート
- [ ] パフォーマンス最適化
- [ ] 詳細ドキュメント
- [ ] 安定版API

### 🔮 将来の構想 (v1.x)
- [ ] AI支援コード生成 (spec-kit統合検討)
- [ ] 多言語対応 (英語版README等)
- [ ] カスタムテンプレート機能

## 🤝 コントリビューション

goossifyの開発に参加してください！

1. 🍴 このリポジトリをフォーク
2. 🌟 フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 💫 変更をコミット (`git commit -m 'Add amazing feature'`)
4. 📤 ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. 🎉 Pull Requestを作成

詳細は[CONTRIBUTING.md](CONTRIBUTING.md)をご覧ください。

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は[LICENSE](LICENSE)ファイルをご覧ください。

## 💖 サポート

- 🐛 [Issue報告](https://github.com/pigeonworks-llc/goossify/issues)
- 💡 [機能要求](https://github.com/pigeonworks-llc/goossify/issues)
- 💬 [ディスカッション](https://github.com/pigeonworks-llc/goossify/discussions)
- 📧 [メール](mailto:support@pigeonworks.llc)

---

**goossify で、Go言語のOSS開発を次のレベルへ 🚀**