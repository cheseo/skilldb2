package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	db *sql.DB
}

func getTemplate() (*template.Template) {
	t, err := template.ParseGlob("*.tmpl")
	if err != nil {
		log.Fatal("couldn't parse templates", err)
	}
	return t
}

func startServer(db *sql.DB){
	app := App{db: db}
	http.HandleFunc("/", app.index)
	http.HandleFunc("/employee", app.employee)
	http.HandleFunc("GET /newemployee", app.newEmpPage)
	http.HandleFunc("POST /newemployee", app.newEmp)
	http.HandleFunc("GET /editEmployee", app.editEmpPage)
	http.HandleFunc("POST /editEmployee", app.editEmp)
	http.HandleFunc("POST /deleteEmployee", app.delEmp)
	http.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir("static/"))))
	log.Println("listening at :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func (a App) Ctx() (*context){
	return &context{db: a.db}
}

func (a App) index(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		getTemplate().ExecuteTemplate(w, "404.tmpl", nil)
		return
	}
	getTemplate().ExecuteTemplate(w, "index.tmpl", GetAllEmployees(a.Ctx()))
}

func (a App) employee(w http.ResponseWriter, r *http.Request){
	ctx := a.getCtx(w, r, false)
	if ctx.err != nil {
		w.WriteHeader(http.StatusNotFound)
		getTemplate().ExecuteTemplate(w, "404.tmpl", nil)
		return
	}
	e := GetEmployee(ctx, true)
	if ctx.err != nil {
		w.WriteHeader(http.StatusNotFound)
		getTemplate().ExecuteTemplate(w, "404.tmpl", ctx.err)
	}
	getTemplate().ExecuteTemplate(w, "employee.tmpl", e)
}

func (a App) newEmpPage(w http.ResponseWriter, r *http.Request){
	getTemplate().ExecuteTemplate(w, "newemployee.tmpl", Employee{})
}

func (a App) newEmp(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	log.Println(r.PostForm)
	e := employeeForm(r)
	ctx := a.Ctx()
	InsertEmployee(ctx, e)
	if ctx.err != nil {
		handleError(w, r, ctx)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func employeeForm(r *http.Request) (Employee){
	e := Employee{}
	e.Name = r.FormValue("Name")
	e.Email = r.FormValue("Email")
	e.PhoneNo = r.FormValue("PhoneNo")
	
	pp := []Project{}
	for k, v := range r.PostForm["Project"] {
		pp = append(pp, Project{Name: v,
			Skills: func() (p []ProjectSkill) {
				for _, v := range strings.Split(r.PostForm["ProjectSkill"][k], ",") {
					p = append(p, ProjectSkill{Name: v})
				}
				return
			}(),
		})
	}
	e.Projects = pp
	
	ww := []WorkExperience{}
	for k, v := range r.PostForm["WorkCompanyName"] {
		ww = append(ww, WorkExperience{CompanyName: v,
			Title: r.PostForm["WorkTitle"][k],
			Duration: r.PostForm["WorkDuration"][k],
			Skills: func() (p []WorkSkill) {
				for _, v := range strings.Split(r.PostForm["WorkSkills"][k], ",") {
					p = append(p, WorkSkill{Name: v})
				}
				return
			}(),
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
	
	
	e.Skills = func() (p []Skill) {
		for _, v := range strings.Split(r.FormValue("Skills"), ",") {
			p = append(p, Skill{Name: v})
		}
		return
	}()
	return e
}

func handleError(w http.ResponseWriter, r *http.Request, ctx *context){
	log.Println(ctx.err)
	w.WriteHeader(http.StatusInternalServerError)
	getTemplate().ExecuteTemplate(w, "error.tmpl", ctx.err)
}

func (a App) editEmpPage(w http.ResponseWriter, r *http.Request){
	ctx := a.getCtx(w, r, true)
	if ctx.err != nil {
		return
	}
	e := GetEmployee(ctx, true)
	if ctx.err != nil {
		handleError(w, r, ctx)
		return
	}

	getTemplate().ExecuteTemplate(w, "editEmployee.tmpl", e)
}

func (a App) editEmp(w http.ResponseWriter, r *http.Request){
	ctx := a.getCtx(w, r, true)
	if ctx.err != nil {
		return
	}
	e := employeeForm(r)

	BeginTransaction(ctx)
	DeleteEmployee(ctx)
	InsertEmployee(ctx, e)
	Commit(ctx)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a App) getCtx(w http.ResponseWriter, r *http.Request, showErr bool) (*context){
	ctx := a.Ctx()
	eid := r.FormValue("eid")
	ctx.eid, ctx.err = strconv.Atoi(eid)
	if ctx.err != nil && showErr {
		handleError(w, r, ctx)
	}
	return ctx
}

func (a App) delEmp(w http.ResponseWriter, r *http.Request){
	ctx := a.getCtx(w, r, true)
	if ctx.err != nil {
		return
	}
	DeleteEmployee(ctx)
	if ctx.err != nil {
		handleError(w, r, ctx)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
