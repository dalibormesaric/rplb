package dashboard

import (
	"html/template"
	"time"
)

var (
	since = time.Now()
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"printSince": func() string {
			return since.Format("2006-01-02 15:04:05")
		},
		"menuSelected": func(selectedMenu, menuItem string) string {
			if selectedMenu == menuItem {
				return " menu-selected"
			}
			return ""
		}}
}
