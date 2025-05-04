package auth

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (a *Auth) createToken(user *entity.User) (string, error) {
	scope := "auth.createToken"
	claim := Claim{
		UserID:    user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		LastName:  user.LastName,
		FirstName: user.FirstName,
	}
	jsonData, err := json.Marshal(claim)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Errorf("error to marshal claim the input is %+v", *user)
	}
	encoded, err := a.crypt.Encrypt(string(jsonData))
	if err != nil {
		return "", err
	}
	encryptedClaim := EncryptedClaim{
		Data: encoded,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.config.ExpireTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "momo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, encryptedClaim)
	return token.SignedString([]byte(a.config.SecretKey))
}

func (a *Auth) DecodeToken(tokenString string) (*Claim, bool, error) {
	scope := "auth.DecodeToken"

	token, err := jwt.ParseWithClaims(tokenString, &EncryptedClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.SecretKey), nil
	})
	if err != nil {
		return nil, false, momoError.Wrap(err).Scope(scope).Errorf("error to marshal claim the input is %s", tokenString)
	}
	encryptedClaim, ok := token.Claims.(*EncryptedClaim)

	if !ok {
		return nil, false, momoError.Scope(scope).Errorf("the token isn't valid the input is %s", tokenString)
	}

	data, err := a.crypt.Decrypt(encryptedClaim.Data)
	if err != nil {
		return nil, false, err
	}

	var claim Claim
	err = json.Unmarshal([]byte(data), &claim)
	if err != nil {
		return nil, false, momoError.Wrap(err).Scope(scope).Errorf("error to unmarshal encrypted claim the input is %s", tokenString)
	}

	return &claim, token.Valid, nil
}
