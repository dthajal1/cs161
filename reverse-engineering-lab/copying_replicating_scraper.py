import requests

# fixed across requests
cookies = {
    'jwt': 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6IjMwMzQ4NjY3ODQiLCJhY2Nlc3NfdG9rZW4iOiJmZDFlYTI1OTg5OTA1ZWNlYjM2ZmZjNWMyYTRjZTRmNWM1N2MyZGJjOWQyMjg2NmZjOTgxMWE4YTM2ZDZjMDQ1In0.IM4e_agVEEjPzNGhjow9BFOOUNNCp7OJvjTdrOgcyVo',
}
# fixed across requests
headers = {
    'Accept': '*/*',
    'Accept-Language': 'en-US,en;q=0.9',
    'Connection': 'keep-alive',
    # 'Cookie': 'jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6IjMwMzQ4NjY3ODQiLCJhY2Nlc3NfdG9rZW4iOiJmZDFlYTI1OTg5OTA1ZWNlYjM2ZmZjNWMyYTRjZTRmNWM1N2MyZGJjOWQyMjg2NmZjOTgxMWE4YTM2ZDZjMDQ1In0.IM4e_agVEEjPzNGhjow9BFOOUNNCp7OJvjTdrOgcyVo',
    'Referer': 'https://generic.cs161.org/dashboard',
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'same-origin',
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36',
    'sec-ch-ua': '" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'x-access-token': 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6IjMwMzQ4NjY3ODQiLCJhY2Nlc3NfdG9rZW4iOiJmZDFlYTI1OTg5OTA1ZWNlYjM2ZmZjNWMyYTRjZTRmNWM1N2MyZGJjOWQyMjg2NmZjOTgxMWE4YTM2ZDZjMDQ1In0.IM4e_agVEEjPzNGhjow9BFOOUNNCp7OJvjTdrOgcyVo',
}

# varies across requests
# response = requests.get('https://generic.cs161.org/api/data/from/1/to/10', cookies=cookies, headers=headers)
# print(response.json())

def get_data(fromIdx: str, toIdx: str):
    dynamic_url = 'https://generic.cs161.org/api/data/from/' + fromIdx + '/to/' + toIdx
    response = requests.get(dynamic_url, cookies=cookies, headers=headers)
    return response.json()

def get_100th_value():
    result = get_data('1', '100')
    return result['Data'][99]

print(get_100th_value())