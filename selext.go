// +build OMIT

package main

import (
	"errors"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var eol = func() string {
	switch runtime.GOOS {
	case "linux":
		return "\n"
	case "windows":
		return "\r"
	case "darwin":
		return "\n\r"
	default:
		return "\n"
	}
}()

func reportError() {

}

func stringsUnique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func commandUniq(input string) (string, error) {
	lines := strings.Split(input, eol)
	lines = stringsUnique(lines)
	return strings.Join(lines, eol), nil
}

func commandTrim(input string) (string, error) {
	lines := strings.Split(input, eol)
	var trimmed []string
	for _, v := range lines {
		trimmed = append(trimmed, strings.Trim(v, " "))
	}
	return strings.Join(trimmed, eol), nil
}

func commandCount(input string) (string, error) {
	return strconv.Itoa(len(strings.Split(input, eol))), nil
}

func commandRe(input string, regexPattern string) (string, error) {
	r, err := regexp.Compile(regexPattern)
	if err != nil {
		return "", errors.New("re: invalid regexp")
	}
	res := r.FindAllString(input, -1)
	return strings.Join(res, eol), nil
}

func commandEmail(input string) (string, error) {
	r, _ := regexp.Compile("([a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+)")
	res := r.FindAllString(input, -1)
	return strings.Join(res, eol), nil
}

func commandIpv4(input string) (string, error) {
	r, _ := regexp.Compile("((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)")
	res := r.FindAllString(input, -1)
	return strings.Join(res, eol), nil
}

func commandUpper(input string) (string, error) {
	return strings.ToUpper(input), nil
}

func commandLower(input string) (string, error) {
	return strings.ToLower(input), nil
}

func commandPrefix(input string, prefix string) (string, error) {
	lines := strings.Split(input, eol)
	var res []string
	for _, v := range lines {
		res = append(res, prefix+v)
	}
	return strings.Join(res, eol), nil
}

func commandPostfix(input string, postfix string) (string, error) {
	lines := strings.Split(input, eol)
	var res []string
	for _, v := range lines {
		res = append(res, v+postfix)
	}
	return strings.Join(res, eol), nil
}

func commandSum(input string) (string, error) {
	lines := strings.Split(input, eol)
	sum := 0.0
	for _, v := range lines {
		if s, err := strconv.ParseFloat(v, 64); err == nil {
			sum += s
		} else {
			return "", errors.New(fmt.Sprintf("re: invalid float format: %s", v))
		}
	}
	return strconv.FormatFloat(sum, 'f', 6, 64), nil
}

func commandAsc(input string) (string, error) {
	lines := strings.Split(input, eol)
	sort.Strings(lines)
	return strings.Join(lines, eol), nil
}

func commandDesc(input string) (string, error) {
	lines := strings.Split(input, eol)
	sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	return strings.Join(lines, eol), nil
}

func processCommands(input string, code string) (string, error) {
	var err error
	code, err = commandTrim(code)
	output := input
	codeLines := strings.Split(code, eol)
	for _, line := range codeLines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		commandSlice := strings.Split(line, ",")
		switch commandSlice[0] {
		case "count":
			output, err = commandCount(output)
		case "trim":
			output, err = commandTrim(output)
		case "uniq":
			output, err = commandUniq(output)
		case "re":
			output, err = commandRe(output, strings.Trim(strings.Join(commandSlice[1:], ""), "\"'"))
		case "email":
			output, err = commandEmail(output)
		case "ipv4":
			output, err = commandIpv4(output)
		case "upper":
			output, err = commandUpper(output)
		case "lower":
			output, err = commandLower(output)
		case "prefix":
			output, err = commandPrefix(output, strings.Trim(strings.Join(commandSlice[1:], ""), "\"'"))
		case "postfix":
			output, err = commandPostfix(output, strings.Trim(strings.Join(commandSlice[1:], ""), "\"'"))
		case "sum":
			output, err = commandSum(output)
		case "asc":
			output, err = commandAsc(output)
		case "desc":
			output, err = commandDesc(output)
		}
	}
	return output, err
}

func main() {
	err := ui.Main(func() {
		window := ui.NewWindow("selext v1", 800, 600, false)

		inputTextArea := ui.NewMultilineEntry()
		inputTextArea.SetText("please@he\n" +
			"12\n" +
			"dcdc\n" +
			"$80\n" +
			"$60\n" +
			"dcdc\n" +
			"23\n" +
			"ceacaec34\n" +
			"ceacaec34\n" +
			"please@help.me")
		vboxInput := ui.NewVerticalBox()
		vboxInput.Append(ui.NewLabel("Input"), false)
		vboxInput.Append(inputTextArea, true)

		outputTextArea := ui.NewMultilineEntry()
		outputTextArea.SetReadOnly(true)
		vboxOutput := ui.NewVerticalBox()
		vboxOutput.Append(ui.NewLabel("Output"), false)
		vboxOutput.Append(outputTextArea, true)

		hbox1 := ui.NewHorizontalBox()
		hbox1.Append(vboxInput, true)
		hbox1.Append(vboxOutput, true)

		buttonRun := ui.NewButton("Run")

		messagesTextArea := ui.NewMultilineEntry()
		messagesTextArea.SetReadOnly(true)
		vbox1 := ui.NewVerticalBox()
		vbox1.Append(ui.NewLabel("Messages"), false)
		vbox1.Append(messagesTextArea, true)

		codeTextArea := ui.NewMultilineEntry()
		vboxCode := ui.NewVerticalBox()
		vboxCode.Append(ui.NewLabel("Code"), false)
		vboxCode.Append(codeTextArea, true)
		vboxCode.Append(buttonRun, false)

		hbox2 := ui.NewHorizontalBox()
		hbox2.Append(vboxCode, true)
		hbox2.Append(vbox1, true)

		vboxmain := ui.NewVerticalBox()
		vboxmain.Append(hbox1, true)
		vboxmain.Append(hbox2, true)

		buttonRun.OnClicked(func(*ui.Button) {
			output, err := processCommands(inputTextArea.Text(), codeTextArea.Text())
			if err != nil {
				ui.MsgBoxError(window, "Command error", err.Error())
			} else {
				outputTextArea.SetText(output)
			}
		})

		window.SetMargined(true)
		window.SetChild(vboxmain)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
