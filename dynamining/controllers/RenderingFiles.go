package controllers

import (
	"net/http"
	"html/template"
	"dynamining/dtools"
	"dynamining/setting"
	"encoding/json"
)

var templates map[string]*template.Template

func renderingJSON(w http.ResponseWriter, data interface{}) {

	j, err := json.Marshal(data)
	if err != nil {
		dtools.Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(j)
}

func rendering(w http.ResponseWriter, r *http.Request, runPage string, backgroundPage string, param interface{}) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	switch backgroundPage {
	case basicPage:
		templates[runPage] = template.Must(template.ParseFiles(tempRoot + basicPage +".tpl", tempRoot + runPage + ".tpl"))
		templ, ok := templates[runPage]
		if !ok {
			dtools.Debug(ok)
			http.Error(w, "The template does not exit.", http.StatusInternalServerError)
		}
		err := templ.ExecuteTemplate(w, basicPage, param)
		if err != nil {
			dtools.Debug(ok)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case basePage:
		templates[runPage] = template.Must(template.ParseFiles(tempRoot + basePage +".tpl", tempRoot + runPage + ".tpl"))
		templ, ok := templates[runPage]
		if !ok {
			dtools.Debug(ok)
			http.Error(w, "The template does not exit.", http.StatusInternalServerError)
		}
		err := templ.ExecuteTemplate(w, basePage, param)
		if err != nil {
			dtools.Debug(ok)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		runPage = setting.TemplateConfig.Error404Page
		templates[runPage] = template.Must(template.ParseFiles(tempRoot + basePage +".tpl", tempRoot + runPage + ".tpl"))
		templ, ok := templates[runPage]
		if !ok {
			dtools.Debug(ok)
			http.Error(w, "The template does not exit.", http.StatusInternalServerError)
		}
		err := templ.ExecuteTemplate(w, basePage, param)
		if err != nil {
			dtools.Debug(ok)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// 자주 호출하게 되는 페이지
// 시작 페이지
func startRendering(w http.ResponseWriter, r *http.Request) {
	rendering(w, r, setting.TemplateConfig.StartPage, setting.TemplateConfig.BasicPage, nil)
}

// 가입 페이지
func joinRendering(w http.ResponseWriter, r *http.Request) {
	rendering(w, r, setting.TemplateConfig.StartPage, setting.TemplateConfig.BasicPage, nil)
}
/*
	switch viewModel.(type) {
	case *models.TSysUser:
		fmt.Println("*TSysUser:", viewModel)
		p := viewModel.(*models.TSysUser)
		fmt.Println("Id:", p.Id)
		fmt.Println("Token:", p.IdToken)
	}
*/
