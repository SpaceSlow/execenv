package metrics

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricWorkers_Err(t *testing.T) {
	mw, err := NewMetricWorkers(1, "", "", "", "", []time.Duration{})
	require.NoError(t, err)
	mw.errorsCh <- errors.New("some error")
	mw.Close()
	assert.Error(t, <-mw.Err())
	assert.Nil(t, <-mw.Err()) // checking for the absence of a deadlock
}

func TestMetricWorkers_getGopsutilMetrics(t *testing.T) {
	mw, err := NewMetricWorkers(1, "", "", "", "", []time.Duration{})
	require.NoError(t, err)
	metrics := <-mw.getGopsutilMetrics()
	assert.Greater(t, len(metrics), 0)

	for _, m := range metrics {
		assert.NotNil(t, m)
	}
}

func TestMetricWorkers_getRuntimeMetrics(t *testing.T) {
	mw, err := NewMetricWorkers(1, "", "", "", "", []time.Duration{})
	require.NoError(t, err)
	metrics := <-mw.getRuntimeMetrics()
	assert.Greater(t, len(metrics), 0)

	for _, m := range metrics {
		assert.NotNil(t, m)
	}
}

func writeRandomCertificateToFile(filename string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: publicKeyBytes,
	}
	pemPublicKey := pem.EncodeToMemory(block)

	return os.WriteFile(filename, pemPublicKey, 0600)
}

func Test_getPublicKey(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "certificate.*.pem")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	certificate, err := getPublicKey(f.Name()) // getting empty certificate
	require.ErrorIs(t, err, ErrDecodePEMBlock)
	assert.Nil(t, certificate)

	err = writeRandomCertificateToFile(f.Name())
	require.NoError(t, err)

	certificate, err = getPublicKey(f.Name()) // get valid public key
	require.NoError(t, err)

	assert.NotNil(t, certificate)
	switch interface{}(certificate).(type) {
	case *rsa.PublicKey:
	default:
		t.Errorf("incorrect type, expected *rsa.PublicKey, got: %T", certificate)
	}
	err = os.Remove(f.Name())
	require.NoError(t, err)

	certificate, err = getPublicKey(f.Name()) // get public key from missing file
	assert.Error(t, err)
	assert.Nil(t, certificate)

}
