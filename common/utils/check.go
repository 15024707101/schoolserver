package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//身份证号验证算法
/**
公民身份证号码按照 GB11643—1999《公民身份证号码》国家标准编制，由18位数字组成：前6位为行政区划分代码，第7位至14位为出生日期码，第15位至17位为顺序码，第18位为校验码。
在上世纪（二十世纪）办的身份证为15位数字码。原来7、8位的年份号到2000年后攺为全称，如1985年过去7、8位码是85，现在增改为1985，而又在最后一位增加校验码，如后三位原来601，加一个5成为6015。身份证一经编定不作改变，派出所会在户口资料中给你加上，你要换新证时就是18位的新码了。

18身份证号码的结构
公民身份号码是特征组合码，由十七位数字本体码和一位校验码组成。

排列顺序从左至右依次为：六位数字地址码，八位数字出生日期码，三位数字顺序码和一位校验码。

1、地址码
表示编码对象常住户口所在县(市、旗、区)的行政区域划分代码，按GB/T2260的规定执行。

2、出生日期码
表示编码对象出生的年、月、日，按GB/T7408的规定执行，年、月、日代码之间不用分隔符。

3、顺序码
表示在同一地址码所标识的区域范围内，对同年、同月、同日出生的人编定的顺序号，顺序码的奇数分配给男性，偶数分配给女性。

4、校验码计算步骤
(1)十七位数字本体码加权求和公式

S = Sum(Ai * Wi), i = 0, … , 16 ，先对前 17 位数字的权求和
eg:
```golang
int s = 0;
for(i=1;i<17;i++) {
 s +=(a[i] * w[i])
}
```
Ai：表示第i位置上的身份证号码数字值(0~9) （17位）
Wi：7 9 10 5 8 4 2 1 6 3 7 9 10 5 8 4 2 （表示第 i 位置上的加权因子）（17位）
eg:
Ai-> a1 , a2, a3, a4, a5, a6... a17 身份证前17位对应(ai) (a18 是校验码)
Wi-> 7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2 (17位)
(2)计算模

Y = mod(S, 11)
eg:
```golang
int y = s % 11
```
(3)根据模，查找得到对应的校验码

Y: 0 1 2 3 4 5 6 7 8 9 10
校验码: 1 0 X 9 8 7 6 5 4 3 2
*/

// /(^\d{15}$)|(^\d{17}(\d|X)$)/
const regularCidCard = `(^\d{15}$)|(^\d{17}(\d|X)$)`
const re_fifteen = `^(\d{6})(\d{2})(\d{2})(\d{2})(\d{3})$`

func CheckString(str string) string {
	var pattern = "~`!@#$%^*()_+{}[]:;&'<>\"|\\/?.,"
	var s = ""

	strArr := strings.Fields(str)
	//fmt.Printf("strArr: %s", strArr)
	for i := 0; i < len(strArr); i++ {
		//fmt.Printf("strArr %s\n", strArr[i])
		//fmt.Printf("Indexof %d\n", strings.Index(pattern, strArr[i]))
		if strings.Index(pattern, strArr[i]) != -1 {
			s += strArr[i]
			//fmt.Printf("s %s\n", s)
		}
	}
	return s
}

// 检查身份证号s
func CheckCardImport(str string) (bool, string) {
	// 检查身份证号尾数
	if CheckCardNO(str) == false {
		//fmt.Println("CheckCardNO 错误\n")
		return false, "请正确填写身份证号码"
	}
	if CheckBirthday(str) == false {
		//fmt.Println("CheckBirthday 错误\n")
		return false, "请正确填写身份证号码"
	}
	if CheckParity(str) == false {
		//fmt.Println("CheckParity 错误\n")
		return false, "身份证号码校验位不正确"
	}
	return true, ""
}
func CheckCard(str string) bool {
	// 检查身份证号尾数
	if CheckCardNO(str) == false {
		//fmt.Println("CheckCardNO 错误\n")
		return false
	}
	if CheckBirthday(str) == false {
		//fmt.Println("CheckBirthday 错误\n")
		return false
	}
	if CheckParity(str) == false {
		//fmt.Println("CheckParity 错误\n")
		return false
	}
	return true
}

// CheckCardNO 检查号码是否符合规范，包括长度，类型
func CheckCardNO(idcard string) bool { // //身份证号码为15位或者18位，15位时全为数字，18位前17位为数字，最后一位是校验位，可能为数字或字符X
	reg := regexp.MustCompile(regularCidCard)
	return reg.MatchString(idcard)
}

var vcity = map[string]string{
	"11": "北京",
	"12": "天津",
	"13": "河北",
	"14": "山西",
	"15": "内蒙古",
	"21": "辽宁",
	"22": "吉林",
	"23": "黑龙江",
	"31": "上海",
	"32": "江苏",
	"33": "浙江",
	"34": "安徽",
	"35": "福建",
	"36": "江西",
	"37": "山东",
	"41": "河南",
	"42": "湖北",
	"43": "湖南",
	"44": "广东",
	"45": "广西",
	"46": "海南",
	"50": "重庆",
	"51": "四川",
	"52": "贵州",
	"53": "云南",
	"54": "西藏",
	"61": "陕西",
	"62": "甘肃",
	"63": "青海",
	"64": "宁夏",
	"65": "新疆",
	"71": "台湾",
	"81": "香港",
	"82": "澳门",
	"91": "国外",
	"88": "未知",
}

