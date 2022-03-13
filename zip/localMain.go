package zip

func LocalZip(base64PdfContent string) {
	//创建文件夹
	//floder := time.Now().Unix()
	//path := fmt.Sprintf("../../../output/%d", floder)
	//err := os.Mkdir(path, 0777)
	//if err != nil {
	//	fmt.Println("创建文件夹错误")
	//}
	////将pdf写到文件，base64PdfContent为pdf文件的base64编码
	//decodeBytes, err := base64.StdEncoding.DecodeString(base64PdfContent)
	//if err != nil {
	//	fmt.Println("\nDecodeString err: ", err.Error())
	//}
	//path = fmt.Sprintf("../../../output/%d/my_pdf.pdf", floder)
	//err1 := ioutil.WriteFile(path, decodeBytes, 0666)
	//if err1 != nil {
	//	fmt.Println("存储到本地pdf错误")
	//}
	//decodeBytes2, err2 := base64.StdEncoding.DecodeString(base64XmlContent)
	//if err2 != nil {
	//	fmt.Println("\nDecodeString err: ", err2.Error())
	//}
	////将xml写到文件
	//path = fmt.Sprintf("../../../output/%d/my_xml.xml", floder)
	//err3 := ioutil.WriteFile(path, decodeBytes2, 0666)
	//if err3 != nil {
	//	fmt.Println("存储到本地xml错误")
	//}
	//
	////将response写到文件
	//decodeBytes3, err4 := base64.StdEncoding.DecodeString(base64ResponseContent)
	//if err4 != nil {
	//	fmt.Println("\nDecodeString err: ", err4.Error())
	//}
	////写到本地
	//path = fmt.Sprintf("../../../output/%d/my_responsexml.xml", floder)
	//err5 := ioutil.WriteFile(path, decodeBytes3, 0666)
	//if err5 != nil {
	//	fmt.Println("存储到本地response错误")
	//}
	////将文件夹打包压缩
	////获取源文件列表
	//target := fmt.Sprintf("../../../output/%d/", floder)
	//f, err := ioutil.ReadDir(target)
	//if err != nil {
	//	fmt.Println("压缩目录读取错误")
	//}
	////buf := new(bytes.Buffer)
	//
	//zipTarget := fmt.Sprintf("../../../output/%d/myzip.zip", floder)
	//fZip, _ := os.Create(zipTarget)
	//w := zip.NewWriter(fZip)
	//defer w.Close()
	//for _, file := range f {
	//	fw, err := w.Create(file.Name())
	//	if err != nil {
	//		fmt.Println("创建失败", fw)
	//	}
	//
	//	fileContent, err := ioutil.ReadFile(target + file.Name())
	//	if err != nil {
	//		fmt.Println("读取文件错误")
	//	}
	//	_, err = fw.Write(fileContent)
	//	if err != nil {
	//		fmt.Println("写文件错误")
	//	}
	//}
}
