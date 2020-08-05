package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go/doc"
	"hash"
	"hash/crc32"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func IsHas(i interface{}) bool {
	if i != nil {
		return true
	}
	return false
}

//func JwtToken(signingKey string, claims map[string]interface{}) string {
//	token := jwt.New(jwt.SigningMethodHS256)
//
//	token.Header["typ"] = "JWT"   //加密方式
//	token.Header["alg"] = "HS256" //加密类型
//	if claims == nil {
//		claims = map[string]interface{}{}
//	}
//	if _, ok := claims["exp"]; !ok {
//		claims["exp"] = time.Now().Add(time.Hour * 96).Unix() //过期时间
//	}
//	token.Claims = claims
//	t, _ := token.SignedString([]byte(signingKey))
//	return t
//}

const (
	PwSaltBytes = 32
	PwHashBytes = 64
)
const (
	//regularMobile = `^(13[0-9]|14[57]|15[0-35-9]|18[06-9])\\d{8}$`
	regularMobile = `^(13[0-9]{9})|(18[0-9]{9})|(14[0-9]{9})|(17[0-9]{9})|(15[0-9]{9})$`
	regularIdCard = `^(\d{17})([0-9]|X)$`
	//regularEmail      = `[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	//regularEmail   = `^[a-z0-9._%+-]+@([a-z0-9-]+\\.)+[a-z]+$`
	regularEmail   = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	regularPGP     = `^[0-9A-Fa-f]{8}{0,2}$`
	regularChinese = `^[\\x{4e00}-\\x{9fa5}]+$`
	// 需要修改
	regularPasswd = `^[0-9A-Za-z]{8,}$`
	//regularPasswd = `((?=.*[a-z])(?=.*\d)|(?=[a-z])(?=.*[#@!~%^&*])|(?=.*\d)(?=.*[#@!~%^&*]))[a-z\d#@!~%^&*]{8,30}`
	//regularIPV4RE     = `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	//regularIPV4CIDRRE = `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\/(?:3[0-2]|[1-2][0-9]|[0-9])`
	regularIPV4CIDRRE = `^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$`
	regularIP         = `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
	regularIPV6RE     = `s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*`
	regularIPV6CIDRRE = `s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*\/(12[0-8]|1[0-1][0-9]|[1-9][0-9]|[0-9])`

	regularAlphaNumberic  = `^[0-9]+$`
	regularDataURI        = `^data:.+\\/(.+);base64$`
	regularBase64         = `^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`
	regularWhiteSpace     = `\\s+`
	regularUUID           = `(?i)^[A-F0-9]{8}-[A-F0-9]{4}-[A-F0-9]{4}-[A-F0-9]{4}-[A-F0-9]{12}$`
	regularYYYY           = `^\d{4}$`
	regularYYYYMM         = `^\b[1-3]\d{3}(0[1-9]|1[0-2])$`
	regularYYYYMMDD       = `^(\d{4})(\d{2})(\d{2})$`
	regularYYYYMMDDHHMMSS = `^((?:19|20)\d\d)(0[1-9]|1[012])(0[1-9]|[12][0-9]|3[01])(0\d|1\d|2[0-3])(0\d|[1-5]\d)(0\d|[1-5]\d)$`
	regularZO             = `^[01]$`
	regularZ012           = `^[012]$`
	regularZO2            = `^[02]$`
	regularZ12            = `^[12]$`
	regularZPN01          = `^[-+]?[01]$`
	regular4AND25         = `^[a-z0-9A-Z\s\p{Han}]{4,25}$`
	regular1AND30         = `^[a-z0-9A-Z\s\p{Han}]{1,30}$`
	regularZ123           = `^[123]$`
	regularZ234           = `^[234]$`
	regularWechat         = `^[a-zA-Z]([-_a-zA-Z0-9]{5,19})+$`
	regularQQ             = `^[1-9][0-9]{4,10}$`
	regular4AND24         = `^[a-z0-9A-Z\s\p{Han}]{4,24}$`
	regular12AND24        = `^[a-z0-9A-Z\s\p{Han}]{12,24}$`
	regular24AND10000     = `^[a-z0-9A-Z\s\p{Han}]{24,8000}$`
)

// Salt 生成加密盐
func Salt() []byte {
	salt := make([]byte, PwSaltBytes)
	_, err := io.ReadFull(crand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	return salt
}

// Hash 密码加密
/*
func Hash(salt, password []byte) []byte {
	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PwHashBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x\n", hash)
	return hash
}
*/

// RandNum return length 6
/*func RandNum() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100000)
}*/
func StructToMap(i interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{})
	vt := reflect.TypeOf(i)
	vv := reflect.ValueOf(i)
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		m[f.Name] = vv.FieldByName(f.Name).String()
	}
	return
}

func FailWith(status int, message string, c echo.Context) error {
	result, _ := json.Marshal(map[string]interface{}{
		"success": false,
		"reason":  message,
	})

	return c.String(status, string(result))
}

// GetRequestBody 从echo.Context获取body
func GetRequestBody(c echo.Context) (string, error) {
	bodyCache := c.Get("requestBody")
	if bodyCache != nil {
		return bodyCache.(string), nil
	}
	//body := c.Request().Body
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return "", err
	}
	c.Set("requestBody", string(b))
	return string(b), nil
}

func GetRequestBodyBytes(c echo.Context) ([]byte, error) {
	bodyCache := c.Get("requestBody")
	if bodyCache != nil {
		return bodyCache.([]byte), nil
	}
	body := c.Request().Body
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	c.Set("requestBody", b)
	return b, nil
}

func GetRequestJSON(s interface{}, c echo.Context) error {
	body, err := GetRequestBody(c)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), s)
	if err != nil {
		return err
	}
	return nil
}

// EncodeMD5 encodes string to md5 hex value.
func EncodeMD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// Encode string to sha1 hex value.
func EncodeSha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ShortSha(sha1 string) string {
	if len(sha1) == 40 {
		return sha1[:10]
	}
	return sha1
}

func BasicAuthDecode(encoded string) (string, string, error) {
	s, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	auth := strings.SplitN(string(s), ":", 2)
	return auth[0], auth[1], nil
}

func BasicAuthEncode(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

// GetRandomString generate random string by specify chars.
func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(bytes)
}

// http://code.google.com/p/go/source/browse/pbkdf2/pbkdf2.go?repo=crypto
// FIXME: use https://godoc.org/golang.org/x/crypto/pbkdf2?
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

const TimeLimitCodeLength = 12 + 6 + 40

// HashEmail hashes email address to MD5 string.
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	h := md5.New()
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 second"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 minute"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 hour"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 day"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 week"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 month"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 year"
	default:
		diffStr = fmt.Sprintf("%d years", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

// TimeSincePro calculates the time interval and generate full user-friendly string.
func TimeSincePro(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

const (
	Byte  = 1
	KByte = Byte * 1024
	MByte = KByte * 1024
	GByte = MByte * 1024
	TByte = GByte * 1024
	PByte = TByte * 1024
	EByte = PByte * 1024
)

var bytesSizeTable = map[string]uint64{
	"b":  Byte,
	"kb": KByte,
	"mb": MByte,
	"gb": GByte,
	"tb": TByte,
	"pb": PByte,
	"eb": EByte,
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+"%s", val, suffix)
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

// Subtract deals with subtraction of all types of number.
func Subtract(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool = true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	if isInt {
		return rleft - rright
	} else {
		return fleft + float64(rleft) - (fright + float64(rright))
	}
}

// EllipsisString returns a truncated short string,
// it appends '...' in the end of the length of string is too large.
func EllipsisString(str string, length int) string {
	if len(str) < length {
		return str
	}
	return str[:length-3] + "..."
}

// TruncateString returns a truncated string with given limit,
// it returns input string if length is not reached limit.
func TruncateString(str string, limit int) string {
	if len(str) < limit {
		return str
	}
	return str[:limit]
}

// Int64sToMap converts a slice of int64 to a int64 map.
func Int64sToMap(ints []int64) map[int64]bool {
	m := make(map[int64]bool)
	for _, i := range ints {
		m[i] = true
	}
	return m
}

// IsLetter reports whether the rune is a letter (category L).
// https://github.com/golang/go/blob/master/src/go/scanner/scanner.go#L257
func IsLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

func IsTextFile(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	if strings.Index(contentType, "text/") != -1 {
		return contentType, true
	}
	return contentType, false
}

func IsImageFile(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	if strings.Index(contentType, "image/") != -1 {
		return contentType, true
	}
	return contentType, false
}

func IsPDFFile(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	if strings.Index(contentType, "application/pdf") != -1 {
		return contentType, true
	}
	return contentType, false
}

// MatchMobile 匹配手机号
func MatchMobile(mobileNum string) bool {
	reg := regexp.MustCompile(regularMobile)
	return reg.MatchString(mobileNum)
}

// MatchIdCard 匹配身份证号
func MatchIdCard(idCard string) bool {
	reg := regexp.MustCompile(regularIdCard)
	return reg.MatchString(idCard)
}

// MatchEmail 匹配电子邮箱
func MatchEmail(email string) bool {
	reg := regexp.MustCompile(regularEmail)
	return reg.MatchString(email)
}

// MatchChinese 匹配中文
func MatchChinese(chinese string) bool {
	reg := regexp.MustCompile(regularChinese)
	return reg.MatchString(chinese)
}

// MatchPwdLength 匹配密码长度8
func MatchPwdLength(str string) bool {
	reg := regexp.MustCompile(regularPasswd)
	return reg.MatchString(str)
}

// MatchIPV4 匹配ip地址
func MatchIPV4(ip string) bool {
	reg := regexp.MustCompile(regularIPV4CIDRRE)
	return reg.MatchString(ip)
}

// MatchAlphaNumberic 匹配是否是数字
func MatchAlphaNumberic(number string) bool {
	reg := regexp.MustCompile(regularAlphaNumberic)
	return reg.MatchString(number)
}

// MatchUUID 匹配UUID
func MatchUUID(str string) bool {
	reg := regexp.MustCompile(regularUUID)
	return reg.MatchString(str)
}

//  MatchYYYYMM 匹配年月
func MatchYYYYMM(str string) bool {
	reg := regexp.MustCompile(regularYYYYMM)
	return reg.MatchString(str)
}

// MatchYYYYMMDD 匹配年月日
func MatchYYYYMMDD(str string) bool {
	reg := regexp.MustCompile(regularYYYYMMDD)
	return reg.MatchString(str)
}

// MatchYYYYMMDDHHMMSS 匹配年月日时分秒
func MatchYYYYMMDDHHMMSS(str string) bool {
	reg := regexp.MustCompile(regularYYYYMMDDHHMMSS)
	return reg.MatchString(str)
}

// MatchYYYY 匹配年份
func MatchYYYY(str string) bool {
	reg := regexp.MustCompile(regularYYYY)
	return reg.MatchString(str)
}

// MatchZO 匹配0 或者1
func MatchZO(str string) bool {
	reg := regexp.MustCompile(regularZO)
	return reg.MatchString(str)
}

// MatchZ012 匹0或者1或者2
func MatchZ012(str string) bool {
	reg := regexp.MustCompile(regularZ012)
	return reg.MatchString(str)
}

func MatchZ02(str string) bool {
	reg := regexp.MustCompile(regularZO2)
	return reg.MatchString(str)
}

// MatchZ12 匹1或者2
func MatchZ12(str string) bool {
	reg := regexp.MustCompile(regularZ12)
	return reg.MatchString(str)
}

// MatchZPN01 匹配 -1 0 1
func MatchZPN01(str string) bool {
	reg := regexp.MustCompile(regularZPN01)
	return reg.MatchString(str)
}

// MatchZ123 匹配 1 2 3 任意一个字符串
func MatchZ123(str string) bool {
	reg := regexp.MustCompile(regularZ123)
	return reg.MatchString(str)
}

// MatchZ123 匹配 2 3 4 任意一个字符串
func MatchZ234(str string) bool {
	reg := regexp.MustCompile(regularZ234)
	return reg.MatchString(str)
}

// Match4And25 配中文，英文字母和数字及_:  4-25
func Match4And25(str string) bool {
	reg := regexp.MustCompile(regular4AND25)
	return reg.MatchString(str)
}

// Match4And24 配中文，英文字母和数字及_:  4-24
func Match4And24(str string) bool {
	reg := regexp.MustCompile(regular4AND24)
	return reg.MatchString(str)
}

// Match12And24 配中文，英文字母和数字及_:  12-24
func Match12And24(str string) bool {
	reg := regexp.MustCompile(regular12AND24)
	return reg.MatchString(str)
}

func Match24And10000(str string) bool {
	reg := regexp.MustCompile(regular24AND10000)
	return reg.MatchString(str)
}

// MatchWechat 微信号正则，6至20位，以字母开头，字母，数字，减号，下划线
func MatchWechat(str string) bool {
	reg := regexp.MustCompile(regularWechat)
	return reg.MatchString(str)
}

// MatchQQ QQ号正则，5至11位
func MatchQQ(str string) bool {
	reg := regexp.MustCompile(regularQQ)
	return reg.MatchString(str)
}

// Match1And30 配中文，英文字母和数字及_:  1-30
func Match1And30(str string) bool {
	reg := regexp.MustCompile(regular1AND30)
	return reg.MatchString(str)

}

func Subexp(exp *regexp.Regexp, matches []string, subexp string) string {
	for index, name := range exp.SubexpNames() {
		if index >= len(matches) {
			continue
		}

		if name == subexp {
			return matches[index]
		}
	}

	return ""
}

func MatchString(exp string, data string) ([]string, error) {
	re, err := regexp.Compile(exp)
	if err != nil {
		return nil, err
	}

	return re.FindStringSubmatch(data), nil
}

func Match(exp string, data []byte) ([][]byte, error) {
	re, err := regexp.Compile(exp)
	if err != nil {
		return nil, err
	}

	return re.FindSubmatch(data), nil
}

/*

func EncodePassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}
*/

// IsEmpty 判断一个值是否为空(0, "", false, 空数组等)。
// []string{""}空数组里套一个空字符串，不会被判断为空。
func IsEmptys(expr interface{}) bool {
	if expr == nil {
		return true
	}

	switch v := expr.(type) {
	case bool:
		return !v
	case int:
		return 0 == v
	case int8:
		return 0 == v
	case int16:
		return 0 == v
	case int32:
		return 0 == v
	case int64:
		return 0 == v
	case uint:
		return 0 == v
	case uint8:
		return 0 == v
	case uint16:
		return 0 == v
	case uint32:
		return 0 == v
	case uint64:
		return 0 == v
	case string:
		return len(v) == 0
	case float32:
		return 0 == v
	case float64:
		return 0 == v
	case time.Time:
		return v.IsZero()
	case *time.Time:
		return v.IsZero()
	}

	// 符合IsNil条件的，都为Empty
	if IsNil(expr) {
		return true
	}

	// 长度为0的数组也是empty
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return 0 == v.Len()
	}

	return false
}

// The IsNum judges string is number or not.
func IsNum(a string) bool {
	reg, _ := regexp.Compile("^\\d+$")
	return reg.MatchString(a)
}

// DeleteWhitespace returns a version of the passed-in string with all
// whitespaces (as by unicode.IsSpace) removed.
func DeleteWhitespace(s string) string {
	if len(s) == 0 {
		return s
	}

	var hasSpace bool
	var buf bytes.Buffer
	for _, r := range s {
		if unicode.IsSpace(r) {
			hasSpace = true
		} else {
			buf.WriteRune(r)
		}
	}

	if !hasSpace {
		return s
	}
	return buf.String()
}

// ShrinkWhitespace returns a version of the passed-in string with whitespace
// removed from the end and beginning as well as all whitespace in-between
// shrunken to a single space.
func ShrinkWhitespace(s string) string {
	reg := regexp.MustCompile(regularWhiteSpace)
	return reg.ReplaceAllString(strings.TrimSpace(s), " ")
}

// IndexAfter returns the index directly after the first instance of sep in s,
// or -1 if sep is not present in s.
func IndexAfter(s, sep string) int {
	pos := strings.Index(s, sep)
	if pos == -1 {
		return -1
	}
	return pos + len(sep)
}

// Reverse returns a reversed version of s.
func Reverse(s string) (rev string) {
	for _, r := range s {
		rev = string(r) + rev
	}
	return rev
}

// Lines returns sequence of lines of passed-in str.
func Lines(s string) []string {
	return strings.Split(s, "\n")
}

// SplitCamelCase splits the string s at each run of upper case runes and
// returns and array of slices of s.
//func SplitCamelCase(s string) (words []string) {
//	fieldStart := 0
//	for i, r := range s {
//		if i != 0 && unicode.IsUpper(r) {
//			words = append(words, s[fieldStart:i])
//			fieldStart = i
//		}
//	}
//	if fieldStart != len(s) {
//		words = append(words, s[fieldStart:])
//	}
//	return words
//}

// IsBlank returns whether a string is whitespace or empty.
func IsBlank(s string) bool {
	strLen := len(s)
	if strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(s[i])) == false {
			return false
		}
	}
	return true
}

// IsEmpty returns whether a string is empty ("").
func IsEmpty(s string) bool {
	return len(s) == 0
}

// WordWrap wraps paragraphs of text to width.
func WordWrap(s string, width int) string {
	buf := new(bytes.Buffer)
	doc.ToText(buf, s, "", "", width)
	return buf.String()
}

// Surround appends and prepends a string to another.
func Surround(s, around string) string {
	return around + s + around
}

// IsNil 判断一个值是否为 nil。
// 当特定类型的变量，已经声明，但还未赋值时，也将返回 true
func IsNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	return (k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Ptr ||
		k == reflect.Slice) &&
		v.IsNil()
}

// IsEqual 判断两个值是否相等。
//
// 除了通过 reflect.DeepEqual() 判断值是否相等之外，一些类似
// 可转换的数值也能正确判断，比如以下值也将会被判断为相等：
//  int8(5)                     == int(5)
//  []int{1,2}                  == []int8{1,2}
//  []int{1,2}                  == [2]int8{1,2}
//  []int{1,2}                  == []float32{1,2}
//  map[string]int{"1":"2":2}   == map[string]int8{"1":1,"2":2}
//
//  // map的键值不同，即使可相互转换也判断不相等。
//  map[int]int{1:1,2:2}        != map[int8]int{1:1,2:2}
func IsEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	vv1 := reflect.ValueOf(v1)
	vv2 := reflect.ValueOf(v2)

	// NOTE: 这里返回false，而不是true
	if !vv1.IsValid() || !vv2.IsValid() {
		return false
	}

	if vv1 == vv2 {
		return true
	}

	vv1Type := vv1.Type()
	vv2Type := vv2.Type()

	// 过滤掉已经在reflect.DeepEqual()进行处理的类型
	switch vv1Type.Kind() {
	case reflect.Struct, reflect.Ptr, reflect.Func, reflect.Interface:
		return false
	case reflect.Slice, reflect.Array:
		// vv2.Kind()与vv1的不相同
		if vv2.Kind() != reflect.Slice && vv2.Kind() != reflect.Array {
			// 虽然类型不同，但可以相互转换成vv1的，如：vv2是string，vv2是[]byte，
			if vv2Type.ConvertibleTo(vv1Type) {
				return IsEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
			}
			return false
		}

		// reflect.DeepEqual()未考虑类型不同但是类型可转换的情况，比如：
		// []int{8,9} == []int8{8,9}，此处重新对slice和array做比较处理。
		if vv1.Len() != vv2.Len() {
			return false
		}

		for i := 0; i < vv1.Len(); i++ {
			if !IsEqual(vv1.Index(i).Interface(), vv2.Index(i).Interface()) {
				return false
			}
		}
		return true // for中所有的值比较都相等，返回true
	case reflect.Map:
		if vv2.Kind() != reflect.Map {
			return false
		}

		if vv1.IsNil() != vv2.IsNil() {
			return false
		}
		if vv1.Len() != vv2.Len() {
			return false
		}
		if vv1.Pointer() == vv2.Pointer() {
			return true
		}

		// 两个map的键名类型不同
		if vv2Type.Key().Kind() != vv1Type.Key().Kind() {
			return false
		}

		for _, index := range vv1.MapKeys() {
			vv2Index := vv2.MapIndex(index)
			if !vv2Index.IsValid() {
				return false
			}

			if !IsEqual(vv1.MapIndex(index).Interface(), vv2Index.Interface()) {
				return false
			}
		}
		return true // for中所有的值比较都相等，返回true
	case reflect.String:
		if vv2.Kind() == reflect.String {
			return vv1.String() == vv2.String()
		}
		if vv2Type.ConvertibleTo(vv1Type) { // 考虑v1是string，v2是[]byte的情况
			return IsEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
		}

		return false
	}

	if vv1Type.ConvertibleTo(vv2Type) {
		return vv2.Interface() == vv1.Convert(vv2Type).Interface()
	} else if vv2Type.ConvertibleTo(vv1Type) {
		return vv1.Interface() == vv2.Convert(vv1Type).Interface()
	}

	return false
}

// HasPanic 判断 fn 函数是否会发生 panic
// 若发生了 panic，将把 msg 一起返回。
func HasPanic(fn func()) (has bool, msg interface{}) {
	defer func() {
		if msg = recover(); msg != nil {
			has = true
		}
	}()
	fn()

	return
}

// IsContains 判断 container 是否包含了 item 的内容。若是指针，会判断指针指向的内容，
// 但是不支持多重指针。
//
// 若 container 是字符串(string、[]byte和[]rune，不包含 fmt.Stringer 接口)，
// 都将会以字符串的形式判断其是否包含 item。
// 若 container是个列表(array、slice、map)则判断其元素中是否包含 item 中的
// 的所有项，或是 item 本身就是 container 中的一个元素。
func IsContains(container, item interface{}) bool {
	if container == nil { // nil不包含任何东西
		return false
	}

	cv := reflect.ValueOf(container)
	iv := reflect.ValueOf(item)

	if cv.Kind() == reflect.Ptr {
		cv = cv.Elem()
	}

	if iv.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	if IsEqual(container, item) {
		return true
	}

	// 判断是字符串的情况
	switch c := cv.Interface().(type) {
	case string:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(c, i)
		case []byte:
			return strings.Contains(c, string(i))
		case []rune:
			return strings.Contains(c, string(i))
		case byte:
			return bytes.IndexByte([]byte(c), i) != -1
		case rune:
			return bytes.IndexRune([]byte(c), i) != -1
		}
	case []byte:
		switch i := iv.Interface().(type) {
		case string:
			return bytes.Contains(c, []byte(i))
		case []byte:
			return bytes.Contains(c, i)
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return bytes.IndexByte(c, i) != -1
		case rune:
			return bytes.IndexRune(c, i) != -1
		}
	case []rune:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(string(c), string(i))
		case []byte:
			return strings.Contains(string(c), string(i))
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return strings.IndexByte(string(c), i) != -1
		case rune:
			return strings.IndexRune(string(c), i) != -1
		}
	}

	if (cv.Kind() == reflect.Slice) || (cv.Kind() == reflect.Array) {
		if !cv.IsValid() || cv.Len() == 0 { // 空的，就不算包含另一个，即使另一个也是空值。
			return false
		}

		if !iv.IsValid() {
			return false
		}

		// item 是 container 的一个元素
		for i := 0; i < cv.Len(); i++ {
			if IsEqual(cv.Index(i).Interface(), iv.Interface()) {
				return true
			}
		}

		// 开始判断 item 的元素是否与 container 中的元素相等。

		// 若 item 的长度为 0，表示不包含
		if (iv.Kind() != reflect.Slice) || (iv.Len() == 0) {
			return false
		}

		// item 的元素比 container 的元素多，必须在判断完 item 不是 container 中的一个元素之
		if iv.Len() > cv.Len() {
			return false
		}

		// 依次比较 item 的各个子元素是否都存在于 container，且下标都相同
		ivIndex := 0
		for i := 0; i < cv.Len(); i++ {
			if IsEqual(cv.Index(i).Interface(), iv.Index(ivIndex).Interface()) {
				if (ivIndex == 0) && (i+iv.Len() > cv.Len()) {
					return false
				}
				ivIndex++
				if ivIndex == iv.Len() { // 已经遍历完 iv
					return true
				}
			} else if ivIndex > 0 {
				return false
			}
		}
		return false
	} // end cv.Kind == reflect.Slice and reflect.Array

	if cv.Kind() == reflect.Map {
		if cv.Len() == 0 {
			return false
		}

		if (iv.Kind() != reflect.Map) || (iv.Len() == 0) {
			return false
		}

		if iv.Len() > cv.Len() {
			return false
		}

		// 判断所有 item 的项都存在于 container 中
		for _, key := range iv.MapKeys() {
			cvItem := iv.MapIndex(key)
			if !cvItem.IsValid() { // container 中不包含该值。
				return false
			}
			if !IsEqual(cvItem.Interface(), iv.MapIndex(key).Interface()) {
				return false
			}
		}
		// for 中的所有判断都成立，返回 true
		return true
	}

	return false
}

// 判断一个值是否可转换为数值。不支持全角数值的判断。
func Number(val interface{}) bool {
	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case []byte:
		_, err := strconv.ParseFloat(string(v), 32)
		return err == nil
	case string:
		_, err := strconv.ParseFloat(v, 32)
		return err == nil
	case []rune:
		_, err := strconv.ParseFloat(string(v), 32)
		return err == nil
	default:
		return false
	}
}

func RuneToString(r rune) string {
	if r >= 0x20 && r < 0x7f {
		return fmt.Sprintf("'%c'", r)
	}
	switch r {
	case 0x07:
		return "'\\a'"
	case 0x08:
		return "'\\b'"
	case 0x0C:
		return "'\\f'"
	case 0x0A:
		return "'\\n'"
	case 0x0D:
		return "'\\r'"
	case 0x09:
		return "'\\t'"
	case 0x0b:
		return "'\\v'"
	case 0x5c:
		return "'\\\\\\'"
	case 0x27:
		return "'\\''"
	case 0x22:
		return "'\\\"'"
	}
	if r < 0x10000 {
		return fmt.Sprintf("\\u%04x", r)
	}
	return fmt.Sprintf("\\U%08x", r)
}

/* Interface */

/*
Convert the literal value of a scanned token to rune
*/
func RuneValue(lit []byte) rune {
	if lit[1] == '\\' {
		return escapeCharVal(lit)
	}
	r, size := utf8.DecodeRune(lit[1:])
	if size != len(lit)-2 {
		panic(fmt.Sprintf("Error decoding rune. Lit: %s, rune: %d, size%d\n", lit, r, size))
	}
	return r
}

/*
Convert the literal value of a scanned token to int64
*/
func IntValue(lit []byte) (int64, error) {
	return strconv.ParseInt(string(lit), 10, 64)
}

/*
Convert the literal value of a scanned token to uint64
*/
func UintValue(lit []byte) (uint64, error) {
	return strconv.ParseUint(string(lit), 10, 64)
}

/* Util */

func escapeCharVal(lit []byte) rune {
	var i, base, max uint32
	offset := 2
	switch lit[offset] {
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'v':
		return '\v'
	case '\\':
		return '\\'
	case '\'':
		return '\''
	case '0', '1', '2', '3', '4', '5', '6', '7':
		i, base, max = 3, 8, 255
	case 'x':
		i, base, max = 2, 16, 255
		offset++
	case 'u':
		i, base, max = 4, 16, unicode.MaxRune
		offset++
	case 'U':
		i, base, max = 8, 16, unicode.MaxRune
		offset++
	default:
		panic(fmt.Sprintf("Error decoding character literal: %s\n", lit))
	}

	var x uint32
	for ; i > 0 && offset < len(lit)-1; i-- {
		ch, size := utf8.DecodeRune(lit[offset:])
		offset += size
		d := uint32(digitVal(ch))
		if d >= base {
			panic(fmt.Sprintf("charVal(%s): illegal character (%c) in escape sequence. size=%d, offset=%d", lit, ch, size, offset))
		}
		x = x*base + d
	}
	if x > max || 0xD800 <= x && x < 0xE000 {
		panic(fmt.Sprintf("Error decoding escape char value. Lit:%s, offset:%d, escape sequence is invalid Unicode code point\n", lit, offset))
	}

	return rune(x)
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch) - '0'
	case 'a' <= ch && ch <= 'f':
		return int(ch) - 'a' + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch) - 'A' + 10
	}
	return 16 // larger than any legal digit val
}

// String int32 to string
// see https://stackoverflow.com/questions/39442167/convert-int32-to-string-in-golang
func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

//
func FormatInt32(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

// RemoveNewLine 移除新行
func RemoveNewLine(value string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(value, "")
}

// JsonpToJson modify jsonp string to json string
// Example: forbar({a:"1",b:2}) to {"a":"1","b":2}
func JsonpToJson(json string) string {
	start := strings.Index(json, "{")
	end := strings.LastIndex(json, "}")
	start1 := strings.Index(json, "[")
	if start1 > 0 && start > start1 {
		start = start1
		end = strings.LastIndex(json, "]")
	}
	if end > start && end != -1 && start != -1 {
		json = json[start : end+1]
	}
	json = strings.Replace(json, "\\'", "", -1)
	regDetail, _ := regexp.Compile("([^\\s\\:\\{\\,\\d\"]+|[a-z][a-z\\d]*)\\s*\\:")
	return regDetail.ReplaceAllString(json, "\"$1\":")
}

// The GetWDPath gets the work directory path.
func GetWDPath() string {
	wd := os.Getenv("GOPATH")
	if wd == "" {
		panic("GOPATH is not setted in env.")
	}
	return wd
}

// The IsDirExists judges path is directory or not.
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

// The IsFileExists judges path is file or not.
func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// simple xml to string  support utf8
/*
func XML2mapstr(xmldoc string) map[string]string {
	var t xml.Token
	var err error
	inputReader := strings.NewReader(xmldoc)
	decoder := xml.NewDecoder(inputReader)
	decoder.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		return charset.NewReader(r, s)
	}
	m := make(map[string]string, 32)
	key := ""
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			key = token.Name.Local
		case xml.CharData:
			content := string([]byte(token))
			m[key] = content
		default:
			// ...
		}
	}

	return m
}
*/

//string to hash
func MakeHash(s string) string {
	const IEEE = 0xedb88320
	var IEEETable = crc32.MakeTable(IEEE)
	hash := fmt.Sprintf("%x", crc32.Checksum([]byte(s), IEEETable))
	return hash
}

func HashString(encode string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(encode))
	return hash.Sum64()
}

// MakeUnique 制作特征值方法一
func MakeUnique(obj interface{}) string {
	baseString, _ := json.Marshal(obj)
	return strconv.FormatUint(HashString(string(baseString)), 10)
}

// MakeMd5 制作特征值方法二
func MakeMd5(obj interface{}, length int) string {
	if length > 32 {
		length = 32
	}
	h := md5.New()
	baseString, _ := json.Marshal(obj)
	h.Write([]byte(baseString))
	s := hex.EncodeToString(h.Sum(nil))
	return s[:length]
}
