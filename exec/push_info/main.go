package main

import (
	"io"
	"log"
	"os"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"bytes"

	mongodb "annstudio.com/spiders/crawler/db" 
)

func get_send_body_json(to_source string, news_media string, news_title string, news_link string, news_date string) []byte {
	// Line Flex 內容(json body)
	send_info_body := `{
		"to": "{{TO_SOURCE_ID}}",
		"messages": [
			{
				"type": "flex",
				"altText": "羅波特來聊訊息！",
				"contents": {
					"type": "bubble",
					"body": {
						"type": "box",
						"layout": "vertical",
						"contents": [
						  {
							"type": "text",
							"text": "哇！一則有趣的新聞重點，您看看～",
							"weight": "bold",
							"size": "md"
						  },
						  {
							"type": "box",
							"layout": "baseline",
							"margin": "md",
							"contents": [
							  {
								"type": "text",
								"text": "{{NEWS_MEDIA}}",
								"size": "sm",
								"weight": "bold",
								"color": "#4444ff",
								"margin": "md",
								"flex": 0
							  }
							]
						  },
						  {
							"type": "box",
							"layout": "vertical",
							"margin": "lg",
							"spacing": "sm",
							"contents": [
							  {
								"type": "box",
								"layout": "baseline",
								"spacing": "sm",
								"contents": [
								  {
									"type": "text",
									"text": "{{NEWS_DATE}}",
									"color": "#FF0000",
									"size": "sm",
									"flex": 1,
									"weight": "bold",
									"align": "end"
								  },
								  {
									"type": "text",
									"text": "{{NEWS_TITLE}}",
									"wrap": true,
									"color": "#4444FF",
									"size": "sm",
									"flex": 5
								  }
								]
							  }
							]
						  }
						]
					  },
					"footer": {
						"type": "box",
						"layout": "vertical",
						"spacing": "sm",
						"contents": [
							{
								"type": "button",
								"style": "link",
								"height": "sm",
								"action": {
									"type": "uri",
									"label": "點擊閱讀全文",
									"uri": "{{NEWS_LINK}}"
								}
							}
						],
						"flex": 0
					}
				}
			}
		]
	}`

	news_date_array := strings.Split(news_date, " ")		// 只取年月日
	new_news_date := news_date_array[0][5:]
	send_info_body = strings.Replace(send_info_body, "{{NEWS_DATE}}", new_news_date, 1)
	send_info_body = strings.Replace(send_info_body, "{{TO_SOURCE_ID}}", to_source, 1)
	send_info_body = strings.Replace(send_info_body, "{{NEWS_MEDIA}}", news_media, 1)
	send_info_body = strings.Replace(send_info_body, "{{NEWS_TITLE}}", news_title, 1)
	send_info_body = strings.Replace(send_info_body, "{{NEWS_LINK}}", news_link, 1)

	jsonStr := []byte(send_info_body)
	return jsonStr
}

func refine_ulr (url string, media string) string {
	media_prefix := "https://tw.news.yahoo.com"
	switch media {
	case "ettoday":
		media_prefix = "https://www.ettoday.net/"
	case "tvbs":
		media_prefix = "https://news.tvbs.com.tw"
	case "setn":
		media_prefix = "https://www.setn.com"
	case "chinatimes":
		media_prefix = "https://www.chinatimes.com"
	case "mirrormedia":
		media_prefix = "https://www.mirrormedia.mg"
	default:
		media_prefix = "https://tw.news.yahoo.com"
	}
	if strings.Contains(url, "https://") {
		return url
	} else {
		return media_prefix+url
	}

}


func main() {
	logrus.SetFormatter(&logrus.TextFormatter{}) // or .JSONFormatter{}
	logrus.SetLevel(logrus.InfoLevel)			 // Can be DebugLevel, ErrorLevel, InfoLevel, WarnLevel

	writer1 := os.Stdout
	writer2, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(writer1, writer2))
	// Line access token 令牌
	access_key := "WBo4wq23pGaso8DKeC3jwcU2YXqxeU+7an2pvltq2MBLiKS3IYj/4QBqqw9DZrMjTH923J9A1vDCu0JrGsH58lhefmCzKTRMVdozVWYLJQP/8PdHksqnji2C2jFyA3TtCtApixSynKFzfQ8BJ734jgdB04t89/1O/w1cDnyilFU="
	send_url := "https://api.line.me/v2/bot/message/push"

	source_id, news_title, news_link, news_time, news_media := mongodb.Get_times_press()
	logrus.Info("send data check list:", source_id, news_title, news_link, news_time, news_media)

	send_body_json := get_send_body_json(source_id, news_media, news_title, refine_ulr(news_link, news_media), news_time.String())

	req, err := http.NewRequest("POST", send_url, bytes.NewBuffer(send_body_json))
	
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ access_key)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

	logrus.Info("response Status:", resp.Status)
	// logrus.Info("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
    // logrus.Info("response Body:", string(body))
	logrus.Info("request process is done.")

}
