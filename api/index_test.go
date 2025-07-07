package handler

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	commit, err := getGithubCommit("kpretty")
	fmt.Println(commit, err)
}
