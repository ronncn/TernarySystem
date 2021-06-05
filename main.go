package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	zero = byte('0')
	one  = byte('1')
)

func init() {
	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1

	code = make(map[int]string)
	code[0] = "0"
	code[1] = "1"
	code[2] = "2"
	code[3] = "3"
	code[4] = "4"
	code[5] = "5"
	code[6] = "6"
	code[7] = "7"
	code[8] = "8"
	code[9] = "9"
	code[10] = "□"
	code[11] = "A"
	code[12] = "B"
	code[13] = "C"
	code[14] = "D"
	code[15] = "E"
	code[16] = "F"
	code[17] = "G"
	code[18] = "H"
	code[19] = "I"
	code[20] = "J"
	code[21] = "K"
	code[22] = "L"
	code[23] = "M"
	code[24] = "N"
	code[25] = "O"
	code[26] = "P"
	code[27] = "Q"
	code[28] = "R"
	code[29] = "S"
	code[30] = "T"
	code[31] = "U"
	code[32] = "V"
	code[33] = "W"
	code[34] = "X"
	code[35] = "Y"
	code[36] = "Z"
	code[41] = "|"
	code[42] = ">"
	code[43] = "<"
	code[44] = "*"
	code[45] = "#"
	code[48] = "+"
	code[49] = "-"
	code[50] = "="
	code[51] = "a"
	code[52] = "b"
	code[53] = "c"
	code[54] = "d"
	code[55] = "e"
	code[56] = "f"
	code[57] = "g"
	code[58] = "h"
	code[59] = "i"
	code[60] = "j"
	code[61] = "k"
	code[62] = "l"
	code[63] = "m"
}

var code map[int]string

func main() {
	fmt.Println("请使用浏览器访问 http://127.0.0.1:8000 ")
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/encode", handlerEncode)
	// http.HandleFunc("/decode", handlerDecode)
	http.HandleFunc("/upload", handlerUpload)
	http.HandleFunc("/download", handlerDownload)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(200)
	tem, err := template.ParseFiles("index.html")
	if err == nil {
		tem.Execute(w, nil)
	}
}

// 处理编码
func handlerEncode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	databs := ternaryEncode(data)
	databin := bytesToBinaryString(databs)
	datater := bytesToTernaryString(databs)
	rsdata := make(map[string]string)
	rsdata["bin"] = databin
	rsdata["ter"] = datater

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	dataByte, err := json.Marshal(rsdata)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(dataByte))
}

// 处理解码
// func handlerDecode(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("content-type", "text/html")
// 	w.WriteHeader(200)
// 	tem, err := template.ParseFiles("index.html")
// 	if err == nil {
// 		tem.Execute(w, nil)
// 		return
// 	}
// 	fmt.Println(err)
// }

// 处理上传
func handlerUpload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	defer file.Close()
	if err == nil {
		buf := make([]byte, fileHeader.Size)
		file.Read(buf)
		str := ternaryDecode(buf)
		databin := bytesToBinaryString(buf)
		datater := bytesToTernaryString(buf)
		rsdata := make(map[string]string)
		rsdata["str"] = str
		rsdata["bin"] = databin
		rsdata["ter"] = datater

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		dataByte, err := json.Marshal(rsdata)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(dataByte))
	}
}

// 处理下载
func handlerDownload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue("data")
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+MD5(data)+".bin")
	databs := ternaryEncode(data)
	w.WriteHeader(200)
	n, err := w.Write(databs)
	if err != nil {
		fmt.Printf("n=%d err:%v\n", n, err)
	}
}

// MD5 算法
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func findkey(str string) int {
	result := 10
	for k, v := range code {
		if str == v {
			result = k
		}
	}
	return result
}

func ternaryEncode(str string) []byte {
	strs := strings.Split(str, "")
	bin := ""
	for _, v := range strs {
		n := findkey(v)
		bin += convertToBin(n, 6)
	}
	fmt.Println("bin:", bin)
	return binaryStringToBytes(bin)
}

func ternaryDecode(bs []byte) string {
	t := ""
	rs := ""
	bin := bytesToBinaryString(bs)
	l := len(bin)
	mo := l % 6
	for i := 0; i < l-mo; i++ {
		t += string(bin[i])
		if (i+1)%6 == 0 {
			n, _ := strconv.ParseInt(t, 2, 32)
			if _, ok := code[int(n)]; ok {
				if !(i+1 == l && int(n) == 10) {
					rs += code[int(n)]
				}
			} else {
				rs += code[10]
			}
			t = ""
		}

	}
	return rs
}

// 10进制转任意进制
func decimalToAny(num, n int) string {
	any := "0123456789abcdefghijklmnopqrstuvwxyz"
	newStr := ""
	var remainder int
	var remString string
	for num != 0 {
		remainder = num % n
		if 76 > remainder && remainder > 9 {
			remString = string(any[remainder])
		} else {
			remString = strconv.Itoa(remainder)
		}
		newStr = remString + newStr
		num = num / n
	}
	return newStr
}

// 将十进制数字转化为二进制字符串 补齐0
func convertToBin(num int, dig int) string {
	s := ""
	if num == 0 {
		s = "0"
	}
	for ; num > 0; num /= 2 {
		lsb := num % 2
		s = strconv.Itoa(lsb) + s
	}
	l := len(s)
	for i := 0; i < dig-l; i++ {
		s = "0" + s
	}
	return s
}

var uint8arr [8]uint8

func binaryStringToBytes(s string) (bs []byte) {
	l := len(s)
	if l == 0 {
		return bs
	}
	mo := l % 8
	l /= 8
	if mo != 0 {
		l++
	}
	bs = make([]byte, 0, l)
	mo = 8 - mo
	var n uint8
	for i, b := range []byte(s) {
		m := i % 8
		switch b {
		case one:
			n += uint8arr[m]
		}
		if i == len([]byte(s))-1 {
			if m == 1 {
				n += uint8arr[4]
				n += uint8arr[6]
			}
		}
		if m == 7 || i == len([]byte(s))-1 {
			bs = append(bs, n)
			n = 0
		}
	}
	return
}

func bytesToBinaryString(bs []byte) string {
	bin := ""
	for _, v := range bs {
		a := int(v)
		bin += convertToBin(a, 8)
	}
	return bin
}

func bytesToTernaryString(bs []byte) string {

	b := ""
	ter := ""
	bin := bytesToBinaryString(bs)
	for i := 0; i < len(bin); i++ {
		b += string(bin[i])
		if (i+1)%6 == 0 {
			n, _ := strconv.ParseInt(b, 2, 32)
			if _, ok := code[int(n)]; ok {
				t := decimalToAny(int(n), 3)
				l := len(t)
				for i := 0; i < 4-l; i++ {
					t = "0" + t
				}
				ter += t
			} else {
				ter += "0101"
			}
			b = ""
		}
	}
	return ter
}
