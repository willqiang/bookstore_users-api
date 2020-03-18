package users

import (
	"fmt"
	"github.com/willqiang/bookstore_users-api/datasources/mysql/users_db"
	"github.com/willqiang/bookstore_users-api/utils/date_utils"
	"github.com/willqiang/bookstore_users-api/utils/errors"
	"github.com/willqiang/bookstore_users-api/utils/mysql_utils"
)

const (
	errorNoRows     = "no rows in result set"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?)"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?;"
	queryUpdateUser = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
	queryDeleteUser = "DELETE FROM users WHERE id = ?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when try get user id %d", user.Id))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	err = result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close() // Must close. Very Important.

	user.DateCreated = date_utils.GetNowString()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}
