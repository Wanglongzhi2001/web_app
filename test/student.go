package main

type student struct {
	Id    int
	Name  string
	Sex   string
	Class string
}

func newStudent(id int, name, sex, class string) *student {
	stu := student{
		Id:    id,
		Name:  name,
		Sex:   sex,
		Class: class,
	}
	return &stu
}
