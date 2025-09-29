package cell

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

// GetTextColor returns the text color of a cell as a hex string
func GetTextColor(f *excelize.File, sheetName, cellReference string) (hex string, err error) {
	styleIndex, err := f.GetCellStyle(sheetName, cellReference)
	if err != nil {
		return "", err
	}

	// Check if Styles and its nested fields exist
	if f.Styles == nil || f.Styles.CellXfs == nil || len(f.Styles.CellXfs.Xf) <= styleIndex {
		return "000000", nil // Default black text
	}

	xf := f.Styles.CellXfs.Xf[styleIndex]
	if xf.FontID == nil {
		return "000000", nil
	}

	fontID := *xf.FontID

	// Check if Fonts exist and fontID is within bounds
	if f.Styles.Fonts == nil || len(f.Styles.Fonts.Font) <= fontID {
		return "000000", nil
	}

	font := f.Styles.Fonts.Font[fontID]
	if font.Color == nil {
		return "000000", nil
	}

	fontColor := font.Color // Font color is a pointer to xlsxColor

	// First try to get RGB color directly
	if fontColor.RGB != "" {
		return strings.TrimPrefix(fontColor.RGB, "FF"), nil
	}

	// If theme color is specified, try to resolve it
	if fontColor.Theme != nil && f.Theme != nil {
		// Use a safer approach to get the theme color similar to getBgColor.go
		func() {
			defer func() {
				if r := recover(); r != nil {
					// If we panic accessing theme colors, just continue to default
				}
			}()

			clrScheme := f.Theme.ThemeElements.ClrScheme
			themeVal := *fontColor.Theme
			var val *string

			// Safely access theme colors
			switch themeVal {
			case 0:
				if clrScheme.Lt1.SysClr.LastClr != "" {
					val = &clrScheme.Lt1.SysClr.LastClr
				}
			case 1:
				if clrScheme.Dk1.SysClr.LastClr != "" {
					val = &clrScheme.Dk1.SysClr.LastClr
				}
			case 2:
				if clrScheme.Lt2.SrgbClr.Val != nil {
					val = clrScheme.Lt2.SrgbClr.Val
				}
			case 3:
				if clrScheme.Dk2.SrgbClr.Val != nil {
					val = clrScheme.Dk2.SrgbClr.Val
				}
			case 4:
				if clrScheme.Accent1.SrgbClr.Val != nil {
					val = clrScheme.Accent1.SrgbClr.Val
				}
			case 5:
				if clrScheme.Accent2.SrgbClr.Val != nil {
					val = clrScheme.Accent2.SrgbClr.Val
				}
			case 6:
				if clrScheme.Accent3.SrgbClr.Val != nil {
					val = clrScheme.Accent3.SrgbClr.Val
				}
			case 7:
				if clrScheme.Accent4.SrgbClr.Val != nil {
					val = clrScheme.Accent4.SrgbClr.Val
				}
			case 8:
				if clrScheme.Accent5.SrgbClr.Val != nil {
					val = clrScheme.Accent5.SrgbClr.Val
				}
			case 9:
				if clrScheme.Accent6.SrgbClr.Val != nil {
					val = clrScheme.Accent6.SrgbClr.Val
				}
			}

			if val != nil && *val != "" {
				hex = strings.TrimPrefix(excelize.ThemeColor(*val, fontColor.Tint), "FF")
			}
		}()

		if hex != "" {
			return hex, nil
		}
	}

	return "000000", nil // Default black text
}
