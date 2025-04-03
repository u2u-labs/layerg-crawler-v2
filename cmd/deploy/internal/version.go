package internal

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
)

type GitCommitVerifier struct {
	Client              *github.Client
	AllowedRepositories []string
	AllowedBranches     []string
}

// NewGitCommitVerifier creates a new verifier with GitHub API access
func NewGitCommitVerifier(githubToken string) *GitCommitVerifier {
	ctx := context.Background()

	// Create authenticated client if token provided
	var client *github.Client
	if githubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &GitCommitVerifier{
		Client: client,
		AllowedRepositories: []string{
			"u2u-labs/layerg-crawler-v2",
		},
		AllowedBranches: []string{
			"master",
			"main",
			"develop",
			"dev",
			"ft/code_fetching",
		},
	}
}

func (gcv *GitCommitVerifier) VerifyCommit(commitHash string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Verifying commit:", commitHash)
	// Parse repository owner and name
	for _, repoFullName := range gcv.AllowedRepositories {
		parts := strings.Split(repoFullName, "/")
		if len(parts) != 2 {
			continue
		}
		owner, repo := parts[0], parts[1]

		// Check if commit exists and is in allowed branches
		if gcv.isCommitInAllowedBranches(ctx, owner, repo, commitHash) {
			return true
		}
	}

	return false
}

func (gcv *GitCommitVerifier) isCommitInAllowedBranches(
	ctx context.Context,
	owner, repo, commitHash string,
) bool {
	// Check each allowed branch
	for _, branch := range gcv.AllowedBranches {
		logger.Info("Checking branch:", branch, " for commit:", commitHash, "in repository:", owner+"/"+repo)
		if gcv.isCommitInBranch(ctx, owner, repo, branch, commitHash) {
			return true
		}
	}
	return false
}

func (gcv *GitCommitVerifier) isCommitInBranch(
	ctx context.Context,
	owner, repo, branch, commitHash string,
) bool {
	// Fetch branch reference
	branchRef, _, err := gcv.Client.Repositories.GetBranch(
		ctx,
		owner,
		repo,
		branch,
		1,
	)
	if err != nil {
		return false
	}

	// Compare commit hash with branch's latest commit
	if branchRef.Commit == nil {
		return false
	}

	// If the commit is part of the branch's history
	comparison, _, err := gcv.Client.Repositories.CompareCommits(
		ctx,
		owner,
		repo,
		*branchRef.Commit.SHA,
		commitHash,
		&github.ListOptions{},
	)
	if err != nil {
		return false
	}

	// Check if commit is reachable from branch head
	return comparison.GetStatus() == "behind" ||
		comparison.GetStatus() == "identical"
}

// Verification Service Integration
type BinaryVerificationService struct {
	gitVerifier *GitCommitVerifier
}

// BuildMetadata stores build-time Git information
type BuildMetadata struct {
	GitCommit    string
	GitBranch    string
	BuildTime    time.Time
	BuildVersion string
}

func (bvs *BinaryVerificationService) VerifyBinaryMetadata(metadata *BuildMetadata) bool {
	// Validate Git commit
	return bvs.gitVerifier.VerifyCommit(metadata.GitCommit)
}

// get current git commit hash
func getGitCommitHash() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

type VersionChecker struct {
}

func NewVersionChecker() *VersionChecker {
	return &VersionChecker{}

}

func (vc *VersionChecker) CheckVersion(binaryPath string) bool {
	v := NewGitCommitVerifier("")

	// Execute the binary with "version" argument
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Convert output to string
	outputStr := string(output)

	// Extract Git Commit hash using regex
	re := regexp.MustCompile(`Git Commit:\s*([a-f0-9]+)`)
	matches := re.FindStringSubmatch(outputStr)
	logger.Info("Matches:", matches)
	if len(matches) < 2 {
		return false // Commit hash not found
	}

	commitHash := matches[1]
	return v.VerifyCommit(commitHash)
}
