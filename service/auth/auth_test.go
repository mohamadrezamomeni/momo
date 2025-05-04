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
	auth := New(
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
	token, err := auth.createToken(user)
	if err != nil {
		t.Fatalf("error to create token the error is %v", err)
	}
	claim, err := auth.DecodeToken(token)
	if err != nil {
		t.Fatalf("error to decode token the error is %v", err)
	}
	if claim.FirstName != user.FirstName ||
		claim.UserID != user.ID ||
		claim.LastName != user.LastName ||
		user.IsAdmin != claim.IsAdmin {
		t.Fatalf("the error is %v", err)
	}
}
