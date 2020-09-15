package main

import (
	"TwitchChat/wikijohn"
	"fmt"
	"html/template"

	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/zmb3/spotify"
)

var messagesm = make(map[string]int)
var warning = make(map[string]int)
var startTime = time.Now()
var wordRes []string
var CountInt int

var UrlGet string

// PathFile тип стуктуры
type PathFile struct {
	Gifs  string
	Count []CounterUsers
}

// CounterUsers принимаем колличество юзеров
type CounterUsers struct {
	CountUser string
}

// TextOut перевод масиива в строку плюс склеивание слов
func TextOut(text []string) string {
	t := strings.Join(text[:], "")
	return fmt.Sprintf("%s", t)
}

// ReadWord func bad word
func ReadWord(text string, check string) (c bool, err error) {
	c = strings.Contains(text, check)
	return c, nil
}

// Error обработка ошибок
func Error(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func update() {
	// integer := CountInt
	// go httpserverobs.OpenServer(CountInt)

}

func main() {
	update()

	client := twitch.NewClient("mrJohnBot", "oauth:u5sfiw5i6cawt7kgkmwsnrk55c6g9h")
	// guichat.InitGui(client * twitch.)
	chaArg := "" // os.Args[1:] наименование канала нужно прописать.

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		logChat(message)
		countMessages(client, chaArg, message)
		badWords(client, message, chaArg)
		go sayTalk(client, message, chaArg)
		router(client, message, chaArg)
		shareNetwork(client, message, chaArg)
		go mesWar(client, message, chaArg)
		wikijohn.Wiki(client, message, chaArg)
		// go outGifs(client, message)
		// userCounMes(client, user, message)
	})
	// lib spotify
	spotifySound()

	// sub, resub and raids
	client.OnNoticeMessage(func(message twitch.NoticeMessage) {
		// subResub(client, message, user, chaArg)
	})

	client.OnWhisperMessage(func(message twitch.WhisperMessage) {
		whisper(message)
	})

	go countOnline(client, chaArg)

	client.Join(chaArg) // (chaArg[0])

	err := client.Connect()
	Error(err)

}

// лог сообщений из чата
func logChat(message twitch.PrivateMessage) {
	usChat, err := os.OpenFile("./logs/Лог чата"+startTime.Format(" 02-01-2006")+".txt", os.O_APPEND|os.O_CREATE, 0644)
	Error(err)

	defer usChat.Close()
	_, err = usChat.Write([]byte(message.Time.Format("15:04:05") + " | " + message.User.DisplayName + " - " + message.Message + "\n"))
	Error(err)
}

// подсчет зрителей
func countOnline(client *twitch.Client, chaArg string) {
	// var IntCountUsers int
	timer := time.NewTicker(20 * time.Second)
	for range timer.C {
		count, err := client.Userlist(chaArg)
		Error(err)
		fmt.Println("cout", len(count))
		CountInt = len(count)
	}

	// CountInt = strconv.Itoa(IntCountUsers)
}

func badWords(client *twitch.Client, message twitch.PrivateMessage, chaArg string) {
	// ПРОВЕРКА НА ПЛОХИЕ СЛОВА!!!
	rer, err := ioutil.ReadFile("word.bad")
	Error(err)

	// разбиваем по запятые
	spl := strings.Split(string(rer), ", ")

	// проверка на матерные слова
stop:
	for _, setWords := range spl {
		word, err := ReadWord(message.Message, setWords)
		Error(err)

		if word == true {
			warning[message.User.DisplayName]++

			if warning[message.User.DisplayName] == 1 {
				client.Say(chaArg, "Mr. John не любит плохие слова! :-(")
			} else if warning[message.User.DisplayName] == 2 {
				client.Say(chaArg, message.User.DisplayName+" не ругайся, 300 сек в студию!")
				client.Say(chaArg, "/timeout "+message.User.DisplayName+" 300")
			} else {
				client.Say(chaArg, message.User.DisplayName+", твои дни сочтены!!!")
				// 259200 -
				client.Say(chaArg, "/ban "+message.User.DisplayName)
			}
			break stop
		}
	}
}

