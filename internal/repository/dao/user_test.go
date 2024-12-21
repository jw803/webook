package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	tests := []struct {
		name string
		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user Users

		wantErr     error
		wantErrCode int
		wantId      int64
	}{
		{
			name: "insert successfully",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: Users{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
		},
		{
			name: "duplicate email",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)
				return mockDB
			},
			user:    Users{},
			wantErr: errorx.WithCode(errcode.ErrUserDuplicated, "email has already been registered"),
		},
		{
			name: "db error",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是，只要是 INSERT 到 users 的语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("db error"))
				require.NoError(t, err)
				return mockDB
			},
			user:    Users{},
			wantErr: errorx.WithCode(errcode.ErrDatabase, "db error"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn: tc.mock(t),
				// 停用gorm初始化 發起 SELECT VERSION的調用
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// no need to ping mock db
				DisableAutomaticPing: true,
				// if false, even if a simple CU clause, gorm also create a transaction
				SkipDefaultTransaction: true,
			})
			d := NewGORMUserDAO(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.True(t, true, errorx.IsEqual(tc.wantErr, err))
		})
	}
}
