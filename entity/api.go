package entity

type ProfileCreds struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type ProfileInfo struct {
	Model
	Nickname    string   `gorm:"varchar(20);unique" json:"nickname" binding:"gte=2,lte=20"`
	Description string   `gorm:"varchar(120)" json:"description" binding:"lte=120"`
	GamesListed uint     `json:"games_listed"`
	Socials     []Social `gorm:"foreignKey:ProfileID"`
}

type LoginProfile struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `gorm:"varchar(70);not null" json:"password" binding:"gte=6,lte=70"`
}
