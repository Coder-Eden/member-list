import json
import requests

def get_github_user_info(username):
    response = requests.get(f"https://api.github.com/users/{username}")
    if response.status_code == 200:
        user_data = response.json()
        return user_data["name"], user_data["avatar_url"]  # 実名とアバターURLを返す
    else:
        print(f"ユーザー情報の取得に失敗しました。ステータスコード: {response.status_code}")
        return None, None

# JSONファイルを読み込む
with open('members.json', 'r', encoding='utf-8') as file:
    members = json.load(file)

markdown_table = "## 概要\n\n"
markdown_table += "![image](https://github.com/Coder-Eden/.github-private/assets/83957178/50505e63-2fba-4733-b825-b9b7e3615ad0)\n\n"
markdown_table += "#### CODE EDENは25卒限定の「完全招待制」のオンラインコミュニティです。\n\n"
markdown_table += "### 参加メンバー\n\n"
markdown_table += "| Number | ユーザー（GitHub） | アバター | 招待された人 | 招待した人（1）| 招待した人（2）| 招待した人（3）| 招待数 |\n"
markdown_table += "|-------|------------|------|------|----------------|----------------|----------------|------|\n"

for index, member in enumerate(members, start=1):
    github_username = member['github']
    real_name, avatar_url = get_github_user_info(github_username)  # 実名とアバターURLを取得
    github_link = f"https://github.com/{github_username}"
    
    if not real_name:
        real_name = github_username

    markdown_link = f"[{real_name}]({github_link})"
    avatar_link = f"![image]({avatar_url})"

    # 招待された人のGitHubリンク
    invited_person_username = member['invited_person']
    invited_person_link = f"[{invited_person_username}](https://github.com/{invited_person_username})"

    invitation_count = len([inv for inv in [member['invited_by_1'], member['invited_by_2'], member['invited_by_3']] if inv != "null"])

    markdown_table += f"| {index} | {markdown_link} | {avatar_link} | {invited_person_link} | {member['invited_by_1']} | {member['invited_by_2']} | {member['invited_by_3']} | {invitation_count} |\n"

# README.mdファイルにMarkdown表を書き込む
with open('README.md', 'w', encoding='utf-8') as file:
    file.write(markdown_table)
