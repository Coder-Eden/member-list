package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GitHubUser は GitHub ユーザー情報を格納する構造体です。
type GitHubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// Member はメンバー情報を格納する構造体です。
type Member struct {
	GitHub        string `json:"github"`
	InvitedPerson string `json:"invited_person"`
	InvitedBy1    string `json:"invited_by_1"`
	InvitedBy2    string `json:"invited_by_2"`
	InvitedBy3    string `json:"invited_by_3"`
}

// getGitHubUserInfo は GitHub API を使用してユーザー情報を取得します。
func getGitHubUserInfo(username string) (*GitHubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	var user GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// createMarkdownLink は Markdown 形式のリンクを生成します。
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

// createAvatarLink はアバターの画像リンクを生成します。
func createAvatarLink(user *GitHubUser) string {
	return fmt.Sprintf("![image](%s)", user.AvatarURL)
}

// createMarkdownTable はメンバーのリストから Markdown テーブルを生成します。
func createMarkdownTable(members []Member) (string, error) {
	var builder strings.Builder
	builder.WriteString("| Number | ユーザー（GitHub） | アバター | 招待した人 | 招待数 |\n")
	builder.WriteString("|-------|------------------|--------|------------|------|\n")

	for i, member := range members {
		user, err := getGitHubUserInfo(member.GitHub)
		invitedPerson, err := getGitHubUserInfo(member.InvitedPerson)
		if err != nil {
			return "", err
		}
		userGitHubLink := createMarkdownLink(user)
		invitedPersonGitHubLink := createMarkdownLink(invitedPerson)
		avatarLink := createAvatarLink(user)
		invitationCount := 0
		for _, invite := range []string{member.InvitedBy1, member.InvitedBy2, member.InvitedBy3} {
			if invite != "" {
				invitationCount++
			}
		}

		builder.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %d |\n",
			i+1, userGitHubLink, avatarLink, invitedPersonGitHubLink, invitationCount))
	}

	return builder.String(), nil
}

func main() {
	// ファイルからメンバー情報を読み込む
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

	// 結果を Markdown ファイルに書き出す
	err = ioutil.WriteFile("README.md", []byte(markdown), 0644)
	if err != nil {
		panic(err)
	}
}
