package handlers

import (
	"html/template"
	"sync"
)

var (
	rootTemplates *template.Template
	viewTemplate  *template.Template
	once          sync.Once
)

func InitTemplates() {
	once.Do(func() {
		rootTemplates = template.Must(template.New("").Funcs(template.FuncMap{
			"seq": func(start, end int) []int {
				if start > end {
					return []int{}
				}
				s := make([]int, end-start+1)
				for i := range s {
					s[i] = start + i
				}
				return s
			},
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
			"max": func(a, b int) int {
				return map[bool]int{true: a, false: b}[a > b]
			},
			"min": func(a, b int) int {
				return map[bool]int{true: a, false: b}[a < b]
			},
		}).ParseFiles(
			"/app/web/templates/index.html",
			"/app/web/templates/pagination.html",
			"/app/web/templates/profile.html",
		))
		viewTemplate = template.Must(template.ParseFiles("/app/web/templates/view.html"))
	})
}
