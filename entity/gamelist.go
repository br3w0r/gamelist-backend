package entity

import (
	"time"

	"gorm.io/gorm"
)

type GameProperties struct {
	Model
	Name         string     `gorm:"unique;uniqueIndex" json:"name" binding:"required"`
	Platforms    []Platform `gorm:"many2many:game_platforms" json:"-"`
	ImageURL     string     `json:"image_url" binding:"required,url"`
	YearReleased uint16     `json:"year_released" binding:"required,gte=1000"`
	Genres       []Genre    `gorm:"many2many:game_genres" json:"-"`
}

type Genre struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"varchar(50);unique;" json:"name"`
}

type Platform struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"varchar(20);unique;" json:"name"`
}

type Profile struct {
	ProfileInfo
	Email         string         `gorm:"unique;not null" json:"email" binding:"required"`
	Password      string         `json:"password" binding:"gte=6,lte=70"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:ProfileID" json:"-"`
}

type RefreshToken struct {
	Model
	ProfileID uint64 `json:"-"`
	Token     string `json:"-"`
}

type Social struct {
	ProfileID uint64     `gorm:"primaryKey" json:"-"`
	Type      SocialType `gorm:"foreignKey:TypeID" json:"-"`
	TypeID    uint64     `gorm:"primaryKey" json:"type" binding:"required"`
	Data      string     `gorm:"varchar(70)" json:"data" binding:"gte=2,lte=70"`
}

type SocialType struct {
	Model
	Name string `gorm:"varchar(20);unique" json:"name"`
}

type ProfileGame struct {
	Profile    ProfileInfo    `gorm:"foreignKey:ProfileID" json:"-"`
	ProfileID  uint64         `gorm:"primaryKey" json:"-"`
	Game       GameProperties `gorm:"foreignKey:GameID" json:"game"`
	GameID     uint64         `gorm:"primaryKey" json:"-"`
	ListType   ListType       `gorm:"foreignKey:ListTypeID" json:"-"`
	ListTypeID uint64         `json:"list_type"`
}

type ListType struct {
	Model
	Name string `gorm:"varchar(20);unique" json:"name"`
}
