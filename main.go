package main

import (
	"database/sql"
	"github.com/Syfaro/telegram-bot-api"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"time"
)

const ADMIN int64 = 375806606

func main() {
	bot := BotStart()
	my_db := DBStart()
	defer my_db.Close()
	go ChatMaker(my_db, bot)
	BotUpdateLoop(bot, my_db)
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("New chat"),
		tgbotapi.NewKeyboardButton("Leave chat"),
	),
)

func BotStart() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI("1022122500:AAFy8sDJFUlgw0e7JURelghBPv_is5kG7ck") //1057128816:AAE3MrZxSXnMPV1UNYuLbOQobd-sxUIhGw4 - AnonStud 1022122500:AAFy8sDJFUlgw0e7JURelghBPv_is5kG7ck - Freedom
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Autorised on account %s", bot.Self.UserName)

	return bot
}
func BotUpdateLoop(my_bot *tgbotapi.BotAPI, database *sql.DB) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := my_bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			switch update.Message.Text {
			case "New chat":
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

				if CheckReg(update.Message.From.ID, database, my_bot) {
					if IsFree(update.Message.From.ID, database, my_bot) {
						if !IsSearch(update.Message.From.ID, database, my_bot) {
							ChangeSearch(database, update.Message.From.ID, 1, my_bot)
							msg.Text = "Search started"
						} else {
							msg.Text = "You are searching already"
						}
					} else {
						msg.Text = "You are chatting already"
					}
				} else {
					msg.Text = "You need /start first"
				}

				_, err := my_bot.Send(msg)
				if err != nil {
					ErrorCatch(err.Error(), my_bot)
				}
				continue

			case "Leave chat":
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

				if CheckReg(update.Message.From.ID, database, my_bot) {
					if FindChat(update.Message.From.ID, database, my_bot) != 0 {
						chat_id := FindChat(update.Message.From.ID, database, my_bot)
						DeleteChat(update.Message.From.ID, database, my_bot)
						ChangeState(database, update.Message.From.ID, 0, my_bot)
						msg.Text = "You leaved a chat"

						DeleteChat(chat_id, database, my_bot)
						ChangeState(database, chat_id, 0, my_bot)
						_, err := my_bot.Send(tgbotapi.NewMessage(int64(chat_id), "The stranger leave the chat"))
						if err != nil {
							ErrorCatch(err.Error(), my_bot)
						}
					} else {
						msg.Text = "You are not chatting now!"
					}
				} else {
					msg.Text = "You need to /start first"
				}

				_, err := my_bot.Send(msg)
				if err != nil {
					ErrorCatch(err.Error(), my_bot)
				}
				continue
			}

			if FindChat(update.Message.From.ID, database, my_bot) != 0 {
				chat_id := FindChat(update.Message.From.ID, database, my_bot)
				msg := tgbotapi.NewMessage(int64(chat_id), "")

				msg.Text = update.Message.Text
				if msg.Text != "" {
					_, err := my_bot.Send(msg)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Photo != nil {
					photo := tgbotapi.NewPhotoShare(int64(chat_id), "")
					photo.FileID = (*update.Message.Photo)[2].FileID
					_, err := my_bot.Send(photo)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Voice != nil {
					voice := tgbotapi.NewVoiceShare(int64(chat_id), "")
					voice.FileID = update.Message.Voice.FileID
					_, err := my_bot.Send(voice)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Animation != nil {
					voice := tgbotapi.NewAnimationShare(int64(chat_id), "")
					voice.FileID = update.Message.Animation.FileID
					_, err := my_bot.Send(voice)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Audio != nil {
					audio := tgbotapi.NewAudioShare(int64(chat_id), "")
					audio.FileID = update.Message.Audio.FileID
					_, err := my_bot.Send(audio)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Sticker != nil {
					sticker := tgbotapi.NewStickerShare(int64(chat_id), "")
					sticker.FileID = update.Message.Sticker.FileID
					_, err := my_bot.Send(sticker)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Document != nil {
					doc := tgbotapi.NewDocumentShare(int64(chat_id), "")
					doc.FileID = update.Message.Document.FileID
					_, err := my_bot.Send(doc)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Video != nil {
					video := tgbotapi.NewVideoShare(int64(chat_id), "")
					video.FileID = update.Message.Video.FileID
					_, err := my_bot.Send(video)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.VideoNote != nil {
					video_note := tgbotapi.NewVideoNoteShare(int64(chat_id),0, "")
					video_note.FileID = update.Message.VideoNote.FileID
					video_note.Length = update.Message.VideoNote.Length
					_, err := my_bot.Send(video_note)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else {
					msg.Text = "Bot cannot send this yet! Please, contact with creator"
					_, err := my_bot.Send(msg)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				}
				continue
			} else {
				continue
			}
		}

		chat_id := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chat_id, "")

		switch update.Message.Command() {
		case "start":
			if !CheckReg(update.Message.From.ID, database, my_bot) {
				FirstStart(update.Message.From.ID, database, my_bot)
				msg.Text = "Hello, this is Freedom chat, where you can freely express your minds and talk with other strangers.\n\n" +
					"To start the chat, send /go_chat command or press \"New chat\" button\n\n" +
					"To leave the chat, send /leave_chat command or press \"Leave chat\" button\n\n" +
					"Bot doesn't store any personal data, so chats are fully anonymous.\n\n" +
					"If You want to check how the bot works - check my video (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
					"If You have some questions or suggestions, please, feel free to contact with me, @YUART\n\n" +
					"Also, check my Patreon page (https://www.patreon.com/artemkakun) if you want receive some bonuses from me :)\n"
				msg.ReplyMarkup = numericKeyboard
			} else {
				msg.Text = "Hello, this is Freedom chat, where you can freely express your minds and talk with other strangers\n" +
					"To start the chat, send /go_chat command or press \"New chat\" button\n" +
					"To leave the chat, send /leave_chat command or press \"Leave chat\" button\n" +
					"Bot doesn't store any personal data, so chats are fully anonymous" +
					"If You want to check how the bot works - check my video (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)" +
					"If You have some questions or suggestions, please, feel free to contact with me, @YUART\n" +
					"Also, check my Patreon page (https://www.patreon.com/artemkakun) if you want receive some bonuses from me :)\n"
				msg.ReplyMarkup = numericKeyboard
			}
		case "go_chat":
			if CheckReg(update.Message.From.ID, database, my_bot) {
				if IsFree(update.Message.From.ID, database, my_bot) {
					if !IsSearch(update.Message.From.ID, database, my_bot) {
						ChangeSearch(database, update.Message.From.ID, 1, my_bot)
						msg.Text = "Search started"
					} else {
						msg.Text = "You are searching already"
					}
				} else {
					msg.Text = "You are chatting already"
				}
			} else {
				msg.Text = "You need to /start first"
			}
		case "leave_chat":
			if CheckReg(update.Message.From.ID, database, my_bot) {
				if FindChat(update.Message.From.ID, database, my_bot) != 0 {
					chat_id := FindChat(update.Message.From.ID, database, my_bot)
					DeleteChat(update.Message.From.ID, database, my_bot)
					ChangeState(database, update.Message.From.ID, 0, my_bot)
					msg.Text = "You leaved a chat"

					DeleteChat(chat_id, database, my_bot)
					ChangeState(database, chat_id, 0, my_bot)
					_, err := my_bot.Send(tgbotapi.NewMessage(int64(chat_id), "The stranger leave the chat"))
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else {
					msg.Text = "You are not chatting now!"
				}
			} else {
				msg.Text = "You need to /start first"
			}
		}

		_, err := my_bot.Send(msg)
		if err != nil {
			ErrorCatch(err.Error(), my_bot)
		}
	}
}

func ChatMaker(database *sql.DB, my_bot *tgbotapi.BotAPI) {
	for true {
		free_users := FindFree(database, my_bot)
		users_amount := len(free_users)
		if users_amount > 1 {
			rand.Seed(time.Now().UnixNano())
			first_user := rand.Intn(users_amount)
			second_user := rand.Intn(users_amount)

			for second_user == first_user {
				second_user = rand.Intn(users_amount)
			}

			ChangeSearch(database, free_users[first_user], 0, my_bot)
			ChangeSearch(database, free_users[second_user], 0, my_bot)
			ChangeState(database, free_users[first_user], 1, my_bot)
			ChangeState(database, free_users[second_user], 1, my_bot)
			AddChat(free_users[first_user], free_users[second_user], database, my_bot)
			AddChat(free_users[second_user], free_users[first_user], database, my_bot)

			msg := tgbotapi.NewMessage(int64(free_users[first_user]), "")
			msg.Text = "Now you can chat"
			_, err := my_bot.Send(msg)
			if err != nil {
				ErrorCatch(err.Error(), my_bot)
			}

			msg = tgbotapi.NewMessage(int64(free_users[second_user]), "")
			msg.Text = "Now you can chat"
			_, err = my_bot.Send(msg)
			if err != nil {
				ErrorCatch(err.Error(), my_bot)
			}
		}
		amt := time.Duration(1000)
		time.Sleep(time.Millisecond * amt)
	}
}
func DBStart() *sql.DB {
	my_db, err := sql.Open("mysql", "root:11hahozeGood!@/anonstudchat")
	if err != nil {
		log.Panic(err)
	} else {
		err = my_db.Ping()
		if err != nil {
			log.Panic(err)
		}
	}
	return my_db
}
func FirstStart(user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) {
	stmtIns, err := my_db.Prepare("INSERT INTO users_info VALUES (?, ?, ?)")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	_, err = stmtIns.Exec(user_id, 0, 0)
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}
}
func CheckReg(user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) bool {
	stmtOut, err := my_db.Prepare("SELECT user_id FROM users_info WHERE user_id = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	var is_reg int
	err = stmtOut.QueryRow(user_id).Scan(&is_reg)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			ErrorCatch(err.Error(), my_bot)
			panic(err.Error())
		}
		return false
	}

	err = stmtOut.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	if is_reg != 0 {
		return true
	} else {
		return false
	}

}
func IsFree(user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) bool {
	stmtOut, err := my_db.Prepare("SELECT user_free FROM users_info WHERE user_id = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	var is_free int
	err = stmtOut.QueryRow(user_id).Scan(&is_free)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			ErrorCatch(err.Error(), my_bot)
			panic(err.Error())
		}
		return false
	}

	err = stmtOut.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	if is_free == 0 {
		return true
	} else {
		return false
	}

}
func IsSearch(user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) bool {
	stmtOut, err := my_db.Prepare("SELECT is_search FROM users_info WHERE user_id = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	var is_free int
	err = stmtOut.QueryRow(user_id).Scan(&is_free)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			ErrorCatch(err.Error(), my_bot)
			panic(err.Error())
		}
		return false
	}

	err = stmtOut.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	if is_free == 1 {
		return true
	} else {
		return false
	}

}
func FindFree(my_db *sql.DB, my_bot *tgbotapi.BotAPI) []int {
	stmtOut, err := my_db.Query("SELECT user_id FROM users_info WHERE user_free = 0 AND is_search = 1")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	var user_free []int
	var one_user int

	for stmtOut.Next() {
		err = stmtOut.Scan(&one_user)
		if err != nil {
			err = stmtOut.Close()
			if err != nil {
				ErrorCatch(err.Error(), my_bot)
				panic(err.Error())
			}
			user_free = append(user_free, 0)
			return user_free
		}
		user_free = append(user_free, one_user)
	}

	err = stmtOut.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	return user_free
}
func ChangeSearch(my_db *sql.DB, user_id int, status int, my_bot *tgbotapi.BotAPI) {
	stmtIns, err := my_db.Prepare("UPDATE users_info SET is_search = ? WHERE user_id = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	_, err = stmtIns.Exec(status, user_id)
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}
}
func ChangeState(my_db *sql.DB, user_id int, status int, my_bot *tgbotapi.BotAPI) {
	stmtIns, err := my_db.Prepare("UPDATE users_info SET user_free = ? WHERE user_id = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	_, err = stmtIns.Exec(status, user_id)
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}
}
func AddChat(first_user_id int, second_user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) {
	stmtIns, err := my_db.Prepare("INSERT INTO chat_buffer VALUES (?, ?)")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	_, err = stmtIns.Exec(first_user_id, second_user_id)
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}
}
func FindChat(user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) int {
	stmtOut, err := my_db.Prepare("SELECT second_user FROM chat_buffer WHERE first_user = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	var second_user int
	err = stmtOut.QueryRow(user_id).Scan(&second_user)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			ErrorCatch(err.Error(), my_bot)
			panic(err.Error())
		}
		return 0
	}

	err = stmtOut.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	return second_user
}
func DeleteChat(user_id int, my_db *sql.DB, my_bot *tgbotapi.BotAPI) {
	stmtIns, err := my_db.Prepare("DELETE FROM chat_buffer WHERE first_user = ?")
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	_, err = stmtIns.Exec(user_id)
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		ErrorCatch(err.Error(), my_bot)
		panic(err.Error())
	}
}

func ErrorCatch(err string, my_bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(ADMIN, err)
	my_bot.Send(msg)
}