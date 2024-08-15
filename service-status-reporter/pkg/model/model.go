package model

type Service struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	SteOwnerID  *int     `json:"ste_owner_id" gorm:"index"`
	SteOwner    *Contact `json:"ste_owner" gorm:"foreignKey:SteOwnerID;references:ID"`
	DeveloperID *int     `json:"developer_id" gorm:"index"`
	Developer   *Contact `json:"developer" gorm:"foreignKey:DeveloperID;references:ID"`
	TeamID      *int     `json:"team_id" gorm:"index"`
	Team        *Team    `json:"team" gorm:"foreignKey:TeamID;references:ID"`
	ProdEnv     string   `json:"prod_env"`
	PreprodEnv  string   `json:"preprod_env"`
	NonprodEnv  string   `json:"nonprod_env"`
}

type Contact struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Slack  string `json:"slack"`
	RoleID *int   `json:"role_id" gorm:"index"`
	Role   Role   `json:"role" gorm:"foreignKey:RoleID;references:ID"`
}

type Role struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

type Team struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}
