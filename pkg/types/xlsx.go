package types

import "github.com/xuri/excelize/v2"

type StyledCell struct {
	Value string          `json:"value"`
	Style *excelize.Style `json:"style,omitempty"`
}

type Sheet struct {
	SheetName string         `json:"sheetName"`
	Cells     [][]StyledCell `json:"cells"`
}

type SheetMinimal struct {
	SheetName string `json:"sheetName"`
	SheetTabColor string `json:"sheetTabColor,omitempty"`
}

type Attributes struct {
	SheetMinimals []SheetMinimal `json:"sheetMinimals"`
	TextColors   []string        `json:"textColors"`
	BgColors     []string        `json:"bgColors"`
}