// CheckProvince 检查省份是否正确
func CheckProvince(idcard string) bool {
	// 取身份证前两位,校验省份
	province := substr(idcard, 0, 2)
	if vcity[province] == "" {
		return false
	}
	return true
}

// 检查生日是否正确
func CheckBirthday(idcard string) bool {
	var len = len(idcard)
	if len == 15 {
		dateY := idcard[6:8]
		dateM := idcard[8:10]
		dateD := idcard[10:12]
		birthday := "19" + dateY + "/" + dateM + "/" + dateD
		return verifyBirthday("19"+dateY, dateM, dateD, birthday)
	}
	if len == 18 {
		dateY := idcard[6:10]
		dateM := idcard[10:12]
		dateD := idcard[12:14]
		birthday := dateY + "/" + dateM + "/" + dateD
		return verifyBirthday(dateY, dateM, dateD, birthday)
	}
	return false
}

// 验证生日
func verifyBirthday(y, m, d, birthday string) bool {
	now_year, _ := strconv.Atoi(fmt.Sprintf("%d", time.Now().Year()))
	dateY, _ := strconv.Atoi(birthday[0:4])
	dateM, _ := strconv.Atoi(birthday[5:7])
	dateD, _ := strconv.Atoi(birthday[8:10])

	pY, _ := strconv.Atoi(y)
	pM, _ := strconv.Atoi(m)
	pD, _ := strconv.Atoi(d)

	if (dateY == pY) && ((dateM) == pM) && (dateD == pD) { //判断年份的范围（3岁到100岁之间)
		time := now_year - pY
		if time >= 0 && time <= 130 {
			return true
		}
		return false
	}
	return false
}

func CheckParity(idcard string) bool {
	idcardStr := changeFivteeenToEighteen(idcard)
	len := len(idcardStr)

	if len == 18 {
		var s int
		var wi [17]int = [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		var arrCh = [11]int{1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2}
		arr := make([]int, 17)
		for i := 0; i < 17; i++ {
			arr[i], _ = strconv.Atoi(string(idcardStr[i]))
		}
		for i := 0; i < 17; i++ {
			s += arr[i] * wi[i]
		}
		valnum := arrCh[s%11]

		if valnum == byte2int(idcardStr[17:]) {
			//fmt.Println("校验码正确")
			return true
		}
		//fmt.Println("校验码错误")
		return false
	}
	//fmt.Println("校验码正确End")
	return true
}

func changeFivteeenToEighteen(idcard string) string {
	if len(idcard) == 15 {
		var wi [17]int = [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		var arrCh = [11]int{1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2}
		dateOf15 := idcard[6:12]                   // 出生日期码: 6-12位位生日 700702,即为70年，07月，02日
		head := idcard[:6]                         // 地址码: 身份证前6位 411327
		tail := idcard[12:]                        // 身份证的后3位 [合起来一共15位]
		newIDCard := head + "19" + dateOf15 + tail // 在年份前加上19，即修改为 411327 19700702 001 ，变成17位

		arr := make([]int, 17)
		for i := 0; i < 17; i++ {
			arr[i], _ = strconv.Atoi(string(newIDCard[i]))
		}

		var s int // 默认值为0
		for i := 0; i < 17; i++ {
			s += arr[i] * wi[i]
		}
		s += arrCh[s%11]
		return string(s)
	}
	return idcard
}

func byte2int(s string) int {
	if s == "X" {
		return 88
	}

	res, _ := strconv.Atoi(s)
	return res
}

func IdentityCardNoGetAge(idCard string) int {
	var mapmou = map[string]int{"January": 1, "february": 2, "March": 3, "April": 4, "May": 5, "June": 6, "July": 7, "August": 8, "September": 9, "October": 10, "November": 11, "December": 12}
	now := time.Now()
	now_year := now.Year()                 // 年
	now_mo := mapmou[now.Month().String()] // 月
	now_day := now.Day()                   // 日
	//fmt.Println(now_year, now_mo, now_day)
	idcard_year, _ := strconv.Atoi(substr(idCard, 6, 4)) // 年
	idcard_mo, _ := strconv.Atoi(substr(idCard, 10, 2))  // 月
	idcard_day, _ := strconv.Atoi(substr(idCard, 12, 2)) // 日
	//fmt.Println(idcard_year, idcard_mo, idcard_day)
	age := now_year - idcard_year // 如果计算虚岁需这样：age := now_year - idcard_year+1
	if now_year < idcard_year {
		age = 0
	} else {
		if now_mo < idcard_mo {
			age = age - 1
		} else if now_mo == idcard_mo {
			if now_day < idcard_day {
				age = age - 1
			}
		}
	}
	return age
}

func substr(str string, start, length int) string {

	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
