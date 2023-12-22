package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GitHubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type Member struct {
	Host   string   `json:"host"`
	GitHub string   `json:"github"`
	Guests []string `json:"guests"`
}

func getGitHubUserInfo(username string) (*GitHubUser, error) {
	user := GitHubUser{Login: username}
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GitHub token is not set")
	}

	// http.NewRequestを使ってリクエストを作成
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
		return &user, nil
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func createMarkdownLink(user *GitHubUser) string {
	if user == nil {
		return ""
	}
	name := user.Name
	if name == "" {
		name = user.Login
	}
	return fmt.Sprintf("[%s](https://github.com/%s)", name, user.Login)
}

func createAvatarLink(user *GitHubUser) string {
	return fmt.Sprintf("![image](%s)", user.AvatarURL)
}

func createMarkdownTable(members []Member) (string, error) {
	var builder strings.Builder
	builder.WriteString("| Number | ユーザー（GitHub） | アバター | 招待された人 | 招待数 | nil  | nil | nil |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")

	for i, member := range members {
		user, err := getGitHubUserInfo(member.GitHub)
		Host, err := getGitHubUserInfo(member.Host)
		if err != nil {
			return "", err
		}
		userGitHubLink := createMarkdownLink(user)
		HostGitHubLink := createMarkdownLink(Host)
		avatarLink := createAvatarLink(user)
		invitationCount := len(member.Guests)

		invitedGuests := make([]string, 0, invitationCount)
		for _, guestUsername := range member.Guests {
			guest, err := getGitHubUserInfo(guestUsername)
			if err != nil {
				return "", err
			}
			invitedGuests = append(invitedGuests, createMarkdownLink(guest))
		}

		builder.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %d | %s | %s | %s |\n",
			i+1, userGitHubLink, avatarLink, HostGitHubLink, invitationCount, "", "", ""))
	}

	return builder.String(), nil
}

func main() {
	data, err := ioutil.ReadFile("members.json")
	if err != nil {
		panic(err)
	}

	var members []Member
	err = json.Unmarshal(data, &members)
	if err != nil {
		panic(err)
	}

	markdownTable, err := createMarkdownTable(members)
	if err != nil {
		panic(err)
	}

	markdown := "## 概要\n\n![image](https://github.com/Coder-Eden/.github-private/assets/83957178/50505e63-2fba-4733-b825-b9b7e3615ad0)\n\n" +
		"#### CODE EDENは25卒限定の「完全招待制」のオンラインコミュニティです。\n\n" +
		"### 参加メンバー\n\n" +
		"以下の表は、招待された人、ユーザーのGitHubプロフィール、何人目に招待されたか、および各ユーザーによって招待された人のリストを示しています。\n\n" +
		markdownTable

	err = ioutil.WriteFile("README.md", []byte(markdown), 0644)
	if err != nil {
		panic(err)
	}
}
