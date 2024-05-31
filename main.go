package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

type Proxy struct {
	IP   string
	Port string
	User string
	Pass string
}

var (
	PostgresqlHost     = os.Getenv("POSTGRESQL_HOST")
	PostgresqlPort     = os.Getenv("POSTGRESQL_PORT")
	PostgresqlUser     = os.Getenv("POSTGRESQL_USER")
	PostgresqlPassword = os.Getenv("POSTGRESQL_PASSWORD")
	PostgresqlDbname   = os.Getenv("POSTGRESQL_DBNAME")
	BotToken           = os.Getenv("BotToken")
)

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		PostgresqlHost, PostgresqlPort, PostgresqlUser, PostgresqlPassword, PostgresqlDbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, chat_id BIGINT UNIQUE)")
	if err != nil {
		log.Fatalf("Ошибка создания таблицы users: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS proxy (id SERIAL PRIMARY KEY, ip VARCHAR(15), port VARCHAR(15), username VARCHAR(50), pass VARCHAR(50))")
	if err != nil {
		log.Fatalf("Ошибка создания таблицы proxy: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil && strings.ToLower(update.Message.Text) == "/start" {
				chatID := update.Message.Chat.ID
				_, err := db.Exec("INSERT INTO users (chat_id) VALUES ($1) ON CONFLICT (chat_id) DO NOTHING", chatID)
				if err != nil {
					log.Printf("Ошибка добавления chat_id в БД: %v", err)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Ты подписан на алерты.")
					_, err := bot.Send(msg)
					if err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
					}
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ticker.C:
				proxies, err := getProxiesFromDB(db)
				if err != nil {
					log.Printf("Ошибка получения прокси из БД: %v", err)
					continue
				}
				checkProxies(bot, db, proxies)
			}
		}
	}()

	go func() {
		router := gin.Default()
		router.GET("/", func(context *gin.Context) {
			context.String(http.StatusOK, "Работает")
		})
		err := router.Run()
		if err != nil {
			log.Fatalf("[Error] failed to start Gin server due to: %v", err)
		}
	}()

	select {}
}

func getProxiesFromDB(db *sql.DB) ([]Proxy, error) {
	rows, err := db.Query("SELECT ip, port, username, pass FROM proxy")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proxies []Proxy
	for rows.Next() {
		var proxy Proxy
		if err := rows.Scan(&proxy.IP, &proxy.Port, &proxy.User, &proxy.Pass); err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}

func checkProxies(bot *tgbotapi.BotAPI, db *sql.DB, proxies []Proxy) {
	for _, proxy := range proxies {
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxy.User, proxy.Pass, proxy.IP, proxy.Port)
		proxyFunc := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyURL)
		}

		transport := &http.Transport{Proxy: proxyFunc}
		client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

		resp, err := client.Get("https://steampowered.com/")
		if err != nil || resp.StatusCode != http.StatusOK {
			sendAlert(bot, db, proxy, err)
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func sendAlert(bot *tgbotapi.BotAPI, db *sql.DB, proxy Proxy, err error) {
	message := fmt.Sprintf("Прокся %s:%s отъебнула: %v", proxy.IP, proxy.Port, err)
	rows, err := db.Query("SELECT chat_id FROM users")
	if err != nil {
		log.Printf("Ошибка выбора chat_id из БД: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var chatID int64
		if err := rows.Scan(&chatID); err != nil {
			log.Printf("Ошибка сканирования chat_id: %v", err)
			continue
		}
		msg := tgbotapi.NewMessage(chatID, message)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Ошибка отправки сообщения в чат: %v", err)
		}
	}
}
