package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	err := os.Remove("tfile.js")
	if err != nil {
		var pathError *os.PathError

		if !errors.As(err, &pathError) {
			log.Fatal(err)
		}
	}

	tfile, err := os.Create("tfile.js")
	if err != nil {
		log.Fatal(err)
	}
	defer tfile.Close()

	file, err := os.Open("file.js")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	defer file.Close()

	//Текст или нет на начало считывания данных
	tcomaa := false
	//Последний символ набора данных
	tlast := ""

	data := make([]byte, 50)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			if _, err := tfile.WriteString(tlast); err != nil {
				panic(err)
			}
			break
		}
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}
		//Переводим считанные данные в строку
		tstring := tlast + string(data[:n])
		if strings.Contains(tstring, "(_0x1dfdc0);") {
			a := 1
			_ = a
		}
		//Обрабатываемые символы
		tbrackes := [6]string{"'", "{", "}"}
		//меняется ли что-то в строке или нужно что-то учесть
		tchange := false
		for _, v := range tbrackes {
			if v == "" || tchange {
				continue
			}
			if strings.Contains(tstring, v) {
				tchange = true
			}
		}
		//если нет ничего из нужного - пишшем в файл
		if !tchange {
			tlast = getlast(tstring)
			if _, err := tfile.Write(data[:n]); err != nil {
				panic(err)
			}
		} else {
			tsave := ""
			//Если что-то есть, то проверяем
			for i, v := range data {
				if i >= n {
					break
				}
				s := string(v)
				//Обрабатываемый символ

				if contains(tbrackes[:], s) {
					switch s {
					case "'":
						if tlast != "\\" {
							tcomaa = !tcomaa
						}
					case "}":
						if !tcomaa {
							tsave += "\n"
						}
					case "{":
						r := []rune(tlast)
						if !tcomaa && (tlast == "" || !unicode.IsLetter(r[:][0])) {
							tsave += "\n"
						}
					}
					tlast = ""
				} else {
					tlast = getlast(s)
				}
				tsave += s
			}
			if _, err := tfile.WriteString(tsave); err != nil {
				panic(err)
			}
		}
	}
}

func getlast(tstring string) string {
	s := []rune(tstring)
	tlast := string(s[len(s)-1:])

	if tlast != "\\" && !unicode.IsLetter(s[len(s)-1:][0]) {
		tlast = ""
	}
	return tlast
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
