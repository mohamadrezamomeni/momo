package usertier

import (
	"os"
	"testing"

	tierRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/tier"
	usertierRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/user_tier"
	"github.com/mohamadrezamomeni/momo/repository/migrate"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"
	tierRepository "github.com/mohamadrezamomeni/momo/repository/sqllite/tier"
)

var (
	userTierRepo *UserTier
	tierRepo     *tierRepository.Tier
)

func TestMain(m *testing.M) {
	config := &sqllite.DBConfig{
		Dialect:    "sqlite3",
		Path:       "test-user.db",
		Migrations: "./repository/sqllite/migrations",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqllite.New(config)

	userTierRepo = New(db)
	tierRepo = tierRepository.New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreating(t *testing.T) {
	defer userTierRepo.DeleteAll()
	err := userTierRepo.Create(usertierRepoDto.Create{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
		Tier:   "silver",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
}

func TestDeleteing(t *testing.T) {
	defer userTierRepo.DeleteAll()
	err := userTierRepo.Create(usertierRepoDto.Create{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
		Tier:   "silver",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	err = userTierRepo.Delete(&usertierRepoDto.IdentifyUserTier{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
		Tier:   "silver",
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
}

func TestUserTiers(t *testing.T) {
	defer userTierRepo.DeleteAll()

	tierRepo.Create(&tierRepositoryDto.CreateTier{
		IsDefault: true,
		Name:      "silver",
	})
	tierRepo.Create(&tierRepositoryDto.CreateTier{
		IsDefault: false,
		Name:      "gold",
	})
	tierRepo.Create(&tierRepositoryDto.CreateTier{
		IsDefault: false,
		Name:      "platinum",
	})
	tierRepo.Create(&tierRepositoryDto.CreateTier{
		IsDefault: false,
		Name:      "diamond",
	})

	userTierRepo.Create(usertierRepoDto.Create{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
		Tier:   "platinum",
	})
	userTierRepo.Create(usertierRepoDto.Create{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c3e",
		Tier:   "gold",
	})
	userTierRepo.Create(usertierRepoDto.Create{
		UserID: "0393ed06-29bb-41c2-b3f4-6382a6729c33",
		Tier:   "gold",
	})

	tiers, err := userTierRepo.FilterTiersBelongToUser("0393ed06-29bb-41c2-b3f4-6382a6729c3e")
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}
	if len(tiers) != 3 {
		t.Fatalf("error to compare data we expected %d records but we got %d", 3, len(tiers))
	}

	seen := map[string]struct{}{}
	for _, tier := range tiers {
		seen[tier.Name] = struct{}{}
	}
	if _, isExist := seen["silver"]; !isExist {
		t.Fatal("we expected silver be exist")
	}

	if _, isExist := seen["gold"]; !isExist {
		t.Fatal("we expected silver be exist")
	}

	if _, isExist := seen["platinum"]; !isExist {
		t.Fatal("we expected silver be exist")
	}
}
