package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type context struct {
	err error
	db *sql.DB
	eid int
}

func BeginTransaction(ctx *context){
	_, ctx.err = ctx.db.Exec("begin transaction;")
}

func Commit(ctx *context){
	if ctx.err != nil {
		ctx.db.Exec("rollback;")
	}else {
		_, ctx.err = ctx.db.Exec("commit;")
	}
}

func GetAllEmployees(ctx *context, filterEid []int) (ee []Employee){
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	if len(filterEid) > 0{
		v := "?" + strings.Repeat(",?", len(filterEid)-1)
		a := make([]any, len(filterEid))
		for i, v := range filterEid {
			a[i] = v
		}
		query := `select eid, name, email, phoneno from employee where eid in (` + v + `);`
		r, ctx.err = ctx.db.Query(query, a...)
	} else {
		query := `select eid, name, email, phoneno from employee;`
		r, ctx.err = ctx.db.Query(query)
	}
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetAllEmployees: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		e := Employee{}
		if ctx.err =r.Scan(&e.Eid, &e.Name, &e.Email, &e.PhoneNo); ctx.err != nil {
			ctx.err = fmt.Errorf("GetAllEmployees: %w", ctx.err)
			break
		}
		ee = append(ee, e)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err = r.Err()
		return
	}
	return
}

func InsertEmployee(ctx *context, e Employee) {
	if ctx.err != nil {
		return
	}
	query := `insert into employee(Name, Email, PhoneNo) values ( ?, ?, ? ) returning eid;`
	ctx.err = ctx.db.QueryRow(query, e.Name, e.Email, e.PhoneNo).Scan(&ctx.eid)
	log.Println("got eid:", ctx.eid)
	e.Eid = ctx.eid
	InsertProjects(ctx, e.Projects)
	InsertWorkExperience(ctx, e.WorkExp)
	InsertTraining(ctx, e.Training)
	InsertEducation(ctx, e.Education)
	InsertSkill(ctx, e.Skills)
}

func DeleteEmployee(ctx *context) {

	if ctx.err != nil {
		return
	}
	_, ctx.err = ctx.db.Exec("delete from employee where eid = ?", ctx.eid)
}

func GetEmployee(ctx *context)(e Employee) {

	if(ctx.err != nil){
		return
	}
	query := `select Name, Email, PhoneNo from employee where eid = ?;`
	e.Eid = ctx.eid
	ctx.db.QueryRow(query, ctx.eid).Scan(&e.Name, &e.Email, &e.PhoneNo)
	e.Projects = GetProjects(ctx)
	e.WorkExp = GetWorkExperience(ctx)
	e.Training = GetTraining(ctx)
	e.Education = GetEducation(ctx)
	e.Skills = GetSkills(ctx)
	log.Println("got employee: ", e)
	return
}

func InsertProjects(ctx *context, pp []Project) {
	for _, p := range pp {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertProjects: %w", ctx.err)
			return
		}
		query := `insert into project(eid, name, url) values (?, ?, ?) returning pid;`
		ctx.db.QueryRow(query, ctx.eid, p.Name, p.Url).Scan(&p.Pid)
		InsertProjectSkill(ctx, p.Pid, p.Skills)
	}
}

func InsertProjectSkill(ctx *context, pid int, skills []string){
	for _, ps := range skills {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertProjectSkill: %w", ctx.err)
			return
		}
		query := `insert into projectskill(eid, pid, name) values (?, ?, ?);`
		_, ctx.err = ctx.db.Exec(query, ctx.eid, pid, ps)
	}
}

func GetProjects(ctx *context) (p []Project) {
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select pid, name, url from project where eid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetProjects: %w", ctx.err)
		return
	}
	defer r.Close()
	log.Println("rows is", r)
	for r.Next() {
		pp := Project{}
		if ctx.err =r.Scan(&pp.Pid, &pp.Name, &pp.Url); ctx.err != nil {
			log.Println("ERROR GetProjects: %w", ctx.err)
			ctx.err = fmt.Errorf("GetProjects: %w", ctx.err)
			break
		}
		pp.Skills = GetProjectSkills(ctx, pp.Pid)
		log.Print("GetProjects: appending", pp, "to", p)
		p = append(p, pp)
	}
	log.Println("r.Next: ", r.Next(), "r.Err(): ",r.Err())
	if r.Err() != nil && ctx.err == nil{
		ctx.err =r.Err()
		return
	}
	log.Println("getProjects returning:", p)
	return
}

func GetProjectSkills(ctx *context, pid int) (s []string){
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select name from projectskill where eid = ? and  pid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid, pid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetProjectSkills: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		ps := ""
		if ctx.err =r.Scan(&ps); ctx.err != nil {
			ctx.err = fmt.Errorf("GetProjectSkills: %w", ctx.err)
			break
		}
		s = append(s, ps)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err =r.Err()
		return
	}
	return
}

func InsertWorkExperience(ctx *context, ww []WorkExperience) {
	for _, w := range ww {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertWorkExperience: %w", ctx.err)
			return
		}
		query := `insert into workexperience(eid, CompanyName, Title, Duration) values ( ?, ?, ?, ? ) returning wid;`
		ctx.db.QueryRow(query, ctx.eid, w.CompanyName, w.Title, w.Duration).Scan(&w.Wid)
		InsertWorkSkill(ctx, w.Wid, w.Skills)
	}
}

