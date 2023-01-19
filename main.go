package main

import (
	//"encoding/json"

	"context"
	"crypto/sha1"
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"

	"github.com/gocolly/colly"
)

type Message struct {
	TEXT_des string    `json:"text_des"`
	LINK     []string  `json:"link"`
	Date     time.Time `json:"date"`
}
type Profile struct {
	ROLLNUMBER string `json:"id"`
	NAME       string `json:"name"`
	PROGRAM    string `json:"program"`
	BRANCH     string `json:"branch"`
	SEMESTER   string `json:"semester"`
	EMAIL      string `json:"email"`
	DOB        string `json:"dob"`
	ADDRESS    string `json:"address"`
	PHONE      string `json:"phonenumber"`
}

func delChar(s string) string {
	ss := ""
	for i := 1; i < len(s); i++ {
		ss += string(s[i])
	}
	return ss
}

func set_time_date(date string) time.Time {
	if date == "" {

		t, _ := time.Parse("2006/01/02", "2000/01/01")
		return t
	}

	temp := strings.Split(date, ".")
	year := temp[2]
	month := temp[1]
	day := temp[0]
	t, _ := time.Parse("2006/01/02", year+"/"+month+"/"+day)
	return t

}

var prev_alltext = make([]Message, 0)
var current_profile Profile

func webscraping_latest_notices() {
	for range time.Tick(time.Second * 19) {

		alltext := make([]Message, 0)
		site_web_address := "http://www.dtu.ac.in/"
		c := colly.NewCollector()
		c.UserAgent = "Go program"

		c.OnHTML(".latest_tab li", func(e *colly.HTMLElement) {
			links_array := e.ChildAttrs("a", "href")

			for i := 0; i < len(links_array); i++ {
				links_array[i] = site_web_address + delChar(links_array[i])
			}
			notice_date := e.ChildText("small em i")

			notice_date_time := set_time_date(notice_date)

			temp := Message{
				TEXT_des: e.ChildText("h6 a"),
				LINK:     links_array,
				Date:     notice_date_time,
			}

			alltext = append(alltext, temp)

			// fmt.Printf("%T", e.DOM.First().Text())
			// fmt.Println(e.DOM.First().Text())
		})

		c.Visit("http://www.dtu.ac.in/")
		// enc := json.NewEncoder(os.Stdout)
		// enc.SetIndent("", " ")
		// enc.Encode(alltext)
		sort.Slice(alltext, func(i, j int) bool { return alltext[i].Date.After(alltext[j].Date) })
		prev_alltext = alltext

	}
}

func getNotice(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, prev_alltext)
}

func getuserinfo(c *gin.Context) {
	erp_login_page(c.Query("id"), c.Query("password"))
	c.IndentedJSON(http.StatusOK, current_profile)
}

func erp_login_page(s1, s2 string) {
	temp_h := ""
	current_profile.ROLLNUMBER = s1
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()
	err := chromedp.Run(ctx, chromedp.Navigate("https://cumsdtu.in/student_dtu/login/login.jsp"),

		chromedp.Text(`body`, &temp_h),
		chromedp.SetValue("usernameId", s1),
		chromedp.SetValue("passwordId", s2),
		chromedp.Click("submitButton"),
		chromedp.Sleep(time.Second*3),
		chromedp.Click("Link145"),
		chromedp.Text(`ListItem248`, &current_profile.ADDRESS),
		chromedp.Text(`ListItem248`, &current_profile.NAME),
		chromedp.Text(`ListItem248`, &current_profile.PROGRAM),
		chromedp.Text(`ListItem248`, &current_profile.BRANCH),
		chromedp.Text(`Label621`, &current_profile.SEMESTER),
		chromedp.Text(`ListItem193`, &current_profile.EMAIL),
		chromedp.Click("Layer_1-2"),
	)
	if err != nil {

		log.Fatal(err)
	}
	log.Fatal(temp_h)

}

func main() {

	go webscraping_latest_notices()

	router := gin.Default()
	router.GET("/notices", getNotice)
	router.GET("/test", getuserinfo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8081"
	} else {
		port = ":" + port
	}
	router.Run(port)

}
func calsha1(str string) string {

	sha := sha1.New()
	sha.Write([]byte(str))
	digest := sha.Sum(nil)
	return string(digest)

}

func calsha256(str string) string {

	sha := sha256.New()
	sha.Write([]byte(str))
	digest := sha.Sum(nil)
	return string(digest)
}
