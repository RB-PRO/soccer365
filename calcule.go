package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Сохранить результаты
func save_calcule(calcules []calcule, ssheet string) {
	var offset int = 1
	f := excelize.NewFile()
	f.NewSheet(ssheet)
	f.DeleteSheet("Sheet1")
	f.SetCellValue(ssheet, "A1", "К-т 1")
	f.SetCellValue(ssheet, "B1", "Команда 1")
	f.SetCellValue(ssheet, "C1", "Голы 1")
	f.SetCellValue(ssheet, "D1", "Голы 2")
	f.SetCellValue(ssheet, "E1", "Команда 2")
	f.SetCellValue(ssheet, "F1", "К-т 1")
	for ind, val := range calcules {
		f.SetCellValue(ssheet, "A"+strconv.Itoa(ind+1+offset), val.koef_left)
		f.SetCellValue(ssheet, "B"+strconv.Itoa(ind+1+offset), val.game.left.name)
		f.SetCellValue(ssheet, "C"+strconv.Itoa(ind+1+offset), val.game.left.gols)
		f.SetCellValue(ssheet, "D"+strconv.Itoa(ind+1+offset), val.game.right.gols)
		f.SetCellValue(ssheet, "E"+strconv.Itoa(ind+1+offset), val.game.right.name)
		f.SetCellValue(ssheet, "F"+strconv.Itoa(ind+1+offset), val.koef_right)
	}
	if err := f.SaveAs(file_calc + ".xlsx"); err != nil {
		fmt.Println(err)
	}
	f.Close()
}

// Сохранить результаты в файл по ТЗ
func save_calcule_other_file(f *excelize.File, calcules []calcule) {
	if _, err := os.Stat(file_out); err == nil { // Файл существует
		writeserDatasCalc(f, calcules)
		f.Save()
	} else { // Файл не существует
		writeserDatasCalc(f, calcules)
		if err := f.SaveAs(file_out); err != nil {
			fmt.Println(err)
		}
	}
}

// Внести всё
func writeserDatasCalc(f *excelize.File, calcules []calcule) {
	for _, val := range calcules {
		tecal_ssheet := sheet_name(val.koef_left, val.koef_right)
		//fmt.Println(">>>>>>>>>>>", tecal_ssheet, "---------", val.koef_left, val.koef_right)
		if len(tecal_ssheet) == 5 {
			f.NewSheet(tecal_ssheet)
			writeData(f, tecal_ssheet, val)
		}
		//fmt.Println(ind, tecal_ssheet, val.koef_left, val.koef_right)
	}
	f.DeleteSheet("Sheet1")
	f.Close()
}

// Составить название листа
func sheet_name(val1, val2 float64) string {
	k1 := int(math.Round(val1 * 10.0))
	k2 := int(math.Round(val2 * 10.0))
	tecal_ssheet := zero_string(k1) + "-" + zero_string(k2)
	return tecal_ssheet
}

// Преобразовать в строку с нулём
func zero_string(input int) string {
	if input < 10 {
		return "0" + strconv.Itoa(input)
	} else {
		return strconv.Itoa(input)
	}
}

// Ввести одну строку
func writeData(f *excelize.File, sSheet string, calc calcule) {
	r, err := f.GetRows(sSheet)
	if err != nil {
		panic(err)
	}
	lens := len(r)
	lens++
	f.SetCellValue(sSheet, "A"+strconv.Itoa(int(lens)), calc.game.left.gols-calc.game.right.gols)
	f.SetCellValue(sSheet, "B"+strconv.Itoa(int(lens)), calc.game.left.gols+calc.game.right.gols)
	f.SetCellValue(sSheet, "C"+strconv.Itoa(int(lens)), calc.game.left.name)
	f.SetCellValue(sSheet, "D"+strconv.Itoa(int(lens)), calc.game.right.name)
}

// Рассчёт всех данных calcule
func calcule_res_itog(results []result, itogs []itog) []calcule {
	var calcules []calcule
	var calcules_tecal calcule
	m := make(map[string]itog)
	for _, itog_val := range itogs {
		m[itog_val.name] = itog_val
	}
	for _, res_val := range results {
		calcules_tecal = calcule_single(res_val, m[res_val.left.name], m[res_val.right.name])
		calcules = append(calcules, calcules_tecal)
	}
	return calcules
}

// Одиночный расчёт
func calcule_single(res result, itog_left itog, itog_right itog) calcule {
	var calc calcule
	calc.stats_left = itog_left
	calc.stats_right = itog_right
	calc.game = res

	calc.koef_left = float64(calc.stats_left.count_in+calc.stats_right.count_out) / float64(calc.stats_left.count_games+calc.stats_right.count_games)
	calc.koef_right = float64(calc.stats_right.count_in+calc.stats_left.count_out) / float64(calc.stats_right.count_games+calc.stats_left.count_games)
	return calc
}
