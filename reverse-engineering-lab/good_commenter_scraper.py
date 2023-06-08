# fetch content
import requests
import json

# Fixed across requests.
cookies = {
    'session_id': '4NiMSKbxN9nymooG6Kycqb6c',
    '__gads': 'ID=41d9fdb87b15c2ef-22a1ae8affc7007f:T=1623198100:RT=1623198100:S=ALNI_MbA9HHn3tezyaeqysXkuXgt8JW2Cw',
    '_ga': 'GA1.2.869757251.1623864961',
    'last_piaz_user': 'jx6qr6jbq1ou6',
    'AWSELB': '732545B312943CE3A5A159C88B79D008FE4AE6530DEA647084583F9DF0336A93698216B8DDBDE1783DE61A9BBCED6DF4B66026BB7C5D23AC87D809FE3C293EA15EEB6C40E7',
    'AWSELBCORS': '732545B312943CE3A5A159C88B79D008FE4AE6530DEA647084583F9DF0336A93698216B8DDBDE1783DE61A9BBCED6DF4B66026BB7C5D23AC87D809FE3C293EA15EEB6C40E7',
    'piazza_session': '1.eyJleHBpcmVzIjoxNjUyNzU4MTU0LCJwZXJtIjoiMCIsImxvZ2dpbmdfaW4iOnRydWUsImNyZWF0ZWRfYXQiOjE2NTE1NDg1NTQsIndoZW4iOjE2NTE1NDg1NTEsImhvbWUiOiIvIiwicmVtZW1iZXIiOiJvbiIsInNlc3Npb25fdG9rZW4iOiJzdF80bXk4S3FDdWVpZEVtVlVJWVBJSyIsIm5pZHMiOiJreTllOGNxODY4NzJ1OjA7aHlxMGJyMXUza3g3ZGc6MDtqeml5a3U1Z29teTdhcTowO2swZzlqM3IwNGVkY246MDtrNWVldnhlYnpwajI1YjowO2s1b2xpenRtenF5NXp2OjA7azVxNDd5cGQ4N2EyaXI6MDtrM3FicHMxaHFkbjdhdjowO2s2azlicm82NG5mb3o6MDtrZDBpZmFjdW1tZDU2YTowO2tlM2FqOHo3NWFwYjA6MDtrZjAwMXViMzgyZ2tvOjA7a2VzMzdpYXYxMWkzcXE6MDtrajlkdW50bWM0dzdpYjowO2tqemM1bXBseDdrMjZtOjA7a2sxcmk0MzBpenM3N2s6MDtrcDNwbTl1MGtmODVuMjowO2twZnhpNHBhMzlqNDZrOjA7a3NjbjRieG41aW4yc246MDtrc3Frb3YxNmo4ajZiNDowO2tzcGgzYmMxZDBoNXk5OjA7a3N3OWtwcm00MDg0enI6MDtrdGdyaDdoeXQwbjRwazowO2t4ejN3aGFrNG1hM3J0OjA7IiwidGFnIjoiIiwidXNlciI6Imp4NnFyNmpicTFvdTYiLCJlbWFpbCI6ImRpcmFqLnRoYWphbGlAYmVya2VsZXkuZWR1In0%3D.hjG7yK93TzjdlKLN%2F49InjF1HaLSoyt4WW4K7WSOkHA%3D',
}
# Fixed across requests.
headers = {
    'Accept': 'application/json, text/javascript, */*; q=0.01',
    'Accept-Language': 'en-US,en;q=0.9',
    'CSRF-Token': '4NiMSKbxN9nymooG6Kycqb6c',
    'Connection': 'keep-alive',
    'Content-Type': 'application/json; charset=UTF-8',
    # Requests sorts cookies= alphabetically
    # 'Cookie': 'session_id=4NiMSKbxN9nymooG6Kycqb6c; __gads=ID=41d9fdb87b15c2ef-22a1ae8affc7007f:T=1623198100:RT=1623198100:S=ALNI_MbA9HHn3tezyaeqysXkuXgt8JW2Cw; _ga=GA1.2.869757251.1623864961; last_piaz_user=jx6qr6jbq1ou6; AWSELB=732545B312943CE3A5A159C88B79D008FE4AE6530DEA647084583F9DF0336A93698216B8DDBDE1783DE61A9BBCED6DF4B66026BB7C5D23AC87D809FE3C293EA15EEB6C40E7; AWSELBCORS=732545B312943CE3A5A159C88B79D008FE4AE6530DEA647084583F9DF0336A93698216B8DDBDE1783DE61A9BBCED6DF4B66026BB7C5D23AC87D809FE3C293EA15EEB6C40E7; piazza_session=1.eyJleHBpcmVzIjoxNjUyNzU4MTU0LCJwZXJtIjoiMCIsImxvZ2dpbmdfaW4iOnRydWUsImNyZWF0ZWRfYXQiOjE2NTE1NDg1NTQsIndoZW4iOjE2NTE1NDg1NTEsImhvbWUiOiIvIiwicmVtZW1iZXIiOiJvbiIsInNlc3Npb25fdG9rZW4iOiJzdF80bXk4S3FDdWVpZEVtVlVJWVBJSyIsIm5pZHMiOiJreTllOGNxODY4NzJ1OjA7aHlxMGJyMXUza3g3ZGc6MDtqeml5a3U1Z29teTdhcTowO2swZzlqM3IwNGVkY246MDtrNWVldnhlYnpwajI1YjowO2s1b2xpenRtenF5NXp2OjA7azVxNDd5cGQ4N2EyaXI6MDtrM3FicHMxaHFkbjdhdjowO2s2azlicm82NG5mb3o6MDtrZDBpZmFjdW1tZDU2YTowO2tlM2FqOHo3NWFwYjA6MDtrZjAwMXViMzgyZ2tvOjA7a2VzMzdpYXYxMWkzcXE6MDtrajlkdW50bWM0dzdpYjowO2tqemM1bXBseDdrMjZtOjA7a2sxcmk0MzBpenM3N2s6MDtrcDNwbTl1MGtmODVuMjowO2twZnhpNHBhMzlqNDZrOjA7a3NjbjRieG41aW4yc246MDtrc3Frb3YxNmo4ajZiNDowO2tzcGgzYmMxZDBoNXk5OjA7a3N3OWtwcm00MDg0enI6MDtrdGdyaDdoeXQwbjRwazowO2t4ejN3aGFrNG1hM3J0OjA7IiwidGFnIjoiIiwidXNlciI6Imp4NnFyNmpicTFvdTYiLCJlbWFpbCI6ImRpcmFqLnRoYWphbGlAYmVya2VsZXkuZWR1In0%3D.hjG7yK93TzjdlKLN%2F49InjF1HaLSoyt4WW4K7WSOkHA%3D',
    'Origin': 'https://piazza.com',
    'Referer': 'https://piazza.com/class/ky9e8cq86872u?cid=8_f11',
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'same-origin',
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36',
    'X-Requested-With': 'XMLHttpRequest',
    'sec-ch-ua': '" Not A;Brand";v="99", "Chromium";v="100", "Google Chrome";v="100"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
}
# Varies across requests
params = {
    'method': 'content.get',
    'aid': 'l2pld84puu6w',
    
}
# Varies across requests
json_data = {
    'method': 'content.get',
    'params': {
        'cid': '8',
        'nid': 'ky9e8cq86872u',
    },
}

