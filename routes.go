package main

import (
	"fmt"
	"time"
	"bytes"
	"io/ioutil"
	"net/smtp"
	"net/http"
	"text/template"
	"encoding/json"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"

)

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		
	encoder := json.NewEncoder(w)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil { panic(err) }	

	signUp := SignUPBody {}
	json.Unmarshal(body, &signUp)

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			Dbhost, Dbuser, Dbpassword, Dbname, Dbport,
		),
	)

	defer db.Close()

	if err != nil { panic(err) }

	user := User {}

	err = db.QueryRow(
		SQL_INSERT_USER, 
		signUp.Email, 
		signUp.Password,
	).Scan(
		&user.Id, 
		&user.Email,
	)

	if err != nil { panic(err) }

	var uuid string

	err = db.QueryRow(
		SQL_SELECT_ACTIVATION,
		user.Id,
	).Scan(
		&uuid,
	)

	if err != nil { panic(err) }

	auth := smtp.PlainAuth(
		"", 
		"sardorgolangdev@gmail.com", 
		"sardor@11",
		"smtp.gmail.com", 
	)

	var buffer bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buffer.Write([]byte(fmt.Sprintf("Subject: Welcome to our site \n%s\n\n", mimeHeaders)))

	t, err := template.ParseFiles("mail-template.html")

	if err != nil { panic(err) }

	t.Execute(&buffer, struct {Email, UUID string}{
		Email: user.Email,
		UUID: uuid,
	})

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"sardorgolangdev@gmail.com",
		[]string {user.Email},
		buffer.Bytes(),
	)

	if err != nil { panic(err) }

	encoder.Encode(user)
}

func Verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			Dbhost, Dbuser, Dbpassword, Dbname, Dbport,
		),
	)

	defer db.Close()

	row, err := db.Exec(
		UPDATE_ACTIVATE,
		vars["uuid"],
	)

	if err != nil {
		panic(err)
	}

	affected, _ := row.RowsAffected()

	if affected > 0 {
		w.WriteHeader(http.StatusAccepted)
	
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

var jwtKey = []byte("secret-key")

type Payload struct {
		Email string
		jwt.StandardClaims
}

func CoursesCtrl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")		
	
	// encoder := json.NewEncoder(w)

	courseBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	accessToCourse := AccessToCourse {}

	json.Unmarshal(courseBody, &accessToCourse)

	// courses := []Course {}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			Dbhost, Dbuser, Dbpassword, Dbname, Dbport,
		),
	)

	defer db.Close()

	if err != nil {
		panic(err)
	}

	var activated_at string

	err = db.QueryRow(
		SELECT_COURSES_ACTIVATED_AT,
		accessToCourse.Email,
	).Scan(&activated_at)

	if activated_at != "" {
		// encoder.Encode(courses)
		expirationTime := time.Now().Add(43200*time.Minute)

		payload := Payload {
			Email: accessToCourse.Email,
			StandardClaims: jwt.StandardClaims {
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			panic(err)
		}
		w.Write([]byte(tokenString))

	} else {
		w.Write([]byte("Signup first or verify your account"))
	}
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("token")
	encoder := json.NewEncoder(w)

	payload := &Payload {}

	tkn, err := jwt.ParseWithClaims(token, payload, func (token *jwt.Token) (interface {}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	courses := []Course {}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			Dbhost, Dbuser, Dbpassword, Dbname, Dbport,
		),
	)

	defer db.Close()

	if err != nil {
		panic(err)
	}

	rows, err := db.Query(
		SELECT_COURSES,
	)

	for rows.Next() {
		var course Course
		rows.Scan(
			&course.Name,
			&course.Price,
		)

		courses = append(courses, course)
	}

	encoder.Encode(courses)

}
