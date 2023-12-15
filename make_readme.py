import json
import requests

def get_github_user_info(username):
    response = requests.get(f"https://api.github.com/users/{username}")
    if response.status_code == 200:
        return response.json()
    else:
        print(f"ユーザー情報の取得に失敗しました。ステータスコード: {response.status_code}")
        return None

def create_markdown_link(username, user_data):
    real_name = user_data["name"] if user_data["name"] else username
    github_link = f"https://github.com/{username}"
    return f"[{real_name}]({github_link})"

def create_avatar_link(user_data):
    return f"![image]({user_data['avatar_url']})" if user_data else "![image](default_avatar_url)"

def generate_markdown_table(members):
    markdown_table = "| Number | ユーザー（GitHub） | アバター | 招待された人 | 招待した人（1）| 招待した人（2）| 招待した人（3）| 招待数 |\n"
    markdown_table += "|-------|------------------|--------|------------|----------------|----------------|----------------|------|\n"

    for index, member in enumerate(members, start=1):
        user_data = get_github_user_info(member['github'])
        markdown_link = create_markdown_link(member['github'], user_data)
        avatar_link = create_avatar_link(user_data)
        invited_person_link = f"[{member['invited_person']}](https://github.com/{member['invited_person']})"
        invitation_count = sum(inv != None for inv in [member['invited_by_1'], member['invited_by_2'], member['invited_by_3']])

        markdown_table += f"| {index} | {markdown_link} | {avatar_link} | {invited_person_link} | {member['invited_by_1']} | {member['invited_by_2']} | {member['invited_by_3']} | {invitation_count} |\n"

    return markdown_table

def main():
    with open('members.json', 'r', encoding='utf-8') as file:
        members = json.load(file)

    markdown = "## 概要\n\n![image](https://github.com/Coder-Eden/.github-private/assets/83957178/50505e63-2fba-4733-b825-b9b7e3615ad0)\n\n"
    markdown += "#### CODE EDENは25卒限定の「完全招待制」のオンラインコミュニティです。\n\n"
    markdown += "### 参加メンバー\n\n" 
    markdown += "以下の表は、招待された人、ユーザーのGitHubプロフィール、何人目に招待されたか、および各ユーザーによって招待された人のリストを示しています。\n\n"
    markdown += generate_markdown_table(members)

    with open('README.md', 'w', encoding='utf-8') as file:
        file.write(markdown)

if __name__ == "__main__":
    main()
