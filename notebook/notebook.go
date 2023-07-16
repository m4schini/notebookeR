package notebook

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type CellType string

const (
	CellTypeMarkdown CellType = "markdown"
	CellTypeCode     CellType = "code"
)

type FormatterFunc func(v string) string

type File struct {
	Cells         []Cell   `json:"cells"`
	Metadata      Metadata `json:"metadata"`
	Nbformat      int      `json:"nbformat"`
	NbformatMinor int      `json:"nbformat_minor"`
}

func (f *File) ToR() string {
	lines := make([]string, 0)
	for _, cell := range f.Cells {
		var formatter FormatterFunc
		switch cell.CellType {
		case CellTypeMarkdown:
			formatter = func(v string) string {
				return fmt.Sprintf("# %v", v)
			}
			break
		case CellTypeCode:
			formatter = func(v string) string {
				return v
			}
			break
		default:
			panic(fmt.Sprintf("unknown celltype: %v", cell.CellType))
		}

		for _, s := range cell.Source {
			lines = append(lines, formatter(s))
		}
	}
	return strings.Join(lines, "")
}

type Metadata struct {
}

type Cell struct {
	CellType CellType `json:"cell_type"`
	Source   []string `json:"source"`
}

func Parse(reader io.Reader) (*File, error) {

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var f File
	err = json.Unmarshal(b, &f)
	return &f, err
}
