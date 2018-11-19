package model

import "time"

type TelegramUserGroup struct {
	Id int64
	UerId int
	ChatId int64
	IsBot  bool
	UserName string
	FirstName string
	LastName string
	LanguageCode string
	IsInGroup bool
	IsCompleteTest bool
	CreateTime time.Time
}