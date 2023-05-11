package office

import (
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/office/excel"
	"testing"
)

func TestExcel(t *testing.T) {
	// 创建ExcelUtil实例
	excelUtil := excel.NewExcelUtil()

	//设置
	headerRow := 0
	keyRow := 1

	// 解析Excel文件
	e := excelUtil.ParseFile("C:/Users/daida/Desktop/大型仪器设备字段导入到正式环境的数据表.xlsx", "Sheet1", &headerRow, &keyRow)
	if e != nil {
		fmt.Println(e)
		return
	}

	// 组装数据并打印结果
	data, e := excelUtil.AssembleDataKey()
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(data)

	// 组装数据并打印结果
	data, e = excelUtil.AssembleDataDef()

	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(data)
}
