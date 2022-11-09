package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
)

func god_of_link() string {
	fmt.Print("Введите номер интересующей Вас год\n(Например: 22 для 2021/2022):\n> ")
	var inp int
	_, err := fmt.Scanf("%d", &inp)
	if err != nil {
		panic(err)
	}
	if inp == 23 {
		return ""
	} else if inp < 23 {
		return strconv.Itoa(inp+2000-1) + "-" + strconv.Itoa(inp+2000) + "/"
	}
	return ""
}

// Получить название страны из ссылки
func country_ligs(link_img string) string {

	if strings.Contains(link_img, "https://s.scr365.net/img/flags/16/") {
		link_img = strings.Replace(link_img, "https://s.scr365.net/img/flags/16/", "", -1)
		link_img = strings.Replace(link_img, ".png", "", -1)
	} else {
		link_img = ""
	}
	return link_img
}

// Получить массив ссылок на все лиги
func list_of_ligs() []lig {
	var ligs []lig
	var tecal_lig lig
	var exits bool
	c := colly.NewCollector()
	c.OnHTML("div[class=season_item]", func(e *colly.HTMLElement) {
		tecal_lig.name = e.DOM.Find("a span").Text()
		tecal_lig.img, _ = e.DOM.Find("img").Attr("src")
		tecal_lig.link, _ = e.DOM.Find("a").Attr("href")
		tecal_lig.link = site + tecal_lig.link
		if len(ligs) > 1 {
			if ligs[0].link == tecal_lig.link {
				exits = true
				return
			}
		}
		ligs = append(ligs, tecal_lig)
	})
	c.Visit("https://soccer365.ru/index.php?c=competitions&a=champs_list_data&tp=0&cn_id=0&st=0&ttl=&p=1")
	/*
		for i := 1; ; i++ {
			c.Visit("https://soccer365.ru/index.php?c=competitions&a=champs_list_data&tp=0&cn_id=0&st=0&ttl=&p=" + strconv.Itoa(i))
			if exits {
				break
			}
		}
	*/
	fmt.Println(exits)
	return ligs
}

// Сохранить лиги в файл
func save_ligs(ligs []lig) {
	f := excelize.NewFile()
	f.NewSheet("main")
	f.DeleteSheet("Sheet1")
	for ind, val := range ligs {
		f.SetCellValue("main", "A"+strconv.Itoa(ind+1), val.name)
		f.SetCellValue("main", "B"+strconv.Itoa(ind+1), val.link)
		f.SetCellValue("main", "C"+strconv.Itoa(ind+1), val.img)
		f.SetCellValue("main", "D"+strconv.Itoa(ind+1), country_ligs(val.img))
	}
	if err := f.SaveAs("ligs.xlsx"); err != nil {
		fmt.Println(err)
	}
}
