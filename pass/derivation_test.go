package pass_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/pass"
)

func TestDerivation(t *testing.T) {
	root, err := pass.Identity(pass.Authentication, []byte("badadmin"), "admin")
	require.NoError(t, err)

	keyA0 := pass.SiteKey(pass.Authentication, root, "a.com", 0)
	keyB0 := pass.SiteKey(pass.Authentication, root, "b.com", 0)
	require.NotEqual(t, keyA0, keyB0)

	keyA1 := pass.SiteKey(pass.Authentication, root, "a.com", 1)
	require.NotEqual(t, keyA0, keyA1)
	require.NotEqual(t, keyB0, keyA1)

	passA0 := pass.SitePassword(keyA0, pass.MaximumSecurity)
	require.Equal(t, "V8-^r&YEi2kq4w6nfDa7", string(passA0))

	passA1 := pass.SitePassword(keyA1, pass.MaximumSecurity)
	require.Equal(t, "aG84VX0l6%UqL5pD^o7+", string(passA1))
}
