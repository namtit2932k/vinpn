package brand

import "testing"

func TestAbout(t *testing.T) {
	if Author != "sickyturtlez" || RepoURL != "https://github.com/sickyturtlez/vinpn" {
		t.Fatalf("Author=%q RepoURL=%q", Author, RepoURL)
	}
}
