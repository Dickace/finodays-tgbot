package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	telegraph "github.com/meinside/telegraph-go"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	telegramBotApiToken string
}

//var html = `<h3 >Машина %s водителя %s</h3>
//    <img style="width: 70%%; " src="https://i.imgur.com/IekblLL.png"/>
//    <p>Редкость: <span style="color:#1E90FF;">%s</span></p>
//            <p>Мощность мотора: %d<br/>
//Аэродинамика: %d<br/>
//Прижимная сила: %d<br/>
//Срок действия до %s<br/>
//</p>
//
//`

var raceChoise = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Один забирает всё", "one winner"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Гонка на интерес", "training"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Рейтинговая гонка", "rating race"),
	),
)

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

var marketplaceMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Купить машину", "buy cars"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Мои машины на продаже", "my cars on sales"),
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
		tgbotapi.NewInlineKeyboardButtonData("Common машина - 5 RRTC", "buyCommonCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Uncommon машина - 10 RRTC", "buyUncommonCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Rare машина - 50 RRTC", "buyRareCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Epic машина - 100 RRTC", "buyEpicCar"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Legendary машина - 500 RRTC", "buyLegendCar"),
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
var imagePool = []string{
	"https://i.imgur.com/RE8ytlh.png",
	"https://i.imgur.com/S4aEIUl.png",
	"https://i.imgur.com/Yr51R0q.jpg",
	"https://i.imgur.com/O4eJEP9.png",
	"https://i.imgur.com/4Hcjkqn.png",
	"https://i.imgur.com/eu74VGK.png",
	"https://i.imgur.com/GTSorY8.png",
	"https://i.imgur.com/ZJOfsRc.png",
	"https://i.imgur.com/YJwRvhr.png",
	"https://i.imgur.com/IekblLL.png",
	"https://i.imgur.com/IekblLL.png",
}
var marketplace = []carSale{
	{salingCar: car{
		owner:       users[0],
		name:        "RX 9004",
		rarity:      "Common",
		image:       "https://i.imgur.com/IekblLL.png",
		power:       64,
		ergo:        34,
		downforce:   78,
		expiredData: "04.03.2022",
		isBusy:      false,
	},
		price: "120",
	},
}

var trainingLobbys = []TrainingLobby{}

type TrainingLobby struct {
	gameId     string
	isOver     bool
	players    []user
	cars       []car
	maxPlayers int
	chances    []float64
	result     []string
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
	isBusy      bool
}

type carSale struct {
	salingCar car
	price     string
}

type user struct {
	username       string
	hash           string
	walletAddress  string
	walletSeed     string
	balance        string
	inputSalePrice bool
	usingCar       string
	carForRace     string
	chatId         int64
}

var users = []user{
	{
		username:       "Hamilton",
		hash:           "",
		walletAddress:  "xasfwegfwegwegweg",
		walletSeed:     "awfawgresgseg",
		balance:        "200",
		inputSalePrice: false,
		usingCar:       "",
		carForRace:     "",
		chatId:         0,
	},
}

func getUserIndex(username string) int {
	for i, v := range users {
		if v.username == username {
			return i
		}
	}
	return -1
}
func getCarIndexFromGarage(carName string) int {
	for i, v := range garage {
		if v.name == carName {
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

func game(bot tgbotapi.BotAPI) {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				for ti, trainingLobby := range trainingLobbys {
					if len(trainingLobby.players) == trainingLobby.maxPlayers && !trainingLobby.isOver {
						for _, usr := range trainingLobby.players {
							var msg = tgbotapi.NewMessage(usr.chatId, fmt.Sprintf("Тренировочная гонка №%s началась!", trainingLobby.gameId))
							_, err := bot.Send(msg)
							if err != nil {
								fmt.Print(err)
							}
						}
						trainingLobbys[ti] = trainingLobby
						var wins []car
						for count := 0; count < len(trainingLobby.cars); count++ {
							var statArray []int
							var statSum = 0
							var chansDyn []float64
							for _, crs := range trainingLobby.cars {
								var carFind = false
								for _, wincar := range wins {
									if wincar.name != crs.name {
										carFind = true
									}
								}
								if !carFind {
									statArray = append(statArray, crs.power+crs.ergo+crs.downforce)
									statSum += crs.power + crs.ergo + crs.downforce
								}
							}
							for _, stat := range statArray {
								chansDyn = append(chansDyn, float64(stat)*100/float64(statSum))
							}
							if count == 0 {
								trainingLobby.chances = chansDyn
							}
							rand.Seed(time.Now().UnixNano())
							var winnerTicket = rand.Float64() * 100
							var sumTotal = 0.0
							for j, chanse := range chansDyn {
								if sumTotal < winnerTicket && winnerTicket < sumTotal+chanse {
									wins = append(wins, trainingLobby.cars[j])
									break
								} else {
									sumTotal += chanse
								}
							}
						}
						var msgText = fmt.Sprintf("Тренировочная гонка №%s закончилась\nРезультаты следующие:\n", trainingLobby.gameId)
						for p, win := range wins {
							var str = fmt.Sprintf("%d место - @%s, на машине %s", p+1, win.owner.username, win.name)
							trainingLobby.result = append(trainingLobby.result, str)
							msgText += str
						}
						trainingLobby.isOver = true
						for _, crs := range trainingLobby.cars {
							crs.isBusy = false
						}
						for _, usr := range trainingLobby.players {
							var msg = tgbotapi.NewMessage(usr.chatId, msgText)
							_, err := bot.Send(msg)
							if err != nil {
								fmt.Print(err)
							}
						}
						trainingLobbys[ti] = trainingLobby
					}
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	telegraph.Verbose = true

	content, err := ioutil.ReadFile("../config.json")
	if err != nil {
		fmt.Print(err)
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Print(err)
	}
	bot, err := tgbotapi.NewBotAPI("5061864092:AAFkaxtv-zwDKyPv5xL9bW0PtPjPdo0helM")
	if err != nil {
		log.Panic(err)
	}

	if client, err = telegraph.Load("5061864092:AAFkaxtv-zwDKyPv5xL9bW0PtPjPdo0helM"); err == nil {
		log.Printf("> Created client: %#+v", client)

		config.telegramBotApiToken = client.AccessToken

		if account, err := client.GetAccountInfo(nil); err == nil {
			log.Printf("> GetAccountInfo result: %#+v", account)
		} else {
			log.Printf("* GetAccountInfo error: %s", err)
		}

		//if page, err := client.CreatePageWithHTML("Test page", "Telegraph Test", "", html, true); err == nil {
		//	log.Printf("> CreatePage: %#+v", page)
		//	log.Printf("> CreatedPage url: %s", page.URL)
		//	if page, err := client.GetPage(page.Path, true); err == nil {
		//		log.Printf("> GetPage result: %#+v", page)
		//	} else {
		//		log.Printf("* GetPage error: %s", err)
		//	}
		//}

	} else {
		log.Printf("* CreatePage error: %s", err)
	}

	bot.Debug = true
	log.Printf("Auth on accouns %s", bot.Self.UserName)
	go game(*bot)
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
					username:       update.Message.Chat.UserName,
					hash:           "",
					walletAddress:  "",
					balance:        "200",
					inputSalePrice: false,
					usingCar:       "",
					chatId:         update.Message.Chat.ID,
				}
				users = append(users, usr)
				rand.Seed(time.Now().UnixNano())
				garage = append(garage, car{
					owner:       usr,
					name:        update.Message.Chat.UserName + "car",
					rarity:      "Common",
					image:       imagePool[rand.Intn(len(imagePool)-0)+0],
					power:       rand.Intn(100-0) + 0,
					ergo:        rand.Intn(100-0) + 0,
					downforce:   rand.Intn(100-0) + 0,
					expiredData: "04.03.2022",
					isBusy:      false,
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
						image:       imagePool[rand.Intn(len(imagePool)-0)+0],
						power:       rand.Intn(100-0) + 0,
						ergo:        rand.Intn(100-0) + 0,
						downforce:   rand.Intn(100-0) + 0,
						expiredData: "04.03.2022",
						isBusy:      false,
					})
				}

			}

			log.Printf("[%s] %d %s", UserName, ChatId, Text)

			msg := tgbotapi.NewMessage(ChatId, "Добро пожаловать в NFT игру Royal Races")
			var usr = getUserIndex(update.Message.Chat.UserName)
			if users[usr].inputSalePrice == true {
				var value, err = strconv.Atoi(update.Message.Text)
				if err != nil {
					msg.Text = "Некорректный ввод цены"
				} else {
					var salecar car
					for i, v := range garage {
						if v.name == users[usr].usingCar {
							salecar = v
							garage = append(garage[:i], garage[i+1:]...)
							break
						}
					}
					marketplace = append(marketplace, carSale{
						salingCar: salecar,
						price:     strconv.Itoa(value),
					})

					msg.Text = fmt.Sprintf(`%s выставлена за %d RRCT`, salecar.name, value)
					msg.ReplyMarkup = mainMenuKeyboard
				}
				users[usr].inputSalePrice = false
			}
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
				var usr = getUserIndex(update.Message.Chat.UserName)
				if users[usr].walletAddress != "" {
					var msgText = fmt.Sprintf("Водитель: @%s\nАдресс кошелька: %s\nБаланс: %s RRTC", users[usr].username, users[usr].walletAddress, users[usr].balance)
					msg.Text = msgText
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Пополнить баланс", "payment"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Подключить кошельек", "payment"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Гараж", "garage"),
						),
					)
				} else {
					var msgText = fmt.Sprintf("Водитель: @%s\nБаланс: %s RRTC", users[usr].username, users[usr].balance)
					msg.Text = msgText
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Пополнить баланс", "payment"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Изменить адресс кошелька", "payment"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Гараж", "garage"),
						),
					)
				}

			case "Саппорт":
				msg.Text = "Наша поддержка ответит на все вопросы, но перед этим прочтите раздел <<Информация>>"
				msg.ReplyMarkup = supportMenu
			case "Гонка":
				msg.Text = "Выбирите гонку в которой хотите участвовать"
				msg.ReplyMarkup = raceChoise
			case "Маркетплейс":
				msg.Text = "На маркетплейсы вы можете купить NFT-машины от других игроков"
				msg.ReplyMarkup = marketplaceMenu
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

			for ti, trainLobby := range trainingLobbys {

				if "check lobby "+trainLobby.gameId == msg.Text {
					var textMsg = fmt.Sprintf("Участники тренировочной гонки №%s:\n", trainLobby.gameId)
					for _, car := range trainLobby.cars {
						textMsg += fmt.Sprintf("@%s на машине %s\nРедкость: %s \nСтаты: \nP:%d E:%d D:%d\n\n", car.owner.username, car.name, car.rarity, car.power, car.ergo, car.downforce)
					}
					msg.Text = textMsg
					if !trainLobby.isOver {
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Выйти из лобби", "left train race "+trainLobby.gameId)))
					}
				} else if "left train race "+trainLobby.gameId == msg.Text {
					for i, cr := range trainLobby.cars {
						if cr.owner.username == update.CallbackQuery.Message.Chat.UserName {
							var crs = getCarIndexFromGarage(cr.name)
							garage[crs].isBusy = false
							trainLobby.cars = append(trainLobby.cars[:i], trainLobby.cars[i+1:]...)
						}
					}
					for j, pls := range trainLobby.players {
						if pls.username == update.CallbackQuery.Message.Chat.UserName {
							trainLobby.players = append(trainLobby.players[:j], trainLobby.players[j+1:]...)
						}
					}
					trainingLobbys[ti] = trainLobby
					msg.Text = "Вы успешно покинули гонку"
				}
			}

			for _, car := range garage {
				if "show car "+car.name == msg.Text {
					//var carPage = fmt.Sprintf(carPageTemplate, car.name, car.owner.username, car.image, car.rarity, car.power, car.ergo, car.downforce, car.expiredData)
					//page, err := client.CreatePageWithHTML("Машина: "+car.name, "Royal Races", "", carPage, true)
					if err == nil {
						var msgText = fmt.Sprintf("<a href='%s'>%s</a>\nРедкость:%s\nСтаты: \nМощность - %d\nАэродинамика - %d,\nПрижимная сила - %d", car.image, car.name, car.rarity, car.power, car.ergo, car.downforce)

						msg.Text = msgText
						msg.ParseMode = tgbotapi.ModeHTML
						var carAction = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Выставить машину на маркетплейс", "sell car "+car.name)),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Улучшить машину", "upgrade car "+car.name)),
						)
						msg.ReplyMarkup = carAction
					}
				} else if "sell car "+car.name == msg.Text {
					var usr = getUserIndex(car.owner.username)
					users[usr].usingCar = car.name
					users[usr].inputSalePrice = true
					msg.Text = "Введите цену продажи"
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				} else if "choose training "+car.name == msg.Text {
					var usr = getUserIndex(update.CallbackQuery.Message.Chat.UserName)
					users[usr].carForRace = car.name
					msg.Text = car.name + " выбрана"
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Найти тренировочную игру", "find training game"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Выбрать другую машину", "choose car for training"),
						),
					)
				}
				//else if "upgrade car "+car.name == msg.Text {
				//
				//}
			}

			for _, saleCar := range marketplace {
				if "buy market "+saleCar.salingCar.name == msg.Text {
					var msgText = fmt.Sprintf("<a href='%s'>%s</a>\nРедкость:%s\nСтаты: \nМощность - %d\nАэродинамика - %d,\nПрижимная сила - %d\nЦена:%s RRCT", saleCar.salingCar.image, saleCar.salingCar.name, saleCar.salingCar.rarity, saleCar.salingCar.power, saleCar.salingCar.ergo, saleCar.salingCar.downforce, saleCar.price)
					msg.Text = msgText
					msg.ParseMode = tgbotapi.ModeHTML
					var carAction = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Купить машину", "buy marketcar "+saleCar.salingCar.name)),
					)
					msg.ReplyMarkup = carAction

				} else if "buy marketcar "+saleCar.salingCar.name == msg.Text {
					var usr = getUserIndex(update.CallbackQuery.Message.Chat.UserName)
					buyerBalance, err := strconv.Atoi(users[usr].balance)
					if err != nil {
						fmt.Print(err)
					}
					price, err := strconv.Atoi(saleCar.price)
					if err != nil {
						fmt.Print(err)
					}
					if buyerBalance > price {
						var recipientUsr = getUserIndex(saleCar.salingCar.owner.username)
						sellerBalance, err := strconv.Atoi(users[recipientUsr].balance)
						if err != nil {
							fmt.Print(err)
						}
						newSellerBalance := strconv.Itoa(sellerBalance + price)
						users[recipientUsr].balance = newSellerBalance
						newBuyerBalance := strconv.Itoa(buyerBalance - price)
						users[usr].balance = newBuyerBalance
						var sellsCar = car{
							name:        saleCar.salingCar.name,
							rarity:      saleCar.salingCar.rarity,
							owner:       users[usr],
							power:       saleCar.salingCar.power,
							ergo:        saleCar.salingCar.ergo,
							downforce:   saleCar.salingCar.downforce,
							expiredData: saleCar.salingCar.expiredData,
							isBusy:      saleCar.salingCar.isBusy,
						}
						garage = append(garage, sellsCar)
						for i, v := range marketplace {
							if v.salingCar.name == sellsCar.name {
								marketplace = append(marketplace[:i], marketplace[i+1:]...)
								break
							}
						}
						msg.Text = fmt.Sprintf("Вы куплили %s, проверьте её в гараже", sellsCar.name)
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("В маркетплейс", "marketplaceMenu"),
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("В гараж", "garage"),
							),
						)
					} else {
						msg.Text = "Недостаточно средств для покупки"
					}

				} else if "upgrade car "+saleCar.salingCar.name == msg.Text {

				} else if "show market "+saleCar.salingCar.name == msg.Text {
					var msgText = fmt.Sprintf("<a href='%s'>%s</a>\nРедкость:%s\nСтаты: \nМощность - %d\nАэродинамика - %d,\nПрижимная сила - %d\nЦена:%s RRCT", saleCar.salingCar.image, saleCar.salingCar.name, saleCar.salingCar.rarity, saleCar.salingCar.power, saleCar.salingCar.ergo, saleCar.salingCar.downforce, saleCar.price)
					msg.Text = msgText
					msg.ParseMode = tgbotapi.ModeHTML
					var carAction = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Снять с продажи", "cancel sell "+saleCar.salingCar.name)),
					)
					msg.ReplyMarkup = carAction

				} else if "cancel sell "+saleCar.salingCar.name == msg.Text {
					garage = append(garage, saleCar.salingCar)
					for i, v := range marketplace {
						if v.salingCar.name == saleCar.salingCar.name {
							marketplace = append(marketplace[:i], marketplace[i+1:]...)
							break
						}
					}
					msg.Text = "Машина снята с продажи"

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
				msg.Text = "В разделе <<Мой аккаунт>> выбать пункт Гараж, там выбрать машину, которую вы хотите улучшить и далее выбираете машину <<донор>>, которая уничтожится.\nТакже надо будет заплатить механикам за улучшение машины. Не забудте про правило улучшения: в зависимости от редкости машины апгрейд работает по разному"
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
							image:       imagePool[rand.Intn(len(imagePool)-0)+0],
							power:       rand.Intn(400-1) + 1,
							ergo:        rand.Intn(400-1) + 1,
							downforce:   rand.Intn(400-1) + 1,
							expiredData: "04.03.2022",
							isBusy:      false,
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							//s := strings.Split(page.URL, "/")
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					} else {
						msg.Text = "Недостаточно средств для покупки машины"
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Пополнить", "payment")))
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
							rarity:      "Uncommon",
							image:       imagePool[rand.Intn(len(imagePool)-0)+0],
							power:       rand.Intn(600-401) + 401,
							ergo:        rand.Intn(600-401) + 401,
							downforce:   rand.Intn(600-401) + 401,
							expiredData: "04.03.2022",
							isBusy:      false,
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					} else {
						msg.Text = "Недостаточно средств для покупки машины"
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Пополнить", "payment")))
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
							rarity:      "Rare",
							image:       imagePool[rand.Intn(len(imagePool)-0)+0],
							power:       rand.Intn(800-601) + 601,
							ergo:        rand.Intn(800-601) + 601,
							downforce:   rand.Intn(800-601) + 601,
							expiredData: "04.03.2022",
							isBusy:      false,
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					} else {
						msg.Text = "Недостаточно средств для покупки машины"
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Пополнить", "payment")))
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
							rarity:      "Epic",
							image:       imagePool[rand.Intn(len(imagePool)-0)+0],
							power:       rand.Intn(900-801) + 801,
							ergo:        rand.Intn(900-801) + 801,
							downforce:   rand.Intn(900-801) + 801,
							expiredData: "04.03.2022",
							isBusy:      false,
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					} else {
						msg.Text = "Недостаточно средств для покупки машины"
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Пополнить", "payment")))
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
							rarity:      "Legendary",
							image:       imagePool[rand.Intn(len(imagePool)-0)+0],
							power:       rand.Intn(999-901) + 901,
							ergo:        rand.Intn(999-901) + 901,
							downforce:   rand.Intn(999-901) + 901,
							expiredData: "04.03.2022",
							isBusy:      false,
						}
						garage = append(garage, randCar)
						var carPage = fmt.Sprintf(carPageTemplate, randCar.name, randCar.owner.username, randCar.image, randCar.rarity, randCar.power, randCar.ergo, randCar.downforce, randCar.expiredData)
						page, err := client.CreatePageWithHTML("Машина: "+randCar.name, "Royal Races", "", carPage, true)
						if err == nil {
							msg.Text = "Вы успешно купили машину! \n" + page.URL
						}
					} else {
						msg.Text = "Недостаточно средств для покупки машины"
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Пополнить", "payment")))
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
					if len(garageCarArray) != 0 {
						var carGarageQuery = tgbotapi.NewInlineKeyboardMarkup(
							garageCarArray...,
						)
						msg.Text = "Выбирете машину для просмотра"
						msg.ReplyMarkup = carGarageQuery
					} else {
						msg.Text = "В вашем гараже нет машин"
					}

				}
			case "buy cars":
				var buyCarsArray = [][]tgbotapi.InlineKeyboardButton{}
				for _, saleCar := range marketplace {
					if saleCar.salingCar.owner.username != update.CallbackQuery.Message.Chat.UserName {
						var carBtn = tgbotapi.NewInlineKeyboardButtonData("Машина: "+saleCar.salingCar.name, "buy market "+saleCar.salingCar.name)
						var carGroupBtn []tgbotapi.InlineKeyboardButton
						carGroupBtn = append(carGroupBtn, carBtn)
						buyCarsArray = append(buyCarsArray, carGroupBtn)
					}
				}
				if len(buyCarsArray) != 0 {
					msg.Text = "Машины на продаже в данный момент"
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buyCarsArray...)
				} else {
					msg.Text = "На данный момент в машины в продаже отсутствуют"
				}
			case "my cars on sales":
				var mySellCarsArray = [][]tgbotapi.InlineKeyboardButton{}
				for _, mySaleCar := range marketplace {
					if mySaleCar.salingCar.owner.username == update.CallbackQuery.Message.Chat.UserName {
						var carBtn = tgbotapi.NewInlineKeyboardButtonData("Машина: "+mySaleCar.salingCar.name, "show market "+mySaleCar.salingCar.name)
						var carGroupBtn []tgbotapi.InlineKeyboardButton
						carGroupBtn = append(carGroupBtn, carBtn)
						mySellCarsArray = append(mySellCarsArray, carGroupBtn)
					}
				}
				if len(mySellCarsArray) != 0 {
					msg.Text = "Мои машины, выставленные на продажу"
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(mySellCarsArray...)
				} else {
					msg.Text = "Ни одна моя машина не выставлена"
				}
			case "training":
				msg.Text = "Тренировочная гонка до 4 учстников, которые играют на интерес. Нет ни ставок, ни рейтинга"
				var usr = getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				if users[usr].carForRace != "" {
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Найти игру", "find training game"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Выбрать другую машину", "choose car for training"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Текущие игры", "current training game"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Прошедшие игры", "last training game"),
						))
				} else {
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Выбрать машину", "choose car for training"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Текущие игры", "current training game"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Прошедшие игры", "last training game"),
						))
				}

			case "choose car for training":
				var garageCarArray = [][]tgbotapi.InlineKeyboardButton{}

				for _, car := range garage {
					if car.owner.username == update.CallbackQuery.Message.Chat.UserName {
						var carText = fmt.Sprintf("%s P:%d E:%d D%d", car.name, car.power, car.ergo, car.downforce)
						var carBtn = tgbotapi.NewInlineKeyboardButtonData(carText, "choose training "+car.name)
						var carGroupBtn []tgbotapi.InlineKeyboardButton
						carGroupBtn = append(carGroupBtn, carBtn)
						garageCarArray = append(garageCarArray, carGroupBtn)

					}
				}
				if len(garageCarArray) != 0 {
					var carGarageQuery = tgbotapi.NewInlineKeyboardMarkup(
						garageCarArray...,
					)
					msg.Text = "Выбирете машину для участия в тренировке"
					msg.ReplyMarkup = carGarageQuery
				} else {
					msg.Text = "В вашем гараже нет машин"
				}
			case "last training game":
				var lobbyArray = [][]tgbotapi.InlineKeyboardButton{}
				for _, lobby := range trainingLobbys {
					for _, usr := range lobby.players {
						if usr.username == update.CallbackQuery.Message.Chat.UserName && lobby.isOver {
							var lobbyBtn = tgbotapi.NewInlineKeyboardButtonData("Посмотреть прошедшую игру №"+lobby.gameId, "check lobby "+lobby.gameId)
							var lobbyGroupBtn []tgbotapi.InlineKeyboardButton
							lobbyGroupBtn = append(lobbyGroupBtn, lobbyBtn)
							lobbyArray = append(lobbyArray, lobbyGroupBtn)
						}
					}
				}
				if len(lobbyArray) != 0 {
					var lobbyQuery = tgbotapi.NewInlineKeyboardMarkup(
						lobbyArray...,
					)
					msg.Text = "Ваши прошлые тренировочные гонки"
					msg.ReplyMarkup = lobbyQuery
				} else {
					msg.Text = "У вас ещё не было тренировочных гонок"
				}
			case "current training game":
				var lobbyArray = [][]tgbotapi.InlineKeyboardButton{}
				for _, lobby := range trainingLobbys {
					for _, usr := range lobby.players {
						if usr.username == update.CallbackQuery.Message.Chat.UserName && !lobby.isOver {
							var lobbyBtn = tgbotapi.NewInlineKeyboardButtonData("Посмотреть лобби №"+lobby.gameId, "check lobby "+lobby.gameId)
							var lobbyGroupBtn []tgbotapi.InlineKeyboardButton
							lobbyGroupBtn = append(lobbyGroupBtn, lobbyBtn)
							lobbyArray = append(lobbyArray, lobbyGroupBtn)
						}
					}
				}
				if len(lobbyArray) != 0 {
					var lobbyQuery = tgbotapi.NewInlineKeyboardMarkup(
						lobbyArray...,
					)
					msg.Text = "Ваши текущие тренировочные гонки"
					msg.ReplyMarkup = lobbyQuery
				} else {
					msg.Text = "Вы пока не участвуете ни в одной из гонок"
				}
			case "find training game":
				var usr = getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				var cr = getCarIndexFromGarage(users[usr].carForRace)
				if !garage[cr].isBusy {
					if len(trainingLobbys) > 0 {
						var lobbyFind = false
						for ti, lobby := range trainingLobbys {
							if !lobby.isOver && len(lobby.players) < 4 {
								var userFind = false
								for _, p := range lobby.players {
									if p.username == update.CallbackQuery.Message.Chat.UserName {
										userFind = true
										break
									}
								}
								if !userFind {
									lobby.players = append(lobby.players, users[usr])
									lobby.cars = append(lobby.cars, garage[cr])
									garage[cr].isBusy = true
									msg.Text = "Игра найдена - игроков в лобби " + strconv.Itoa(len(lobby.players)) + "/" + strconv.Itoa(lobby.maxPlayers)
									msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("Посмотреть лобби", "check lobby "+lobby.gameId),
									))
									trainingLobbys[ti] = lobby
									lobbyFind = true
									break
								} else {
									continue
								}
							} else {
								continue
							}
						}
						if !lobbyFind {
							newGameId, err := strconv.Atoi(trainingLobbys[len(trainingLobbys)-1].gameId)
							if err != nil {
								fmt.Print(err)
							}
							newGameId += 1
							trainingLobbys = append(trainingLobbys, TrainingLobby{
								gameId:     strconv.Itoa(newGameId),
								isOver:     false,
								players:    []user{users[usr]},
								cars:       []car{garage[cr]},
								maxPlayers: 4,
							})
							garage[cr].isBusy = true
							msg.Text = "Игра найдена - игроков в лобби " + strconv.Itoa(len(trainingLobbys[newGameId-1].players)) + "/" + strconv.Itoa(trainingLobbys[newGameId-1].maxPlayers)
							msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Посмотреть лобби", "check lobby "+strconv.Itoa(newGameId)),
							))
						}
					} else {
						trainingLobbys = []TrainingLobby{
							{gameId: "1",
								isOver:     false,
								players:    []user{users[usr]},
								cars:       []car{garage[cr]},
								maxPlayers: 4},
						}
						garage[cr].isBusy = true
						msg.Text = "Игра найдена - игроков в лобби " + strconv.Itoa(len(trainingLobbys[0].players)) + "/" + strconv.Itoa(trainingLobbys[0].maxPlayers)
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Посмотреть лобби", "check lobby 1"),
						))
					}
				} else {
					msg.Text = "Выбранная машина уже участвует в гонке, выберите другую "
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Выбрать другую машину", "choose car for training"),
					))
				}
			case "payment":
				msg.Text = "Получите бесплатную 1000 RRCT в рамках работы прототипа"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Получить 1000 RRCT", "bonus"),
					))
			case "bonus":
				var usr = getUserIndex(update.CallbackQuery.Message.Chat.UserName)
				oldBalance, err := strconv.Atoi(users[usr].balance)
				if err != nil {
					fmt.Print(err)
				}
				users[usr].balance = strconv.Itoa(oldBalance + 1000)
				msg.Text = "Ваши 1000 RRCT уже на вашем балансе, удачных гонок!"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("В магазин машин", "buyCar"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("В маркеталейс", "marketplaceMenu"),
					),
				)

			}

			var phrases = strings.Split(msg.Text, " ")
			if len(phrases) > 2 && phrases[0] == "sell" && phrases[1] == "car" {
				break
			} else {
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}

		} else {
			continue
		}
	}
}
