package github

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func PostPRComment(owner, repo string, prNumber string, body string) error {
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        return errors.New("GITHUB_TOKEN is required")
    }

    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s/comments", owner, repo, prNumber)
    payload := map[string]string{"body": body}
    data, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("marshal comment payload: %w", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewReader(data))
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Accept", "application/vnd.github+json")
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("failed to post PR comment: %s", resp.Status)
    }
    return nil
}

func DetectPullRequest() (owner string, repo string, prNumber string, err error) {
    repoSlug := os.Getenv("GITHUB_REPOSITORY")
    if repoSlug == "" {
        return "", "", "", errors.New("GITHUB_REPOSITORY not set")
    }

    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        return "", "", "", errors.New("GITHUB_TOKEN not set")
    }

    eventPath := os.Getenv("GITHUB_EVENT_PATH")
    if eventPath == "" {
        return "", "", "", errors.New("GITHUB_EVENT_PATH not set")
    }

    data, err := os.ReadFile(filepath.Clean(eventPath))
    if err != nil {
        return "", "", "", fmt.Errorf("read event file: %w", err)
    }

    var event struct {
        PullRequest struct {
            Number int `json:"number"`
        } `json:"pull_request"`
    }
    if err := json.Unmarshal(data, &event); err != nil {
        return "", "", "", fmt.Errorf("unmarshal event: %w", err)
    }

    if event.PullRequest.Number == 0 {
        return "", "", "", errors.New("not a pull_request event")
    }

    parts := strings.Split(repoSlug, "/")
    if len(parts) != 2 {
        return "", "", "", fmt.Errorf("invalid GITHUB_REPOSITORY: %q", repoSlug)
    }
    owner, repo = parts[0], parts[1]
    prNumber = fmt.Sprintf("%d", event.PullRequest.Number)
    return owner, repo, prNumber, nil
}
