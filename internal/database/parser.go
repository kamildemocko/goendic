package database

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type LexicalResource struct {
	Lexicons []Lexicon `xml:"Lexicon"`
}

type Lexicon struct {
	LexicalEntries []LexicalEntry `xml:"LexicalEntry"`
	Synsets        []Synset       `xml:"Synset"`
}

type LexicalEntry struct {
	XMLName xml.Name `xml:"LexicalEntry"`
	ID      string   `xml:"id,attr"`
	Lemma   Lemma    `xml:"Lemma"`
	Senses  []Sense  `xml:"Sense"`
}

type Lemma struct {
	WrittenForm  string `xml:"writtenForm,attr"`
	PartOfSpeech string `xml:"partOfSpeech,attr"`
}

type Sense struct {
	Synset string `xml:"synset,attr"`
}

type Synset struct {
	XMLName    xml.Name `xml:"Synset"`
	ID         string   `xml:"id,attr"`
	POS        string   `xml:"partOfSpeech,attr"`
	Definition string   `xml:"Definition"`
	Examples   []string `xml:"Example"`
}

func ParseXML(filePath string) ([][4]string, error) {
	log.Println("parsing")

	lexicalFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer lexicalFile.Close()

	lexicalBytes, err := io.ReadAll(lexicalFile)
	if err != nil {
		return nil, err
	}

	var resource LexicalResource
	err = xml.Unmarshal(lexicalBytes, &resource)
	if err != nil {
		return nil, err
	}

	if len(resource.Lexicons) == 0 {
		return nil, fmt.Errorf("no lexicon found in XML")
	}

	lexicon := resource.Lexicons[0]
	log.Printf("found %d entries", len(lexicon.LexicalEntries))

	synsetMap := make(map[string]Synset)
	for _, s := range lexicon.Synsets {
		synsetMap[s.ID] = s
	}

	// word, pos, defintion, example
	values := make([][4]string, 0)
	for _, entry := range lexicon.LexicalEntries {
		word := entry.Lemma.WrittenForm
		pos := getFullPartOfSpeech(entry.Lemma.PartOfSpeech)

		for _, sense := range entry.Senses {
			syn, ok := synsetMap[sense.Synset]
			if !ok {
				continue
			}

			definition := syn.Definition
			examples := ""
			if len(syn.Examples) > 0 {
				for i, e := range syn.Examples {
					if i > 0 {
						examples += " | "
					}
					examples += e
				}
			}

			values = append(values, [4]string{word, pos, definition, examples})
		}
	}

	log.Println("success")

	return values, nil
}

func getFullPartOfSpeech(abbr string) string {
	switch abbr {
	case "n":
		return "noun"
	case "v":
		return "verb"
	case "a", "s":
		return "adjective"
	case "r":
		return "adverb"
	default:
		return abbr
	}
}
