package utils

import "errors"

var ErrDecodePEMBlock = errors.New("failed to decode PEM block containing public key")
