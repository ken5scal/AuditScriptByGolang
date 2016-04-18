package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

const (
	RSYNC  = "rsync"
	SFTP   = "sftp"
	SCP    = "scp"
	TELNET = "telnet"
)

var inspectingFile, resultFile string
var originalContent []byte
var auditingComands = []string{RSYNC, SFTP, SCP, TELNET}

func main() {
	// Argumentからファイル名を読み込む
	if len(os.Args) == 3 {
		inspectingFile = os.Args[1]
		resultFile = os.Args[2]
	} else {
		panic("Must specify both inspecting file name and result file name")
	}

	// 正規表現の生成
	regexpCommands := auditingComands[len(auditingComands)-1]
	auditingComands = auditingComands[:len(auditingComands)-1] // POP
	for _, command := range auditingComands {
		regexpCommands = regexpCommands + `|` + command
	}

	// 読み込み
	originalContent = findMatchingLines(inspectingFile, regexpCommands, originalContent)

	// 書き込み
	writeToFile(originalContent)
}

// テキストをバイト型で読み込む
func findMatchingLines(fileName string, grepStrings string, contents []byte) []byte {
	regex := regexp.MustCompile(grepStrings)
	regexPrompt := regexp.MustCompile(`\@[A-Za-z0-9\\\._-]+\s[A-Za-z0-9\\\.\/\~_-]+\]\$\s`) //@HOSTNAME DIR]$ のようなhost名からからプロンプトまでにマッチ:w

	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	for {
		line, _, err := reader.ReadLine()
		if regex.MatchString(string(line)) && regexPrompt.MatchString(string(line)) {
			contents = append(contents, (string(line) + "\n")...)
		} else if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return contents
}

// テキスト(バイト型)を書き込む
func writeToFile(extractedStrings []byte) {
	ioutil.WriteFile(resultFile, extractedStrings, os.ModePerm)
}
