package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/goudev/survey/v2/core"
	"github.com/goudev/survey/v2/terminal"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestInputRender(t *testing.T) {

	suggestFn := func(string) (s []string) { return s }

	tests := []struct {
		title    string
		prompt   Input
		data     InputTemplateData
		expected string
	}{
		{
			"Teste de saída da pergunta sem valor padrão",
			Input{Message: "Qual é o seu mês favorito:"},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: ", defaultIcons().Question.Text),
		},
		{
			"Teste de saída da pergunta com valor padrão",
			Input{Message: "Qual é o seu mês favorito:", Default: "Abril"},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: (Abril) ", defaultIcons().Question.Text),
		},
		{
			"Teste de saída da resposta da pergunta",
			Input{Message: "Qual é o seu mês favorito:"},
			InputTemplateData{ShowAnswer: true, Answer: "Outubro"},
			fmt.Sprintf("%s Qual é o seu mês favorito: Outubro\n", defaultIcons().Question.Text),
		},
		{
			"Teste de saída da pergunta sem valor padrão, mas com ajuda oculta",
			Input{Message: "Qual é o seu mês favorito:", Help: "Isso é útil"},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: [%s para ajuda] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Teste de saída da pergunta com valor padrão e ajuda oculta",
			Input{Message: "Qual é o seu mês favorito:", Default: "Abril", Help: "Isso é útil"},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: [%s para ajuda] (Abril) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Teste de saída da pergunta sem valor padrão, mas com ajuda exibida",
			Input{Message: "Qual é o seu mês favorito:", Help: "Isso é útil"},
			InputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s Isso é útil\n%s Qual é o seu mês favorito: ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Teste de saída da pergunta com valor padrão e ajuda exibida",
			Input{Message: "Qual é o seu mês favorito:", Default: "Abril", Help: "Isso é útil"},
			InputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s Isso é útil\n%s Qual é o seu mês favorito: (Abril) ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Teste de saída da pergunta com sugestões",
			Input{Message: "Qual é o seu mês favorito:", Suggest: suggestFn},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: [%s para sugestões] ", defaultIcons().Question.Text, string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Teste de saída da pergunta com sugestões e ajuda oculta",
			Input{Message: "Qual é o seu mês favorito:", Suggest: suggestFn, Help: "Isso é útil"},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: [%s para ajuda, %s para sugestões] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput), string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Teste de saída da pergunta com sugestões, valor padrão e ajuda oculta",
			Input{Message: "Qual é o seu mês favorito:", Suggest: suggestFn, Help: "Isso é útil", Default: "Abril"},
			InputTemplateData{},
			fmt.Sprintf("%s Qual é o seu mês favorito: [%s para ajuda, %s para sugestões] (Abril) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput), string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Teste de saída da pergunta com sugestões exibidas",
			Input{Message: "Qual é o seu mês favorito:", Suggest: suggestFn},
			InputTemplateData{
				Answer:        "Fevereiro",
				PageEntries:   core.OptionAnswerList([]string{"Janeiro", "Fevereiro", "Março", "etc..."}),
				SelectedIndex: 1,
			},
			fmt.Sprintf(
				"%s Qual é o seu mês favorito: Fevereiro [Use as setas para mover, enter para selecionar, digite para continuar]\n"+"  Janeiro\n%s Fevereiro\n  Março\n  etc...\n",
				defaultIcons().Question.Text, defaultPromptConfig().Icons.SelectFocus.Text,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r, w, err := os.Pipe()
			assert.NoError(t, err)

			test.prompt.WithStdio(terminal.Stdio{Out: w})
			test.data.Input = test.prompt

			// define a configuração de tempo de execução
			test.data.Config = defaultPromptConfig()

			err = test.prompt.Render(
				InputQuestionTemplate,
				test.data,
			)
			assert.NoError(t, err)

			assert.NoError(t, w.Close())
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.NoError(t, err)

			assert.Contains(t, buf.String(), test.expected)
		})
	}
}