// подсчет сообщений
func countMessages(client *twitch.Client, chaArg string, message twitch.PrivateMessage) {
	messagesm[message.User.ID]++
	if message.Message == "!use" {
		var couMes string
		users, err := client.Userlist(chaArg) // client.Userlist(channel) - список зрителей
		Error(err)
		for _, uu := range users {
			// сделать проверку на юзера
			if uu == message.User.DisplayName {
				couMes = uu + " - " + strconv.Itoa(messagesm[message.User.ID]) + " сообщение/й за стрим\n"
			}
		}
		client.Say(chaArg, couMes)
	}
}

// функция вывода рандомных сообщений
func mesWar(client *twitch.Client, message twitch.PrivateMessage, chaArg string) {
	// каждые 15 минут выводит сообщение
	ticker := time.NewTicker(time.Minute * 15)

	// карта с фразами
	randomMes := [...]string{
		"Попробуй команду Поиск/ (пример: Поиск/шаман)",
	}

	// рандомный вывод сообщений
rand:
	for range ticker.C {
		leng := len(randomMes)
		randomNumber := rand.Intn(leng)
		randomPhrase := randomMes[randomNumber]
		client.Say(chaArg, randomPhrase)
		break rand
	}
}

func subResub(client *twitch.Client, message twitch.PrivateMessage, user twitch.User, chaArg string) {
	client.Say(chaArg, message.User.DisplayName+" подписался! TwitchVotes "+message.Raw)
	fmt.Println(message.Action, message.Emotes, message.Type, message.Tags)
}

// обращение, команды к боту
func sayTalk(client *twitch.Client, message twitch.PrivateMessage, chaArg string) {
	nameJohn := [...]string{
		"John", "Джон", "mr. John", "mr. Jon", "mrJohnBot", "mrjohnbot", "jon", "john",
	}

	reg := regexp.MustCompile(`[a-zA-Zа-яА-Я]+`)
	saySay := reg.FindAllString(message.Message, -1)

says:
	for _, values := range saySay {
		for _, findJohn := range nameJohn {
			if values == findJohn {
				client.Say(chaArg, message.User.DisplayName+", да, сэ-э-эр!")
				break says
			}
		}
	}
}

func router(client *twitch.Client, message twitch.PrivateMessage, chaArg string) {
	switch message.Message {
	case "!Время":
		client.Say(chaArg, "С начала срима прошло - "+printElapsedTime())
		break
		// ссылка в телеграм
	case "!tel":
		client.Say(chaArg, message.User.DisplayName+" https://t.me/joinchat/AAAAAFGAVk9hZ7vAci-mNQ")
		break
	}
}

// вывод вркмени в чат с начала трансляции отсчет
// const (
// 	pars = "72h3m0.5s"
// 	now  = "15:04:05"
// )

func printElapsedTime() string {
	elapsed := time.Since(startTime)
	elapsed = elapsed.Round(time.Minute)
	h := elapsed / time.Hour
	elapsed -= h * time.Hour
	m := elapsed / time.Minute
	spl := strings.Replace(elapsed.String(), "h", ":", 2)

	return fmt.Sprintf("Время стрима = %s:%s", spl, m)
}

