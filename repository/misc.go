package repository

import (
	"fmt"

	"bitbucket.org/br3w0r/gamelist-backend/entity"
	"gorm.io/gorm"
)

func CheckSocialTypes(db *gorm.DB, profile *entity.Profile) error {
	for i := range profile.Socials {
		err := db.First(&entity.SocialType{}, profile.Socials[i].TypeID).Error
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("couldn't find social of type %d", profile.Socials[i].TypeID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
