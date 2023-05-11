package excel

import (
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/fmt"
	"github.com/xuri/excelize/v2"
)

//func main() {
//	f := excelize.NewFile()
//	// 创建一个工作表
//	index, _ := f.NewSheet("Sheet2")
//	// 设置单元格的值
//	f.SetCellValue("Sheet2", "A2", "Hello world.")
//	f.SetCellValue("Sheet1", "B2", 100)
//	// 设置工作簿的默认工作表
//	f.SetActiveSheet(index)
//	// 根据指定路径保存文件
//	if err := f.SaveAs("Book1.xlsx"); err != nil {
//		fmt.Println(err)
//	}
//}

// ExcelData 包含解析Excel文件的方法返回的数据
type ExcelData struct {
	//表头
	Headers []string
	//对应到数据库中的key
	Keys []string
	//多行数据
	Rows [][]string
}

// ExcelUtil 包含解析Excel文件和组装数据的方法
type ExcelUtil struct {
	data *ExcelData
}

// NewExcelUtil 返回一个ExcelUtil实例
func NewExcelUtil() *ExcelUtil {
	return &ExcelUtil{
		data: &ExcelData{},
	}
}

// ParseFile 从指定的Excel文件中解析数据 传入指针阐述可以数据nil
func (excelUtil *ExcelUtil) ParseFile(filePath string, sheetName string, headerRow, keyRow *int) error {
	// 打开Excel文件
	f, e := excelize.OpenFile(filePath)
	if e != nil {
		return e
	}

	// 获取指定工作表的所有行
	rows, e := f.GetRows(sheetName)
	if e != nil {
		return e
	}

	//判断是否传入了headerRow, keyRow
	count, headerIndex, keyIndex := 0, 0, 0

	if headerRow != nil {
		count++
		headerIndex = *headerRow
		// 将第一行作为表头
		excelUtil.data.Headers = rows[headerIndex]
	}

	if keyRow != nil {
		count++
		// 将第二行作为key
		keyIndex = *keyRow
		excelUtil.data.Keys = rows[keyIndex]
	}

	if len(rows) < count {
		return err.New("没有发现数据")
	}

	headersLen := len(excelUtil.data.Headers)

	// 将每一行的数据保存到二维切片中
	for i, row := range rows {
		if headerIndex == i || keyIndex == i {
			continue
		}

		rowData := make([]string, headersLen)
		copy(rowData, row)
		for j := len(row); j < headersLen; j++ {
			rowData[j] = ""
		}

		excelUtil.data.Rows = append(excelUtil.data.Rows, rowData)
	}

	return nil
}

// AssembleData 使用指定方法将Excel数据组装成map
func (excelUtil *ExcelUtil) AssembleData(f func(headers, Keys, row []string) (map[string]string, error)) ([]map[string]string, error) {
	if excelUtil.data == nil {
		return nil, err.New("没有发现数据")
	}

	var result []map[string]string

	// 遍历每一行的数据，并将其组装成map
	for _, row := range excelUtil.data.Rows {
		if len(row) != len(excelUtil.data.Headers) {
			return nil, err.New("行数据长度与标头长度不匹配")
		}

		m, err := f(excelUtil.data.Headers, excelUtil.data.Keys, row)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}

// AssembleDataDef 使用默认方法进行key传输
func (excelUtil *ExcelUtil) AssembleDataDef() ([]map[string]string, error) {
	return excelUtil.AssembleData(func(headers, keys, row []string) (map[string]string, error) {
		m := make(map[string]string)
		for i, _ := range headers {
			m[fmt.Sprintf("key%d", i)] = row[i]
		}
		return m, nil
	})
}

// AssembleDataKey 使用默认方法进行key传输
func (excelUtil *ExcelUtil) AssembleDataKey() ([]map[string]string, error) {
	return excelUtil.AssembleData(func(headers, keys, row []string) (map[string]string, error) {
		m := make(map[string]string)
		for i, _ := range headers {
			m[keys[i]] = row[i]
		}
		return m, nil
	})
}
