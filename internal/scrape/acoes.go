package scrape

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func AcHandler(config *Config) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		papel := r.URL.Query().Get("papel")
		if len(papel) == 0 {
			response(w, nil, http.StatusBadRequest)
		}

		log.Println("Scraping: ", config.TargerAddr)

		acao := NewACAO(papel)

		c := colly.NewCollector(
			colly.CacheDir(config.CacheDir),
		)

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0")
		})

		c.OnResponse(func(r *colly.Response) {
			log.Printf("(%s) Response Code: %d", papel, r.StatusCode)
		})

		c.OnError(func(r *colly.Response, err error) {
			log.Printf("(%s) Response Code: %d, error: %v, url: %s", papel, r.StatusCode, err, r.Request.URL)
		})

		c.OnHTML("div.conteudo.clearfix", func(h *colly.HTMLElement) {
			var key string

			valor := make(map[string]string)

			h.ForEach("table.w728:nth-child(3) > tbody:nth-child(1) > tr", func(_ int, hl *colly.HTMLElement) {

				key = hl.ChildText("td:nth-child(1) > span.txt")
				valor[key] = hl.ChildText("td:nth-child(2) > span.txt")

				key = hl.ChildText("td:nth-child(3) > span.txt")
				valor[key] = hl.ChildText("td:nth-child(4) > span.txt")

			})
			acao.Valor = valor

			h.ForEach("table.w728:nth-child(5) > tbody:nth-child(1) > tr", func(i int, hl *colly.HTMLElement) {

				if i < 1 {
					return
				}

				key = hl.ChildText("td:nth-child(1) > span.txt")
				acao.Balanco[key] = hl.ChildText("td:nth-child(2) > span.txt")

				key = hl.ChildText("td:nth-child(3) > span.txt")
				acao.Balanco[key] = hl.ChildText("td:nth-child(4) > span.txt")

			})

			r3 := make(Resultados)
			r12 := make(Resultados)
			h.ForEach("table.w728:nth-child(6) > tbody:nth-child(1) > tr", func(i int, hl *colly.HTMLElement) {

				if i < 2 {
					return
				}

				key = hl.ChildText("td:nth-child(1) > span.txt")
				r12[key] = hl.ChildText("td:nth-child(2) > span.txt")

				key = hl.ChildText("td:nth-child(3) > span.txt")
				r3[key] = hl.ChildText("td:nth-child(4) > span.txt")

			})
			acao.Demonstrativos[3] = r3
			acao.Demonstrativos[12] = r12

			h.ForEach("table.w728:nth-child(4) > tbody:nth-child(1) > tr", func(i int, hl *colly.HTMLElement) {
				if i < 1 {
					return
				}

				key = hl.ChildText("td:nth-child(1) > span.txt")
				if len(key) > 0 {
					acao.Ocilacao[key] = hl.ChildText("td:nth-child(2) > span.oscil")
				}

				key = hl.ChildText("td:nth-child(3) > span.txt")
				acao.Indicadores[key] = hl.ChildText("td:nth-child(4) > span.txt")

				key = hl.ChildText("td:nth-child(5) > span.txt")
				acao.Indicadores[key] = hl.ChildText("td:nth-child(6) > span.txt")

			})
		})

		if err := c.Visit(fmt.Sprintf("%s?papel=%s", config.TargerAddr, papel)); err != nil {
			log.Println("Visit error: ", err)
		}

		response(w, acao, http.StatusOK)

	}

}
