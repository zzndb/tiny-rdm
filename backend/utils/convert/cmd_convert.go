package convutil

import (
	"encoding/base64"
	"log"
	"strconv"
	"time"
	"strings"
	sliceutil "tinyrdm/backend/utils/slice"
)

type CmdConvert struct {
	Name       string
	Auto       bool
	DecodePath string
	DecodeArgs []string
	EncodePath string
	EncodeArgs []string
}

const replaceholder = "{VALUE}"

func (c CmdConvert) Enable() bool {
	return true
}

func (c CmdConvert) Encode(str string) (string, bool) {
	base64Content := base64.StdEncoding.EncodeToString([]byte(str))
	var containHolder bool
	args := sliceutil.Map(c.EncodeArgs, func(i int) string {
		arg := strings.TrimSpace(c.EncodeArgs[i])
		if strings.Contains(arg, replaceholder) {
			arg = strings.ReplaceAll(arg, replaceholder, base64Content)
			containHolder = true
		}
		return arg
	})
	if len(args) <= 0 || !containHolder {
		args = append(args, base64Content)
	}
	output, err := runCommand(c.EncodePath, args...)
	var fileFallback bool = false
	if err != nil || len(output) <= 0 || string(output) == "[RDM-ERROR]" {
		if err != nil {
			if strings.Contains(err.Error(), "argument list too long") {
				fileFallback = true
				log.Println("fallback to file encode")
			}
		}
		if (!fileFallback) {
			return str, false
		}
	}
	if (fileFallback) {
		var filePath string
		filePath, err = writeTempFile([]byte(base64Content), strconv.Itoa(int(time.Now().UnixMilli())))

		args = args[:len(args)-1]
		args = append(args, "fallback", filePath)
		output, err = runCommand(c.EncodePath, args...)
		err = cleanTempFile(filePath)

		if err != nil || len(output) <= 0 || string(output) == "[RDM-ERROR]" {
			return str, false
		}
	}

	outputContent := make([]byte, base64.StdEncoding.DecodedLen(len(output)))
	n, err := base64.StdEncoding.Decode(outputContent, output)
	if err != nil {
		return str, false
	}
	return string(outputContent[:n]), true
}

func (c CmdConvert) Decode(str string) (string, bool) {
	base64Content := base64.StdEncoding.EncodeToString([]byte(str))
	var containHolder bool
	args := sliceutil.Map(c.DecodeArgs, func(i int) string {
		arg := strings.TrimSpace(c.DecodeArgs[i])
		if strings.Contains(arg, replaceholder) {
			arg = strings.ReplaceAll(arg, replaceholder, base64Content)
			containHolder = true
		}
		return arg
	})
	if len(args) <= 0 || !containHolder {
		args = append(args, base64Content)
	}
	output, err := runCommand(c.DecodePath, args...)
	var fileFallback bool = false
	if err != nil || len(output) <= 0 || string(output) == "[RDM-ERROR]" {
		if err != nil {
			if strings.Contains(err.Error(), "argument list too long") {
				fileFallback = true
				log.Println("fallback to file decode")
			}
		}
		if (!fileFallback) {
			return str, false
		}
	}
	if (fileFallback) {
		var filePath string
		filePath, err = writeTempFile([]byte(base64Content), strconv.Itoa(int(time.Now().UnixMilli())))

		args = args[:len(args)-1]
		args = append(args, "fallback", filePath)
		output, err = runCommand(c.DecodePath, args...)
		err = cleanTempFile(filePath)

		if err != nil || len(output) <= 0 || string(output) == "[RDM-ERROR]" {
			return str, false
		}
	}

	outputContent := make([]byte, base64.StdEncoding.DecodedLen(len(output)))
	n, err := base64.StdEncoding.Decode(outputContent, output)
	if err != nil {
		return str, false
	}
	return string(outputContent[:n]), true
}
