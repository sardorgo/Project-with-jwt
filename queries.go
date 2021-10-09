package main

var SQL_INSERT_USER = `

	insert into users (
		email, password
	) values ($1, $2)
	returning
		user_id,
		email

`

var SQL_SELECT_ACTIVATION = `
	select id from activation where user_id = $1
`

var UPDATE_ACTIVATE = `
	update users set activated_at = current_timestamp
	where user_id = (
		select user_id from activation where id = $1
	)
`

var SELECT_COURSES = `
	select name, price from courses
`

var SELECT_COURSES_ACTIVATED_AT = `
	select u.activated_at from users as u
	where u.email = $1
`