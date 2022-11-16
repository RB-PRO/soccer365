package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
)

type lig struct {
	name string // Название лиги
	link string // Ссылка
	img  string // link of image
}
type comand struct {
	name string // Название команды
	gols int    // Количество команд
}
type result struct { // результаты матча
	left  comand
	right comand
}
type itog struct {
	name        string // Название команды
	count_games int    // Всего игр
	count_win   int    // Выйграно игр
	count_draw  int    // Ничья
	count_lost  int    // Пройгрышей
	count_in    int    // Забил
	count_out   int    // Пропустил
	koef        int    // +/-
	obsh        int    // Последние данные
}
type calcule struct {
	stats_left  itog
	stats_right itog
	game        result
	koef_left   float64 // Левый к-т домашней команды
	koef_right  float64 // Правый к-т приезжей команды
}
type linkid struct {
	link  string
	label string
}

const site string = "https://soccer365.ru"
const file_itog string = "Итоги турниров"   //"itog"
const file_lig string = "Все лиги"          //"lig"
const file_result string = "Результаты игр" //"result"
const file_out string = "Вывод"             //"out"

func main() {
	fmt.Println("Загрузка списка всех лиг")
	//  Создаём или открываем файлы
	f_itog, f_itog_err := createXLSX(file_itog)
	if f_itog_err != nil {
		fmt.Println(f_itog_err)
	}

	f_lig, f_lig_err := createXLSX(file_lig)
	if f_lig_err != nil {
		fmt.Println(f_lig_err)
	}

	f_result, f_result_err := createXLSX(file_result)
	if f_result_err != nil {
		fmt.Println(f_result_err)
	}

	f_out, f_out_err := createXLSX(file_out)
	if f_out_err != nil {
		fmt.Println(f_out_err)
	}

	f_out_tz, f_out_tz_err := openOrCreateXLSX(file_out + "_all")
	if f_out_tz_err != nil {
		fmt.Println(f_out_tz_err)
	}

	startParse(f_itog, f_lig, f_result, f_out, f_out_tz)

	saveCloseExit(f_itog)
	saveCloseExit(f_lig)
	saveCloseExit(f_result)
	saveCloseExit(f_out)
	saveCloseExit(f_out_tz)

	fmt.Println("\nГотово\nНажмите на Enter")
	var input string
	fmt.Scanf("%v", &input)
}

func startParse(f_itog, f_lig, f_result, f_out, f_out_tz *excelize.File) {
	ligs := list_of_ligs()
	for ind, val := range ligs {
		fmt.Printf("%v\t%v - %v", ind+1, val.name, country_ligs(val.img))
		fmt.Println()
	}
	save_ligs(f_lig, ligs, "main")

	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	text := scanner.Text()
	text = strings.Replace(text, "  ", " ", -1)

	// Получить год
	link_god := god_of_link()

	//strs := strings.Split(text, " ")
	strs := allTournirs(text)

	var sheetName string

	for _, str := range strs {
		str = strings.Replace(str, " ", "", -1)
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, "\t", "", -1)
		input_ligs, _ := strconv.Atoi(str)

		if input_ligs != 0 {

			fmt.Printf("\n%v - %v, %v\n", input_ligs, ligs[input_ligs-1].name, country_ligs(ligs[input_ligs-1].img))
			//fmt.Printf("%v: %v - %v\n\n", input_ligs, ligs[input_ligs-1].name, country_ligs(ligs[input_ligs-1].img))
			link_thil_lig := ligs[input_ligs-1].link

			// Составляем ссылку
			link_thil_lig = makeLinkOfYear(link_thil_lig, link_god)
			fmt.Println("Загружаю:", link_thil_lig)

			if link_thil_lig != "" {
				sheetName = ligs[input_ligs-1].name + " " + country_ligs(ligs[input_ligs-1].img)
				sheetName = createNameSheet(sheetName)
				parseLig(link_thil_lig, sheetName, f_itog, f_lig, f_result, f_out, f_out_tz)
			}
		}
	}
}

