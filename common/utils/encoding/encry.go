package encoding

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

var privateKey = "QWERTYUIhgfedcba87654321"

type NonceBlock struct {
	Hash           []byte
	TimeStamp      int64
	ExpiredSeconds int64
	Nonce          int64
	Data           []byte
}

func Int322Byte(num int32) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	CheckErr(err)
	return buffer.Bytes()
}

func Int642Byte(num int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	CheckErr(err)
	return buffer.Bytes()
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("err:", err)
	}
}

func NewBlock(data string, expiredSeconds int64) *NonceBlock {
	block := &NonceBlock{
		TimeStamp:      time.Now().Unix(),
		ExpiredSeconds: expiredSeconds,
		Nonce:          time.Now().UnixNano(),
		Data:           []byte(data)}
	block.SetHash()

	return block
}

func (block *NonceBlock) SetHash() {
	tmp := [][]byte{
		//实现将int转成byte数组的函数
		Int642Byte(block.TimeStamp),
		Int642Byte(block.ExpiredSeconds),
		Int642Byte(block.Nonce),
		block.Data}
	//将区块个字段连接成一个切片，使用[]byte{}进行连接

	data := bytes.Join(tmp, []byte{})

	//算出hash的值
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

func (block *NonceBlock) ToBase64() string {
	tmp := [][]byte{
		block.Hash,
		Int642Byte(block.TimeStamp),
		Int642Byte(block.ExpiredSeconds),
		Int642Byte(block.Nonce),
		block.Data}
	data := bytes.Join(tmp, []byte{})
	encrypt, err := AesCBCEncrypt(data, []byte(privateKey))
	if err != nil {
		//log.Error(err)
		return ""
	}
	encodeString := base64.StdEncoding.EncodeToString(encrypt)

	return encodeString
}

func NewFromBase64(encodeString string) (*NonceBlock, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		return nil, err
	}
	decodeBytes, err = AesCBCDecrypt(decodeBytes, []byte(privateKey))
	if err != nil || len(decodeBytes) < 56 {
		return nil, err
	}

	var part1 = decodeBytes[0:32]
	var part2 = decodeBytes[32:40]
	var part3 = decodeBytes[40:48]
	var part4 = decodeBytes[48:56]
	var part5 = decodeBytes[56:(len(decodeBytes))]
	block := &NonceBlock{
		Hash:           part1,
		TimeStamp:      int64(binary.BigEndian.Uint64(part2)),
		ExpiredSeconds: int64(binary.BigEndian.Uint64(part3)),
		Nonce:          int64(binary.BigEndian.Uint64(part4)),
		Data:           part5}
	if (time.Now().Unix() - block.TimeStamp) > block.ExpiredSeconds {
		return nil, errors.New("请求数据已失效,请刷新后操作")
	}
	return block, nil
}

func EncryptWithExpire(in string, ex int64) string {
	return NewBlock(in, ex).ToBase64()
}

func DecryptExpireData(in string) (string, error) {
	b, err := NewFromBase64(in)
	if err != nil {
		return "", err
	}
	return string(b.Data), nil
}

func (block *NonceBlock) Check() bool {
	var retValue = false

	if block != nil && ((time.Now().Unix() - block.TimeStamp) <= block.ExpiredSeconds) {
		retValue = true
	}

	return retValue
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, len(rawData))
	//block大小 16
	iv := make([]byte, blockSize)

	//iv = []byte("1234567890123456")
	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, rawData)

	return cipherText, nil
}

func AesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		return []byte{}, errors.New("data too short")
	}
	iv := make([]byte, blockSize)

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return []byte{}, errors.New("data is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func EncryptByUid(in, key string) string {
	if len(in) == 0 {
		return ""
	}
	out, err := AesCBCEncrypt([]byte(in), PrepareKey(key))
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(out)
}

func DecryptByUid(in, key string) (string, error) {
	if len(in) == 0 {
		return "", errors.New("data is null or '' ")
	}
	decodeBytes, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}
	out, err := AesCBCDecrypt(decodeBytes, PrepareKey(key))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// 截取 32 位, key 加 uid
func PrepareKey(uid string) []byte {

	var buffer bytes.Buffer
	buffer.WriteString(privateKey)
	buffer.WriteString(uid)
	buffer.WriteString("It was the best of times,it was the worst of times,")

	return buffer.Bytes()[:32]
}





///////////////同一个字符串加密后的结果不一样（不变的key），现用于登录返回给vue，根据此解码出 userId_leagueId ///////////////////
func AesRandEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	cipherText := make([]byte, blockSize+len(rawData))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)
	return cipherText, nil
}
func AesRandDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

var key = "321ffFDs$#25GDS+!fwGf554"

func EncryptRand(rawData string) string {
	data, err := AesRandEncrypt([]byte(rawData), []byte(key))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

func DecryptRand(rawData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := AesRandDecrypt(data, []byte(key))
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

///////////////////////////////////////////////////////////////////////
