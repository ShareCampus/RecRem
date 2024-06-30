package models

type OpenAI struct {
	Id    int    `gorm:"column:id;primary_key"`
	Token string `gorm:"column:token"`
	Model string `gorm:"column:model"`
}
