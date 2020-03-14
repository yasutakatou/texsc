package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"github.com/yasutakatou/ishell"
)

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var label string
var filename string
var regexpStr string
var predict string
var learn string
var predictChar int
var learnChar int
var modelName string
var bufStr string

var annotationFlag bool
var Debug bool

func main() {
	_debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_label := flag.String("label", "__label__1", "[-label=classifier label]")
	_modelName := flag.String("model", "model", "[-model=define #MODEL#.]")
	_file := flag.String("file", "test.txt", "[-file=target text file]")
	_regexp := flag.String("regexp", "[0-9][0-9]:[0-9][0-9]", "[-regexp=target word to empty]")
	_predict := flag.String("predict", "fasttext.exe predict-prob #MODEL#.bin #DATA#", "[-predict=predict command. #DATA# is clipboard text.]")
	_learn := flag.String("learn", "fasttext.exe supervised -input #DATA# -output #MODEL# -epoch 1000", "[-learn=learn command. #DATA# is target text file.]")
	_predictChar := flag.Int("predictChar", 88, "[-predictChar=predict keyboard char.]")
	_learnChar := flag.Int("learnChar", 90, "[-learnChar=learn keyboard char.]")

	flag.Parse()

	label = string(*_label)
	filename = string(*_file)
	regexpStr = string(*_regexp)
	predict = string(*_predict)
	learn = string(*_learn)
	predictChar = int(*_predictChar)
	learnChar = int(*_learnChar)
	modelName = string(*_modelName)
	Debug = bool(*_debug)

	annotationFlag = true

	showConfig()

	bufStr, _ = clipboard.ReadAll()

	go clipboardToExport()
	go backgroundExecute()

	var shell = ishell.New()

	shell.AddCmd(&ishell.Cmd{Name: "setLabel",
		Help: "label setting",
		Func: CliHandler})

	shell.AddCmd(&ishell.Cmd{Name: "showConfig",
		Help: "print config",
		Func: CliHandler})

	shell.AddCmd(&ishell.Cmd{Name: "setFile",
		Help: "file setting",
		Func: CliHandler})

	shell.AddCmd(&ishell.Cmd{Name: "default",
		Help: "default is print config",
		Func: CliHandler})

	shell.Run()
}

func clipboardToExport() {
	for {
		text, _ := clipboard.ReadAll()

		if bufStr != text {
			str := cleansingStr(text)
			if annotationFlag == true {
				fileWrite(filename, str)
				fmt.Println(" -- -- -- -- clip! -- -- -- -- ")
				if Debug == true {
					fmt.Println(str)
				}
			} else {
				fmt.Println(" -- -- -- -- predict! -- -- -- -- ")
				fmt.Println(str)
			}

			bufStr = text
		}
		time.Sleep(time.Duration(250) * time.Millisecond)
	}
}

func cleansingStr(text string) string {
	rep := regexp.MustCompile(regexpStr)
	str := rep.ReplaceAllString(text, "")

	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "　", "", -1)

	tmpFile := RandStr(8) + ".txt"

	fileWrite(tmpFile, str)

	str = Execmd("mecab -Owakati " + tmpFile)

	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)

	if annotationFlag == false {
		tmp2File := RandStr(8) + ".txt"

		fileWrite(tmp2File, str)

		rule := strings.Replace(predict, "#MODEL#", modelName, 1)
		rule = strings.Replace(rule, "#DATA#", tmp2File, 1)
		str = Execmd(rule)

		//if Debug == false {
		Execmd("del " + tmp2File)
		//}
	} else {
		str = label + ", " + str
	}

	//if Debug == false {
	Execmd("del " + tmpFile)
	//}

	return str
}

func CliHandler(c *ishell.Context) {
	switch c.Cmd.Name {
	case "setLabel":
		if len(c.Args) == 0 {
			return
		}
		setStr(&label, c.Args[0])
	case "setFile":
		if len(c.Args) == 0 {
			return
		}
		setStr(&filename, c.Args[0])
	case "showConfig":
		showConfig()
	default:
		showConfig()
	}
}

func showConfig() {
	fmt.Printf(" - - predict (%s) key press ascii code ctrl+(%s[%d]) - - \n", predict, string(predictChar), predictChar)
	fmt.Printf(" - - learn (%s) key press ascii code  ctrl+(%s[%d]) - - \n", learn, string(learnChar), learnChar)
	fmt.Println("File: ", filename)
	fmt.Println("Label: ", label)
	fmt.Printf("annotation: %t\n", annotationFlag)
}

func setStr(str *string, params string) bool {
	if len(params) == 0 {
		return false
	}
	*str = params
	return true
}

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

func fileWrite(filename, str string) {
	if Exists(filename) == false {
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		fmt.Fprintln(file, str)
	} else {
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		fmt.Fprintln(file, str)
	}

}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Execmd(command string) string {
	if Debug == true {
		fmt.Println("command: ", command)
	}
	out, err := exec.Command("cmd", "/C", command).Output()
	if err != nil {
		fmt.Println(err)
	}
	if Debug == true {
		fmt.Println("output: ", string(out))
	}
	return string(out)
}

func backgroundExecute() {
	ctlFlag := false

	EvChan := hook.Start()
	defer hook.End()

	for ev := range EvChan {
		strs := ""

		//KeyDown = 3
		if ev.Kind == 3 {
			if ctlFlag == true {
				strs = string(ev.Keychar)

				switch int(ev.Rawcode) {
				case learnChar:
					fmt.Printf(" learning..\n")
					strs = strings.Replace(learn, "#MODEL#", modelName, 1)
					strs = strings.Replace(strs, "#DATA#", filename, 1)
					Execmd(strs)
					fmt.Printf("Done! \n")
				case predictChar:
					if annotationFlag == false {
						annotationFlag = true
					} else {
						annotationFlag = false
					}
					fmt.Printf(" annotation: %t\n", annotationFlag)
				}
				ctlFlag = false
			}
		}

		//KeyHold = 4,KeyUp   = 5
		if ev.Kind == 4 || ev.Kind == 5 {
			if ev.Rawcode == 162 {
				ctlFlag = true
			}
		}
	}
}