# response = requests.post('https://piazza.com/logic/api', params=params, cookies=cookies, headers=headers, json=json_data)
# print(response.json())


# Fixed across `content.get` requests.
def get_content(cid: str, nid: str, aid: str):
    params = {
        'method': 'content.get',
        'aid': aid,
    }

    json_data = {
        'method': 'content.get',
        'params': {
            'cid': cid,
            'nid': nid,
        },
    }

    response = requests.post('https://piazza.com/logic/api', headers=headers, params=params, cookies=cookies, json=json_data)
    return response.json()


aid = 'l2wmev078xau'
cid = '8'
nid = 'ky9e8cq86872u'
# print(get_content(cid, nid, aid))
# with open("temp.json", "w") as file1:
#     # Writing data to a file
#     file1.write(json.dumps(get_content(cid, nid, aid)))

result = get_content(cid, nid, aid)
i_like_cheese = result['result']['children'][10]
with open("temp.json", "w") as file1:
    # Writing data to a file
    file1.write(json.dumps(i_like_cheese))


good_commenters = i_like_cheese['tag_good']
ta_or_instructor_good_commenters = []
for i in range(len(good_commenters)):
    # print(good_commenters[i]['role'])
    # print(good_commenters[i]['name'])
    curr_comment = good_commenters[i]
    if (curr_comment['role'] == 'ta' or curr_comment['role'] == 'instructor'):
        ta_or_instructor_good_commenters.append(curr_comment['name'])

print(ta_or_instructor_good_commenters)




