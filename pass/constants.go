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
	"strings"
)

// Scope defines an enumeration of possible scopes used in key derivation.
type Scope string

const (
	// Authentication is used when generating a key that is used for authenticating the user, such as a password.
	Authentication Scope = "com.lyndir.masterpassword"

	// Identification is used when generating a key that is intended for the purpose of identifying the user.
	// Identification keys are not necessarily private.
	Identification Scope = "com.lyndir.masterpassword.login"

	// Recovery is used for generating fallback identifiers for use in access recovery when the primary mechanism has
	// failed.
	Recovery Scope = "com.lyndir.masterpassword.answer"
)

// TemplateClass defines an enumeration of password templates to choose from.
type TemplateClass string

const (
	// MaximumSecurity defines a set of templates used to generate passwords with the strongest security.
	MaximumSecurity TemplateClass = "max"
	// Long defines a set of templates used to generate long passwords.
	Long TemplateClass = "long"
	// Medium defines a set of templates used to generate medium-length passwords.
	Medium TemplateClass = "medium"
	// Short defines a set of templates used to generate short-length passwords.
	Short TemplateClass = "short"
	// Basic defines a set of templates used to generate basic passwords.
	Basic TemplateClass = "basic"
	// PIN generates a pin.
	PIN TemplateClass = "pin"
	// VerificationCode provides a template for generating a 6-digit verification code.
	VerificationCode TemplateClass = "code"
)

// characterClasses defines a mapping of a character code to it's associated character set.
type characterClasses map[byte]string

// templateClasses defines a mapping of templates that yield different strengths.
type templateClasses map[TemplateClass][][]byte

var (
	vowels          = "aeiou"
	vowelsUpper     = strings.ToUpper(vowels)
	consonants      = "bcdfghjklmnpqrstvwxyz"
	consonantsUpper = strings.ToUpper(consonants)
	numeric         = "123456789"
	other           = "@&%?,=[]_:-+*$#!'^~;()/."
	alphabeticUpper = vowelsUpper + consonantsUpper
	alphabetic      = vowelsUpper + vowels + consonantsUpper + consonants
	union           = alphabetic + "0" + numeric + "!@#$%^&*()"

	// defaultCharacterClasses defines a default set of characters to pick from when a template class sees a specific
	// byte. This is used to ensure the password is a good mix of values.
	defaultCharacterClasses = characterClasses{
		'v': vowels,
		'V': vowelsUpper,
		'c': consonants,
		'C': consonantsUpper,
		'A': alphabeticUpper,
		'a': alphabetic,
		'n': numeric,
		'o': other,
		'x': union,
	}

	// defaultTemplateClasses defines a common set of templates for passwords to use. They're keyed by the class of
	// password formats.
	defaultTemplateClasses = templateClasses{
		// according to specification

		MaximumSecurity: {
			[]byte("anoxxxxxxxxxxxxxxxxx"), []byte("axxxxxxxxxxxxxxxxxno"),
		},
		Long: {
			[]byte("CvcvnoCvcvCvcv"), []byte("CvcvCvcvCvccno"),
			[]byte("CvcvCvcvnoCvcv"), []byte("CvccnoCvccCvcv"),
			[]byte("CvcvCvcvCvcvno"), []byte("CvccCvccnoCvcv"),
			[]byte("CvccnoCvcvCvcv"), []byte("CvccCvccCvcvno"),
			[]byte("CvccCvcvnoCvcv"), []byte("CvcvnoCvccCvcc"),
			[]byte("CvccCvcvCvcvno"), []byte("CvcvCvccnoCvcc"),
			[]byte("CvcvnoCvccCvcv"), []byte("CvcvCvccCvccno"),
			[]byte("CvcvCvccnoCvcv"), []byte("CvccnoCvcvCvcc"),
			[]byte("CvcvCvccCvcvno"), []byte("CvccCvcvnoCvcc"),
			[]byte("CvcvnoCvcvCvcc"), []byte("CvccCvcvCvccno"),
			[]byte("CvcvCvcvnoCvcc"),
		},
		Medium: {[]byte("CvcnoCvc"), []byte("CvcCvcno")},
		Short:  {[]byte("Cvcn")},
		Basic:  {[]byte("aaanaaan"), []byte("aannaaan"), []byte("aaannaaa")},
		PIN:    {[]byte("nnnn")},

		// custom formats

		VerificationCode: {[]byte("nnnnnn")},
	}
)
