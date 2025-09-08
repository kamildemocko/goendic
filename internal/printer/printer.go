package printer

import (
	"fmt"
	"goendic/internal/data/model"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

func PrintEmpty() {
	fmt.Println("No result found.")
}

func PrintResult(values []model.UpdateEntry) {
	printerWord := color.New(color.FgHiMagenta).Add(color.Underline)
	printerDef := color.New(color.FgHiYellow)
	printerGray := color.New(color.FgHiBlack)

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
