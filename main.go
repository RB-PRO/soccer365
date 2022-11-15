package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
const file_out string = "out.xlsx"

func main() {
	ligs := list_of_ligs()
	for ind, val := range ligs {
		fmt.Printf("%v\t%v - %v", ind+1, val.name, country_ligs(val.img))
		fmt.Println()
	}
	//save_ligs(ligs)

	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	text := scanner.Text()
	text = strings.Replace(text, "  ", " ", -1)
	strs := strings.Split(text, " ")

	// Получить год
	link_god := god_of_link()

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
			main2(link_thil_lig)
		}
	}

}

func main2(link_thil_lig string) {

	// получить результаты всех матчей
	results := result_of_lig_god(link_thil_lig)

	// получить данные с итога
	itogs := itog_of_lig_god(link_thil_lig)

	// Расчёт по ТЗ
	calcules := calcule_res_itog(results, itogs)

	// Сохранение данных
	//save_res(results)
	//save_itog(itogs)
	//save_calcule(calcules)
	save_calcule_other_file(calcules)
}
