package setn

// 抓取SETN 的最新新聞文章數據
// Author: 江謝廸
import (
	// "encoding/json"
	"fmt"
	"time"
	"strings"
	"strconv"
	"crypto/md5"
	"encoding/hex"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/sirupsen/logrus"

	mongodb "annstudio.com/spiders/crawler/db"
)

var data_field_list []mongodb.Data_field

var HEADER = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"

func init() {
	fmt.Println("spider setn init")
}

func get_detail(detail_url string) string {

	content_string := ""

	c2 := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		// colly.Async(true),
	)

	c2.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		RandomDelay: 3 * time.Second,
	})

	c2.OnHTML("body", func(e *colly.HTMLElement) {
		// logrus.Info(">>>", e.ChildText("article"))
		// logrus.Info(e.ChildText(".title"))
		// fmt.Println(e.ChildText(".story"))
		e.ForEach("article > div > p", func(_ int, el *colly.HTMLElement) {
			// logrus.Info("*->", el.Text)
			if ! strings.Contains(el.Text, "圖／") {
				content_string += el.Text
			}
		})
		// logrus.Info(e.ChildAttr(".date", "datetime"))
	})

	c2.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", HEADER)
	})

	c2.OnError(func(r *colly.Response, err error) {
		logrus.Info("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	logrus.Info("fetching url:", detail_url)
	c2.Visit("https://www.setn.com"+detail_url)

	// logrus.Info("end of fetch detail")
	return content_string
}

func Get_list(get_count int) [] mongodb.Data_field {
	cloop := 0
	

	logrus.Info("job started")

	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		RandomDelay: 3 * time.Second,
	})

	c.OnHTML("div[id='NewsList']", func(e *colly.HTMLElement) {
		// fmt.Println(e.Text)
		logrus.Info(e.ChildText(".title"))
		// fmt.Println(e.ChildText(".story"))
		e.ForEach("div[class='col-sm-12 newsItems']", func(_ int, el *colly.HTMLElement) {
			if cloop >= get_count {
				return
			}
			// <div class="newslabel-tab " style="background-color:#4f6128;"><a href="?PageGroupID=41">社會</a></div>
			// fmt.Println("***:", el.ChildAttr("div","class"))
			p_tag := el.ChildText("div > div > div > a")
			p_date_str := el.ChildText("div > div > time")
			p_article := el.ChildText("div > div > h3")
			p_link := el.ChildAttr("div > div > h3 > a", "href")
			p_content := get_detail(p_link)
			year, _, _ := time.Now().Date()
			p_date_str = strconv.Itoa(year)+"/"+p_date_str+":00"
			p_date_str = strings.Replace(p_date_str,"/","-", -1)
			p_date, err := time.Parse("2006-01-02 15:04:05", p_date_str)
			if err != nil {
				logrus.Error("*** Error:", err)
			}
			// fmt.Println("***1:",  p_date, p_tag, p_article, p_link)
			// fmt.Println("***2:",  p_content)
			if len(p_content) > 0{
				pid_ecrpt := md5.Sum([]byte(p_article+p_link))		// 使用md5 製造unique key,方便管理
				fetch_data := mongodb.Data_field {
					Pid:		  		"setn-"+hex.EncodeToString(pid_ecrpt[:]),
					Pub_date:  			p_date,
					Title:  			p_article,
					Page_link: 			"https://www.setn.com"+p_link,
					Search_key: 		"*sys-即時新聞",
					Author: 	  		"-",
					Content:   			p_content,
					Score:				0,
					Cat_tag:   			p_tag,
					Channel_url: 		"setn",
					Update_datetime: 	time.Now(),
				}
				if err != nil {
					fmt.Println(err)
					return
				}
				data_field_list = append(data_field_list, fetch_data)
				// data_field_json, err := json.MarshalIndent(&fetch_data, "", "\t")
				// logrus.Info("*->", string(data_field_json))
				// logrus.Info("===============================")
				cloop += 1
			}
			// type Data_field struct {
			// 	Pid		  	string    	`json:"pid"`
			// 	Pub_date  	time.Time 	`json:"post_date"`
			// 	Title   	string    	`json:"title"`
			// 	Page_link 	string    	`json:"url"`
			// 	Search_key 	string 		`json:"search_key"`
			// 	Author 	  	string 		`json:"author"`
			// 	Content   	string    	`json:"content"`
			// 	Score		int			`json:"dcore"`
			// 	Cat_tag   	string    	`json:"cat_tag"`
			// 	Channel_url string		`json:"channel_url"`
			//  Update_datetime 	time.Time 	`json:"update_datetime"`
			// }

		})

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", HEADER)
	})

	c.OnError(func(r *colly.Response, err error) {
		logrus.Info("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://www.setn.com/ViewAll.aspx")

	logrus.Info("end of fetch list")
	return data_field_list
}
