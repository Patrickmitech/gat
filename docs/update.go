package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/koki-develop/gat/pkg/printer"
)

var (
	src = `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}`
)

func String(s string) *string {
	return &s
}

func main() {
	updateLanguages()
	updateThemes()
	updateFormats()
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func OrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func Map[T, I any](ts []T, f func(t T) I) []I {
	is := []I{}
	for _, t := range ts {
		is = append(is, f(t))
	}
	return is
}

func updateLanguages() {
	f := Must(os.Create("docs/languages.md"))
	defer f.Close()

	Must(f.WriteString("# Languages\n\n"))

	Must(f.WriteString("| Language | Aliases |\n"))
	Must(f.WriteString("| --- | --- |\n"))

	for _, l := range lexers.GlobalLexerRegistry.Lexers {
		cfg := l.Config()
		Must(f.WriteString(fmt.Sprintf("| `%s` ", cfg.Name)))

		if len(cfg.Aliases) > 0 {
			Must(f.WriteString(
				fmt.Sprintf(
					"| %s |",
					strings.Join(
						Map(
							cfg.Aliases,
							func(a string) string { return fmt.Sprintf("`%s`", a) },
						),
						", ",
					),
				),
			))
		} else {
			Must(f.WriteString("| |"))
		}
		Must(f.WriteString("\n"))
	}
}

func updateFormats() {
	f := Must(os.Create("docs/formats.md"))
	defer f.Close()

	Must(f.WriteString("# Output Formats\n\n"))

	for _, format := range formatters.Names() {
		Must(f.WriteString(fmt.Sprintf("- [`%s`](#%s)\n", format, format)))
	}
	Must(f.WriteString("\n"))

	for _, format := range formatters.Names() {
		Must(f.WriteString(fmt.Sprintf("## `%s`\n\n", format)))

		p := printer.New(&printer.PrinterConfig{
			Format: format,
			Theme:  printer.DefaultTheme,
		})

		b := new(bytes.Buffer)
		OrPanic(p.Print(&printer.PrintInput{
			In:       strings.NewReader(src),
			Out:      b,
			Filename: String("main.go"),
		}))

		Must(f.WriteString(fmt.Sprintf("```%s\n", strings.TrimSuffix(format, "-min"))))
		if strings.HasPrefix(format, "terminal") {
			Must(f.WriteString(strings.Trim(strings.ReplaceAll(fmt.Sprintf("%#v", strings.TrimSpace(b.String())), "\\n", "\n"), "\"")))
		} else {
			Must(f.WriteString(strings.TrimSpace(b.String())))
		}
		Must(f.WriteString("\n"))
		Must(f.WriteString("```\n"))

		Must(f.WriteString("\n"))
	}
}

func updateThemes() {
	f := Must(os.Create("docs/themes.md"))
	defer f.Close()

	Must(f.WriteString("# Highlight Themes\n\n"))

	for _, s := range styles.Names() {
		Must(f.WriteString(fmt.Sprintf("- [`%s`](#%s)\n", s, s)))
	}
	Must(f.WriteString("\n"))

	for _, s := range styles.Names() {
		Must(f.WriteString(fmt.Sprintf("## `%s`\n\n", s)))

		p := printer.New(&printer.PrinterConfig{
			Format: "svg",
			Theme:  s,
		})

		b := new(bytes.Buffer)
		OrPanic(p.Print(&printer.PrintInput{
			In:       strings.NewReader(src),
			Out:      b,
			Filename: String("main.go"),
		}))

		img := Must(os.Create(fmt.Sprintf("./docs/themes/%s.svg", s)))
		defer img.Close()
		Must(img.Write(b.Bytes()))

		Must(f.WriteString(fmt.Sprintf("![%s](./themes/%s.svg)\n\n", s, s)))
	}
}
