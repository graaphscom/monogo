package sandbox

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFa2TsSuccess(t *testing.T) {
	resultPath, initialPath, nextPath, err := preparePaths()
	if err != nil {
		t.Fatal(err)
	}

	err = icons2Ts(initialPath, resultPath)
	if err != nil {
		t.Fatal(err)
	}

	resultFS := os.DirFS(resultPath)
	var actualDirLayout string
	fs.WalkDir(resultFS, ".", func(p string, d fs.DirEntry, err error) error {
		actualDirLayout = fmt.Sprintln(actualDirLayout, strings.Repeat(" ", strings.Count(p, string(os.PathSeparator))), d)
		return nil
	})
	const expectedDirLayout = `  d ./
  d brands/
   - __42Group.ts
   - __500px.ts
   - index.ts
  d regular/
   - addressBook.ts
   - addressCard.ts
   - bellSlash.ts
   - index.ts
  d solid/
   - index.ts
   - nextDeleted.ts
   - nextModified.ts
`

	assert.Equal(t, expectedDirLayout, actualDirLayout)

	err = icons2Ts(nextPath, resultPath)
	if err != nil {
		t.Fatal(err)
	}

	var actualDirLayout1 string
	fs.WalkDir(resultFS, ".", func(p string, d fs.DirEntry, err error) error {
		actualDirLayout1 = fmt.Sprintln(actualDirLayout1, strings.Repeat(" ", strings.Count(p, string(os.PathSeparator))), d)
		return nil
	})
	const expectedDirLayout1 = `  d ./
  d brands/
   - __42Group.ts
   - __500px.ts
   - index.ts
  d regular/
   - addressBook.ts
   - addressCard.ts
   - bellSlash.ts
   - index.ts
  d solid/
   - index.ts
   - nextAdded.ts
   - nextDeleted.ts
   - nextModified.ts
`

	assert.Equal(t, expectedDirLayout1, actualDirLayout1)
}

func preparePaths() (resultPath, initialPath, nextPath string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	resultPath = path.Join(wd, "testdata", "successful-result")
	err = os.RemoveAll(resultPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return
	}
	err = os.Mkdir(resultPath, 0750)
	if err != nil {
		return
	}

	initialPath = path.Join(wd, "testdata", "initial-fontawesome", "svgs")
	nextPath = path.Join(wd, "testdata", "next-fontawesome", "svgs")

	return
}

func TestFa2TsError(t *testing.T) {

}
