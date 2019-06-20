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

	twitch "github.com/gempir/go-twitch-irc"
)

const (
	cahnnel = "monstrum_gear"
)

var messagesm = make(map[string]int) // var arr = []int{1, 2, 3, 4}
var warning = make(map[string]int)
var timers = make(map[int]int)

func main() {
	go htmlHH()

	client := twitch.NewClient("mrJohnBot", "oauth:nwaoopj79z91twfuts32tbnm4pe5d7")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		countMessages(client, message)
		badWords(client, message)
		sayTalk(client, message)
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
		} else {
			client.Say(message.Channel, "Тебя приветствует mr. John, "+message.User+", я слежу за порядком в чатике!!!")
		}
		// Позвольте вас поприветствовать и представиться users.
		// Я mister John и в мои обязанности входит следить за порядком в чате

		// запись статистики в файл
		readText := []byte(message.User)
		read, err := os.OpenFile("Отчет по трансляции.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // да здесь много непонятных букв.
		// os.O_APPEND|os.O_CREATE|os.O_WRONLY, - это раздешения на файл, скоторыми он пытается создать его
		// os.O_APPEND - запись в конец файла, os.O_CREATE - либо создание нового файла, os.O_WRONLY - либо запись в пустой существующий файл

		if err != nil {
			log.Panic(err)
		}
		defer read.Close()

		// вот это будет каждый раз записывать строчку в конец файла!
		// по идее должно работать
		if _, err := read.Write([]byte(", " + message.User + ", ")); err != nil {
			log.Panic(err)
		}
		read.Write(readText)
	})

	client.OnConnect(func() {
		go mesWar(client)
	})

	// for _, oneChanel := range channels {
	client.Join(cahnnel)
	// }

	err := client.Connect()
	if err != nil {
		log.Panic(err)
	}
}

func badWords(client *twitch.Client, message twitch.PrivateMessage) {
	// ПРОВЕРКА НА ПЛОХИЕ СЛОВА!!!
	rer, err := ioutil.ReadFile("mat.txt")
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
	if message.Message == "use!" {
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
func mesWar(client *twitch.Client) {
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
		"Squid1 Squid2 Squid3 Squid4 - Я знал что он водится!!!",
		"mr. John еще сдесь :)",
		"Моя жизнь, это то во что ты её превратил...",
	}

	// рандомный вывод сообщений
	for range ticker.C {
		leng := len(randomMes)
		randomNumber := rand.Intn(leng)
		randomPhrase := randomMes[randomNumber]
		// for _, oneChanel := range channels {
		client.Say(cahnnel, randomPhrase)
		// }
	}
	// go get -u all
	// for i := 0; i < 10; i++ {

	// }
	// fmt.Println(randomMes)
}

func subResub(client *twitch.Client, message twitch.UserNoticeMessage) {
	client.Say(cahnnel, message.User.Name+" подписался! TwitchVotes "+message.SystemMsg)
	fmt.Println(message.MsgID, message.MsgParams, message.SystemMsg, message.Tags)
}

// обращение, команды к боту
func sayTalk(client *twitch.Client, message twitch.PrivateMessage) {
	// обращение к мистеру Джону
	nameJohn := [...]string{
		"John", "Джон", "mr. John", "mr. Jon", "mrJohnBot",
	}

	reg := regexp.MustCompile(`[a-zA-Zа-яА-Я]+`)
	saySay := reg.FindAllString(message.Message, -1)

says:
	for _, values := range saySay {
		for _, findJohn := range nameJohn {
			if values == findJohn {
				client.Say(message.Channel, message.User.Name+", да, сэ-э-эр!")
				break says
			}
		}
	}

	ticker := time.Tick(time.Second)
time:
	for _, quest := range saySay {
		if quest == "время" {
			for {
				select {
				case <-ticker:
					timers[1]++
					fmt.Printf("%s", timers)
					client.Say(message.Channel, "q")
					break time
				}
			}
		}
	}
}

func htmlHH() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	file := "gif.html"

	// 	// data := "100.gif"
	// 	// // tmpl, _ := template.ParseFiles("gif.html")
	// 	// tmpl, _ := template.ParseGlob("gif.html")
	// 	// tmpl.Execute(w, data)
	// })

	fs := http.FileSystem.Open("gif")(`/`)
	http.Handle("/", fs)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
