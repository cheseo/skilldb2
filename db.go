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
	Skills    []string
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
	Skills      []string
}

const ProjectCreate = `
create table project(
pid integer primary key,
eid integer not null,
name text,
foreign key (eid) references employee(eid) on delete cascade
);
`

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
	Skills             []string
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

const SkillCreate = `
create table skill(
eid integer not null,
name text,
foreign key (eid) references employee(eid) on delete cascade,
primary key (eid, name)
);
`

const SkillViewCreate = `
create view allskills as
select eid, name from skill union
select eid, name from projectskill union
select eid, name from workskill;
`
