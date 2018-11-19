package task

import (
	"time"
	"fmt"
	"telegrambot/dao"
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(db *sql.DB, bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			fmt.Printf("ticked at %v \n", time.Now())
			limitTime := time.Now().Add(-time.Minute * 10).Format("2006-01-02 15:04:05")
			userGroups := dao.TimeoutUser(db, limitTime)
			if userGroups != nil {
				for _, userGroup := range userGroups {
					dao.UpdateIsInGroup(db, userGroup.UerId, userGroup.ChatId,false, false)

					chatMemberConfig := tgbotapi.ChatMemberConfig{
						ChatID: userGroup.ChatId,
						UserID: userGroup.UerId,
					}
					// 踢出群组（如果是supergroup,会直接屏蔽此人）
					kickChatMemberConfig := tgbotapi.KickChatMemberConfig{
						ChatMemberConfig: chatMemberConfig,
					}
					bot.KickChatMember(kickChatMemberConfig)
				}
			}
		}
	}()
}
