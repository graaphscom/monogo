package unitree

import (
	"fmt"
	"github.com/graaphscom/monogo/icommon/metadata"
	"github.com/graaphscom/monogo/icommon/tsmakers"
	"os"
	"regexp"
	"slices"
	"strings"
)

var treeBuilders = map[string]treeBuilder{
	"boxicons": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			tagsExtractor: func(_ metadata.Store, rawRootName, rawName string) (IconTags, error) {
				return IconTags{
					Search: []string{
						strings.ReplaceAll(strings.TrimSuffix(firstHyphenRegexp.ReplaceAllString(rawName, ""), ".svg"), "-", " "),
					},
					Visual: []string{rawRootName},
				}, nil
			},
			tsMaker: tsmakers.Boxicons,
		},
	},
	"bytesize": iconsTreeBuilder{
		iconNameConverter: iconNameKebabCaseConverter,
		tagsExtractor:     tagsExtractorKebabCase,
		tsMaker:           tsmakers.Bytesize,
	},
	"fluentui": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: func(in string) string {
				return fixVarNameFirstChar(toCamelCase(strings.TrimPrefix(strings.TrimSuffix(in, ".svg"), "ic_fluent"), snakeCaseRegexp))
			},
			treeNameConverter: treeNameSpaceConverter,
			srcSuffix:         "SVG",
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetFluentui(rawRootName)
				matches := fluentuiTagsRegexp.FindStringSubmatch(rawName)
				if os.IsNotExist(err) {
					return IconTags{
						Search: []string{strings.ReplaceAll(matches[1], "_", " ")},
						Visual: matches[2:],
					}, nil
				}
				if err != nil {
					return IconTags{}, err
				}
				return IconTags{
					Search: append(m.Metaphor, strings.ReplaceAll(matches[1], "_", " ")),
					Visual: matches[2:],
				}, nil
			},
			tsMaker: tsmakers.Fluentui,
		},
	},
	"fontawesome": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetFontawesome()
				if err != nil {
					return IconTags{}, err
				}
				rawNameTrimmed := strings.TrimSuffix(rawName, ".svg")
				if icoMeta, ok := m[rawNameTrimmed]; ok {
					return IconTags{
						Search: append(icoMeta.Search.Terms, strings.ReplaceAll(rawNameTrimmed, "-", " ")),
						Visual: []string{rawRootName},
					}, nil
				}
				return IconTags{}, nil
			},
			tsMaker: tsmakers.Fontawesome,
		},
	},
	"octicons": iconsTreeBuilder{
		iconNameConverter: iconNameKebabCaseConverter,
		tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
			m, err := metadata.GetOcticons()
			if err != nil {
				return IconTags{}, err
			}
			matches := octiconsTagsRegexp.FindStringSubmatch(rawName)
			visualTags := slices.DeleteFunc(matches[2:], func(s string) bool {
				return s == ""
			})
			for idx, visualTag := range visualTags {
				visualTags[idx] = strings.Trim(visualTag, "-")
			}
			searchTags := []string{strings.ReplaceAll(matches[1], "-", " ")}
			if _, ok := m[matches[1]]; ok {
				searchTags = append(searchTags, m[matches[1]]...)
			}
			return IconTags{Search: searchTags, Visual: visualTags}, nil
		},
		tsMaker: tsmakers.Octicons,
	},
	"radixui": iconsTreeBuilder{
		iconNameConverter: iconNameKebabCaseConverter,
		tagsExtractor:     tagsExtractorKebabCase,
		tsMaker:           tsmakers.Radixui,
	},
	"remixicon": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			treeNameConverter: treeNameSpaceConverter,
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetRemixicon()
				if err != nil {
					return IconTags{}, err
				}
				matches := remixiconTagsRegexp.FindStringSubmatch(rawName)

				searchTags := []string{strings.ToLower(rawRootName), strings.ReplaceAll(matches[1], "-", " ")}
				if _, ok := m[rawRootName][matches[1]]; ok {
					searchTags = append(searchTags, strings.Split(m[rawRootName][matches[1]], ",")...)
				}
				return IconTags{
					Search: searchTags,
					Visual: []string{strings.Trim(matches[2], "-")},
				}, nil
			},
			tsMaker: tsmakers.Remixicon,
		},
	},
	"unicons": categoriesTreeBuilder{
		iconsTreeBuilder: iconsTreeBuilder{
			iconNameConverter: iconNameKebabCaseConverter,
			tagsExtractor: func(metadata metadata.Store, rawRootName, rawName string) (IconTags, error) {
				m, err := metadata.GetUnicons(rawRootName)
				var _ = m
				if err != nil {
					return IconTags{}, err
				}

				var searchTags []string
				if tags, ok := m[strings.TrimSuffix(rawName, ".svg")]; ok {
					for _, tag := range tags {
						searchTags = append(searchTags, strings.ReplaceAll(tag, "-", " "))
					}
				} else {
					searchTags = []string{strings.ReplaceAll(strings.TrimSuffix(rawName, ".svg"), "-", " ")}
				}

				return IconTags{Search: searchTags, Visual: []string{rawRootName}}, nil
			},
			tsMaker: tsmakers.Unicons,
		},
	},
}

func iconNameKebabCaseConverter(in string) string {
	return fixVarNameFirstChar(toCamelCase(strings.TrimSuffix(in, ".svg"), kebabCaseRegexp))
}

func treeNameSpaceConverter(in string) string {
	return toCamelCase(in, spaceCaseRegexp)
}

func tagsExtractorKebabCase(_ metadata.Store, _, rawName string) (IconTags, error) {
	return IconTags{Search: []string{strings.ReplaceAll(strings.TrimSuffix(rawName, ".svg"), "-", " ")}}, nil
}

func toCamelCase(varName string, initialCaseRegexp *regexp.Regexp) string {
	converted := initialCaseRegexp.ReplaceAllStringFunc(varName, func(s string) string {
		return strings.ToUpper(string(s[1]))
	})
	return strings.ToLower(string(converted[0])) + converted[1:]
}

var kebabCaseRegexp, _ = regexp.Compile(`-\w{1}`)
var snakeCaseRegexp, _ = regexp.Compile(`_\w{1}`)
var spaceCaseRegexp, _ = regexp.Compile(` \w{1}`)
var firstHyphenRegexp, _ = regexp.Compile(`^\w*-`)
var fluentuiTagsRegexp, _ = regexp.Compile(`ic_fluent_(.*?)_(\d*)?_(regular|filled)?(_ltr)?(_rtl)?.svg`)
var octiconsTagsRegexp, _ = regexp.Compile(`(.*?)(-circle)?(-fill)?(-\d*)?.svg`)
var remixiconTagsRegexp, _ = regexp.Compile(`(.*)(-line|-fill)?.svg`)

func fixVarNameFirstChar(varName string) string {
	notAllowedFirstCharRegexp, _ := regexp.Compile("^[^A-Za-z_]{1}")
	if notAllowedFirstCharRegexp.MatchString(varName) {
		return fmt.Sprintf("__%s", varName)
	}

	return varName
}
