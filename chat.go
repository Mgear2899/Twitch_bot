package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	// twi "github.com/Onestay/go-new-twitch"
	"github.com/PuerkitoBio/goquery"

	twitch "github.com/gempir/go-twitch-irc"
)

var messagesm = make(map[string]int) // var arr = []int{1, 2, 3, 4}
var warning = make(map[string]int)
var startTime = time.Now()

func main() {
	client := twitch.NewClient("mrJohnBot", "oauth:u5sfiw5i6cawt7kgkmwsnrk55c6g9h")

	chaArg := os.Args[1:]

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		countMessages(client, message)
		badWords(client, message)
		go sayTalk(client, message)
		router(client, message)
		shareNetwork(client, message)
		go mesWar(client, message)
		// aiSpeak(message)
	})

	// sub, resub and raids
	client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		subResub(client, message)
	})

	// приветствие зрителя
	client.OnUserJoinMessage(func(message twitch.UserJoinMessage) {
		t := time.Now().Format("15:04:05")
		fmt.Printf("%s = ", t)
		fmt.Println(message.User, "- зашел в чат")
		if message.User == "mrjohnbot" {
			fmt.Println("Yes, se-e-er!")
		} else if message.User == "стример" {
			// client.Say(message.Channel, "Тебя приветствует mr. John, "+message.User+", я слежу за порядком в чатике!!!")
			client.Say(message.Channel, "Здравствуй, "+message.User+", я mr. John. И умею следить за порядком в чатике!")
		}
		// Позвольте вас поприветствовать и представиться users.
		// Я mister John и в мои обязанности входит следить за порядком в чате

		// запись статистики в файл
		read, err := os.OpenFile("Отчет по трансляции.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // да здесь много непонятных букв.
		// os.O_APPEND|os.O_CREATE|os.O_WRONLY, - это раздешения на файл, скоторыми он пытается создать его
		// os.O_APPEND - запись в конец файла, os.O_CREATE - либо создание нового файла, os.O_WRONLY - либо запись в пустой существующий файл

		if err != nil {
			log.Panic(err)
		}
		defer read.Close()

		// вот это будет каждый раз записывать строчку в конец файла!
		if _, err := read.Write([]byte(t + " - " + message.User + "\n ")); err != nil {
			log.Panic(err)
		}

	})

	// client.OnConnect(func() {
	// 	go mesWar(client)
	// })

	client.Join(chaArg[0])

	err := client.Connect()
	if err != nil {
		log.Panic(err)
	}
}

func badWords(client *twitch.Client, message twitch.PrivateMessage) {
	// ПРОВЕРКА НА ПЛОХИЕ СЛОВА!!!
	rer, err := ioutil.ReadFile("word.bad")
	if err != nil {
		log.Panicln(err)
	}

	// убираем запятые
	str := string(rer)
	spl := strings.Split(str, ", ")

	// разбор предложений из чата
	re := regexp.MustCompile(`[a-zA-Zа-яА-Я0-9]+`)
	match := re.FindAllString(message.Message, -1) //

	// проверка на матерные слова
mm:
	for _, value := range match {
		for _, strSplit := range spl {
			if value == strSplit {
				warning[message.User.Name]++
				fmt.Println("warnings -", warning)
				if warning[message.User.Name] == 1 {
					client.Say(message.Channel, "Mr. John не любит плохие слова! :-(")
				} else if warning[message.User.Name] == 2 {
					client.Say(message.Channel, message.User.Name+" не ругайся, 30 сек держи!")
					client.Say(message.Channel, "/timeout "+message.User.Name+" 30")
				} else if warning[message.User.Name] == 3 {
					client.Say(message.Channel, message.User.Name+" не ругайся, 300 сек в студию!")
					client.Say(message.Channel, "/timeout "+message.User.Name+" 300")
				} else {
					client.Say(message.Channel, message.User.Name+", твои дни сочтены!!!")
					// 259200 -
					client.Say(message.Channel, "/ban "+message.User.Name)
				}

				break mm
				// https://help.twitch.tv/s/article/chat-commands?language=en_US
				// /timeout night_delirium 30
			}
		}
	}
}

