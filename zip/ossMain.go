package zip

func OssZip() {
	//创建buf用于存储zip中的文件
	//var zipBuffer *bytes.Buffer = new(bytes.Buffer)
	//var zipWriter *zip.Writer = zip.NewWriter(zipBuffer)
	//var zipEntry io.Writer
	//var err error
	//zipEntry, err = zipWriter.Create("ticket.pdf")
	//if err != nil {
	//	fmt.Println("创建pdf文件错误")
	//	return
	//}
	//decodeBytes, err1 := base64.StdEncoding.DecodeString(base64PdfContent)
	//if err1 != nil {
	//	fmt.Println("\nDecodeString err: ", err1.Error())
	//}
	//_, err = zipEntry.Write(decodeBytes)
	//if err != nil {
	//	fmt.Println("写入pdf错误")
	//}
	//zipEntry, err = zipWriter.Create("ticket.xml")
	//if err != nil {
	//	fmt.Println("创建xml文件错误")
	//	return
	//}
	//decodeBytes2, err2 := base64.StdEncoding.DecodeString(base64XmlContent)
	//if err2 != nil {
	//	fmt.Println("\nDecodeString err: ", err2.Error())
	//}
	//_, err = zipEntry.Write(decodeBytes2)
	//if err != nil {
	//	fmt.Println("写入xml错误")
	//}
	//
	//zipEntry, err = zipWriter.Create("response.xml")
	//if err != nil {
	//	fmt.Println("创建response文件错误")
	//	return
	//}
	//decodeBytes3, err3 := base64.StdEncoding.DecodeString(base64ResponseContent)
	//if err3 != nil {
	//	fmt.Println("\nDecodeString err: ", err3.Error())
	//}
	//_, err = zipEntry.Write(decodeBytes3)
	//if err != nil {
	//	fmt.Println("写入xml错误")
	//}
	//err = zipWriter.Close()
	//if err != nil {
	//	fmt.Println("关闭文件错误")
	//}
	//ctx := context.Background()
	////MD5
	////m := md5.New()
	////m.Write(zipBuffer.Bytes())
	////md5Str := hex.EncodeToString(m.Sum(nil))
	////上传oss
	//randStr := convert.RandStringBytes(20)
	//nowTs := time.Now().Unix()
	//fileName := fmt.Sprintf("upload-base64-test/%d-%s-3c28425b-88f6-0ca7-dbe5-5149db9ac8d7.zip", nowTs, randStr)
	//objectKey := GetTicketOssCountryEnvKey("CR", fileName)
	//fmt.Println("objectkey:", objectKey)
	//bufSize := int64(len(zipBuffer.Bytes()))
	//fmt.Println("文件大小", bufSize/1024)
	////fmt.Println("md5:", md5Str)
	////_, err = UpToDidiOssPublicRead(ctx, objectKey, zipBuffer, bufSize, md5Str)
	//_, err = UpToDidiOssPublicRead(ctx, objectKey, zipBuffer, bufSize)
	//if err != nil {
	//	fmt.Println("UpToDidiOssPublicRead err: ", err.Error())
	//}
	//fmt.Printf("======== %s \n", GetDidiOssCDNUrl(objectKey))
	//fmt.Println("结束！！！")
}
