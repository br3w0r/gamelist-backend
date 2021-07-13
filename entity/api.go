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

type TokenPair struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type GameListRequest struct {
	GameId   uint64 `json:"game_id" binding:"required"`
	ListType uint64 `json:"list_type"`
}

type TypedGameListProperties struct {
	GameProperties
	ListTypeID uint64 `json:"user_list"`
}

type SearchRequest struct {
	Name string `json:"name"`
}

type GameSearchResult struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type GameDetailsRequest struct {
	Id uint64 `json:"id"`
}

type GameDetailsResponse struct {
	Game      TypedGameListProperties `json:"game"`
	Platforms []Platform              `json:"platforms"`
	Genres    []Genre                 `json:"genres"`
}
