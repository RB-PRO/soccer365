package main

import (
	"fmt"
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

const site string = "https://soccer365.ru"

func main() {
	ligs := list_of_ligs()
	for ind, val := range ligs {
		fmt.Printf("%v\t%v - %v", ind+1, val.name, country_ligs(val.img))
		fmt.Println()
	}

	fmt.Print("Введите номер интересующей Вас лиги:\n> ")
	var input_ligs int
	_, err := fmt.Scanf("%d", &input_ligs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Вы выбрали ")
	fmt.Printf("%v: %v - %v\n", input_ligs, ligs[input_ligs-1].name, country_ligs(ligs[input_ligs-1].img))
	var link_thil_lig string = ligs[input_ligs-1].link

	fmt.Print("Введите номер интересующей Вас год\n(Например: 22 для 2021/2022):\n> ")
	var input_god int
	_, err = fmt.Scanf("%d", &input_god)
	if err != nil {
		panic(err)
	}
	link_god := god_of_link(input_god)

	// Составляем ссылку
	link_thil_lig += link_god

	// получить результаты всех матчей
	results := result_of_lig_god(link_thil_lig)

	// получить данные с итога
	itogs := itog_of_lig_god(link_thil_lig)

	// Сохранение данных
	save_res(results)
	save_ligs(ligs)
	save_itog(itogs)
}
