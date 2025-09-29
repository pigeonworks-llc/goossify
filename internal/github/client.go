package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client performs GitHub API operations
type Client struct {
	client *github.Client
	owner  string
	repo   string
	ctx    context.Context
}

// Config is GitHub API configuration
type Config struct {
	Token string
	Owner string
	Repo  string
}

// RepositorySettings is repository settings
type RepositorySettings struct {
	BranchProtection    BranchProtectionSettings
	Labels              []LabelConfig
	RequiredReviews     int
	StatusChecks        []string
	AutoMerge           bool
	DeleteBranchOnMerge bool
}

// BranchProtectionSettings is branch protection settings
type BranchProtectionSettings struct {
	Branch                  string
	RequiredStatusChecks    []string
	RequiredReviews         int
	DismissStaleReviews     bool
	RequireCodeOwnerReviews bool
	RestrictPushes          bool
}

// LabelConfig is label configuration
type LabelConfig struct {
	Name        string
	Color       string
	Description string
}

// NewClient creates a new GitHub API client
func NewClient(config Config) (*Client, error) {
	if config.Token == "" {
		return nil, fmt.Errorf("GitHub token is required")
	}

	if config.Owner == "" || config.Repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Client{
		client: client,
		owner:  config.Owner,
		repo:   config.Repo,
		ctx:    ctx,
	}, nil
}

// SetupRepository configures basic repository settings
func (c *Client) SetupRepository(settings *RepositorySettings) error {
	// 1. Configure labels
	if err := c.setupLabels(settings.Labels); err != nil {
		return fmt.Errorf("failed to setup labels: %w", err)
	}

	// 2. Configure branch protection
	if err := c.setupBranchProtection(settings.BranchProtection); err != nil {
		return fmt.Errorf("failed to setup branch protection: %w", err)
	}

	// 3. Configure general repository settings
	if err := c.updateRepositorySettings(settings); err != nil {
		return fmt.Errorf("failed to setup repository settings: %w", err)
	}

	return nil
}

// setupLabels configures repository labels
func (c *Client) setupLabels(labels []LabelConfig) error {
	// Get existing labels
	existingLabels, _, err := c.client.Issues.ListLabels(c.ctx, c.owner, c.repo, nil)
	if err != nil {
		return fmt.Errorf("failed to get existing labels: %w", err)
	}

	existingLabelMap := make(map[string]*github.Label)
	for _, label := range existingLabels {
		existingLabelMap[*label.Name] = label
	}

	// Create or update labels
	for i := range labels {
		labelConfig := labels[i]
		label := &github.Label{
			Name:        &labelConfig.Name,
			Color:       &labelConfig.Color,
			Description: &labelConfig.Description,
		}

		if existingLabel, exists := existingLabelMap[labelConfig.Name]; exists {
			// Update existing label
			label.ID = existingLabel.ID
			_, _, err := c.client.Issues.EditLabel(c.ctx, c.owner, c.repo, labelConfig.Name, label)
			if err != nil {
				return fmt.Errorf("failed to update label (%s): %w", labelConfig.Name, err)
			}
		} else {
			// Create new label
			_, _, err := c.client.Issues.CreateLabel(c.ctx, c.owner, c.repo, label)
			if err != nil {
				return fmt.Errorf("failed to create label (%s): %w", labelConfig.Name, err)
			}
		}
	}

	return nil
}

// setupBranchProtection configures branch protection
func (c *Client) setupBranchProtection(config BranchProtectionSettings) error {
	if config.Branch == "" {
		config.Branch = "main"
	}

	protection := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   true,
			Contexts: config.RequiredStatusChecks,
		},
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews:          config.DismissStaleReviews,
			RequireCodeOwnerReviews:      config.RequireCodeOwnerReviews,
			RequiredApprovingReviewCount: config.RequiredReviews,
		},
		EnforceAdmins: true,
	}

	if config.RestrictPushes {
		protection.Restrictions = &github.BranchRestrictionsRequest{
			Users: []string{},
			Teams: []string{},
		}
	}

	_, _, err := c.client.Repositories.UpdateBranchProtection(c.ctx, c.owner, c.repo, config.Branch, protection)
	if err != nil {
		return fmt.Errorf("failed to setup branch protection: %w", err)
	}

	return nil
}

// updateRepositorySettings updates general repository settings
func (c *Client) updateRepositorySettings(settings *RepositorySettings) error {
	repo := &github.Repository{
		DeleteBranchOnMerge: &settings.DeleteBranchOnMerge,
	}

	_, _, err := c.client.Repositories.Edit(c.ctx, c.owner, c.repo, repo)
	if err != nil {
		return fmt.Errorf("failed to update repository settings: %w", err)
	}

	return nil
}

// GetDefaultLabels returns default label configuration
func GetDefaultLabels() []LabelConfig {
	return []LabelConfig{
		{Name: "bug", Color: "d73a4a", Description: "Something isn't working"},
		{Name: "documentation", Color: "0075ca", Description: "Improvements or additions to documentation"},
		{Name: "duplicate", Color: "cfd3d7", Description: "This issue or pull request already exists"},
		{Name: "enhancement", Color: "a2eeef", Description: "New feature or request"},
		{Name: "good first issue", Color: "7057ff", Description: "Good for newcomers"},
		{Name: "help wanted", Color: "008672", Description: "Extra attention is needed"},
		{Name: "invalid", Color: "e4e669", Description: "This doesn't seem right"},
		{Name: "question", Color: "d876e3", Description: "Further information is requested"},
		{Name: "wontfix", Color: "ffffff", Description: "This will not be worked on"},
		{Name: "priority: high", Color: "b60205", Description: "High priority issue"},
		{Name: "priority: low", Color: "0e8a16", Description: "Low priority issue"},
		{Name: "type: breaking", Color: "B60205", Description: "Breaking change"},
		{Name: "type: feature", Color: "0e8a16", Description: "New feature"},
		{Name: "type: refactor", Color: "fbca04", Description: "Code refactoring"},
	}
}

// ParseRepositoryURL extracts owner/repo from GitHub URL
func ParseRepositoryURL(url string) (owner, repo string, err error) {
	// Expected format: https://github.com/owner/repo.git
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "git@github.com:")

	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	return parts[0], parts[1], nil
}
