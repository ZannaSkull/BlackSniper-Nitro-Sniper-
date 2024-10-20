package main

import (
    "github.com/lxi1400/GoTitle"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	

	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
	"github.com/valyala/fasthttp"
)

var (
	Token     string
	WebhookURL string
	Regex     = regexp.MustCompile(`(?m)(?:discord\.gift|discord\.com/gifts)/\b([0-9a-zA-Z]{16,24})\b`)
)

func init() {
	if len(os.Args) == 1 {
		fileBytes, err := ioutil.ReadFile("config.toml")

		if err != nil {
			log.Fatalln("💔 Couldn't read config file, is it missing?")
		}

		var config = struct {
			Token     string
			WebhookURL string
		}{}

		err = toml.Unmarshal(fileBytes, &config)

		if err != nil {
			log.Fatalln("💔 Couldn't parse config file, exiting")
		}

		Token = config.Token
		WebhookURL = config.WebhookURL 
	} else {
		flag.StringVar(&Token, "t", "", "Discord Token")
		flag.StringVar(&WebhookURL, "w", "", "Webhook URL")
		flag.Parse()
	}

	if Token == "" {
		log.Fatalln("💔 No Discord token provided, exiting")
	}

	if WebhookURL == "" {
		log.Fatalln("💔 No Webhook URL provided, exiting")
	}
}

func main() {
	title.SetTitle("Black Sniper")
	blue := "\033[34m"   // Blue
	reset := "\033[0m"   // Reset color
	
	bot, err := discordgo.New(Token)

	if err != nil {
		log.Fatalln("💔 Couldn't create Discord session:", err)
	}

	fmt.Printf("\n")

	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)

	err = bot.Open()

	if err != nil {
		log.Fatalln("💔 Couldn't establish WebSocket connection:", err)
	}
	
	fmt.Println(blue + `
	██████  ██       █████   ██████ ██   ██ ███████ ███    ██ ██ ██████  ███████ ██████  
	██   ██ ██      ██   ██ ██      ██  ██  ██      ████   ██ ██ ██   ██ ██      ██   ██ 
	██████  ██      ███████ ██      █████   ███████ ██ ██  ██ ██ ██████  █████   ██████  
	██   ██ ██      ██   ██ ██      ██  ██       ██ ██  ██ ██ ██ ██      ██      ██   ██ 
	██████  ███████ ██   ██  ██████ ██   ██ ███████ ██   ████ ██ ██      ███████ ██   ██ 
	
	` + reset)
	fmt.Println("👑 Crown Sniper Based.")
	fmt.Println("👑 Bot running, press CTRL+C to exit")
	syscalls := make(chan os.Signal, 1)
	signal.Notify(syscalls, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	fmt.Printf("🔺 Signal `%v` detected, disconnecting bot and exiting\n\n", <-syscalls)

	_ = bot.Close()
}

func ready(_ *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("👤 Logged in as", event.User.String())
}

func messageCreate(_ *discordgo.Session, event *discordgo.MessageCreate) {
	start := time.Now()
	matches := Regex.FindAllStringSubmatch(event.Content, -1)

	for i := 0; i < len(matches); i++ {
		code := matches[i][1]

		if length := len(code); length != 16 && length != 24 {
			return
		}
		if code == "" {
			return
		}

		fmt.Println("🛠️  Handler -> Redeem:", time.Since(start))
		go redeemNitroGift(code, event.ChannelID)
	}
}

func redeemNitroGift(code, channelID string) {
	start := time.Now()
	requestBody := "{\"channel_id\": \"" + channelID + "\", \"payment_source_id\":null}"
	url := "https://discordapp.com/api/v9/entitlements/gift-codes/" + code + "/redeem"
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", Token)
	req.SetBodyString(requestBody)
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Println("Error making request:", err)
		return
	}

	body := resp.Body()
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		fmt.Println("✨ Successfully claimed code:", code, string(body))
		go sendWebhook("Functional code: " + code)
	} else {
		fmt.Println("⛔ Couldn't claim code:", code, string(body))
		go sendWebhook("Invalid code: " + code)
	}

	fmt.Println("🛠️  Request :", time.Since(start))
}

func sendWebhook(message string) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.SetBodyString("{\"content\":\"" + message + "\"}")
	req.SetRequestURI(WebhookURL)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Println("Error sending webhook:", err)
	}
}
