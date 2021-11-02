/*
Package ulid provides code for generating variable length unique, lexigraphic identifiers (ULID) with programmable
fills. Currently, there is a RandomGenerator that can be used to generate ULIDs with a randomized payload. To provide
a custom payload, simply extend the BaseGenerator, and override the Generate method. It's important to call the
BaseGenerator's Generate method, otherwise the skew and timestamp bits won't be set properly.
*/
package ulid
