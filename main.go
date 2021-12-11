package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	telegraph "github.com/meinside/telegraph-go"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var html = `<h3 >Машина %s водителя %s</h3>
    <img style="width: 70%%; " src="https://i.imgur.com/IekblLL.png"/>
    <p>Редкость: <span style="color:#1E90FF;">%s</span></p>
            <p>Мощность мотора: %d<br/>
Аэродинамика: %d<br/>
Прижимная сила: %d<br/>
Срок действия до %s<br/>
</p>

`

var mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Купить машину"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мой аккаунт"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Гонка"),
		tgbotapi.NewKeyboardButton("Маркетплейс"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Информация"),
		tgbotapi.NewKeyboardButton("Саппорт"),
	),


)

var marketplaceMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Продать моменты"),
		tgbotapi.NewKeyboardButton("Купить моменты"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мои моменты на продаже"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var infoMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("О нас", "about us"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Как купить или получить машину?", "how to buy car"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Как продать машину?", "how to sale car"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Правила гонок", "race rules"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Как участвовать в гонке?", "how go to race"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Почему машина не может участвовать в гонке?", "invalid car"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Как улучшить машину?", "upgrade car"),
	),
)

var supportMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Написать в поддержку", "mailto:romix0000@gmail.com"),
	),
)

var buyCar = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Common машина - 5$", "buyCommonCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Uncommon машина - 10$", "buyUncommonCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Rare машина - 50$", "buyRareCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Epic машина - 100$", "buyEpicCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Legendary машина - 500$", "buyLegendCar"),
	),
)

var garage = []car{
	//{
	//	owner:       users[0],
	//	name:        "RX 9002",
	//	rarity:      "Common",
	//	image:       "https://i.imgur.com/IekblLL.png",
	//	power:       82,
	//	ergo:        24,
	//	downforce:   63,
	//	expiredData: "04.03.2022",
	//},
	//{
	//	owner:       users[0],
	//	name:        "RX 9003",
	//	rarity:      "Common",
	//	image:       "https://i.imgur.com/IekblLL.png",
	//	power:       99,
	//	ergo:        64,
	//	downforce:   92,
	//	expiredData: "04.03.2022",
	//},
	//{
	//	owner:       users[0],
	//	name:        "RX 9004",
	//	rarity:      "Common",
	//	image:       "https://i.imgur.com/IekblLL.png",
	//	power:       64,
	//	ergo:        34,
	//	downforce:   78,
	//	expiredData: "04.03.2022",
	//},
}

type car struct {
	name        string
	rarity      string
	image       string
	power       int
	ergo        int
	downforce   int
	expiredData string
	owner       user
}

type user struct {
	username      string
	hash          string
	walletAddress string
	balance       string
}

var users = []user{
	//	{
	//	username:      "rZaharov10",
	//	hash:          "",
	//	walletAddress: "",
	//},
}

func getUserIndex(username string) int {
	for i, v := range users {
		if v.username == username {
			return i
		}
	}
	return -1
}

var client *telegraph.Client

var carPageTemplate = `
    <h3 >Машина %s водителя %s</h3>
    <img style="width: 70%%; " src="%s"/>
    <p>Редкость: <span style="color:#1E90FF;">%s</span></p>
            <p>Мощность мотора: %d</p>
            <p>Аэродинамика: %d</p>
            <p>Прижимная сила: %d</p>
            <p>Срок действия до %s</p>`

