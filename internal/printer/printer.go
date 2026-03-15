package printer

import (
	"embed"
	"flag"
	"fmt"
	"strings"
	"unicode"

	"github.com/kamildemocko/goendic/v2/internal/data/model"

	"github.com/fatih/color"
)

func PrintFirstTimeDB() {
	fmt.Println("Preparing DB...")
}

func PrintUpdateDB() {
	fmt.Println("Updating DB...")
}

func PrintDbUpdated() {
	fmt.Println("DB was updated")
}

func PrintOldDB() {
	printerGray := color.New(color.FgHiBlack)
	printerGray.Println("Consider updating DB with the -u flag")
	fmt.Println()
}

func PrintEmpty() {
	printerRed := color.New(color.FgRed)
	printerRed.Println("No result found.")
}

func SetupPrintUsage() {
	printerGray := color.New(color.FgHiBlack)

	flag.Usage = func() {
		fmt.Println("Error: No search word provided")
		fmt.Print("Usage: endic")
		printerGray.Print(" [OPTIONS] ")
		fmt.Println("WORD")
		fmt.Println("Options: ")
		printerGray.Println(" -e  : Use exact matching")
		printerGray.Println(" -l  : Return all results")
		printerGray.Println(" -u  : Update database")
		printerGray.Println(" -v  : Show version")
	}
}

func PrintVersion(version embed.FS) {
	fmt.Println(version)
}

func PrintOfflineState() {
	printerWord := color.New(color.FgHiMagenta).Add(color.Underline)

	printerWord.Println("Please connect to network to create or update database")
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
