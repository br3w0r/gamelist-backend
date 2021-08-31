//go:generate mockgen -destination=mock_service.go -package=service github.com/br3w0r/gamelist-backend/repository GamelistRepository

package service

import (
	"testing"

	"github.com/br3w0r/gamelist-backend/entity"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
)

var (
	mockGameProperty = entity.GameProperties{
		Name: "Test Game",
		Platforms: []entity.Platform{
			{Name: "Test Platform"},
		},
		ImageURL:     "image/url/im.png",
		YearReleased: 1970,
		Genres: []entity.Genre{
			{Name: "Test Genre"},
		},
	}

	mockGamePropertyArray = []entity.GameProperties{mockGameProperty}
	mockProfile           = entity.Profile{
		ProfileInfo: entity.ProfileInfo{
			Nickname:    "test",
			Description: "test desc",
			GamesListed: 100,
		},
		Email:    "test@mail.com",
		Password: "test_pass",
	}
)

func TestCreateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := NewMockGamelistRepository(ctrl)

	repo.EXPECT().
		CreateProfile(gomock.Any()).
		DoAndReturn(func(profile entity.Profile) error {
			return bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(mockProfile.Password))
		}).
		Times(1)

	service := NewGameListService(repo, "")

	convey.Convey("service.CreateProfile() should return nil error", t, func() {
		err := service.CreateProfile(mockProfile)

		convey.So(err, convey.ShouldBeNil)
	})
}
