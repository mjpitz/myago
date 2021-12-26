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

package auth

import (
	"encoding/json"
	"strconv"
)

// UserInfo represents a minimum set of user information.
type UserInfo struct {
	// Subject is the users ID. CACF1875-7B44-4B77-BF52-51A06E52FFDF
	Subject string `json:"sub"`
	// Profile is the users name. "Jane Doe"
	Profile string `json:"profile"`
	// Email is the users' email address. jane@example.com
	Email string `json:"email"`
	// EmailVerified indicates if the user has verified their email address.
	EmailVerified bool `json:"email_verified"`
	// Groups contains a list of groups that the user belongs to.
	Groups []string `json:"groups"`

	claims []byte
}

// Claims provides a convenient way to read additional data from the request.
func (u UserInfo) Claims(v interface{}) error {
	return json.Unmarshal(u.claims, v)
}

// UnmarshalJSON transparently unmarshals the user information structure.
func (u *UserInfo) UnmarshalJSON(data []byte) error {
	raw := &userInfoWire{}
	err := json.Unmarshal(data, raw)
	if err != nil {
		return err
	}

	u.Subject = raw.Subject
	u.Profile = raw.Profile
	u.Email = raw.Email
	u.EmailVerified = bool(raw.EmailVerified)
	u.Groups = append([]string{}, raw.Groups...)
	u.claims = append([]byte{}, data...)
	return nil
}

// userInfoWire is an intermediary data structure for unmarshalling the UserInfo structure.
type userInfoWire struct {
	Subject       string       `json:"sub"`
	Profile       string       `json:"profile"`
	Email         string       `json:"email"`
	EmailVerified boolOrString `json:"email_verified"`
	Groups        []string     `json:"groups"`
}

// boolOrString is used to deserialize a boolean value that may also be represented as a string. Some identity providers
// use this representation rather than a boolean in the JSON response.
type boolOrString bool

func (b *boolOrString) UnmarshalJSON(data []byte) error {
	var stringValue string
	var boolValue bool

	// try deserializing a bool first
	boolErr := json.Unmarshal(data, &boolValue)
	if boolErr == nil {
		*b = boolOrString(boolValue)
		return nil
	}

	// if that fails, decode a string
	stringErr := json.Unmarshal(data, &stringValue)
	if stringErr != nil {
		return stringErr
	}

	// and parse
	boolValue, boolErr = strconv.ParseBool(stringValue)
	if boolErr != nil {
		return boolErr
	}

	*b = boolOrString(boolValue)
	return nil
}
