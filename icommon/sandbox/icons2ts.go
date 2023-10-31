package sandbox

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

const faExtension = ".svg"

var icons2Ts Transformer = func(iconsDir string, destination string) error {

	faStyleToFilenames, iconsCount, err := buildFaStyleToFilenames(iconsDir)

	if err != nil {
		return err
	}

	for faStyle := range faStyleToFilenames {
		if err := os.Mkdir(path.Join(destination, faStyle), 0750); err != nil && !os.IsExist(err) {
			return err
		}
	}

	svgViewBoxRegexp, err1 := regexp.Compile("viewBox=\"([^\"]+)\"")
	svgPathRegexp, err2 := regexp.Compile("d=\"([^\"]+)\"")
	copyrightRegexp, err3 := regexp.Compile("<!--! (.*) -->")
	tsIconFileTmpl, err4 := template.New("tsIconFileTpl").Parse(`// {{.CopyrightNote}}
export const {{.IconName}} = {
    d: "{{.SvgPath}}",
    viewBox: "{{.SvgViewBox}}",
};
`)
	tsIndexFileTmpl, err5 := template.New("tsIndexFileTpl").Parse(`{{range .}}export { {{.}} } from "./{{.}}";
{{end}}`)

	err = errors.Join(err1, err2, err3, err4, err5)

	if err != nil {
		return err
	}

	resultCh := make(chan faResult, iconsCount)

	for faStyle, faNames := range faStyleToFilenames {
		for _, faName := range faNames {
			go func(faStyle string, faName string) {
				iconName := fixVarNameFirstChar(kebabToCamelCase(strings.TrimSuffix(faName, faExtension)))

				contents, err := os.ReadFile(path.Join(iconsDir, faStyle, faName))
				if err != nil {
					resultCh <- faResult{nil, &faError{faResultDetails{faStyle, iconName}, err}}
					return
				}
				svgViewBoxMatch := svgViewBoxRegexp.FindSubmatch(contents)
				if len(svgViewBoxMatch) != 2 {
					resultCh <- faResult{nil, &faError{faResultDetails{faStyle, iconName}, errors.New("invalid svg viewbox regexp match")}}
					return
				}
				svgPathMatch := svgPathRegexp.FindSubmatch(contents)
				if len(svgPathMatch) != 2 {
					resultCh <- faResult{nil, &faError{faResultDetails{faStyle, iconName}, errors.New("invalid svg path regexp match")}}
					return
				}
				copyrightNoteMatch := copyrightRegexp.FindSubmatch(contents)
				if len(copyrightNoteMatch) != 2 {
					resultCh <- faResult{nil, &faError{faResultDetails{faStyle, iconName}, errors.New("invalid copyright note regexp match")}}
					return
				}

				file, err := os.OpenFile(
					path.Join(destination, faStyle, fmt.Sprintf("%s.ts", iconName)),
					os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
					0666,
				)
				if err != nil {
					resultCh <- faResult{nil, &faError{faResultDetails{faStyle, iconName}, err}}
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
					resultCh <- faResult{nil, &faError{faResultDetails{faStyle, iconName}, err}}
					return
				}

				resultCh <- faResult{&faResultDetails{faStyle, iconName}, nil}
			}(faStyle, faName)
		}
	}

	faStylesToNames := map[string][]string{}
	var resultErrors []error
	for i := 0; i < iconsCount; i++ {
		writeResult := <-resultCh
		if writeResult.success != nil {
			if styleNames, ok := faStylesToNames[writeResult.success.faStyle]; ok {
				faStylesToNames[writeResult.success.faStyle] = append(styleNames, writeResult.success.iconName)
			} else {
				faStylesToNames[writeResult.success.faStyle] = []string{writeResult.success.iconName}
			}
		}
		if writeResult.err != nil {
			resultErrors = append(resultErrors, writeResult.err)
		}
	}

	for faStyle, faNames := range faStylesToNames {
		file, err := os.OpenFile(
			path.Join(destination, faStyle, "index.ts"),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0666,
		)

		if err != nil {
			resultErrors = append(resultErrors, err)
			continue
		}

		if err = tsIndexFileTmpl.Execute(file, faNames); err != nil {
			resultErrors = append(resultErrors, err)
		}
	}

	return errors.Join(resultErrors...)
}

type faResult struct {
	success *faResultDetails
	err     *faError
}

type faError struct {
	details faResultDetails
	err     error
}

type faResultDetails struct {
	faStyle  string
	iconName string
}

func (e *faError) Error() string {
	return fmt.Sprintf("Font Awesome icon (%s/%s) error: %s", e.details.faStyle, e.details.iconName, e.err)
}

func buildFaStyleToFilenames(faDir string) (faStyleToFilenames map[string][]string, iconsCount int, err error) {
	faStyleToFilenames = map[string][]string{
		"brands":  {},
		"regular": {},
		"solid":   {},
	}

	for faCategory := range faStyleToFilenames {
		var dirEntries []os.DirEntry
		dirEntries, err = os.ReadDir(path.Join(faDir, faCategory))
		if err != nil {
			return
		}

		for _, dirEntry := range dirEntries {
			faStyleToFilenames[faCategory] = append(faStyleToFilenames[faCategory], dirEntry.Name())
		}
	}

	for _, faPaths := range faStyleToFilenames {
		iconsCount += len(faPaths)
	}

	return
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

type Transformer func(iconsDir, destination string) error