// подсчет сообщений
func countMessages(client *twitch.Client, message twitch.PrivateMessage) {
	messagesm[message.User.Name]++
	fmt.Println(messagesm)
	if message.Message == "!use" {
		var answer string
		users, err := client.Userlist(message.Channel) // client.Userlist(channel) - список зрителей
		if err != nil {
			log.Panic(err)
		}
		for _, user := range users {
			answer = answer + user + " - " + strconv.Itoa(messagesm[message.User.Name]) + " сообщений за стрим\n"
		}
		client.Say(message.Channel, answer)
	}
}

// функция вывода рандомных сообщений
func mesWar(client *twitch.Client, message twitch.PrivateMessage) {
	// каждые 15 минут выводит сообщение
	ticker := time.NewTicker(time.Minute * 20)

	// карта с фразами
	randomMes := [...]string{
		"Подписаться на канал это пол дела, надо сидеть на стриме! (c)",
		"— Это мой первый бой. Что мне делать? — Не умирать.",
		"Война. Война никогда не меняется.",
		"Я не спасся. Меня убили… Обожаю эту шутку.",
		"Эта драка бессмысленна. Как и твоё сопротивление.",
		"BloodTrail",
		"Моя жизнь, это то во что ты её превратил...",
	}

	// рандомный вывод сообщений
	for range ticker.C {
		leng := len(randomMes)
		randomNumber := rand.Intn(leng)
		randomPhrase := randomMes[randomNumber]
		client.Say(message.Channel, randomPhrase)
	}
}

func subResub(client *twitch.Client, message twitch.UserNoticeMessage) {
	client.Say(message.Channel, message.User.Name+" подписался! TwitchVotes "+message.SystemMsg)
	fmt.Println(message.MsgID, message.MsgParams, message.SystemMsg, message.Tags)
}

// обращение, команды к боту
func sayTalk(client *twitch.Client, message twitch.PrivateMessage) {
	// обращение к мистеру Джону
	nameJohn := [...]string{
		"John", "Джон", "mr. John", "mr. Jon", "mrJohnBot", "mrjohnbot", "jon", "john",
	}

	reg := regexp.MustCompile(`[a-zA-Zа-яА-Я]+`)
	saySay := reg.FindAllString(message.Message, -1)

says:
	for _, values := range saySay {
		for _, findJohn := range nameJohn {
			if values == findJohn {
				client.Say(message.Channel, ", да, сэ-э-эр!")
				break says
			}
		}
	}
}

func router(client *twitch.Client, message twitch.PrivateMessage) {
	switch message.Message {
	case "!Время":
		client.Say(message.Channel, "С начала срима прошло - "+printElapsedTime())
		break
		// ссылка в телеграм
	case "!tel":
		client.Say(message.Channel, message.User.Name+" https://t.me/joinchat/AAAAAFGAVk9hZ7vAci-mNQ")
		break
	}
}

// вывод вркмени в чат с начала трансляции отсчет
const (
	pars = "72h3m0.5s"
	now  = "15:04:05"
)

func printElapsedTime() string {
	elapsed := time.Since(startTime)
	// elapsed = elapsed.Round(time.Minute)
	// h := elapsed / time.Hour
	// elapsed -= h * time.Hour
	// m := elapsed / time.Minute
	spl := strings.Replace(elapsed.String(), "h", ":", 2)

	return fmt.Sprintf("%s", spl)
}

func shareNetwork(client *twitch.Client, message twitch.PrivateMessage) {
	reg := regexp.MustCompile(`Поиск\/`)
	saySay := reg.FindAllString(message.Message, -1)

share:
	for _, values := range saySay {
		if values == "Поиск/" {
			regsha := regexp.MustCompile(`\/.+`)

			sayShare := regsha.FindAllString(message.Message, -1)
			for _, val := range sayShare {
				sha, err := http.Get("https://ru.wikipedia.org/wiki" + val)
				if err != nil {
					fmt.Println(err)
				}
				defer sha.Body.Close()

				pars, err := goquery.NewDocumentFromReader(sha.Body)
				if err != nil {
					log.Fatal(err)
				}
				client.Say(message.Channel, pars.Find(".mw-parser-output p").First().Text())
			}
			break share
		}

	}
}
