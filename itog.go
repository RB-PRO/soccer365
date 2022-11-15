package main

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
)

// Сохранить результаты
func save_itog(f *excelize.File, itogs []itog, ssheet string) {
	var offset int = 1
	f.SetCellValue(ssheet, "A1", "Команда")
	f.SetCellValue(ssheet, "B1", "Всего игр")
	f.SetCellValue(ssheet, "C1", "Выйграно игр")
	f.SetCellValue(ssheet, "D1", "Ничья")
	f.SetCellValue(ssheet, "E1", "Пройгрышей")
	f.SetCellValue(ssheet, "F1", "Забил")
	f.SetCellValue(ssheet, "G1", "Пропустил")
	f.SetCellValue(ssheet, "H1", "+/-")
	f.SetCellValue(ssheet, "I1", "Последние данные")
	for ind, val := range itogs {
		f.SetCellValue(ssheet, "A"+strconv.Itoa(ind+1+offset), val.name)        // Название команды
		f.SetCellValue(ssheet, "B"+strconv.Itoa(ind+1+offset), val.count_games) // Всего игр
		f.SetCellValue(ssheet, "C"+strconv.Itoa(ind+1+offset), val.count_win)   // Выйграно игр
		f.SetCellValue(ssheet, "D"+strconv.Itoa(ind+1+offset), val.count_draw)  // Ничья
		f.SetCellValue(ssheet, "E"+strconv.Itoa(ind+1+offset), val.count_lost)  // Пройгрышей
		f.SetCellValue(ssheet, "F"+strconv.Itoa(ind+1+offset), val.count_in)    // Забил
		f.SetCellValue(ssheet, "G"+strconv.Itoa(ind+1+offset), val.count_out)   // Пропустил
		f.SetCellValue(ssheet, "H"+strconv.Itoa(ind+1+offset), val.koef)        // +/-
		f.SetCellValue(ssheet, "I"+strconv.Itoa(ind+1+offset), val.obsh)        // Последние данные
	}
	if err := f.SaveAs(file_itog); err != nil {
		fmt.Println(err)
	}
	f.Close()
}

// Получить результаты всех матчей
func itog_of_lig_god(url string) []itog {
	var itogs []itog
	var tecal_itog itog
	c := colly.NewCollector()
	c.OnHTML("table[class=stngs] tbody tr", func(e *colly.HTMLElement) {
		tecal_itog.name = e.DOM.Find("td:nth-child(2) a").Text()
		tecal_itog.count_games, _ = strconv.Atoi(e.DOM.Find("td:nth-child(3)").Text())
		tecal_itog.count_win, _ = strconv.Atoi(e.DOM.Find("td:nth-child(4)").Text())
		tecal_itog.count_draw, _ = strconv.Atoi(e.DOM.Find("td:nth-child(5)").Text())
		tecal_itog.count_lost, _ = strconv.Atoi(e.DOM.Find("td:nth-child(6)").Text())
		tecal_itog.count_in, _ = strconv.Atoi(e.DOM.Find("td:nth-child(7)").Text())
		tecal_itog.count_out, _ = strconv.Atoi(e.DOM.Find("td:nth-child(8)").Text())
		tecal_itog.koef, _ = strconv.Atoi(e.DOM.Find("td:nth-child(9)").Text())
		tecal_itog.obsh, _ = strconv.Atoi(e.DOM.Find("td:nth-child(10) b").Text())

		itogs = append(itogs, tecal_itog)
	})
	c.Visit(url)

	return itogs
}
