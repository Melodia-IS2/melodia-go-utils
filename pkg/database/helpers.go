package database

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func HandleSqlError(err error) *SqlError {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return NewSqlError(SqlErrorTypeNotFound, nil, err.Error())
	}
	// Handle MySQL-specific errors
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return NewSqlError(SqlErrorTypeUnknown, nil, err.Error())
	}

	switch mysqlErr.Number {
	case 1062: // Duplicate entry for unique constraint
		var duplicateEntryRegex = regexp.MustCompile(`Duplicate entry '([^']*)' for key '([^']*)'`)

		matches := duplicateEntryRegex.FindStringSubmatch(mysqlErr.Message)
		var attributeName *string
		value := ""
		if len(matches) == 3 {
			fullKey := matches[2]
			parts := strings.Split(fullKey, ".")
			attr := ""
			if len(parts) == 2 {
				attr = parts[1]
			} else {
				attr = fullKey
			}
			if idx := strings.Index(attr, "_"); idx != -1 {
				attr = attr[:idx]
			}
			attributeName = &attr
			value = matches[1]
		}
		return NewSqlError(SqlErrorTypeConflict, attributeName, value)
	case 1451: // Cannot delete or update a parent row: a foreign key constraint fails
		var fkConstraintRegex = regexp.MustCompile("FOREIGN KEY \\(`([^`]*)`\\)")
		matches := fkConstraintRegex.FindStringSubmatch(mysqlErr.Message)
		var constraintName *string
		if len(matches) >= 2 {
			cn := matches[1]
			if idx := strings.Index(cn, "_"); idx != -1 {
				cn = cn[:idx]
			}
			constraintName = &cn
		}
		return NewSqlError(SqlErrorTypeFkError, constraintName, mysqlErr.Message)
	case 1452: // Cannot delete or update a child or parent row: a foreign key constraint fails
		var fkConstraintRegex = regexp.MustCompile("FOREIGN KEY \\(`([^`]*)`\\)")
		matches := fkConstraintRegex.FindStringSubmatch(mysqlErr.Message)
		var constraintName *string
		if len(matches) >= 2 {
			cn := matches[1]
			if idx := strings.Index(cn, "_"); idx != -1 {
				cn = cn[:idx]
			}
			constraintName = &cn
		}
		return NewSqlError(SqlErrorTypeFkError, constraintName, mysqlErr.Message)
	}

	return NewSqlError(SqlErrorTypeUnknown, nil, mysqlErr.Message)
}
