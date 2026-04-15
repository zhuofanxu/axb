package database

import (
	"context"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/zhuofanxu/axb/errx"
)

func WrapDBErr(err error, data ...string) error {
	var mysqlErr *mysql.MySQLError
	var pgErr *pgconn.PgError
	var wrappedErr error
	msg := ""
	if len(data) > 0 {
		msg = data[0]
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		wrappedErr = errx.NotFoundError(err)
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		wrappedErr = errx.AlreadyExistsError(err)
	} else if isSQLiteDuplicateErr(err) {
		wrappedErr = errx.AlreadyExistsError(err)
	} else if errors.Is(err, context.Canceled) {
		wrappedErr = errx.CanceledError(err)
	} else if errors.Is(err, context.DeadlineExceeded) {
		wrappedErr = errx.TimeoutError(err)
	} else if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // PostgreSQL 唯一键冲突
			wrappedErr = errx.AlreadyExistsError(err)
		default:
			return errors.Wrap(err, msg+"[数据库错误]")
		}
	} else if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062: // 重复键
			wrappedErr = errx.AlreadyExistsError(err)
		case 1054: // 未知列
			return errors.Wrap(err, msg+"[数据库未知列错误]")
		default:
			return errors.Wrap(err, msg+"[数据库错误]")
		}
	} else {
		return errors.Wrap(err, msg+"[数据库错误]")
	}

	return wrappedErr
}

// isSQLiteDuplicateErr 判断错误是否为 SQLite 唯一键冲突。
func isSQLiteDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "unique constraint failed")
}
