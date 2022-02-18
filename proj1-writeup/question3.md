# Question 3: Polaris Write-up
## Main Idea
#### (A description of the vulnerability)
For this problem, stack canaries were enabled which prevented us from a simple buffer overflow attack. However, the first while loop inside the ```dehexify``` function allowed us to leak the stack canary and after which we were able to attack this program using the same technique we used for a simple buffer overflow attack except by also adding the stack canary onto it to make it seem like the stack canary remained unchanged (which didn't exit the program).

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)
We first determined the address of buffer (0xffffd5ec) and the address of the rip of ```dehexify``` function (0xffffd60c). This was done by invoking GDB and setting a breakpoint at line 17.
```
(gdb) x/16x c.buffer
0xffffd5ec:     0x41414141      0x41414141      0x41414141      0x0800785c
0xffffd5fc:     0xd1cd93b5      0x0804d020      0x00000000      0xffffd618
0xffffd60c:     0x08049341      0x00000000      0xffffd630      0xffffd6ac
0xffffd61c:     0x0804952a      0x00000001      0x08049329      0x0804cfe8

(gdb) i f
Stack level 0, frame at 0xffffd610:
 eip = 0x8049245 in dehexify (dehexify.c:22); saved eip = 0x8049341
 called by frame at 0xffffd630
 source language c.
 Arglist at 0xffffd608, args:
 Locals at 0xffffd608, Previous frame's sp is 0xffffd610
 Saved registers:
  ebp at 0xffffd608, eip at 0xffffd60c
```
By doing so we learned that 

## Exploit Structure
#### (A description of your exploit structure)
This exploit has five parts:
1. Write 32 garbage characters to overwrite the ```struct c``` which includes ```c.buffer``` and ```c.answer``` both 16 bytes each. 
2. Above the ```struct c``` lies the stack canary and we want it to remain the same. We successfully leaked the stack canary by exploiting a vulnerability inside the while loop of ```dehexify``` function and so we use that leaked canary here. 
3. Write 12 more garbage characters to overwrite compiler padding and the sfp of ```dehexify```.
4. Overwrite the rip of dehexify with the address of shellcode which lies directly above rip (rip + 4 = 0xffffd60c + 4 = 0xffffd610).
5. Finally, insert the shellcode directly above the rip.

## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)
When we ran GDB after inputting the malicious exploit string, we got the following output:
```
(gdb) x/16x c.buffer
0xffffd5ec:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5fc:     0x4bb779e1      0x41414141      0x41414141      0x41414141
0xffffd60c:     0xffffd610      0xdb31c031      0xd231c931      0xb05b32eb
0xffffd61c:     0xcdc93105      0xebc68980      0x3101b006      0x8980cddb
```
The address of rip (0xffffd60c) is successfully overwritten with 0xffffd610 which points to the shellcode right next to it. 

