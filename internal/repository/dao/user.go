package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jw803/webook/internal/pkg/errcode"
	"github.com/jw803/webook/pkg/errorx"
	"github.com/jw803/webook/pkg/loggerx"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (Users, error)
	FindById(ctx context.Context, id int64) (Users, error)
	FindByPhone(ctx context.Context, phone string) (Users, error)
	Insert(ctx context.Context, u Users) error
	FindByWechat(ctx context.Context, openID string) (Users, error)
	UpdateExtraInfoById(ctx context.Context, id int64, u Users) error
}

type GORMUserDAO struct {
	db *gorm.DB
	l  loggerx.Logger
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	res := &GORMUserDAO{
		db: db,
	}
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	db, err := gorm.Open(mysql.Open())
	//	pt := unsafe.Pointer(&res.db)
	//	atomic.StorePointer(&pt, unsafe.Pointer(&db))
	//})
	return res
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (Users, error) {
	var u Users
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(&u).Error
	//err := dao.p().WithContext(ctx).Where("wechat_open_id = ?", openID).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (Users, error) {
	var u Users
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dao.l.Error(ctx, "user not found", loggerx.Error(err))
		return u, errorx.WithCode(errcode.ErrUserNotFound, err.Error())
	}
	if err != nil {
		dao.l.Error(ctx, "db error", loggerx.Error(err))
		return u, errorx.WithCode(errcode.ErrDatabase, err.Error())
	}
	return u, nil
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (Users, error) {
	var u Users
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dao.l.Error(ctx, "user not found", loggerx.Error(err))
		return u, errorx.WithCode(errcode.ErrUserNotFound, err.Error())
	}
	if err != nil {
		dao.l.Error(ctx, "db error", loggerx.Error(err))
		return u, errorx.WithCode(errcode.ErrDatabase, err.Error())
	}
	return u, nil
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (Users, error) {
	var u Users
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u Users) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return errorx.WithCode(errcode.ErrUserDuplicated, "email has already been registered")
		}
	}
	return err
}

func (dao *GORMUserDAO) UpdateExtraInfoById(ctx context.Context, id int64, u Users) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	res := dao.db.WithContext(ctx).Where("id = ?", id).
		Updates(map[string]any{
			"nick_name": u.NickName,
			"birthday":  u.Birthday,
			"intro":     u.Intro,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errorx.WithCode(errcode.ErrUserNotFound, "user % not found", id)
	}
	return nil
}

// Users 直接对应数据库表结构
// 有些人叫做 entity，有些人叫做 model，有些人叫做 PO(persistent object)

type Users struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全部用户唯一
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引允许有多个空值
	// 但是不能有多个 ""
	Phone sql.NullString `gorm:"unique"`
	// 最大问题就是，你要解引用
	// 你要判空
	//Phone *string

	// 往这面加
	NickName string
	Birthday string
	Intro    string

	// 索引的最左匹配原则：
	// 假如索引在 <A, B, C> 建好了
	// A, AB, ABC 都能用
	// WHERE A =?
	// WHERE A = ? AND B =?    WHERE B = ? AND A =?
	// WHERE A = ? AND B = ? AND C = ?  ABC 的顺序随便换
	// WHERE 里面带了 ABC，可以用
	// WHERE 里面，没有 A，就不能用

	// 如果要创建联合索引，<unionid, openid>，用 openid 查询的时候不会走索引
	// <openid, unionid> 用 unionid 查询的时候，不会走索引
	// 微信的字段
	WechatUnionID sql.NullString
	WechatOpenID  sql.NullString `gorm:"unique"`

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
