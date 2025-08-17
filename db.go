package main

type Employee struct {
	Eid      int
	Name     string
	Email    string
	PhoneNo  string
	Projects []Project
	WorkExp  []WorkExperience
	Training []Training
	Education []Education
	Skills    []Skill
}

const EmployeeCreate = `
create table employee(
eid integer primary key,
name text,
email text,
phoneno text
);
`
	
type Project struct {
	Pid         int
	Name        string
	Skills      []ProjectSkill
}

const ProjectCreate = `
create table project(
pid integer primary key,
eid integer not null,
name text,
foreign key (eid) references employee(eid) on delete cascade
);
`

type ProjectSkill struct {
	Name      string
}

const ProjectSkillCreate = `
create table projectskill(
eid integer not null,
pid integer not null,
name text,
foreign key (eid) references employee(eid) on delete cascade,
foreign key (pid) references project(pid) on delete cascade,
primary key (eid, pid, name)
);
`

type WorkExperience struct {
	Wid                int
	CompanyName        string
	Title              string
	Duration           string
	Skills             []WorkSkill
}

const WorkExperienceCreate = `
create table workexperience(
wid integer primary key,
eid integer not null,
companyname text,
title text,
duration text,
foreign key (eid) references employee(eid) on delete cascade
);
`

type WorkSkill struct {
	Name      string
}

const WorkSkillCreate = `
create table workskill(
eid integer not null,
wid integer not null,
name text,
foreign key (eid) references employee(eid) on delete cascade,
foreign key (wid) references workexperience(wid) on delete cascade,
primary key (eid, wid, name)
);
`

type Training struct {
	Tid          int
	Name         string
	Institute    string
	Certificate  string
	Duration     string
}

const TrainingCreate = `
create table training(
tid integer primary key,
eid integer not null,
name text,
institute text,
certificate text,
duration text,
foreign key (eid) references employee(eid) on delete cascade
);
`

type Education struct {
	Eid           int
	Name          string
	Duration      string
}
const EducationCreate = `
create table education(
edid integer primary key,
eid integer not null,
name text,
duration text,
foreign key (eid) references employee(eid) on delete cascade
);
`

type Skill struct {
	Name      string
}

const SkillCreate = `
create table skill(
eid integer not null,
name text,
foreign key (eid) references employee(eid) on delete cascade,
primary key (eid, name)
);
`
