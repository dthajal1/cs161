# Tests

## 3.1.2.c Empty Pasword

  * Create one user w/ password as '' and username as 'Bob'
  * Expect a user pointer to be returned, with no error
 
## 3.2.2 Multiple sessions, Same user

  * Bob calls `InitUser()` to create session bobLaptop
  * Bob calls `GetUser()` to create sesssion bobDesktop
  * Bob calls `StorFile('proj2.txt', content)` from bobLaptop
  * Bob calls `LoadFile('proj2.txt')` from BobDesktop which should equal `content`

## 3.6.2 Expect Error when Invite has not been accpeted

  * Bob calls `InitUser()` to create session bob
  * Alice calls `InitUser()` to create session alice
  * Bob calls `StorFile('proj2.txt', content)` from bob
  * Bob calls `CreateInvitation('proj2.txt', 'alice')`
  * Alice calls `LoadFile('proj2.txt')`, which should error

## 3.6.2 Expect Error after Revocation

  * Bob calls `InitUser()` to create session bob
  * Alice calls `InitUser()` to create session alice
  * Bob calls `StorFile('proj2.txt', content)` from bob
  * Bob calls `CreateInvitation('proj2.txt', 'alice')`
  * Alice accepts invitation
  * Alice calls `LoadFile('proj2.txt')`, which should return `content`
  * Bob revoke's Alice's access to `proj2.txt`
  * Alice calls `LoadFile('proj2.txt')`, which should error

## 3.6.2 Unauthorized access

  * Bob calls `InitUser()` to create session bob
  * Alice calls `InitUser()` to create session alice
  * Bob calls `StorFile('proj2.txt', content)` from bob
  * Alice calls `LoadFile('proj2.txt')`, which should error

## 3.6.9 Expect Error in Chain Revocation

  * Bob calls `InitUser()` to create session bob
  * Alice calls `InitUser()` to create session alice
  * Evantbot calls `InitUser()` to create session evanbot
  * Bob calls `StorFile('proj2.txt', content)` from bob
  * Bob calls `CreateInvitation('proj2.txt', 'alice')`
  * Alice accepts invitation
  * Alice calls `CreateInvitation('proj2.txt', 'evanbot')`
  * Evanbot accepts invitation
  * Bob revoke's Alice's access to `proj2.txt`
  * Evanbot calls `LoadFile('proj2.txt')`, which should error
