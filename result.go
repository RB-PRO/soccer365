package main

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly/v2"
)

// Сохранить результаты
func save_res(results []result) {
	var tecal_ssheet string = "main"
	var offset int = 1
	f.NewSheet(ssheet)
	f.SetCellValue(ssheet, "A1", "Команда 1")
	f.SetCellValue(ssheet, "B1", "Голы 1")
	f.SetCellValue(ssheet, "C1", "Голы 2")
	f.SetCellValue(ssheet, "D1", "Команда 2")
	for ind, val := range results {
		f.SetCellValue(ssheet, "A"+strconv.Itoa(ind+1+offset), val.left.name)
		f.SetCellValue(ssheet, "B"+strconv.Itoa(ind+1+offset), val.left.gols)
		f.SetCellValue(ssheet, "C"+strconv.Itoa(ind+1+offset), val.right.gols)
		f.SetCellValue(ssheet, "D"+strconv.Itoa(ind+1+offset), val.right.name)
	}
	if err := f.SaveAs(file_result); err != nil {
		fmt.Println(err)
	}
	f.Close()
}

// Получить результаты всех матчей
func result_of_lig_god(url string) []result {
	var res []result
	var tecal_res result
	c := colly.NewCollector()
	c.OnHTML("div[class^=game_block]", func(e *colly.HTMLElement) {
		tecal_res.left.name = e.DOM.Find("div[class=ht] div[class=name] span").Text()
		tecal_res.left.gols, _ = strconv.Atoi(e.DOM.Find("div[class=ht] div[class=gls]").Text())

		tecal_res.right.name = e.DOM.Find("div[class=at] div[class=name] span").Text()
		tecal_res.right.gols, _ = strconv.Atoi(e.DOM.Find("div[class=at] div[class=gls]").Text())

		res = append(res, tecal_res)
	})
	c.Visit(url + "results/")

	return res
}
