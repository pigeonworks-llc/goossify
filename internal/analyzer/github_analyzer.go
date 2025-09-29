package analyzer

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// GitHubAnalyzer はGitHub固有の分析を行う
type GitHubAnalyzer struct {
	client *github.Client
	owner  string
	repo   string
}

// GitHubCheck はGitHub設定のチェック結果
type GitHubCheck struct {
	BranchProtection  Item `json:"branch_protection"`
	RequiredReviews   Item `json:"required_reviews"`
	StatusChecks      Item `json:"status_checks"`
	AdminEnforcement  Item `json:"admin_enforcement"`
	DefaultBranch     Item `json:"default_branch"`
	SelfReviewAllowed Item `json:"self_review_allowed"`
}

// NewGitHubAnalyzer は新しいGitHubAnalyzerを作成
func NewGitHubAnalyzer(token string) (*GitHubAnalyzer, error) {
	if token == "" {
		return nil, fmt.Errorf("GitHub token is required")
	}

	// Git remote URLからowner/repo取得
	owner, repo, err := getRepositoryInfo()
	if err != nil {
		return nil, fmt.Errorf("リポジトリ情報取得失敗: %w", err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubAnalyzer{
		client: client,
		owner:  owner,
		repo:   repo,
	}, nil
}

// AnalyzeGitHubSettings はGitHub設定を分析
func (g *GitHubAnalyzer) AnalyzeGitHubSettings() (*GitHubCheck, error) {
	ctx := context.Background()

	// デフォルトブランチ取得
	repository, _, err := g.client.Repositories.Get(ctx, g.owner, g.repo)
	if err != nil {
		return nil, fmt.Errorf("リポジトリ情報取得失敗: %w", err)
	}

	defaultBranch := repository.GetDefaultBranch()
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	result := &GitHubCheck{
		DefaultBranch: Item{
			Name:        fmt.Sprintf("デフォルトブランチ: %s", defaultBranch),
			Status:      "present",
			Required:    false,
			Description: "リポジトリのデフォルトブランチ",
		},
	}

	// ブランチ保護設定取得
	protection, resp, err := g.client.Repositories.GetBranchProtection(ctx, g.owner, g.repo, defaultBranch)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// ブランチ保護が設定されていない
			result.BranchProtection = Item{
				Name:        "ブランチ保護",
				Status:      "missing",
				Required:    false,
				Description: "mainブランチの保護設定",
			}
			result.RequiredReviews = Item{
				Name:        "必須レビュー",
				Status:      "missing",
				Required:    false,
				Description: "プルリクエストの必須レビュー設定",
			}
			result.StatusChecks = Item{
				Name:        "ステータスチェック",
				Status:      "missing",
				Required:    false,
				Description: "CI/CDチェックの必須化",
			}
			result.AdminEnforcement = Item{
				Name:        "管理者強制",
				Status:      "missing",
				Required:    false,
				Description: "管理者に対するルール強制",
			}
			result.SelfReviewAllowed = Item{
				Name:        "セルフレビュー設定",
				Status:      "missing",
				Required:    false,
				Description: "ブランチ保護が設定されていません",
			}
		} else {
			return nil, fmt.Errorf("ブランチ保護設定取得失敗: %w", err)
		}
	} else {
		// ブランチ保護が設定されている
		result.BranchProtection = Item{
			Name:        "ブランチ保護",
			Status:      "present",
			Required:    false,
			Description: "mainブランチの保護設定",
		}

		// 必須レビュー設定
		if protection.GetRequiredPullRequestReviews() != nil {
			reviews := protection.GetRequiredPullRequestReviews()
			reviewCount := reviews.RequiredApprovingReviewCount
			result.RequiredReviews = Item{
				Name:        fmt.Sprintf("必須レビュー (%d人)", reviewCount),
				Status:      "present",
				Required:    false,
				Description: "プルリクエストの必須レビュー設定",
			}

			// セルフレビュー許可の確認（RequireCodeOwnerReviewsフィールドを使用）
			if reviews.RequireCodeOwnerReviews {
				result.SelfReviewAllowed = Item{
					Name:        "セルフレビュー制限",
					Status:      "present",
					Required:    false,
					Description: "コードオーナーレビューが必須",
				}
			} else {
				result.SelfReviewAllowed = Item{
					Name:        "セルフレビュー許可",
					Status:      "warning",
					Required:    false,
					Description: "セルフレビューが可能な設定",
				}
			}
		} else {
			result.RequiredReviews = Item{
				Name:        "必須レビュー",
				Status:      "missing",
				Required:    false,
				Description: "プルリクエストの必須レビュー設定",
			}
		}

		// ステータスチェック設定
		if protection.GetRequiredStatusChecks() != nil {
			checks := protection.GetRequiredStatusChecks()
			checkCount := len(checks.Contexts)
			if checkCount > 0 {
				result.StatusChecks = Item{
					Name:        fmt.Sprintf("ステータスチェック (%d個)", checkCount),
					Status:      "present",
					Required:    false,
					Description: "CI/CDチェックの必須化",
				}
			} else {
				result.StatusChecks = Item{
					Name:        "ステータスチェック",
					Status:      "missing",
					Required:    false,
					Description: "CI/CDチェックの必須化",
				}
			}
		} else {
			result.StatusChecks = Item{
				Name:        "ステータスチェック",
				Status:      "missing",
				Required:    false,
				Description: "CI/CDチェックの必須化",
			}
		}

		// 管理者強制設定
		if protection.GetEnforceAdmins() != nil && protection.GetEnforceAdmins().Enabled {
			result.AdminEnforcement = Item{
				Name:        "管理者強制",
				Status:      "present",
				Required:    false,
				Description: "管理者に対するルール強制",
			}
		} else {
			result.AdminEnforcement = Item{
				Name:        "管理者強制",
				Status:      "missing",
				Required:    false,
				Description: "管理者に対するルール強制",
			}
		}
	}

	return result, nil
}

// getRepositoryInfo はGit remote URLからowner/repoを取得
func getRepositoryInfo() (owner, repo string, err error) {
	// git remote get-url origin
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("git remote URL 取得失敗: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))
	return parseRepositoryURL(remoteURL)
}

// parseRepositoryURL はGitHub URLからowner/repoを抽出
func parseRepositoryURL(url string) (owner, repo string, err error) {
	// https://github.com/owner/repo.git の形式を想定
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "git@github.com:")

	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	return parts[0], parts[1], nil
}
