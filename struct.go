package main

type SignUPBody struct {
	Email string 
	Password string
}

type User struct {
	Id int
	Email string
}

type AccessToCourse struct {
	Email string
}

type Course struct {
	Name string
	Price int
}