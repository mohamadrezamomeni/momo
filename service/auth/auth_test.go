package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mohamadrezamomeni/momo/entity"
	mockUserService "github.com/mohamadrezamomeni/momo/mocks/service/user"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
	crypt "github.com/mohamadrezamomeni/momo/service/crypt"
)

func TestTokenProcessing(t *testing.T) {
	auth1 := New(
		mockUserService.New(),
		crypt.New(&crypt.CryptConfig{EncryptKey: "qFQGz0yE0nYBxqi9"}),
		&AuthConfig{
			SecretKey:  "Ve9Bu3A0ZbhunKFd",
			ExpireTime: 20,
		},
	)
	user := &entity.User{
		ID:        uuid.New().String(),
		Username:  utils.RandomString(5),
		FirstName: utils.RandomString(5),
		LastName:  utils.RandomString(5),
	}
	token, err := auth1.createToken(user)
	if err != nil {
		t.Fatalf("error to create token the error is %v", err)
	}
	claim, isValid, err := auth1.DecodeToken(token)
	if err != nil {
		t.Fatalf("error to decode token the error is %v", err)
	}
	if !isValid {
		t.Fatalf("this token must be valid")
	}
	if claim.FirstName != user.FirstName ||
		claim.UserID != user.ID ||
		claim.LastName != user.LastName ||
		user.IsAdmin != claim.IsAdmin {
		t.Fatalf("the error is %v", err)
	}

	auth2 := New(
		mockUserService.New(),
		crypt.New(&crypt.CryptConfig{EncryptKey: "qFQGz0yE0nYBxqi9"}),
		&AuthConfig{
			SecretKey:  "Ve9Bu3A0ZbhunKFn",
			ExpireTime: 20,
		},
	)

	_, isValid, _ = auth2.DecodeToken(token)

	if isValid {
		t.Errorf("the token must be invalid")
	}
}
