package main

import (
	"database/sql"
)

type User struct {
	Uuid string 		`json:"id"`
	Username string 	`json:"username"`
	Name string			`json:"name"`
}

func ListUsers (db *sql.DB) ([]User, error) {
	rows, err:= db.Query("select uuid,username,name from users")
	if err!=nil {
		return nil, err
	}
	defer rows.Close()
	
	users:=[]User{}
	for rows.Next() {
		user:=User{}
		if err:=rows.Scan(&user.Uuid,&user.Username,&user.Name); err!=nil {
			return users, err
		}
		users=append(users,user)
	}
	if err=rows.Err(); err!=nil {
		return users, err
	}
	return users, nil
}