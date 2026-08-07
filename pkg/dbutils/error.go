package dbutils

import (
	"errors"
	"strings"
)

var errorFilters = []func(err error) (bool, error) {
	filterDuplicationUsername,
	filterDuplicationEmail,
	filterRecordNotFound,
}
func CatchDBError(err error) error {
	if err == nil {
		return nil
	}

	for _, filter := range errorFilters {
		match, filteredErr := filter(err)
		if match {
			return filteredErr
		}
	}
	return err
}

var (
ErrDuplicationUsername = errors.New("username already exists") 
ErrDuplicationEmail = errors.New("email already exists")
ErrRecordNotFound = errors. New("record not found")
)
func filterDuplicationUsername(err error) (bool, error) { 
return strings.Contains(strings.ToLower(err.Error()), `"duplicate key value violates unique constraint "uni_users_username"`), nil
}
func filterDuplicationEmail(err error) (bool, error) { 
return strings.Contains(strings.ToLower(err.Error()), "duplicate key value"), nil
}
func filterRecordNotFound(err error) (bool, error) { 
return strings.Contains(strings.ToLower(err.Error()), "record not found"), nil
}