func TestInputPrompt(t *testing.T) {

	tests := []PromptTest{
		{
			"Test Input prompt interaction",
			&Input{
				Message: "What is your name?",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction with default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("")
				c.ExpectEOF()
			},
			"Johnny Appleseed",
		},
		{
			"Test Input prompt interaction overriding default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction and prompt for help",
			&Input{
				Message: "What is your name?",
				Help:    "It might be Satoshi Nakamoto",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("?")
				c.ExpectString("It might be Satoshi Nakamoto")
				c.SendLine("Satoshi Nakamoto")
				c.ExpectEOF()
			},
			"Satoshi Nakamoto",
		},
		{
			// https://en.wikipedia.org/wiki/ANSI_escape_code
			// Device Status Report - Reports the cursor position (CPR) to the
			// application as (as though typed at the keyboard) ESC[n;mR, where n is the
			// row and m is the column.
			"SKIP: Test Input prompt with R matching DSR",
			&Input{
				Message: "What is your name?",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.SendLine("R")
				c.ExpectEOF()
			},
			"R",
		},
		{
			"Test Input prompt interaction when delete",
			&Input{
				Message: "What is your name?",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.Send("Johnny ")
				c.Send(string(terminal.KeyDelete))
				c.SendLine("")
				c.ExpectEOF()
			},
			"Johnny",
		},
		{
			"Test Input prompt interaction when delete rune",
			&Input{
				Message: "What is your name?",
			},
			func(c expectConsole) {
				c.ExpectString("What is your name?")
				c.Send("小明")
				c.Send(string(terminal.KeyBackspace))
				c.SendLine("")
				c.ExpectEOF()
			},
			"小",
		},
		{
			"Test Input prompt interaction when ask for suggestion with empty value",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"January", "February"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("What is your favorite month?")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("January")
				c.ExpectString("February")
				c.SendLine("")
				c.ExpectEOF()
			},
			"January",
		},
		{
			"Test Input prompt interaction when ask for suggestion with some value",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"February"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("What is your favorite month?")
				c.Send("feb")
				c.Send(string(terminal.KeyTab))
				c.SendLine("")
				c.ExpectEOF()
			},
			"February",
		},
		{
			"Test Input prompt interaction when ask for suggestion with some value, choosing the second one",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"January", "February", "March"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("What is your favorite month?")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine("")
				c.ExpectEOF()
			},
			"March",
		},
		{
			"Test Input prompt interaction when ask for suggestion with some value, choosing the second one",
			&Input{
				Message: "What is your favorite month?",
				Suggest: func(string) []string {
					return []string{"January", "February", "March"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("What is your favorite month?")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowUp))
				c.SendLine("")
				c.ExpectEOF()
			},
			"February",
		},
		{
			"Test Input prompt interaction when ask for suggestion, complementing it and get new suggestions",
			&Input{
				Message: "Where to save it?",
				Suggest: func(complete string) []string {
					if complete == "" {
						return []string{"folder1/", "folder2/", "folder3/"}
					}
					return []string{"folder3/file1.txt", "folder3/file2.txt"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("Where to save it?")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("folder1/")
				c.Send(string(terminal.KeyArrowDown))
				c.Send(string(terminal.KeyArrowDown))
				c.Send("f")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("folder3/file2.txt")
				c.Send(string(terminal.KeyArrowDown))
				c.SendLine("")
				c.ExpectEOF()
			},
			"folder3/file2.txt",
		},
		{
			"Test Input prompt interaction when asked suggestions, but abort suggestions",
			&Input{
				Message: "Wanna a suggestion?",
				Suggest: func(string) []string {
					return []string{"suggest1", "suggest2"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("Wanna a suggestion?")
				c.Send("typed answer")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("suggest1")
				c.Send(string(terminal.KeyEscape))
				c.ExpectString("typed answer")
				c.SendLine("")
				c.ExpectEOF()
			},
			"typed answer",
		},
		{
			"Test Input prompt interaction with suggestions, when tabbed with list being shown, should select next suggestion",
			&Input{
				Message: "Choose the special one:",
				Suggest: func(string) []string {
					return []string{"suggest1", "suggest2", "special answer"}
				},
			},
			func(c expectConsole) {
				c.ExpectString("Choose the special one:")
				c.Send("s")
				c.Send(string(terminal.KeyTab))
				c.ExpectString("suggest1")
				c.ExpectString("suggest2")
				c.ExpectString("special answer")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyTab))
				c.SendLine("")
				c.ExpectEOF()
			},
			"special answer",
		},
		{
			"Test Input prompt must allow moving cursor using right and left arrows",
			&Input{Message: "Filename to save:"},
			func(c expectConsole) {
				c.ExpectString("Filename to save:")
				c.Send("essay.txt")
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send("_final")
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyArrowRight))
				c.Send(string(terminal.KeyBackspace))
				c.Send(string(terminal.KeyBackspace))
				c.Send(string(terminal.KeyBackspace))
				c.Send("md")
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.SendLine("2")
				c.ExpectEOF()
			},
			"essay_final2.md",
		},
		{
			"Test Input prompt must allow moving cursor using right and left arrows, even after suggestions",
			&Input{Message: "Filename to save:", Suggest: func(string) []string { return []string{".txt", ".csv", ".go"} }},
			func(c expectConsole) {
				c.ExpectString("Filename to save:")
				c.Send(string(terminal.KeyTab))
				c.ExpectString(".txt")
				c.ExpectString(".csv")
				c.ExpectString(".go")
				c.Send(string(terminal.KeyTab))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send(string(terminal.KeyArrowLeft))
				c.Send("newtable")
				c.SendLine("")
				c.ExpectEOF()
			},
			"newtable.csv",
		},
	}

	for _, test := range tests {
		testName := strings.TrimPrefix(test.name, "SKIP: ")
		t.Run(testName, func(t *testing.T) {
			if testName != test.name {
				t.Skipf("warning: flakey test %q", testName)
			}
			RunPromptTest(t, test)
		})
	}
}
