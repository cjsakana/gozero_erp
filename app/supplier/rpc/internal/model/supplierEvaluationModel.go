package model

import (
	"context"
	"erp/app/supplier/rpc/internal/types"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"
)

var _ SupplierEvaluationModel = (*customSupplierEvaluationModel)(nil)

type (
	// SupplierEvaluationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSupplierEvaluationModel.
	SupplierEvaluationModel interface {
		supplierEvaluationModel
		Search(ctx context.Context, req *types.SearchSupplierEvaluation) ([]*SupplierEvaluation, int64, error)
	}

	customSupplierEvaluationModel struct {
		*defaultSupplierEvaluationModel
	}
)

// NewSupplierEvaluationModel returns a model for the database table.
func NewSupplierEvaluationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SupplierEvaluationModel {
	return &customSupplierEvaluationModel{
		defaultSupplierEvaluationModel: newSupplierEvaluationModel(conn, c, opts...),
	}
}

// 构建评分搜索条件
func (m *customSupplierEvaluationModel) buildScoreConditions(req *types.SearchSupplierEvaluation) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// 质量评分搜索
	if req.QualityOp != "" || req.QualityMin > 0 || req.QualityMax > 0 {
		condition := m.buildScoreCondition("quality_score", req.QualityMin, req.QualityMax, req.QualityOp)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	// 交付评分搜索
	if req.DeliveryOp != "" || req.DeliveryMin > 0 || req.DeliveryMax > 0 {
		condition := m.buildScoreCondition("delivery_score", req.DeliveryMin, req.DeliveryMax, req.DeliveryOp)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	// 服务评分搜索
	if req.ServiceOp != "" || req.ServiceMin > 0 || req.ServiceMax > 0 {
		condition := m.buildScoreCondition("service_score", req.ServiceMin, req.ServiceMax, req.ServiceOp)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	// 综合评分搜索
	if req.OverallOp != "" || req.OverallMin > 0 || req.OverallMax > 0 {
		condition := m.buildScoreCondition("overall_score", req.OverallMin, req.OverallMax, req.OverallOp)
		if condition != "" {
			conditions = append(conditions, condition)
		}
	}

	if len(conditions) > 0 {
		return "(" + strings.Join(conditions, " AND ") + ")", args
	}

	return "", args
}

// 构建单个评分字段的搜索条件
func (m *customSupplierEvaluationModel) buildScoreCondition(field string, min, max float64, op string) string {
	switch op {
	case "gt": // 大于
		return fmt.Sprintf("%s > ?", field)
	case "gte": // 大于等于
		return fmt.Sprintf("%s >= ?", field)
	case "eq": // 等于
		return fmt.Sprintf("%s = ?", field)
	case "lt": // 小于
		return fmt.Sprintf("%s < ?", field)
	case "lte": // 小于等于
		return fmt.Sprintf("%s <= ?", field)
	case "between": // 范围
		if min > 0 && max > 0 {
			return fmt.Sprintf("%s BETWEEN ? AND ?", field)
		}
	case "": // 默认使用范围搜索
		if min > 0 && max > 0 {
			return fmt.Sprintf("%s BETWEEN ? AND ?", field)
		} else if min > 0 {
			return fmt.Sprintf("%s >= ?", field)
		} else if max > 0 {
			return fmt.Sprintf("%s <= ?", field)
		}
	}

	return ""
}

// 完整的搜索方法
func (m *customSupplierEvaluationModel) Search(ctx context.Context, req *types.SearchSupplierEvaluation) ([]*SupplierEvaluation, int64, error) {
	var conditions []string
	var args []interface{}

	// 基础条件
	if req.SupplierId > 0 {
		conditions = append(conditions, "supplier_id = ?")
		args = append(args, req.SupplierId)
	}

	if req.EvaluatorId > 0 {
		conditions = append(conditions, "evaluator_id = ?")
		args = append(args, req.EvaluatorId)
	}

	if req.StartData > 0 {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, time.Unix(req.StartData, 0))
	}

	if req.EndData > 0 {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, time.Unix(req.EndData, 0))
	}

	// 评分搜索条件
	scoreCondition, scoreArgs := m.buildScoreConditions(req)
	if scoreCondition != "" {
		conditions = append(conditions, scoreCondition)
		args = append(args, scoreArgs...)
	}

	// 构建WHERE子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)

	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	var evaluations []*SupplierEvaluation
	dataQuery := fmt.Sprintf(`select %s from %s %s LIMIT ? OFFSET ?`, supplierEvaluationRows, m.table, whereClause)

	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)

	err = m.QueryRowsNoCacheCtx(ctx, &evaluations, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	switch {
	case err == nil:
		return evaluations, total, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}
