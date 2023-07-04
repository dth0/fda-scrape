package scrape

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func FisHandler(config *Config) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		papel := r.URL.Query().Get("papel")

		if len(papel) == 0 {
			response(w, nil, http.StatusBadRequest)
			return
		}

		log.Println("Scraping: ", config.TargerAddr)

		fis := NewFis(papel)

		c := colly.NewCollector(
			colly.CacheDir(config.CacheDir),
		)

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0")
		})

		c.OnResponse(func(r *colly.Response) {
			log.Printf("(%s) Response Code: %d", papel, r.StatusCode)
		})

		c.OnError(func(_ *colly.Response, err error) {
			log.Println("Error: ", err)
		})

		c.OnHTML("div.conteudo.clearfix", func(h *colly.HTMLElement) {

			var key string

			valor := make(map[string]string)
			h.ForEach("table.w728:nth-child(4) > tbody:nth-child(1) > tr", func(_ int, hl *colly.HTMLElement) {

				key = hl.ChildText("td:nth-child(1) > span.txt")
				valor[key] = hl.ChildText("td:nth-child(2) > span.txt")

				key = hl.ChildText("td:nth-child(3) > span.txt")
				valor[key] = hl.ChildText("td:nth-child(4) > span.txt")

			})
			fis.Valor = valor

			r3 := make(Resultados)
			r12 := make(Resultados)
			h.ForEach("table.w728:nth-child(5) > tbody:nth-child(1) > tr", func(i int, hl *colly.HTMLElement) {

				if i > 0 && i < 4 {
					key = hl.ChildText("td:nth-child(3) > span.txt")
					fis.Indicadores[key] = hl.ChildText("td:nth-child(4) > span.txt")

					key = hl.ChildText("td:nth-child(5) > span.txt")
					fis.Indicadores[key] = hl.ChildText("td:nth-child(6) > span.txt")

				}

				if i > 5 && i < 10 {
					key = hl.ChildText("td:nth-child(3) > span.txt")
					r12[key] = hl.ChildText("td:nth-child(4) > span.txt")

					key = hl.ChildText("td:nth-child(5) > span.txt")
					r3[key] = hl.ChildText("td:nth-child(6) > span.txt")

				}

				if i == 11 {
					key = hl.ChildText("td:nth-child(3) > span.txt")
					fis.Balanco[key] = hl.ChildText("td:nth-child(4) > span.txt")

					key = hl.ChildText("td:nth-child(5) > span.txt")
					fis.Balanco[key] = hl.ChildText("td:nth-child(6) > span.txt")
				}

			})
			fis.Resultados[3] = r3
			fis.Resultados[12] = r12

			h.ForEach("table.w728:nth-child(7) > tbody:nth-child(1) > tr", func(i int, hl *colly.HTMLElement) {

				if i == 0 {
					return
				}

				key = hl.ChildText("td:nth-child(1) > span.txt")
				fis.Imoveis[key] = hl.ChildText("td:nth-child(2) > span.txt")

				key = hl.ChildText("td:nth-child(3) > span.txt")
				fis.Imoveis[key] = hl.ChildText("td:nth-child(4) > span.txt")

				key = hl.ChildText("td:nth-child(5) > span.txt")
				fis.Imoveis[key] = hl.ChildText("td:nth-child(6) > span.txt")

			})

		})

		c.Visit(fmt.Sprintf("%s?papel=%s", config.TargerAddr, papel))

		response(w, fis, http.StatusOK)
	}
}
