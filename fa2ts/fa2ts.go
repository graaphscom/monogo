package fa2ts

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

const (
	faDir       = "./src/foundations/fontawesome-free-web/svgs"
	faExtension = ".svg"
	destination = "./src/foundations/fontawesome"
)

func Fa2Ts() {
	faNamesByCategory := map[string][]string{
		"brands":  {},
		"regular": {},
		"solid":   {},
	}

	for faCategory := range faNamesByCategory {
		dirEntries, err := os.ReadDir(path.Join(faDir, faCategory))
		if err != nil {
			log.Fatalln(err)
		}

		for _, dirEntry := range dirEntries {
			faNamesByCategory[faCategory] = append(faNamesByCategory[faCategory], dirEntry.Name())
		}
	}

	for faCategory := range faNamesByCategory {
		if err := os.Mkdir(path.Join(destination, faCategory), 0750); err != nil && !os.IsExist(err) {
			log.Fatalln(err)
		}
	}

	iconsCount := 0
	for _, faPaths := range faNamesByCategory {
		iconsCount += len(faPaths)
	}

	svgViewBoxRegexp, err1 := regexp.Compile("viewBox=\"([^\"]+)\"")
	svgPathRegexp, err2 := regexp.Compile("d=\"([^\"]+)\"")
	copyrightRegexp, err3 := regexp.Compile("<!--! (.*) -->")
	tsIconFileTmpl, err4 := template.New("tsIconFileTpl").Parse("// {{.CopyrightNote}}\nexport const {{.IconName}} = {\n  d: \"{{.SvgPath}}\",\n  viewBox: \"{{.SvgViewBox}}\",\n};\n")

	for _, err := range []error{err1, err2, err3, err4} {
		if err != nil {
			log.Fatalln(err)
		}
	}

	resultCh := make(chan FaResult, iconsCount)

	for faCategory, faNames := range faNamesByCategory {
		for _, faName := range faNames {
			go func(faCategory string, faName string) {
				contents, err := os.ReadFile(path.Join(faDir, faCategory, faName))
				if err != nil {
					resultCh <- FaResult{false, &FaError{faCategory, faName, err}}
					return
				}
				svgViewBoxMatch := svgViewBoxRegexp.FindSubmatch(contents)
				if len(svgViewBoxMatch) != 2 {
					resultCh <- FaResult{false, &FaError{faCategory, faName, errors.New("invalid svg viewbox regexp match")}}
					return
				}
				svgPathMatch := svgPathRegexp.FindSubmatch(contents)
				if len(svgPathMatch) != 2 {
					resultCh <- FaResult{false, &FaError{faCategory, faName, errors.New("invalid svg path regexp match")}}
					return
				}
				copyrightNoteMatch := copyrightRegexp.FindSubmatch(contents)
				if len(copyrightNoteMatch) != 2 {
					resultCh <- FaResult{false, &FaError{faCategory, faName, errors.New("invalid copyright note regexp match")}}
					return
				}

				iconName := fixVarNameFirstChar(kebabToCamelCase(strings.TrimSuffix(faName, faExtension)))

				file, err := os.OpenFile(
					path.Join(destination, faCategory, fmt.Sprintf("%s.ts", iconName)),
					os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
					0666,
				)
				if err != nil {
					resultCh <- FaResult{false, &FaError{faCategory, faName, err}}
					return
				}

				if err := tsIconFileTmpl.Execute(file, struct {
					CopyrightNote string
					IconName      string
					SvgPath       string
					SvgViewBox    string
				}{
					string(copyrightNoteMatch[1]),
					iconName,
					string(svgPathMatch[1]),
					string(svgViewBoxMatch[1]),
				}); err != nil {
					resultCh <- FaResult{false, &FaError{faCategory, faName, err}}
					return
				}

				resultCh <- FaResult{true, nil}
			}(faCategory, faName)
		}
	}

	errSummary, hasErrors := "", false
	for i := 0; i < iconsCount; i++ {
		writeResult := <-resultCh
		if !writeResult.success {
			hasErrors = true
			errSummary = strings.Join([]string{errSummary, writeResult.err.Error()}, "\n")
		}
	}

	if hasErrors {
		log.Fatalln(errSummary)
	}

	log.Printf("Successfully transformed %d icons", iconsCount)
}

type FaResult struct {
	success bool
	err     *FaError
}

type FaError struct {
	faCategory string
	faName     string
	err        error
}

func (e *FaError) Error() string {
	return fmt.Sprintf("Font Awesome icon (%s/%s) error: %s", e.faCategory, e.faName, e.err)
}

func kebabToCamelCase(varName string) string {
	hyphenPartRegexp, _ := regexp.Compile(`-\w{1}`)
	return hyphenPartRegexp.ReplaceAllStringFunc(varName, func(s string) string {
		return strings.ToUpper(string(s[1]))
	})
}

func fixVarNameFirstChar(varName string) string {
	notAllowedFirstCharRegexp, _ := regexp.Compile("^[^A-Za-z_]{1}")
	if notAllowedFirstCharRegexp.MatchString(varName) {
		return fmt.Sprintf("__%s", varName)
	}

	return varName
}
