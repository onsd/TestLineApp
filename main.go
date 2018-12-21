package main

import (
    "fmt"
    "github.com/line/line-bot-sdk-go/linebot"
    "log"
    "math/rand"
    "net/http"
    "os"
    "regexp"
    "time"
)

func main(){
    bot ,err := linebot.New(
        os.Getenv("ChannelSecret"),
        os.Getenv("AccessToken"),
        )
    //err
    if err != nil{
        log.Fatalf("Error at linebot.new: %v",err)
    }

    http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
        //fmt.Printf("ping\n")
        events, err := bot.ParseRequest(req)
        if err != nil {
            if err == linebot.ErrInvalidSignature {
                w.WriteHeader(400)
            } else {
                w.WriteHeader(500)
            }
            return
        }
        for _, event := range events {
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type) {
                case *linebot.TextMessage:
                    log.Printf("%v", message)
                    re,_ := regexp.Compile(`.*こんにちは.*`)
                    if re.MatchString(message.Text){
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("こんにちは！")).Do(); err != nil {
                            log.Print(err)
                        }
                    }
                    re, _ = regexp.Compile(`.*占い.*`)
                    if re.MatchString(message.Text){
                        rand.Seed(time.Now().UnixNano())
                        luck := []string{"大凶","凶","末吉","小吉","中吉","吉","大吉"}
                        content := "あなたの運勢は" + luck[rand.Intn(7)] + "!"
                        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
                            log.Print(err)
                        }
                    }
                case *linebot.StickerMessage:
                    log.Printf("Sticker: %v", message)
                    if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("1","1")).Do(); err != nil{
                        log.Print(err)
                }
                }
            }
        }
    })
    fmt.Println("Serve at " + os.Getenv("PORT"))
    if err := http.ListenAndServe(":" + os.Getenv("PORT"), nil); err != nil {
        log.Fatalf("Error http.ListenAndServe: %v",err)
    }
}
