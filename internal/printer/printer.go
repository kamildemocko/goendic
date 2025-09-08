package printer

import (
	"endic/internal/data/model"
	"fmt"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

func PrintEmpty() {
	printerRed := color.New(color.FgRed)
	printerRed.Println("No result found.")
}

func PrintUsage() {
	printerGray := color.New(color.FgHiBlack)

	fmt.Println("Error: No search word provided")
	fmt.Print("Usage: endic")
	printerGray.Print(" [OPTIONS] ")
	fmt.Println("WORD")
	fmt.Println("Options: ")
	printerGray.Println(" -e  : Use exact matching")
	printerGray.Println(" -l  : Return all results")
	printerGray.Println(" -d  : Debug mode")
}

func PrintResult(values []model.UpdateEntry, allResults bool) {
	printerWord := color.New(color.FgHiMagenta).Add(color.Underline)
	printerDef := color.New(color.FgHiYellow)
	printerGray := color.New(color.FgHiBlack)

	if len(values) > 10 && !allResults {
		printerGray.Println("over 10 results, printing first 10...")
		fmt.Println()
		values = values[:9]
	}

	for _, value := range values {
		wordR := []rune(value.Word)
		wordR[0] = unicode.ToUpper(wordR[0])
		word := string(wordR)

		printerWord.Print(word)
		printerGray.Printf(" [%s]\n", value.Pos)
		printerDef.Println(fmt.Sprintf(" » %s", value.Definition))

		if value.Examples == "" {
			continue
		}
		examples := strings.Split(value.Examples, " | ")
		if len(examples) > 3 {
			examples = examples[:2]
		}
		for _, example := range examples {
			printerGray.Println(fmt.Sprintf(" … %s", example))
		}
	}
}