func allTournirs(str string) []string {
	reference := strings.Split(str, " ")
	var output []string
	for _, val := range reference {
		//val_int, _ := strconv.Atoi(val)
		if strings.Contains(val, "-") {
			container := strings.Split(val, "-")
			if len(container) == 2 {
				val_before, _ := strconv.Atoi(container[0])
				val_after, _ := strconv.Atoi(container[1])
				if val_before > val_after {
					val_before, val_after = val_after, val_before
				}
				for i := val_before; i <= val_after; i++ {
					output = append(output, strconv.Itoa(i))
				}
			}
		} else {
			output = append(output, val)
		}

	}
	return output
}

func makeLinkOfYear(link string, year int) string {
	//fmt.Println(link)
	lil := makeAllYear(link)
	//fmt.Println(lil)

	yearStr := strconv.Itoa(year)
	yearT := yearStr[len(yearStr)-2:]

	for _, val := range lil {
		tecalYear := val.label[len(val.label)-2:]
		if tecalYear == yearT {
			//fmt.Println(tecalYear, yearT)
			return site + val.link
		}
	}

	return ""
}
func makeAllYear(link string) []linkid {
	var label string
	var tecal_year string
	var linkids []linkid
	c := colly.NewCollector()
	c.OnHTML("div[class=breadcrumb] div[class^=selectbox]:nth-of-type(2) a", func(e *colly.HTMLElement) {
		label = e.DOM.Text()
		tecal_year, _ = e.DOM.Attr("href")
		linkids = append(linkids, linkid{link: tecal_year, label: label})
	})
	c.Visit(link)
	return linkids
}

func createNameSheet(str string) string {
	if len([]rune(str)) <= 31 {
		return strings.TrimSpace(str)
	} else {
		return strings.TrimSpace(str[:31])
	}
}

func parseLig(link_thil_lig, ssheet string, f_itog, f_lig, f_result, f_out, f_out_tz *excelize.File) {

	// получить результаты всех матчей
	results := result_of_lig_god(link_thil_lig)

	// получить данные с итога
	itogs := itog_of_lig_god(link_thil_lig)

	// Расчёт по ТЗ
	calcules := calcule_res_itog(results, itogs)

	// Сохранение данных
	save_res(f_result, results, ssheet)
	save_itog(f_itog, itogs, ssheet)
	save_calcule(f_out, calcules, ssheet)
	save_calcule_other_file(f_out_tz, calcules)
}

// System function

func saveCloseExit(f *excelize.File) {
	f.DeleteSheet("DeleteMe")
	// Close the spreadsheet.
	if err := f.Save(); err != nil {
		fmt.Println(err)
	}
	// Close the spreadsheet.
	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}
func openOrCreateXLSX(filename string) (*excelize.File, error) {
	var f *excelize.File
	var err_create_open error
	if _, err_create_open = os.Stat(filename + ".xlsx"); err_create_open == nil {
		// Файл существует
		if err_create_open != nil {
			return nil, err_create_open
		}
		f, err_create_open = openXLSX(filename)
		if err_create_open != nil {
			return nil, err_create_open
		}
	} else {
		// файл не существует
		f, err_create_open = createXLSX(filename)
		if err_create_open != nil {
			return nil, err_create_open
		}
	}
	return f, nil
}

// Открыть xlsx
func openXLSX(filename string) (*excelize.File, error) {
	return excelize.OpenFile(filename + ".xlsx")
}

// Создать xlsx
func createXLSX(filename string) (*excelize.File, error) {
	f := excelize.NewFile()
	f.NewSheet("DeleteMe")
	f.DeleteSheet("Sheet1")
	if err_save := f.SaveAs(filename + ".xlsx"); err_save != nil {
		return nil, err_save
	}
	return f, nil
}