func main() {
	telegraph.Verbose = true
	var savedAccessToken string = "eb5a3dea9b91dfe7fa268b9e4e0641cf46fa7a8f86dc900bb147e1f8a2e9"

	bot, err := tgbotapi.NewBotAPI("5006119341:AAFFzqR3kQacXOpO7SoQvrzVlln2m9ISiag")
	if err != nil {
		log.Panic(err)
	}

	if client, err = telegraph.Load(savedAccessToken); err == nil {
		log.Printf("> Created client: %#+v", client)

		savedAccessToken = client.AccessToken

		if account, err := client.GetAccountInfo(nil); err == nil {
			log.Printf("> GetAccountInfo result: %#+v", account)
		} else {
			log.Printf("* GetAccountInfo error: %s", err)
		}

		if page, err := client.CreatePageWithHTML("Test page", "Telegraph Test", "", html, true); err == nil {
			log.Printf("> CreatePage: %#+v", page)
			log.Printf("> CreatedPage url: %s", page.URL)
			if page, err := client.GetPage(page.Path, true); err == nil {
				log.Printf("> GetPage result: %#+v", page)
			} else {
				log.Printf("* GetPage error: %s", err)
			}
		}

	} else {
		log.Printf("* CreatePage error: %s", err)
	}

	bot.Debug = true
	log.Printf("Auth on accouns %s", bot.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates := bot.GetUpdatesChan(ucfg)
	for update := range updates {
		if update.Message != nil {
			UserName := update.Message.From.UserName
			ChatId := update.Message.Chat.ID
			Text := update.Message.Text
			Command := update.Message.Command()
			var userFind = false
			for _, u := range users {
				if update.Message.Chat.UserName == u.username {
					userFind = true
					break
				}
			}
			if !userFind {
				log.Print(update.Message.Chat.UserName)
				var usr = user{
					username:      update.Message.Chat.UserName,
					hash:          "",
					walletAddress: "",
					balance:       "200",
				}
				users = append(users, usr)
				rand.Seed(time.Now().UnixNano())
				garage = append(garage, car{
					owner:       usr,
					name:        update.Message.Chat.UserName + "car",
					rarity:      "Common",
					image:       "https://i.imgur.com/IekblLL.png",
					power:       rand.Intn(100-0) + 0,
					ergo:        rand.Intn(100-0) + 0,
					downforce:   rand.Intn(100-0) + 0,
					expiredData: "04.03.2022",
				})

			}
			if len(update.Message.NewChatMembers) != 0 {
				log.Print(update.Message.NewChatMembers)
				for _, v := range update.Message.NewChatMembers {
					var usr = user{
						username:      v.UserName,
						hash:          "",
						walletAddress: "",
						balance:       "200",
					}
					users = append(users, usr)
					rand.Seed(time.Now().UnixNano())
					garage = append(garage, car{
						owner:       usr,
						name:        v.UserName + "car",
						rarity:      "Common",
						image:       "https://i.imgur.com/IekblLL.png",
						power:       rand.Intn(100-0) + 0,
						ergo:        rand.Intn(100-0) + 0,
						downforce:   rand.Intn(100-0) + 0,
						expiredData: "04.03.2022",
					})
				}

			}

			log.Printf("[%s] %d %s", UserName, ChatId, Text)

			var accountInfo = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Подключить кошелёк", "input wallet"),
					tgbotapi.NewInlineKeyboardButtonData("Гараж", "garage"),
					tgbotapi.NewInlineKeyboardButtonData("Как продать момент", "how to sell moment"),
					tgbotapi.NewInlineKeyboardButtonData("Как купить момент на маркетплейсе", "how to buy moment"),
				),
			)

			msg := tgbotapi.NewMessage(ChatId, "Добро пожаловать в NFT игру Royal Races")

			if Command == "start" {
				msg.ReplyMarkup = mainMenuKeyboard
			}

			switch Text {
			case "Информация":
				msg.Text = "Информация о Royal Races"
				msg.ReplyMarkup = infoMenu

			case "Купить машину":
				msg.Text = "Выбирите редкость покупаемой машины"
				msg.ReplyMarkup = buyCar
			case "Мой аккаунт":
				msg.Text = "Водитель: @" + update.Message.Chat.UserName
				msg.ReplyMarkup = accountInfo
			case "Саппорт":
				msg.Text = "Наша поддержка ответит на все вопросы, но перед этим прочтите раздел <<Информация>>"
				msg.ReplyMarkup = supportMenu
			case "Гонка":
				msg.Text = "Информация о NFT BOX"
			case "Маркетплейс":
				msg.Text = "Информация о NFT BOX"
			}
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			for _, car := range garage {
				if "show car "+car.name == msg.Text {
					var carPage = fmt.Sprintf(carPageTemplate, car.name, car.owner.username, car.image, car.rarity, car.power, car.ergo, car.downforce, car.expiredData)
					page, err := client.CreatePageWithHTML("Машина: "+car.name, "Royal Races", "", carPage, true)
					if err == nil {
						msg.Text =  page.URL
					}

				}
			}
			switch msg.Text {
			case "about us":
				msg.Text = "Проект на Finodays, сделанный с любовью))"

			case "how to buy car":
				msg.Text = "Купить машину со случайными показателями можно в разделе маркетплейс, а также участвуя в различных эвентах"
			case "how to sale car":
				msg.Text = "Продать машину можно выставив её на торговую площадку"
			case "race rules":
				msg.Text = "Правила для гонок описаны в разделе <<Гонка>> в описании"
			case "how go to race":
				msg.Text = "Надо зайти в раздел <<Гонка>>, выбрать интересующий тип гонки и начать участвовать"
			case "invalid car":
				msg.Text = "Одной из ключевых особенностей этой игры является срок службы кузова, после истечения которого, машина не сможет участвовать в гонке, а сможет пойти только на детали"
			case "upgrade car":
				msg.Text = "В разделе <<Мой аккаунт>> выбать пункт Гараж, там выбрать машину, которую вы хотите улучшить и далее выбираете машину <<донор>>, которая уничтожится.\n Также надо будет заплатить механикам за улучшение машины. Не забудте про правило улучшения: в зависимости от редкости машины апгрейд работает по разному"
			case "buyCommonCar":
				usr := getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				if usr != -1 {
					i, err := strconv.Atoi(users[usr].balance)
					if err != nil {
						log.Printf("Convert error: %s", err)
					}
					if i > 5 {
						i -= 5
						users[usr].balance = strconv.Itoa(i)
						rand.Seed(time.Now().UnixNano())
						var randCar = car{
							name:        "RX" + strconv.Itoa(rand.Intn(4000-1001)+1001),
							owner:       users[usr],
							rarity:      "Common",
							image:       "https://i.imgur.com/IekblLL.png",
							power:       rand.Intn(400-1) + 1,
							ergo:        rand.Intn(400-1) + 1,
							downforce:   rand.Intn(400-1) + 1,
							expiredData: "04.03.2022",
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							//s := strings.Split(page.URL, "/")
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					}
				} else {
					msg.Text = "Пользователь не найден, обратитесь в поддержку"
				}

			case "buyUncommonCar":
				usr := getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				if usr != -1 {
					i, err := strconv.Atoi(users[usr].balance)
					if err != nil {
						log.Printf("Convert error: %s", err)
					}
					if i > 10 {
						i -= 10
						users[usr].balance = strconv.Itoa(i)
						rand.Seed(time.Now().UnixNano())
						var randCar = car{
							name:        "RX" + strconv.Itoa(rand.Intn(6000-4001)+4001),
							owner:       users[usr],
							rarity:      "Common",
							image:       "https://i.imgur.com/IekblLL.png",
							power:       rand.Intn(600-401) + 401,
							ergo:        rand.Intn(600-401) + 401,
							downforce:   rand.Intn(600-401) + 401,
							expiredData: "04.03.2022",
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					}
				} else {
					msg.Text = "Пользователь не найден, обратитесь в поддержку"
				}
			case "buyRareCar":
				usr := getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				if usr != -1 {

					i, err := strconv.Atoi(users[usr].balance)
					if err != nil {
						log.Printf("Convert error: %s", err)
					}
					if i > 50 {
						i -= 50
						users[usr].balance = strconv.Itoa(i)
						rand.Seed(time.Now().UnixNano())
						var randCar = car{
							name:        "RX" + strconv.Itoa(rand.Intn(8000-6001)+6001),
							owner:       users[usr],
							rarity:      "Common",
							image:       "https://i.imgur.com/IekblLL.png",
							power:       rand.Intn(800-601) + 601,
							ergo:        rand.Intn(800-601) + 601,
							downforce:   rand.Intn(800-601) + 601,
							expiredData: "04.03.2022",
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					}
				} else {
					msg.Text = "Пользователь не найден, обратитесь в поддержку"
				}
			case "buyEpicCar":
				usr := getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				if usr != -1 {
					i, err := strconv.Atoi(users[usr].balance)
					if err != nil {
						log.Printf("Convert error: %s", err)
					}
					if i > 100 {
						i -= 100
						users[usr].balance = strconv.Itoa(i)
						rand.Seed(time.Now().UnixNano())
						var randCar = car{
							name:        "RX" + strconv.Itoa(rand.Intn(9000-8001)+8001),
							owner:       users[usr],
							rarity:      "Common",
							image:       "https://i.imgur.com/IekblLL.png",
							power:       rand.Intn(900-801) + 801,
							ergo:        rand.Intn(900-801) + 801,
							downforce:   rand.Intn(900-801) + 801,
							expiredData: "04.03.2022",
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					}
				} else {
					msg.Text = "Пользователь не найден, обратитесь в поддержку"
				}
			case "buyLegendCar":
				usr := getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				if usr != -1 {
					i, err := strconv.Atoi(users[usr].balance)
					if err != nil {
						log.Printf("Convert error: %s", err)
					}
					if i > 500 {
						i -= 500
						users[usr].balance = strconv.Itoa(i)
						rand.Seed(time.Now().UnixNano())
						var randCar = car{
							name:        "RX" + strconv.Itoa(rand.Intn(9999-9001)+9001),
							owner:       users[usr],
							rarity:      "Common",
							image:       "https://i.imgur.com/IekblLL.png",
							power:       rand.Intn(999-901) + 901,
							ergo:        rand.Intn(999-901) + 901,
							downforce:   rand.Intn(999-901) + 901,
							expiredData: "04.03.2022",
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					}
				} else {
					msg.Text = "Пользователь не найден, обратитесь в поддержку"
				}
			case "garage":
				{
					var garageCarArray = [][]tgbotapi.InlineKeyboardButton{}

					for _, car := range garage {
						if car.owner.username == update.CallbackQuery.Message.Chat.UserName {
							var carBtn = tgbotapi.NewInlineKeyboardButtonData("Машина: "+car.name, "show car "+car.name)
							var carGroupBtn []tgbotapi.InlineKeyboardButton
							carGroupBtn = append(carGroupBtn, carBtn)
							garageCarArray = append(garageCarArray, carGroupBtn)
							log.Print("Error create car page")
						}
					}

					var carGarageQuery = tgbotapi.NewInlineKeyboardMarkup(
						garageCarArray...,
					)
					msg.Text = "Выбирете машину для просмотра"
					msg.ReplyMarkup = carGarageQuery

				}
			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		} else {
			continue
		}
	}
}
