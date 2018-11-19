package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"telegrambot/config"
	"telegrambot/dao"
	"strings"
)

func main() {
	configuration, err := config.GetConfiguration()
	if err != nil {
		log.Panic(err)
	}

	db, err := dao.Init()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(configuration.BotToken)
	if err != nil {
		log.Panic(err)
	}

	//task.Start(db, bot)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			back1 := "/goonUZtAUcd1LiJBgU"
			back2 := "/erroroneUZtAUcd1LiJBgU"
			back3 := "/successUZtAUcd1LiJBgU"
			back4 := "/errortwoUZtAUcd1LiJBgU"
			chatId := update.CallbackQuery.Message.Chat.ID
			if strings.EqualFold(update.CallbackQuery.Data, back1) {
				msg := tgbotapi.NewMessage(chatId, "回答正确。即使是币安的官方工作人员也不会向您索要这些信息")
				bot.Send(msg)

				markup := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{
							Text:                         "直接报警",
							CallbackData: &back3,
						},
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{
							Text:                         "输入账号",
							CallbackData: &back4,
						},
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{
							Text:                         "不知道",
							CallbackData: &back4,
						},
					),
				)
				reply := tgbotapi.NewMessage(chatId, "当自称币安客服的人员主动私聊您,并给您链接让您登陆并输入账号邮箱密码等信息时，您会怎么做？")
				reply.ReplyMarkup = &markup
				bot.Send(reply)
			}

			if strings.EqualFold(update.CallbackQuery.Data, back2) {
				msg := tgbotapi.NewMessage(chatId, "很遗憾，回答错误，如果您这样做，账户可能不安全。")
				bot.Send(msg)

				msg = tgbotapi.NewMessage(chatId, "很遗憾，您未能通过测试。请再次输入 /start 重新尝试。")
				bot.Send(msg)
			}

			if strings.EqualFold(update.CallbackQuery.Data, back3) {
				msg := tgbotapi.NewMessage(chatId, "回答正确。这是骗子常用的骗术之一，请注意防范哦。")
				bot.Send(msg)

				msg = tgbotapi.NewMessage(chatId, "恭喜您已通过入群测试，现在可以加入 @BinanceChinese 发言咯！")
				bot.Send(msg)
			}

			if strings.EqualFold(update.CallbackQuery.Data, back4) {
				msg := tgbotapi.NewMessage(chatId, "很遗憾，回答错误，您遇到的是骗子，还请提高警惕，切勿上当受骗！")
				bot.Send(msg)

				msg = tgbotapi.NewMessage(chatId, "很遗憾，您未能通过测试。请再次输入 /start 重新尝试。")
				bot.Send(msg)
			}
		}
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("UserId is [%d], userName is [%s], text is [%s]", update.Message.From.ID, update.Message.From.UserName, update.Message.Text )

		chat := update.Message.Chat

		if chat.IsSuperGroup() || chat.IsGroup() {
			// 群组或超级群组
			if len(update.Message.Text) > 0 {
				chatId := chat.ID
				for _, cat := range configuration.Chats {
					if chatId == cat.Id {
						fromId := update.Message.From.ID
						for _, bannedNews := range configuration.BannedNews {
							if strings.Contains(update.Message.Text, bannedNews) {
								// 删除消息
								deleteMessageConfig := tgbotapi.DeleteMessageConfig{
									ChatID: chatId,
									MessageID: update.Message.MessageID,
								}
								bot.DeleteMessage(deleteMessageConfig)

								chatMemberConfig := tgbotapi.ChatMemberConfig{
									ChatID: chatId,
									UserID: fromId,
								}
								// 踢出群组（如果是supergroup,会直接屏蔽此人）
								kickChatMemberConfig := tgbotapi.KickChatMemberConfig{
									ChatMemberConfig: chatMemberConfig,
								}
								bot.KickChatMember(kickChatMemberConfig)
							}
						}
					}
				}
			} else {
				chatId := chat.ID
				newChatMembers := update.Message.NewChatMembers
				if newChatMembers != nil {
					for _, user := range *newChatMembers {
						log.Printf("New chat user id is %d", user.ID)
						if dao.HasTgu(db, user.ID, chatId) {
							dao.UpdateTgu(db, user, chatId)
						} else {
							dao.InsertTgu(db, user, chatId)
						}
					}
				}

				leftChatMember := update.Message.LeftChatMember
				if leftChatMember != nil && leftChatMember.ID != 0 {
					log.Printf("Left chat member id is %d", leftChatMember.ID)
					dao.UpdateIsInGroup(db, leftChatMember.ID, chatId, false, false)
				}
			}
		} else if chat.IsPrivate() {
			userId := update.Message.From.ID
			back1 := "/goonUZtAUcd1LiJBgU"
			back2 := "/erroroneUZtAUcd1LiJBgU"
			if len(update.Message.Text) > 0 {
				if strings.Contains(update.Message.Text, "/start") {
					if dao.HasNoNewGroup(db, userId) {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You've no new group test pending")
						msg.ReplyToMessageID = update.Message.MessageID

						bot.Send(msg)
					} else {
						markup := tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.InlineKeyboardButton{
									Text:                         "会",
									CallbackData: &back2,
								},
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.InlineKeyboardButton{
									Text:                         "不会,这是一个钓鱼网站",
									CallbackData: &back1,
								},
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.InlineKeyboardButton{
									Text:                         "不知道",
									CallbackData: &back2,
								},
							),
						)
						reply := tgbotapi.NewMessage(update.Message.Chat.ID, "当打开 ‘币安’ 网站时，要求输入2FA 的密钥和邮箱密码时，您会输入么？")
						reply.ReplyMarkup = &markup
						bot.Send(reply)
					}
				}

			}

		}
	}
}
