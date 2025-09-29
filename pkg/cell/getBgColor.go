package cell

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

// GetBgColor returns the background color of a cell as a hex string
func GetBgColor(f *excelize.File, sheetName, cellReference string) (hex string, err error) {
	styleIndex, err := f.GetCellStyle(sheetName, cellReference)
	if err != nil {
		return "", err
	}

	// Check if Styles and its nested fields exist
	if f.Styles == nil || f.Styles.CellXfs == nil || len(f.Styles.CellXfs.Xf) <= styleIndex {
		return "FFFFFF", nil
	}

	xf := f.Styles.CellXfs.Xf[styleIndex]
	if xf.FillID == nil {
		return "FFFFFF", nil
	}

	fillID := *xf.FillID

	// Check if Fills exist and fillID is within bounds
	if f.Styles.Fills == nil || len(f.Styles.Fills.Fill) <= fillID {
		return "FFFFFF", nil
	}

	fill := f.Styles.Fills.Fill[fillID]
	if fill.PatternFill == nil {
		return "FFFFFF", nil
	}

	fgColor := fill.PatternFill.FgColor
	if fgColor != nil {
		// First try to get RGB color directly
		if fgColor.RGB != "" {
			return strings.TrimPrefix(fgColor.RGB, "FF"), nil
		}

		// If theme color is specified, try to resolve it
		if fgColor.Theme != nil && f.Theme != nil {
			// Had panics when using the theme reference so creating a safer approach to get the theme color
			func() {
				defer func() {
					if r := recover(); r != nil {
						// If we panic accessing theme colors, just continue to default
					}
				}()

				clrScheme := f.Theme.ThemeElements.ClrScheme
				themeVal := *fgColor.Theme
				var val *string

				// Safely access theme colors using reflection-like approach
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
					hex = strings.TrimPrefix(excelize.ThemeColor(*val, fgColor.Tint), "FF")
				}
			}()

			if hex != "" {
				return hex, nil
			}
		}
	}
	return "FFFFFF", nil
}
