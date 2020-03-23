package users

import (
	"fmt"
	"github.com/willqiang/bookstore_users-api/datasources/mysql/users_db"
	"github.com/willqiang/bookstore_users-api/logger"
	"github.com/willqiang/bookstore_users-api/utils/date_utils"
	"github.com/willqiang/bookstore_users-api/utils/errors"
)

const (
	errorNoRows            = "no rows in result set"
	queryInsertUser        = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?)"
	queryGetUser           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryUpdateUser        = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
	queryDeleteUser        = "DELETE FROM users WHERE id = ?;"
	queryFindUsersByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	err = result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
	if err != nil {
		logger.Error("error when trying to get user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare create user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close() // Must close. Very Important.

	user.DateCreated = date_utils.GetNowDbFormat()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to create user", saveErr)
		return errors.NewInternalServerError("database error")
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get user id after created user", saveErr)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUsersByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan row to user", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s:", status))
	}

	return results, nil
}
