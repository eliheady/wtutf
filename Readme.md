## What the UTF

A simple utility to help me out of my ASCII-centric shell

This program just prints out the code points of the string you feed into it. It can also show you the punycode conversion of your string and failure reasons if conversion isn't possible.

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
$ wtutf -ts $PINATA1
could not punycode-convert input
total bytes:	8
 characters:	7
----------------------------------
       code point |  bytes (len) | conversion rules violated
  p:         0x70 |       70 (1) | 
  i:         0x69 |       69 (1) | 
  n:         0x6e |       6e (1) | 
  ◌̃:       0x0303 |     cc83 (2) | CheckJoiners (RFC 5892), ValidateForRegistration (RFC 5891), ValidateLabels (RFC 5891), UseSTD3ASCIIRules (RFC 1034, 5891, UTS 46)
  a:         0x61 |       61 (1) | 
  t:         0x74 |       74 (1) | 
  a:         0x61 |       61 (1) | 

$ wtutf -ts $PINATA2
   punycode:	xn--piata-pta
total bytes:	7
 characters:	6
----------------------------------
       code point |  bytes (len)
  p:         0x70 |       70 (1) | 
  i:         0x69 |       69 (1) | 
  ñ:       0x00f1 |     c3b1 (2) | 
  a:         0x61 |       61 (1) | 
  t:         0x74 |       74 (1) | 
  a:         0x61 |       61 (1) |
```

Decode punycode strings

```shell
$ wtutf -p xn--piata-pta
   punycode:	xn--piata-pta
      utf-8:	piñata
total bytes:	7
 characters:	6
 ```

Care is taken to avoid echoing control characters in the output

```shell
$ wtutf -trs "$(printf '🔔bell\u07')"  
could not punycode-convert input
   total bytes:	9
    characters:	6
unicode ranges:
    Common: 2
    Latin: 4
----------------------------------
       code point |  bytes (len) | conversion rules violated
 🔔:   0x0001f514 | f09f9494 (4) | 
  b:         0x62 |       62 (1) | 
  e:         0x65 |       65 (1) | 
  l:         0x6c |       6c (1) | 
  l:         0x6c |       6c (1) | 
 ^G:         0x07 |       07 (1) | ValidateForRegistration (RFC 5891)
```

### Why make this?

I was interested in punycode and IDNA standards and wanted to make a simple utility to run locally to test conversion of various Unicode characters.

The Unicode Transformation Format – 8-bit (UTF-8) encoding allows for some difficult to interpret strings even when your rendering environment doesn't garble the characters with question marks or boxes.

The combining characters are a good example. The usage example above shows two strings that look identical on my system: "piñata" and "piñata". Only if I examine the bytes of those strings can I see that the second one uses 0x6e (n) and the UTF-8 "combining tilde" character 0x0303 ( ̃) to create the Spanish eñe. The first uses the single 0xf1 (ñ) "precomposed character".

Many combining characters aren't allowed in IDN domain registrations because they would provide a way to register names that are visually indistinguishable but comprised of different bytes, making things confusing for online piñata shopping. This is an example of general problem of homoglyphs in the DNS. Though combining characters are disallowed in the INDA specs, homoglyphs can still be used in various underhanded ways in and this tool could be useful to examine suspect strings.

Use the `--check`,`-c` flag to get a simple ok/caution validation of a string:

```shell
$ wtutf -c www.ցooցlе.com || echo 'WARNING'
WARNING
```

And to see a summary of the Unicode script families found in the input, use `--show-ranges`,`-r`:
```shell
$ wtutf --check --show-ranges www.ցooցlе.com
Cyrillic: 1
Latin: 9
Armenian: 2
```

```shell
$ wtutf -tr www.ցooցlе.com             
      punycode:	www.xn--ool-tdd07nca.com
   total bytes:	17
    characters:	14
unicode ranges:
    Latin: 9
    Common: 2
    Armenian: 2
    Cyrillic: 1
----------------------------------
       code point |  bytes (len)
  w:         0x77 |       77 (1) | 
  w:         0x77 |       77 (1) | 
  w:         0x77 |       77 (1) | 
  .:         0x2e |       2e (1) | 
  ց:       0x0581 |     d681 (2) | 
  o:         0x6f |       6f (1) | 
  o:         0x6f |       6f (1) | 
  ց:       0x0581 |     d681 (2) | 
  l:         0x6c |       6c (1) | 
  е:       0x0435 |     d0b5 (2) | 
  .:         0x2e |       2e (1) | 
  c:         0x63 |       63 (1) | 
  o:         0x6f |       6f (1) | 
  m:         0x6d |       6d (1) | 
  ```

The 0x0581 (ց) and 0x0435 (е) look slightly different from 'g' and 'e' on my system, but they could easily go unnoticed in many contexts.

This program shows what went into strings that look similar but aren't identical. It is also useful if you need to troubleshoot punycode conversion.

### Useful documents

* https://www.unicode.org/reports/tr46/#Validity_Criteria
* https://datatracker.ietf.org/doc/html/rfc5892
* https://datatracker.ietf.org/doc/html/rfc8753


### Installing

Download a build from the Releases section at right or [here](https://github.com/eliheady/wtutf/releases).

To verify provenance of a release, use the [slsa-verifier utility](https://github.com/slsa-framework/slsa-verifier) provided by the SLSA Framework project.

Example of verifying the [v0.0.1 release](https://github.com/eliheady/wtutf/releases/tag/v0.0.1):
```shell
$ curl -sLo wtutf-darwin-arm64 https://github.com/eliheady/wtutf/releases/download/v0.0.1/wtutf-darwin-arm64
$ curl -sLo wtutf-darwin-arm64.intoto.jsonl https://github.com/eliheady/wtutf/releases/download/v0.0.1/wtutf-darwin-arm64.intoto.jsonl
$ slsa-verifier verify-artifact wtutf-darwin-arm64 --provenance-path wtutf-darwin-arm64.intoto.jsonl --source-uri github.com/eliheady/wtutf --source-tag v0.0.1

Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@refs/tags/v2.1.0" at commit cca466a774d7475d4c4c404d9374f95d09afc6ee
Verifying artifact wtutf-darwin-arm64: PASSED

PASSED: SLSA verification passed
```

### Build from source

```shell
git clone https://github.com/eliheady/wtutf
cd wtutf
go build .
```