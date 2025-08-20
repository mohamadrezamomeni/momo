package user

import (
	"encoding/json"

	userRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/user"
	eventServiceDto "github.com/mohamadrezamomeni/momo/dto/service/event"
	userServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user"
	"github.com/mohamadrezamomeni/momo/entity"
	userEvent "github.com/mohamadrezamomeni/momo/event/user"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	crypt "github.com/mohamadrezamomeni/momo/service/crypt"
)

type User struct {
	userRepo UserRepo
	crypt    *crypt.Crypt
	eventSvc EventService
}

type UserRepo interface {
	FindUserByID(string) (*entity.User, error)
	FindUserByUsername(string) (*entity.User, error)
	Create(*userRepoDto.Create) (*entity.User, error)
	Upsert(*userRepoDto.Create) (*entity.User, error)
	DeleteByUsername(string) error
	DeletePreviousSuperAdmins() error
	FilterUsers(*userRepoDto.FilterUsers) ([]*entity.User, error)
	FindByTelegramID(string) (*entity.User, error)
	Update(string, *userRepoDto.UpdateUser) error
}

type EventService interface {
	Create(*eventServiceDto.CreateEventDto)
}

func New(userRepo UserRepo, crypt *crypt.Crypt, eventSvc EventService) *User {
	return &User{
		eventSvc: eventSvc,
		userRepo: userRepo,
		crypt:    crypt,
	}
}

func (u *User) Create(userDto *userServiceDto.AddUser) (*entity.User, error) {
	return u.userRepo.Create(&userRepoDto.Create{
		IsAdmin:          userDto.IsAdmin,
		TelegramID:       userDto.TelegramID,
		FirstName:        userDto.FirstName,
		TelegramUsername: userDto.TelegramUsername,
		LastName:         userDto.LastName,
		Username:         userDto.Username,
		IsApproved:       userDto.IsApproved,
	})
}

func (u *User) CreateUserAdmin(userDto *userServiceDto.AddUser) (*entity.User, error) {
	err := u.userRepo.DeletePreviousSuperAdmins()
	if err != nil {
		return nil, err
	}

	passwordHashed, err := u.crypt.Hash(userDto.Password)
	if err != nil {
		return nil, err
	}

	return u.userRepo.Upsert(&userRepoDto.Create{
		IsAdmin:   userDto.IsAdmin,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Username:  userDto.Username,
		Password:  passwordHashed,
	})
}

func (u *User) ApproveUser(id string) error {
	scope := "user.service.ApproveUser"
	approve := true
	err := u.userRepo.Update(id, &userRepoDto.UpdateUser{
		IsApproved: &approve,
	})
	if err != nil {
		return err
	}
	eventData := userEvent.UserApproved{
		UserID: id,
	}
	jsonStr, err := json.Marshal(eventData)
	if err != nil {
		return momoError.Wrap(err).Scope(scope)
	}
	u.eventSvc.Create(&eventServiceDto.CreateEventDto{
		Name: "approve_user",
		Data: string(jsonStr),
	})

	return nil
}

func (u *User) FindByID(id string) (*entity.User, error) {
	return u.userRepo.FindUserByID(id)
}

func (u *User) FindByUsername(username string) (*entity.User, error) {
	return u.userRepo.FindUserByUsername(username)
}

func (u *User) DeleteByUsername(username string) error {
	return u.userRepo.DeleteByUsername(username)
}

func (u *User) Filter() ([]*entity.User, error) {
	return u.userRepo.FilterUsers(&userRepoDto.FilterUsers{})
}

func (u *User) FindByTelegramID(tid string) (*entity.User, error) {
	return u.userRepo.FindByTelegramID(tid)
}
