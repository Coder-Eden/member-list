package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type GitHubUser struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

const Organization = "Coder-Eden"

func GetGitHubMember() ([]*GitHubUser, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GitHub token is not set")
	}

	url := fmt.Sprintf("https://api.github.com/orgs/%s/members", Organization)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API status: %s, body: %s", resp.Status, string(body))
	}

	var users []*GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
