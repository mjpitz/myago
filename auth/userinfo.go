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

	claims []byte
}

// Claims provides a convenient way to read additional data from the request.
func (u *UserInfo) Claims(v interface{}) error {
	return json.Unmarshal(u.claims, v)
}

// WithExtra adds additional claims to the raw payload. This is a rather expensive operation and should really
// only need to be done once during authentication.
func (u *UserInfo) WithExtra(v interface{}) error {
	idx := make(map[string]interface{})

	providedData, providedErr := json.Marshal(v)
	if providedErr != nil {
		return providedErr
	}

	userData, userErr := json.Marshal(u)
	if userErr != nil {
		return userErr
	}

	providedErr = json.Unmarshal(providedData, &idx)
	if providedErr != nil {
		return providedErr
	}

	userErr = json.Unmarshal(userData, &idx)
	if userErr != nil {
		return userErr
	}

	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	u.claims = data
	return nil
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
	u.claims = data[:]
	return nil
}

// userInfoWire is an intermediary data structure for unmarshalling the UserInfo structure.
type userInfoWire struct {
	Subject       string       `json:"sub"`
	Profile       string       `json:"profile"`
	Email         string       `json:"email"`
	EmailVerified boolOrString `json:"email_verified"`
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
