package mini_spider

import (
	"fmt"
	"sync"

	"github.com/Fromsko/gouitls/logs"

	"github.com/gocolly/colly"
)

var log = logs.BaseLog

type BaseSpider struct {
	Name      string
	Collector *colly.Collector
	Mux       sync.Mutex
	BaseUrl   string
	UA        string
}

func (sp *BaseSpider) GetIndexPage() {
	sp.Collector.OnRequest(func(r *colly.Request) {
		log.InfoMsg("Spider", sp.Name, "is Running!")
		log.InfoMsg(
			fmt.Sprintf("Visited url => %s", r.URL),
		)
	})
}

func (sp *BaseSpider) ParserHtml(tmp chan map[string]string) {
	sp.Collector.OnHTML("[href]", func(h *colly.HTMLElement) {
		if h.Text != "" {
			result := make(map[string]string)

			result[h.Text] = h.Attr("href")

			tmp <- result

			log.InfoMsg(h.Text)
			h.Request.Visit(h.Attr("href"))
		}
	})
}
