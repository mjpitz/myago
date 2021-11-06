/*
Package ulid provides code for generating variable length unique, lexigraphic identifiers (ULID) with programmable
fills. Currently, there is a RandomGenerator that can be used to generate ULIDs with a randomized payload. To provide
a custom payload, simply extend the BaseGenerator, and override the Generate method. It's important to call the
BaseGenerator's Generate method, otherwise the skew and timestamp bits won't be set properly.

Unlike the canonical [ULID](https://github.com/ulid/spec), this version holds a placeholder byte for major clock skews
which can often occur in distributed systems. The wire format is as follows: `[ skew ][ sec ][ payload ]`

 - `skew` - 1 byte used to handle major clock skews (reserved, unused)
 - `sec` - 6 bytes of a unix timestamp (should give us until the year 10k or so)
 - `payload` - N bytes for the payload
*/
package ulid
