// Copyright (C) The AetherFS Authors - All Rights Reserved
// See LICENSE for more information.

package oidcauth

import (
	"errors"
)

var errUnauthorized = errors.New("unauthorized")

var errInternal = errors.New("internal server error")
