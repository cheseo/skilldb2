package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	db *sql.DB
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

	http.HandleFunc("GET /searchSkill", app.searchSkill)
	
	http.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir("static/"))))

	log.Println("listening at :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func (a App) Ctx() (*context){
	return &context{db: a.db}
}

func (a App) CtxWithEid(w http.ResponseWriter, r *http.Request, showErr bool) (*context){
	ctx := a.Ctx()
	eid := r.FormValue("eid")
	ctx.eid, ctx.err = strconv.Atoi(eid)
	if ctx.err != nil && showErr {
		handleError(w, ctx)
	}
	return ctx
}

func (a App) index(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		sendTemplate(w, "404.tmpl", nil)

		return
	}
	sendTemplate(w, "index.tmpl", GetAllEmployees(a.Ctx(), nil))
}

func (a App) employee(w http.ResponseWriter, r *http.Request){
	ctx := a.CtxWithEid(w, r, false)
	if ctx.err != nil {
		w.WriteHeader(http.StatusNotFound)
		sendTemplate(w, "404.tmpl", nil)
		return
	}
	e := GetEmployee(ctx)
	if ctx.err != nil {
		w.WriteHeader(http.StatusNotFound)
		sendTemplate(w, "404.tmpl", ctx.err)
		return
	}
	sendTemplate(w, "employee.tmpl", e)
}

func (a App) newEmpPage(w http.ResponseWriter, r *http.Request){
	sendTemplate(w, "editEmployee.tmpl", Employee{})
}

func (a App) newEmp(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	log.Println(r.PostForm)
	e := employeeForm(r)
	ctx := a.Ctx()
	InsertEmployee(ctx, e)
	if ctx.err != nil {
		handleError(w, ctx)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a App) editEmpPage(w http.ResponseWriter, r *http.Request){
	ctx := a.CtxWithEid(w, r, true)
	if ctx.err != nil {
		return
	}
	e := GetEmployee(ctx)
	if ctx.err != nil {
		handleError(w, ctx)
		return
	}

	sendTemplate(w, "editEmployee.tmpl", e)
}

func (a App) editEmp(w http.ResponseWriter, r *http.Request){
	ctx := a.CtxWithEid(w, r, true)
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

func (a App) delEmp(w http.ResponseWriter, r *http.Request){
	ctx := a.CtxWithEid(w, r, true)
	if ctx.err != nil {
		return
	}
	DeleteEmployee(ctx)
	if ctx.err != nil {
		handleError(w, ctx)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (a App) searchSkill(w http.ResponseWriter, r *http.Request){
	s := r.FormValue("skills")
	if s == "" {
		sendTemplate(w, "searchFailed.tmpl", "No skills given.")
		return
	}
	ss := strings.Split(s, ",")
	ctx := a.Ctx()
	e := SearchSkills(ctx, ss)
	if ctx.err != nil {
		handleError(w, ctx)
		return
	}
	sendTemplate(w, "index.tmpl", GetAllEmployees(a.Ctx(), e))
}
