package logic

import (
	"context"
	"database/sql"
	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/model"
	"erp/app/supplier/rpc/internal/svc"
	"erp/app/supplier/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSupplierLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSupplierLogic {
	return &AddSupplierLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------supplier-----------------------
func (l *AddSupplierLogic) AddSupplier(in *pb.AddSupplierReq) (*pb.AddSupplierResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.SupplierModel.Insert(l.ctx, &model.Supplier{
		Id:           id,
		Code:         sql.NullString{String: in.Code, Valid: true},
		Uscc:         sql.NullString{String: in.Uscc, Valid: true},
		Name:         in.Name,
		Contact:      sql.NullString{String: in.Contact, Valid: true},
		Phone:        sql.NullString{String: in.Phone, Valid: true},
		Address:      sql.NullString{String: in.Address, Valid: true},
		PaymentTerms: sql.NullString{String: in.PaymentTerms, Valid: true},
		Credit:       sql.NullString{String: in.Credit, Valid: true},
		IsActive:     in.IsActive,
		CreatedBy:    in.CreatedBy,
		UpdatedBy:    in.CreatedBy,
	})
	if err != nil {
		return nil, code.AddSupplierFail
	}

	return &pb.AddSupplierResp{
		Id: id,
	}, nil
}
