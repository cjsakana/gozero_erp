package model

import (
	"context"
	"database/sql"
	"erp/app/user/rpc/internal/types"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		XUpdate(ctx context.Context, newData *User) error
		SearchUsers(ctx context.Context, data *types.SearchUserParams) ([]*User, int64, error)
		XDelete(ctx context.Context, id int64) error
		BulkInsert(ctx context.Context, data []*User) (err error)
	}

	customUserModel struct {
		*defaultUserModel
		conn sqlx.SqlConn
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
		conn:             conn,
	}
}

func (m *customUserModel) XUpdate(ctx context.Context, newData *User) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	var setClauses []string
	var args []interface{}

	if newData.Username != "" {
		setClauses = append(setClauses, "username = ?")
		args = append(args, newData.Username)
	}
	if newData.RealName != "" {
		setClauses = append(setClauses, "real_name = ?")
		args = append(args, newData.RealName)
	}
	if newData.PasswordHash != "" {
		setClauses = append(setClauses, "password_hash = ?")
		args = append(args, newData.PasswordHash)
	}
	if newData.Phone.Valid && newData.Phone.String != "" {
		setClauses = append(setClauses, "phone = ?")
		args = append(args, newData.Phone.String)
	}
	if newData.Email.String != "" {
		setClauses = append(setClauses, "email = ?")
		args = append(args, newData.Email)
	}
	if newData.Avatar.Valid && newData.Avatar.String != "" {
		setClauses = append(setClauses, "avatar = ?")
		args = append(args, newData.Avatar.String)
	}
	if newData.UpdatedBy.Valid {
		setClauses = append(setClauses, "updated_by = ?")
		args = append(args, newData.UpdatedBy.Int64)
	}
	if newData.Resigned != 0 {
		setClauses = append(setClauses, "resigned = ?")
		args = append(args, newData.Resigned)
	}

	// 如果没有字段需要更新，直接返回成功
	if len(setClauses) == 0 {
		return nil
	}

	erpUserUserIdKey := fmt.Sprintf("%s%v", cacheErpUserUserIdPrefix, data.Id)
	erpUserUserPhoneKey := fmt.Sprintf("%s%v", cacheErpUserUserPhonePrefix, data.Phone)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setClauses, ", "))
		args = append(args, newData.Id)
		return conn.ExecCtx(ctx, query, args...)
	}, erpUserUserIdKey, erpUserUserPhoneKey)
	return err
}

func (m *customUserModel) SearchUsers(ctx context.Context, data *types.SearchUserParams) ([]*User, int64, error) {
	var users []*User

	conditions := []string{}
	args := []any{}

	if data.EmployeeId != 0 {
		conditions = append(conditions, "employee_id = ?")
		args = append(args, data.EmployeeId)
	}

	if data.Username != "" {
		conditions = append(conditions, "username LIKE ?")
		args = append(args, "%"+data.Username+"%")
	}

	if data.RealName != "" {
		conditions = append(conditions, "real_name LIKE ?")
		args = append(args, "%"+data.RealName+"%")
	}

	if data.Phone != "" {
		conditions = append(conditions, "phone LIKE ?")
		args = append(args, "%"+data.Phone+"%")
	}

	if data.Email != "" {
		conditions = append(conditions, "email LIKE ?")
		args = append(args, "%"+data.Email+"%")
	}

	// 是否离职筛选（可选）
	if data.Resigned != nil {
		if *data.Resigned {
			conditions = append(conditions, "resigned = ?")
			args = append(args, 1)
		} else {
			conditions = append(conditions, "resigned = ?")
			args = append(args, 0)
		}
	}

	// 构建完整 SQL
	sql := fmt.Sprintf("select %s from %s", userRows, m.table)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM  %s", m.table)
	if len(conditions) > 0 {
		where := " where " + strings.Join(conditions, " AND ")
		sql += where
		countQuery += where
	}

	// 查询总数
	var total int64
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 添加分页
	if data.Limit != -1 { // 约定 -1 表示查询全部
		sql += fmt.Sprintf(" LIMIT %d OFFSET %d", data.Limit, (data.Page-1)*data.Limit)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &users, sql, args...)

	switch {
	case err == nil:
		return users, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customUserModel) XDelete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	erpUserUserIdKey := fmt.Sprintf("%s%v", cacheErpUserUserIdPrefix, id)
	erpUserUserPhoneKey := fmt.Sprintf("%s%v", cacheErpUserUserPhonePrefix, data.Phone)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `resigned` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, 1, id)
	}, erpUserUserIdKey, erpUserUserPhoneKey)
	return err
}

func (m *customUserModel) BulkInsert(ctx context.Context, data []*User) (err error) {
	if len(data) == 0 {
		return nil
	}

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)

	blk, err := sqlx.NewBulkInserter(m.conn, query)
	if err != nil {
		return err
	}
	defer blk.Flush()

	// 统计交给EmployeeDetail
	blk.SetResultHandler(func(result sql.Result, err error) {
		if err != nil {
			logx.Error(err)
			return
		}
	})

	for _, v := range data {
		err := blk.Insert(
			v.Id,
			v.EmployeeId,
			v.EmployeeNo,
			v.Username,
			v.RealName,
			v.PasswordHash,
			v.Phone.String,
			v.Email.String,
			v.Avatar.String,
			v.Resigned,
			v.CreatedBy.Int64,
			v.UpdatedBy.Int64,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
