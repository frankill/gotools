package array

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

// WriteToExcelStringSliceChannel 从一个字符串切片的通道读取数据，并将其写入指定的 Excel 文件。
// 参数:
//
//	ch - 包含字符串切片的通道，每个切片代表一行数据。
//	stop - 用于停止写入过程的信号通道。
//	filename - 输出的 Excel 文件的路径。
//	sheet - Excel 文件中的工作表名称，默认为 "Sheet1"。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToExcelStringSliceChannel(ch chan []string, stop chan struct{}, filename string, sheet string) error {

	defer close(stop)
	file := excelize.NewFile()

	if sheet == "" {
		sheet = "Sheet1"
	}
	stream, err := file.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	defer stream.Flush()
	num := 1

	for strSlice := range ch {
		stream.SetRow(fmt.Sprintf("A%d", num), ArrayToAny(strSlice))
		num++
	}

	if err := file.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

// WriteToExcel 将二维字符串数组 data 写入指定的 Excel 文件。
// 参数:
//
//	data - 二维字符串数组，其中每一行代表 Excel 中的一行数据。
//	filename - 输出的 Excel 文件的路径。
//	sheet - Excel 文件中的工作表名称，默认为 "Sheet1"。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToExcel(data [][]string, filename string, sheet string) error {

	file := excelize.NewFile()

	if sheet == "" {
		sheet = "Sheet1"
	}
	stream, err := file.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	defer stream.Flush()
	num := 1

	for _, row := range data {
		stream.SetRow(fmt.Sprintf("A%d", num), ArrayToAny(row))
		num++
	}

	// 保存文件
	if err := file.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

// WriteToCSVStringSliceChannel 从一个字符串切片的通道读取数据，并将其写入指定的 CSV 文件。
// 参数:
//
//	ch - 包含字符串切片的通道，每个切片代表一行数据。
//	stop - 用于停止写入过程的信号通道。
//	filename - 输出的 CSV 文件的路径。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToCSVStringSliceChannel(ch chan []string, stop chan struct{}, filename string) error {

	defer close(stop)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for row := range ch {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// WriteToCSV 将二维字符串数组 data 写入指定的 CSV 文件。
// 参数:
//
//	data - 二维字符串数组，其中每一行代表 CSV 中的一行数据。
//	filename - 输出的 CSV 文件的路径。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToCSV(data [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ReadFromCsvSliceChannel 从 CSV 文件中读取数据，并将每一行数据作为字符串切片发送到通道。
// 参数:
//
//	filename - CSV 文件的路径。
//	ch - 用于接收字符串切片的通道。
//
// 返回:
//
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromCsvSliceChannel(filename string, ch chan []string) error {

	defer close(ch)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {

		row, err := reader.Read()
		if err != nil {
			break
		}
		ch <- row

	}
	return nil
}

// ReadFromCsv 一次性读取整个 CSV 文件的内容，并返回一个二维字符串数组。
// 参数:
//
//	filename - CSV 文件的路径。
//
// 返回:
//
//	一个二维字符串数组，其中每一行代表 CSV 文件中的一行数据。
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromCsv(filename string) ([][]string, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	return reader.ReadAll()
}

// ReadFromExcel 从 Excel 文件中读取指定工作表的所有行，并返回一个二维字符串数组。
// 参数:
//
//	filename - Excel 文件的路径。
//	sheet - 工作表的名称。
//
// 返回:
//
//	一个二维字符串数组，其中每一行代表 Excel 工作表中的一行数据。
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromExcel(filename string, sheet string) ([][]string, error) {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return file.GetRows(sheet)
}

// ReadFromExcelSliceChannel 从 Excel 文件中读取指定工作表的每一行数据，并通过通道发送出去。
// 参数:
//
//	filename - Excel 文件的路径。
//	sheet - 工作表的名称。
//	ch - 用于发送字符串切片的通道。
//
// 返回:
//
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromExcelSliceChannel(filename string, sheet string, ch chan []string) error {

	defer close(ch)

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	rows, err := file.Rows(sheet)
	if err != nil {
		return err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}
		ch <- row
	}
	return nil
}
