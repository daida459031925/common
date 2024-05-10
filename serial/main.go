package serial

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"strings"
)

// 可以调用剪辑的串口并且可以开锁
func main() {
	//origData := []byte("Hello World") // 待加密的数据
	//key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	//log.Println("原文：", string(origData))

	//enc, _ := gaes.Encrypt(origData, key)

	//log.Println("------------------ CBC模式 --------------------")
	//encrypted := AesEncryptCBC(origData, key)
	//log.Println("密文(hex)：", string(enc))
	//log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	//log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	//decrypted := AesDecryptCBC(encrypted, key)
	//log.Println("解密结果：", string(decrypted))
	//dec, _ := gaes.DecryptCBC(encrypted, key)
	//log.Println(AesDecryptCBC(enc, key))
	//log.Println(dec)
	//
	//log.Println("------------------ ECB模式 --------------------")
	//encrypted = AesEncryptECB(origData, key)
	//log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	//log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	//decrypted = AesDecryptECB(encrypted, key)
	//log.Println("解密结果：", string(decrypted))
	//
	//log.Println("------------------ CFB模式 --------------------")
	//encrypted = AesEncryptCFB(origData, key)
	//log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	//log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	//decrypted = AesDecryptCFB(encrypted, key)
	//log.Println("解密结果：", string(decrypted))
	origData, e := randomBytes(16)
	if e != nil {
		log.Fatal(e)
	}
	origData[0] = 0xA5
	origData[1] = 0x5A
	origData[2] = 0x14
	origData[3] = 0x04
	origData[4] = 0x19

	keys := []byte{0x23, 0x59, 0x74, 0x03, 0x53, 0x6c, 0x1d, 0x4b, 0x33, 0x5e, 0x09, 0x18, 0x7a, 0x62, 0x71, 0x2c}

	// 配置串口参数
	options5 := serial.OpenOptions{
		PortName:              getUartPath("TTYS0"),
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       4,
		InterCharacterTimeout: 2 * 1000, // 设置超时时间
		//FlowControl:     serial.HardwareFlowControl,
	}

	//options6 := serial.OpenOptions{
	//	PortName:        "COM6",
	//	BaudRate:        9600,
	//	DataBits:        8,
	//	StopBits:        1,
	//	MinimumReadSize: 4,
	//	//FlowControl:     serial.HardwareFlowControl,
	//}

	// 打开串口
	port5, err5 := serial.Open(options5)
	if err5 != nil {
		log.Fatal(err5)
	}

	encrypted := AesEncryptECB(origData, keys)

	//port6, err6 := serial.Open(options6)
	//if err6 != nil {
	//	log.Fatal(err6)
	//}

	// 关闭串口
	defer func() {
		port5.Close()
		//port6.Close()
	}()

	//for {
	// 发送数据
	buf := make([]byte, 128)
	n, err := port5.Write(encrypted)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent %d bytes: %v", n, buf[:n])
	//接收数据
	//buf = make([]byte, 128)
	//n, err = port6.Read(buf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if n == 0 {
	//	//continue
	//}
	//
	log.Println("解密结果：", string(AesDecryptECB(encrypted, keys)))
	//log.Printf("Received %d bytes: %v", n, buf[:n])

	//}

}

// 生成的随机 byte 数组的长度
func randomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// 将 uart 转换为大写，以便进行比较
func getUartPath(uart string) string {
	uart = strings.ToUpper(uart)

	// 根据不同的 uart 名称返回对应的路径
	switch uart {
	case "TTYS0":
		return "/dev/ttyS0"
	case "TTYS1":
		return "/dev/ttyS1"
	case "TTYS2":
		return "/dev/ttyS2"
	case "TTYS3":
		return "/dev/ttyS3"
	case "TTYS4":
		return "/dev/ttyS4"
	case "TTYS5":
		return "/dev/ttyS5"
	default:
		return ""
	}
}

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// =================== CFB ======================
func AesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func AesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
