## What the UTF

A simple utility to help me out of my ASCII-centric shell

This program just prints out the Unicode code points of the string you feed into it. It can also show you the punycode conversion.

```shell
$ wtutf piñata
total bytes: 8
characters: 7
    code point  (bytes)
p:      0x70    (1)
i:      0x69    (1)
n:      0x6e    (1)
̃:      0x303   (2)
a:      0x61    (1)
t:      0x74    (1)
a:      0x61    (1)
punycode:
could not punycode-convert input

$ wtutf piñata  
total bytes: 7
characters: 6
    code point  (bytes)
p:      0x70    (1)
i:      0x69    (1)
ñ:      0xf1    (2)
a:      0x61    (1)
t:      0x74    (1)
a:      0x61    (1)
punycode:
xn--piata-pta
.

Usage:
  wtutf [flags]

Flags:
  -h, --help   help for wtutf
```