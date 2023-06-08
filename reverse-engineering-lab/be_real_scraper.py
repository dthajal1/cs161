import requests

cookies = {
    'jwt': 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6IjMwMzQ4NjY3ODQiLCJhY2Nlc3NfdG9rZW4iOiJmZDFlYTI1OTg5OTA1ZWNlYjM2ZmZjNWMyYTRjZTRmNWM1N2MyZGJjOWQyMjg2NmZjOTgxMWE4YTM2ZDZjMDQ1In0.IM4e_agVEEjPzNGhjow9BFOOUNNCp7OJvjTdrOgcyVo',
    # 'jwt': '<redacted>'
}

headers = {
    'Accept': '*/*',
    'x-access-token': 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6IjMwMzQ4NjY3ODQiLCJhY2Nlc3NfdG9rZW4iOiJmZDFlYTI1OTg5OTA1ZWNlYjM2ZmZjNWMyYTRjZTRmNWM1N2MyZGJjOWQyMjg2NmZjOTgxMWE4YTM2ZDZjMDQ1In0.IM4e_agVEEjPzNGhjow9BFOOUNNCp7OJvjTdrOgcyVo',
    # 'x-access-token': '<redacted>',
}

json_data = {
    'toInclude': [
        'alice',
    ],
    'toExclude': [
        'mallory',
    ],
}

response = requests.post('https://generic.cs161.org/api/bereal/getFriendSuggestions', headers=headers, json=json_data, cookies=cookies, verify=False)
print(response.json())