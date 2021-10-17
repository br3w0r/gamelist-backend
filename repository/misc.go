package repository

import (
	"fmt"

	"github.com/br3w0r/gamelist-backend/entity"
	utilErrs "github.com/br3w0r/gamelist-backend/util/errors"
	"gorm.io/gorm"
)

func CheckSocialTypes(db *gorm.DB, profile *entity.Profile) error {
	for i := range profile.Socials {
		res := db.First(&entity.SocialType{}, profile.Socials[i].TypeID)
		if res.Error != nil {
			return utilErrs.FromGORM(res, fmt.Sprint("couldn't find social of type ", profile.Socials[i].TypeID))
		}
	}
	return nil
}