func GetWorkExperience(ctx *context) (ww []WorkExperience) {
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select wid, companyname, title, duration from workexperience where eid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetWorkExperience: db.Query: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		w := WorkExperience{}
		if ctx.err =r.Scan(&w.Wid, &w.CompanyName, &w.Title, &w.Duration); ctx.err != nil {
			ctx.err = fmt.Errorf("GetWorkExperience:r.Scan: %w", ctx.err)
			break
		}
		w.Skills = GetWorkSkills(ctx, w.Wid)
		ww = append(ww, w)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err =r.Err()
		return
	}
	return
}

func InsertWorkSkill(ctx *context, wid int, skills []string) {
	for _, s := range skills {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertWorkSkill: %w", ctx.err)
			return
		}
		query := `insert into workskill(eid, wid, name) values(?, ?, ?);`
		_, ctx.err = ctx.db.Exec(query, ctx.eid, wid, s)
	}
}

func GetWorkSkills(ctx *context, wid int) (s []string){
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select name from workskill where eid = ? and wid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid, wid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetWorkSkills: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		ws := ""
		if ctx.err =r.Scan(&ws); ctx.err != nil {
			ctx.err = fmt.Errorf("GetWorkSkills: %w", ctx.err)
			break
		}
		s = append(s, ws)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err =r.Err()
		return
	}
	return
}

func InsertTraining(ctx *context, tt []Training) {
	for _, t := range tt {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertTraining: %w", ctx.err)
			return
		}
		query := `insert into training( eid, Name, Institute, Certificate, certurl, Duration) values ( ?, ?, ?, ?, ? );`
		_, ctx.err = ctx.db.Exec(query, ctx.eid, t.Name, t.Institute, t.Certificate, t.CertUrl, t.Duration)
	}
}

func GetTraining(ctx *context) (tt []Training) {
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select tid, name, institute, certificate, certurl, duration from training where eid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetTraining: db.Query: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		t := Training{}
		if ctx.err =r.Scan(&t.Tid, &t.Name, &t.Institute, &t.Certificate, &t.CertUrl, &t.Duration); ctx.err != nil {
			ctx.err = fmt.Errorf("GetTraining:r.Scan: %w", ctx.err)
			break
		}
		tt = append(tt, t)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err =r.Err()
		return
	}
	return
}

func InsertEducation(ctx *context, ee []Education) {
	for _, e := range ee {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertEducation: %w", ctx.err)
			return
		}
		query := `insert into education( eid, Name, Duration) values ( ?, ?, ? );`
		_, ctx.err = ctx.db.Exec(query, ctx.eid, e.Name, e.Duration)
	}
}

func GetEducation(ctx *context) (ee []Education){
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select name, duration from education where eid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetEducation: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		e := Education{}
		if ctx.err =r.Scan(&e.Name, &e.Duration); ctx.err != nil {
			ctx.err = fmt.Errorf("GetEducation: %w", ctx.err)
			break
		}
		ee = append(ee, e)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err =r.Err()
		return
	}
	return
}

func InsertSkill(ctx *context, ss []string) {
	for _, s := range ss {
		if(ctx.err != nil){
			ctx.err = fmt.Errorf("InsertSkill: %w", ctx.err)
			return
		}
		query := `insert into skill( eid, Name) values ( ?, ? );`
		_, ctx.err = ctx.db.Exec(query, ctx.eid, s)
	}
}

func GetSkills(ctx *context) (ss []string){
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	query := `select name from skill where eid = ?;`
	r, ctx.err = ctx.db.Query(query, ctx.eid)
	if ctx.err != nil {
		ctx.err = fmt.Errorf("GetSkills: %w", ctx.err)
		return
	}
	defer r.Close()
	for r.Next() {
		s := ""
		if ctx.err =r.Scan(&s); ctx.err != nil {
			ctx.err = fmt.Errorf("GetSkills: %w", ctx.err)
			break
		}
		ss = append(ss, s)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err = fmt.Errorf("GetSkills: %w",r.Err())
		return
	}
	return
}

func SearchSkills(ctx *context, s []string) (eids []int){
	var r *sql.Rows
	if(ctx.err != nil){
		return
	}
	if len(s) < 1 {
		ctx.err = fmt.Errorf("SearchSkills: provide atleast one skill to search")
		return
	}
	v := "?" + strings.Repeat(",?", len(s)-1)
	a := make([]any, len(s))
	for i, v := range s {
		a[i] = v
	}
	
	query := `select eid from allskills where name in (` + v + `);`
	r, ctx.err = ctx.db.Query(query, a...)
	defer r.Close()
	for r.Next(){
		e := 0
		if ctx.err =r.Scan(&e); ctx.err != nil {
			ctx.err = fmt.Errorf("SearchSkills: %w", ctx.err)
			return
		}
		eids = append(eids, e)
	}
	if r.Err() != nil && ctx.err == nil{
		ctx.err = fmt.Errorf("SearchSkills: %w", r.Err())
		return
	}
	return
}
