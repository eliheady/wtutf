## What the UTF

A simple utility to help me out of my ASCII-centric shell

This program just prints out the Unicode code points of the string you feed into it. It can also show you the punycode conversion of your string and failure reasons if conversion isn't possible.

Try this
```shell
export PINATA1="piñata" 
export PINATA2="piñata" 

echo "$PINATA1 has $(echo -n $PINATA1 | wc -c) bytes"
echo "$PINATA2 has $(echo -n $PINATA2 | wc -c) bytes"

if [[ "$PINATA1" != "$PINATA2" ]]
then
    echo 'they do not match!'
fi
```

```
piñata has        8 bytes
piñata has        7 bytes
they do not match!
```

Huh? But they look the same!?

```shell
$ wtutf -s $PINATA1
total bytes:    8
characters:     7
punycode:       could not punycode-convert input
    code point | bytes  | conversion errors
p:        0x70 | (1)    | 
i:        0x69 | (1)    | 
n:        0x6e | (1)    | 
̃:       0x0303 | (2)    | CheckJoiners (RFC 5892), ValidateForRegistration (RFC 5891), ValidateLabels (RFC 5891), UseSTD3ASCIIRules (RFC 1034, 5891, UTS 46)
a:        0x61 | (1)    | 
t:        0x74 | (1)    | 
a:        0x61 | (1)    |  

$ wtutf -s $PINATA2
total bytes:    7
characters:     6
punycode:       xn--piata-pta
    code point | bytes 
p:        0x70 | (1)
i:        0x69 | (1)
ñ:        0xf1 | (2)
a:        0x61 | (1)
t:        0x74 | (1)
a:        0x61 | (1)
```

Care is taken to avoid echoing control characters in the output

```shell
$ wtutf -s "$(echo -n 🔔bell. | tr '.' '\007')" 
total bytes:    9
characters:     6
punycode:       could not punycode-convert input
     code point | bytes | conversion errors
 🔔:    0x1f514 | (4)   | 
 b:        0x62 | (1)   | 
 e:        0x65 | (1)   | 
 l:        0x6c | (1)   | 
 l:        0x6c | (1)   | 
^G:        0x07 | (1)   | ValidateForRegistration (RFC 5891)
```

### Why make this?

I was interested in punycode and IDNA standards and wanted to make a simple utility to run locally to test coversion of various Unicode characters.

The Unicode Transformation Format – 8-bit (UTF-8) encoding allows for some difficult to interpret strings even when your rendering environment doesn't garble the characters with question marks or boxes.

The combining characters are a good example. The usage string above shows two strings that look identical on my system: "piñata" and "piñata". Only if I examine the bytes of those strings can I see that the second one uses 0x6e (n) and the UTF-8 "combining tilde" character 0x0303 ( ̃) to create the Spanish eñe. The first uses the single 0xf1 (ñ) "precomposed character".

The combining characters aren't allowed in IDN domain registrations because they would provide a way to register names that are visually indistinguishable but comprised of different bytes, making things confusing for online piñata shopping.

This program shows what went into strings that look identical but aren't. It is also useful if you need to troubleshoot punycode conversion.

### Useful documents

* https://www.unicode.org/reports/tr46/#Validity_Criteria
* https://datatracker.ietf.org/doc/html/rfc5892
* https://datatracker.ietf.org/doc/html/rfc8753