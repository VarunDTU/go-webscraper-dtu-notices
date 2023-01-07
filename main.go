package main

import (
	//"encoding/json"
	"net/http"

	//"fmt"
	"os"
	//"unsafe"

	"github.com/gin-gonic/gin"

	"github.com/gocolly/colly"
	//"github.com/google/go-cmp/cmp"
)

type Message struct {
	TEXT_des string `json:"text_des"`
	LINK     string `json:"link"`
}

func delChar(s string) string {
	ss := ""
	for i := 1; i < len(s); i++ {
		ss += string(s[i])
	}
	return ss
}

var prev_alltext = make([]Message, 0)

func webscraping() []Message {
	alltext := make([]Message, 0)
	site_web_address := "http://www.dtu.ac.in/"
	c := colly.NewCollector()
	c.UserAgent = "Go program"

	c.OnHTML(".latest_tab li h6 a", func(e *colly.HTMLElement) {

		temp := Message{
			TEXT_des: e.DOM.First().Text(),
			LINK:     site_web_address + delChar(e.DOM.First().AttrOr("href", "")),
		}
		//fmt.Println("-----------------------------")
		alltext = append(alltext, temp)
		//fmt.Printf("%T", e.DOM.First().Text())
		//fmt.Println(e.DOM.First().Text())
	})

	//fmt.Print(alltext[1])

	c.Visit("http://www.dtu.ac.in/")
	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", " ")
	//enc.Encode(alltext)
	// if cmp.Equal(alltext, prev_alltext) {
	// 	for i := 0; i < int(unsafe.Sizeof(prev_alltext)); i++ {
	// 		if(alltext[i]!=prev_alltext[i])
	// 	}

	// }
	prev_alltext = alltext
	return prev_alltext
}
func getNotice(c *gin.Context) {
	x := webscraping()
	c.IndentedJSON(http.StatusOK, x)
}
func main() {
	router := gin.Default()
	router.GET("/notices", getNotice)
	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8080"
	}
	router.Run(port)

}
