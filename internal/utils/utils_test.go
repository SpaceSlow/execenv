package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "medium text",
			data: []byte(`Dolores doloribus est ut qui cumque.
Ratione repudiandae placeat est voluptas sequi aliquam illum.
Voluptates a veniam quidem quidem explicabo doloribus. 
Explicabo perferendis voluptas harum eveniet et dolorem. 
Error tempore perspiciatis sit perferendis. 
At voluptatem soluta quod esse. Ipsa ducimus dolores et quam. 
Voluptates voluptas ea blanditiis placeat dolorem. 
Exercitationem animi deserunt repellat assumenda eveniet quam reprehenderit.
Consequatur fugiat et sequi.`),
		},
		{
			name: "large text",
			data: []byte(`Consequatur ut velit officia repudiandae quas. 
Quam voluptatem voluptatibus nesciunt et ut nobis. 
Molestiae placeat et non atque error omnis.
Accusamus neque quasi consequatur necessitatibus nihil iure.
Exercitationem minima amet tempore ratione aperiam aut in.
Nulla dolor omnis molestiae ex optio.

Maiores dolores placeat sunt odit quidem. 
Distinctio nesciunt rerum porro et. 
Quis quibusdam aut qui. Qui culpa earum sit.
Consequatur consequuntur autem tempore. Autem ex nesciunt minima officia quisquam natus. 
Et odio pariatur et. Dolorem consequatur voluptas nihil necessitatibus. Ut ipsam aut maiores adipisci similique omnis.

Aperiam dolore est molestiae non qui vel. Delectus nostrum quaerat dolor fugiat. 
Unde qui omnis reprehenderit beatae esse. Aliquid dolores tempora cum.
Repellat ducimus quam ut quasi occaecati. Id sint voluptatum non libero ipsum doloremque. 
Nulla rerum at consequatur deleniti officia ut voluptate enim. Aliquid dolorum rerum qui. 
Vel ratione voluptatem commodi voluptatem illo. Exercitationem sunt omnis voluptates consequatur.

Facilis reiciendis magnam explicabo quod repellendus et fuga cumque. Repudiandae sequi ut eos aliquid nemo.
Eum id aliquid delectus ipsum magni qui. Tempore non ab excepturi ut. Minima id est incidunt quaerat qui ratione. 
Qui omnis commodi blanditiis molestiae ut. Non ex ut nesciunt. Et repellat perferendis eos maiores ratione.
Eius vel aut possimus ipsa omnis. Consequatur quae dicta magnam incidunt consequuntur doloremque dolor. 
Omnis cupiditate tempore sit corporis nam et tempore a.
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.data)
			require.NoError(t, err)
			assert.Less(t, len(got), len(tt.data))
		})
	}
}

func TestGetPublicKey(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "certificate.*.pem")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	certificate, err := GetPublicKey(f.Name()) // getting empty certificate
	require.ErrorIs(t, err, ErrDecodePEMBlock)
	assert.Nil(t, certificate)

	err = writeRandomCertificateToFile(f.Name())
	require.NoError(t, err)

	certificate, err = GetPublicKey(f.Name()) // get valid public key
	require.NoError(t, err)

	assert.NotNil(t, certificate)
	switch interface{}(certificate).(type) {
	case *rsa.PublicKey:
	default:
		t.Errorf("incorrect type, expected *rsa.PublicKey, got: %T", certificate)
	}
	err = os.Remove(f.Name())
	require.NoError(t, err)

	certificate, err = GetPublicKey(f.Name()) // get public key from missing file
	assert.Error(t, err)
	assert.Nil(t, certificate)
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
