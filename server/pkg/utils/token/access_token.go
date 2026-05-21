package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unitechio/oss-monorepo/server/pkg/utils/str"
)

type AccessTokenUtil struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	Salt       string
}

func NewAccessTokenUtil(publicKeyB64, privateKeyB64 string) (*AccessTokenUtil, error) {
	pub, err := parsePublicKey(publicKeyB64)
	if err != nil {
		return nil, err
	}

	priv, err := parsePrivateKey(privateKeyB64)
	if err != nil {
		return nil, err
	}

	return &AccessTokenUtil{
		PublicKey:  pub,
		PrivateKey: priv,
		Salt:       "thisI$$alt",
	}, nil
}

func parsePrivateKey(b64 string) (*rsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS8PrivateKey(decoded)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid RSA private key")
	}

	return rsaKey, nil
}

func parsePublicKey(b64 string) (*rsa.PublicKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(decoded)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid RSA public key")
	}

	return rsaPub, nil
}

func (u *AccessTokenUtil) GenerateAccessToken(
	contractID string,
	issuer string,
	audience string,
	phoneNumber string,
	partnerID string,
) (string, error) {

	now := time.Now()
	exp := now.Add(1 * time.Hour)

	claims := jwt.MapClaims{
		"jti":          GenerateUUID(),
		"sub":          contractID,
		"contract_id":  contractID,
		"partner_id":   partnerID,
		"phone_number": phoneNumber,
		"iss":          issuer,
		"aud":          audience,
		"iat":          now.Unix(),
		"nbf":          now.Unix(),
		"exp":          exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(u.PrivateKey)
}

type AccessToken struct {
	ContractID  string `json:"contract_id"`
	PartnerID   string `json:"partner_id"`
	PhoneNumber string `json:"phone_number"`
	Subject     string `json:"sub"`
	Issuer      string `json:"iss"`
	Audience    string `json:"aud"`
}

func (u *AccessTokenUtil) Parse(tokenStr string) (*AccessToken, error) {
	tokenStr = ResolveBearer(tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return u.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &AccessToken{
		ContractID:  str.ToString(claims["contract_id"]),
		PartnerID:   str.ToString(claims["partner_id"]),
		PhoneNumber: str.ToString(claims["phone_number"]),
		Subject:     str.ToString(claims["sub"]),
		Issuer:      str.ToString(claims["iss"]),
		Audience:    str.ToString(claims["aud"]),
	}, nil
}

func (u *AccessTokenUtil) Validate(tokenStr string) bool {
	_, err := u.Parse(tokenStr)
	return err == nil
}
