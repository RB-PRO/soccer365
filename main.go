package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	strs := strings.Split(text, " ")

	var sheetName string

	for _, str := range strs {
		str = strings.Replace(str, " ", "", -1)
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Replace(str, "\t", "", -1)
		input_ligs, _ := strconv.Atoi(str)

		if input_ligs != 0 {

			fmt.Printf("\n%v - %v\n", input_ligs, ligs[input_ligs-1].name)
			//fmt.Printf("%v: %v - %v\n\n", input_ligs, ligs[input_ligs-1].name, country_ligs(ligs[input_ligs-1].img))
			link_thil_lig := ligs[input_ligs-1].link

			// Составляем ссылку
			link_thil_lig += link_god

			// Making
			sheetName = ligs[input_ligs-1].name + " " + country_ligs(ligs[input_ligs-1].img)
			sheetName = createNameSheet(sheetName)
			parseLig(link_thil_lig, sheetName, f_itog, f_lig, f_result, f_out, f_out_tz)
		}
	}
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
