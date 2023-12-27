package memberloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"member-list/internal/api"
)

type Member struct {
	GitHub string   `json:"github"`
	Host   string   `json:"host"`
	Guests []string `json:"guests"`
}

func ReadMembersFromFile(filePath string) ([]Member, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var members []Member
	if err := json.Unmarshal(data, &members); err != nil {
		return nil, fmt.Errorf("failed to unmarshal members: %w", err)
	}
	return members, nil
}

func WriteMembersToFile(members []Member, filePath string) error {
	jsonMembers, err := json.Marshal(members)
	if err != nil {
		return fmt.Errorf("failed to marshal members: %w", err)
	}

	err = ioutil.WriteFile(filePath, jsonMembers, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func UpdateMemberList() ([]Member, error) {
	githubMembers, err := api.GetGitHubMember()
	if err != nil {
		return nil, err
	}

	members, err := ReadMembersFromFile("./members.json")
	if err != nil {
		return nil, err
	}

	membersToUpdate := make(map[string]bool)
	for _, user := range githubMembers {
		membersToUpdate[user.Login] = true
	}

	for _, member := range members {
		if _, exists := membersToUpdate[member.GitHub]; exists {
			delete(membersToUpdate, member.GitHub)
		}
	}

	for githubMember := range membersToUpdate {
		fmt.Printf("Adding new member from GitHub: %s\n", githubMember)
		members = append(members, Member{GitHub: githubMember})
	}

	err = WriteMembersToFile(members, "./members.json")
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		fmt.Println("Member:", member.GitHub)
	}

	return members, nil
}
