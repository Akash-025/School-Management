package utils

import "errors"

type ContextKey string

func AuthorizeUser(userRole string, allowedRoles ...string) (bool, error) {

	for _, allowedrole := range allowedRoles {
		if allowedrole == userRole {
			return true, nil
		}
	}
	return false, errors.New("User not authorized")
}