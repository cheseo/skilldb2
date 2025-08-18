package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func getTemplate() (*template.Template) {
	t, err := template.New("noooooooo").Funcs(template.FuncMap{
		"Join": strings.Join,
	}).ParseGlob("*.tmpl")
	if err != nil {
		log.Fatal("couldn't parse templates", err)
	}
	return t
}

func employeeForm(r *http.Request) (Employee){
	e := Employee{}
	e.Name = r.FormValue("Name")
	e.Email = r.FormValue("Email")
	e.PhoneNo = r.FormValue("PhoneNo")
	
	pp := []Project{}
	for k, v := range r.PostForm["Project"] {
		pp = append(pp, Project{Name: v,
			Skills: strings.Split(r.PostForm["ProjectSkill"][k], ","),
		})
	}
	e.Projects = pp
	
	ww := []WorkExperience{}
	for k, v := range r.PostForm["WorkCompanyName"] {
		ww = append(ww, WorkExperience{CompanyName: v,
			Title: r.PostForm["WorkTitle"][k],
			Duration: r.PostForm["WorkDuration"][k],
			Skills: strings.Split(r.PostForm["WorkSkills"][k], ","),
		})
	}
	e.WorkExp = ww
	
	tt := []Training{}
	for k, v := range r.PostForm["TrainingName"] {
		tt = append(tt, Training{Name: v,
			Institute: r.PostForm["TrainingInstitute"][k],
			Certificate: r.PostForm["TrainingCertificate"][k],
			Duration: r.PostForm["TrainingDuration"][k],
		})
	}
	e.Training = tt

	ee := []Education{}
	for k, v := range r.PostForm["EducationName"] {
		ee = append(ee, Education{Name: v,
			Duration: r.PostForm["EducationDuration"][k],
		})
	}
	e.Education = ee
	
	
	e.Skills = strings.Split(r.FormValue("Skills"), ",")
	return e
}

func handleError(w http.ResponseWriter, ctx *context){
	log.Println(ctx.err)
	w.WriteHeader(http.StatusInternalServerError)
	sendTemplate(w, "error.tmpl", ctx.err)
}

func JoinComma(s []string) (string){
	return strings.Join(s, ",")
}

func sendTemplate(w http.ResponseWriter, tname string, data any){
	err := getTemplate().ExecuteTemplate(w, tname, data)
	if err != nil {
		log.Println("error executing template", tname, ":", err)
		w.WriteHeader(http.StatusInternalServerError)
		getTemplate().ExecuteTemplate(w, "tmplerror.tmpl", err)
	}
}
