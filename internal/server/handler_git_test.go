package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"skillshare/internal/testutil"
)

func TestHandleGitCommit_NoRemote_CreatesLocalCommit(t *testing.T) {
	s, src := newTestServer(t)
	initServerGitRepo(t, src)
	addSkill(t, src, "local-skill")

	req := httptest.NewRequest(http.MethodPost, "/api/git/commit", strings.NewReader(`{"message":"local checkpoint"}`))
	rr := httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp pushResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success response: %+v", resp)
	}

	message := testutil.RunGit(t, src, "log", "-1", "--pretty=%s")
	if message != "local checkpoint" {
		t.Fatalf("commit message = %q, want %q", message, "local checkpoint")
	}

	status := testutil.RunGit(t, src, "status", "--porcelain")
	if status != "" {
		t.Fatalf("expected clean working tree, got %q", status)
	}
}

func TestHandleGitCommit_DryRun_DoesNotCreateCommit(t *testing.T) {
	s, src := newTestServer(t)
	initServerGitRepo(t, src)
	addSkill(t, src, "dry-run-skill")

	req := httptest.NewRequest(http.MethodPost, "/api/git/commit", strings.NewReader(`{"message":"dry run checkpoint","dryRun":true}`))
	rr := httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	count := testutil.RunGit(t, src, "rev-list", "--count", "HEAD")
	if count != "1" {
		t.Fatalf("commit count = %q, want 1", count)
	}
}

func TestHandleGitCommit_NoChanges(t *testing.T) {
	s, src := newTestServer(t)
	initServerGitRepo(t, src)

	req := httptest.NewRequest(http.MethodPost, "/api/git/commit", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	s.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp pushResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Message != "nothing to commit (working tree clean)" {
		t.Fatalf("message = %q, want nothing to commit", resp.Message)
	}
}

func initServerGitRepo(t *testing.T, dir string) {
	t.Helper()

	testutil.RunGit(t, dir, "init")
	testutil.ConfigureGitUser(t, dir)
	testutil.RunGit(t, dir, "commit", "--allow-empty", "-m", "initial")
}
