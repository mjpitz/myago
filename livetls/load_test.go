package livetls_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/livetls"
)

type TestCase struct {
	Name      string
	Config    *livetls.Config
	NilResult bool
	NilErr    bool
}

type Callback func(cfg *livetls.Config) (interface{}, error)

func run(t *testing.T, testCases []TestCase, cb Callback) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result, err := cb(testCase.Config)

			if testCase.NilErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			if testCase.NilResult {
				require.Nil(t, result)
			} else {
				require.NotNil(t, result)
			}
		})
	}
}

func TestLoadCertPool(t *testing.T) {
	ctx := context.Background()

	run(
		t, []TestCase{
			{
				Name: "missing",
				Config: &livetls.Config{
					CertPath: "sslconf",
				},
				NilResult: true,
				NilErr:    true,
			},
			{
				Name: "invalid",
				Config: &livetls.Config{
					CertPath: "sslconf",
					CAFile:   "invalid.pem",
				},
				NilResult: true,
				NilErr:    true,
			},
			{
				Name: "valid",
				Config: &livetls.Config{
					CertPath: "sslconf",
					CAFile:   "ca.pem",
				},
				NilErr: true,
			},
		},
		func(cfg *livetls.Config) (interface{}, error) {
			return livetls.LoadCertPool(ctx, cfg)
		},
	)
}

func TestLoadCertificate(t *testing.T) {
	ctx := context.Background()

	run(
		t, []TestCase{
			{
				Name: "missing cert",
				Config: &livetls.Config{
					CertPath: "sslconf",
					KeyFile:  "key.pem",
				},
				NilResult: true,
				NilErr: true,
			},
			{
				Name: "invalid cert",
				Config: &livetls.Config{
					CertPath: "sslconf",
					CertFile: "invalid.pem",
					KeyFile:  "key.pem",
				},
				NilErr: true,
				NilResult: true,
			},
			{
				Name: "missing key",
				Config: &livetls.Config{
					CertPath: "sslconf",
					CertFile: "cert.pem",
				},
				NilResult: true,
				NilErr: true,
			},
			{
				Name: "invalid key",
				Config: &livetls.Config{
					CertPath: "sslconf",
					CertFile: "cert.pem",
					KeyFile:  "invalid.pem",
				},
				NilErr: true,
				NilResult: true,
			},
			{
				Name: "valid",
				Config: &livetls.Config{
					CertPath: "sslconf",
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
				},
				NilErr: true,
			},
		},
		func(cfg *livetls.Config) (interface{}, error) {
			return livetls.LoadCertificate(ctx, cfg)
		},
	)
}
