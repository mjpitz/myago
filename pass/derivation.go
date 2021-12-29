// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package pass

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"

	"golang.org/x/crypto/scrypt"
)

const (
	n      = 32768
	r      = 8
	p      = 2
	keyLen = 64
)

// Identity computes your identity which is defined by your root password. This key unlocks all doors. The result is a
// cryptographic key that is derived from the scope of the operation, your password (authentication), and your name
// (identification).
func Identity(scope Scope, password []byte, name string) ([]byte, error) {
	seed := bytes.NewBuffer(nil)

	seed.WriteString(string(scope))
	_ = binary.Write(seed, binary.BigEndian, uint32(len(name)))
	seed.WriteString(name)

	return scrypt.Key(password, seed.Bytes(), n, r, p, keyLen)
}

// SiteKey derives a site specific key from your identity key. Use of your identity key ensures only your identity has
// access to this key and your site name scopes the key to the site. The site counter ensures you can easily create new
// keys for the site should a key become compromised.
func SiteKey(scope Scope, identity []byte, site string, counter uint32) []byte {
	seed := bytes.NewBuffer(nil)

	seed.WriteString(string(scope))
	_ = binary.Write(seed, binary.BigEndian, uint32(len(site)))
	seed.WriteString(site)
	_ = binary.Write(seed, binary.BigEndian, counter)

	sig := hmac.New(sha256.New, identity)
	sig.Write(seed.Bytes())

	return sig.Sum(nil)
}

// SitePassword is an identifier derived from your site key in compliance with the site's password policy. This step
// renders the sites cryptographic key into a format that the site's password input will accept.
func SitePassword(siteKey []byte, class TemplateClass) []byte {
	templates := defaultTemplateClasses[class]

	template := templates[int(siteKey[0])%len(templates)]

	password := make([]byte, len(template))
	for i := 0; i < len(template); i++ {
		charset := defaultCharacterClasses[template[i]]
		password[i] = charset[int(siteKey[i+1])%len(charset)]
	}

	return password
}