func shareNetwork(client *twitch.Client, message twitch.PrivateMessage, chaArg string) {
	reg := regexp.MustCompile(`Поиск\/`)
	saySay := reg.FindAllString(message.Message, -1)

share:
	for _, values := range saySay {
		if values == "Поиск/" {
			regsha := regexp.MustCompile(`\/.+`)
			sayShare := regsha.FindString(strings.Replace(message.Message, " ", "_", -1))

			sha, err := http.Get("https://ru.wikipedia.org/wiki" + strings.Replace(sayShare, " ", "_", -1))
			Error(err)
			defer sha.Body.Close()

			pars, err := goquery.NewDocumentFromReader(sha.Body)
			Error(err)

			pars.Find(".mw-parser-output").Each(func(i int, s *goquery.Selection) {
				bandCo := s.Find("p").Text()
				text := bandCo[:strings.IndexByte(bandCo, '\n')]

				leng := utf8.RuneCountInString(text)

				if leng > 500 {
					spl := strings.Split(text, "")
					// вот что нужно сделать чтобы вывести первые 500 символов
					for count := 0; count < 499; count++ {
						// собираем слова обратно в массив
						wordRes = append(wordRes, spl[count])
						if count == 498 {
							client.Say(chaArg, TextOut(wordRes))
						}

						minus := len(spl) - 499
						if count == 498 && minus > 0 {
							// чистим массив
							wordRes = nil
							for count := 498; count < leng; count++ {
								wordRes = append(wordRes, spl[count])
							}
						}
					}
					client.Say(chaArg, TextOut(wordRes))
				} else {
					client.Say(chaArg, text)
				}
			})
			break share
		}
	}
}

func userCounMes(client *twitch.Client, user twitch.User, message twitch.PrivateMessage, chaArg string) {
	// склад юзеров
	var userStorage = make(map[string]int)

	t := time.Now().Format("15:04:05")
	fmt.Printf("%s = ", t)
	fmt.Println(message.User.DisplayName, "- зашел в чат")
	if message.User.DisplayName == "mrjohnbot" {
		fmt.Println("Yes, se-e-er!")
	} else if message.User.DisplayName == "стример" {
		client.Say(chaArg, "Здравствуй, "+message.User.DisplayName+", я mr. John. И умею следить за порядком в чатике!")
	}

	read, err := os.OpenFile("Отчет по трансляции.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Error(err)
	defer read.Close()
	// проверка на повторных юзеров
	userStorage[message.User.DisplayName]++
	if userStorage[message.User.DisplayName] == 1 {
		_, err := read.Write([]byte(t + " - " + message.User.DisplayName + "\n"))
		Error(err)
	}
}

func spotifySound() {
	// URL перенаправления должен точно соответствовать URL, который вы зарегистрировали для своего приложения
	// области действия определяют, какие разрешения пользователю предлагается авторизовать
	auth := spotify.NewAuthenticator("https://developer.spotify.com/dashboard/applications/", "")
	fmt.Println(auth)
	// если вы не сохранили свой идентификатор и секретный ключ в указанных переменных среды,
	// вы можете установить их здесь вручную
	auth.SetAuthInfo("2b21fc85a5494f29ac11d35145d8a141", "fc89960c2dd14ca8afbbb5a064d68385")
}

// PathGif путь до картинок
var PathGif string

// ниже функ для вывода гиф анимаций на экран, не доделана
func outGifs(client *twitch.Client, message twitch.PrivateMessage) {
	const tpl = `
	<!DOCTYPE html>
	<html>
		<body>
			<img src="{{ . }}">
		</body>
	</html>`

	readGif, err := ioutil.ReadDir(`C:\Users\volgi\OneDrive\Рабочий стол\TwitchChat\gifs`)
	Error(err)

	for _, gg := range readGif {
		switch message.Message {
		case "подписка":
			if gg.Name() == "chear.gif" {
				PathGif = gg.Name()
				log.Println(gg.Name())
			}
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := PathFile{
				Gifs: PathGif,
			}
			tmpl, err := template.New("webpafe").Parse(tpl)
			if err != nil {
				log.Println(err.Error())
			}
			tmpl.Execute(w, path)
		})

		log.Fatal(http.ListenAndServe("localhost:8080", nil))
	}
}

func whisper(message twitch.WhisperMessage) {
	if message.User.DisplayName == "monstrumgear" {
		fmt.Println(message.Message + " приватное сообщение")
	}
}
