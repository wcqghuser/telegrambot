package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"telegrambot/model"
	"time"
	"fmt"
)

func Init() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/telegram")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 群组用户是否存在
func HasTgu(db *sql.DB, userId int, chatId int64) (bool) {
	stmtOut, err := db.Prepare("SELECT COUNT(*) FROM telegram_group_user WHERE user_id=? AND chat_id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var count int
	err = stmtOut.QueryRow(userId, chatId).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	return count > 0
}

// 插入群组用户
func InsertTgu(db *sql.DB, user tgbotapi.User, chatId int64) {
	stmtIns, err := db.Prepare("INSERT INTO telegram_group_user(user_id, chat_id, is_bot, user_name, first_name, last_name, language_code, is_in_group, is_complete_test, create_time) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	groupUser := model.TelegramUserGroup{
		UerId: user.ID,
		ChatId: chatId,
		IsBot: user.IsBot,
		UserName: user.UserName,
		FirstName: user.FirstName,
		LastName: user.LastName,
		LanguageCode: user.LanguageCode,
		IsInGroup: true,
		IsCompleteTest: false,
		CreateTime: time.Now(),
	}

	_, err = stmtIns.Exec(groupUser.UerId, groupUser.ChatId, groupUser.IsBot, groupUser.UserName, groupUser.FirstName,
		groupUser.LastName, groupUser.LanguageCode, groupUser.IsInGroup, groupUser.CreateTime)
	if err != nil {
		panic(err.Error())
	}
}

// 修改群组用户
func UpdateTgu(db *sql.DB, user tgbotapi.User, chatId int64) {
	stmtIns, err := db.Prepare("UPDATE telegram_group_user SET is_bot=?, user_name=?, first_name=?, last_name=?, language_code=?, is_in_group=? WHERE user_id=? AND chat_id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	groupUser := model.TelegramUserGroup{
		UerId: user.ID,
		ChatId: chatId,
		IsBot: user.IsBot,
		UserName: user.UserName,
		FirstName: user.FirstName,
		LastName: user.LastName,
		LanguageCode: user.LanguageCode,
		IsInGroup: true,
	}

	_, err = stmtIns.Exec(groupUser.IsBot, groupUser.UserName, groupUser.FirstName,
		groupUser.LastName, groupUser.LanguageCode, groupUser.IsInGroup, groupUser.UerId, groupUser.ChatId)
	if err != nil {
		panic(err.Error())
	}
}

// 修改群组用户是否在群内
func UpdateIsInGroup(db *sql.DB, userId int, chatId int64, isInGroup bool, isCompleteTest bool) {
	stmtUpd, err := db.Prepare("UPDATE telegram_group_user SET is_in_group=?, is_complete_test=? WHERE user_id=? AND chat_id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtUpd.Close()

	_, err = stmtUpd.Exec(isInGroup, isCompleteTest, userId, chatId)
	if err != nil {
		panic(err.Error())
	}
}

// 是否没有加入群组
func HasNoNewGroup(db *sql.DB, userId int) (bool) {
	stmtOut, err := db.Prepare("SELECT COUNT(*) FROM telegram_group_user WHERE user_id=? AND is_in_group=1 AND is_complete_test=0")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var count int
	err = stmtOut.QueryRow(userId).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	return count == 0
}

func TimeoutUser(db *sql.DB, limitTime string) ([]model.TelegramUserGroup) {
	rows, err := db.Query("SELECT user_id, chat_id FROM telegram_group_user WHERE is_in_group=1 AND is_complete_test=0 AND create_time < ?", limitTime)
	if err != nil {
		panic(err.Error())
	}

	var userGroups []model.TelegramUserGroup
	i := 0
	for rows.Next() {
		var userGroup model.TelegramUserGroup
		var userId int
		var chatId int64
		if err := rows.Scan(&userId, &chatId); err != nil {
			panic(err.Error())
		}
		fmt.Printf("userId %d chatId %d", userId, chatId)
		userGroup.UerId = userId
		userGroup.ChatId = chatId
		userGroups = append(userGroups, userGroup)
		i++
	}

	return userGroups
}
