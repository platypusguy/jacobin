/*
	Standard (Java 21):
	https://docs.oracle.com/en/java/javase/21/docs/specs/security/standard-names.html
*/

package security

import (
	"crypto/rand"
	"math/big"
)

type DHParameters struct {
	P *big.Int
	G *big.Int
}

type DHKeyPair struct {
	PrivateKey *big.Int
	PublicKey  *big.Int
	Params     *DHParameters
}

func (dh *DHKeyPair) ComputeShared(peerPub *big.Int) (*big.Int, error) {
	secret := new(big.Int).Exp(peerPub, dh.PrivateKey, dh.Params.P)
	return secret, nil
}

func NewDHParameters(bits int) (*DHParameters, error) {
	p, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &DHParameters{P: p, G: big.NewInt(2)}, nil
}

func GenerateDHKeyPair(params *DHParameters) (*DHKeyPair, error) {
	priv, _ := rand.Int(rand.Reader, new(big.Int).Sub(params.P, big.NewInt(2)))
	priv = priv.Add(priv, big.NewInt(2))
	pub := new(big.Int).Exp(params.G, priv, params.P)
	return &DHKeyPair{PrivateKey: priv, PublicKey: pub, Params: params}, nil
}

func (kp *DHKeyPair) ComputeSharedSecret(peerPub *big.Int) (*big.Int, error) {
	secret := new(big.Int).Exp(peerPub, kp.PrivateKey, kp.Params.P)
	return secret, nil
}
