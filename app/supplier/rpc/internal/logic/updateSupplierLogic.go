package logic

import (
	"context"
	"database/sql"
	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/model"

	"erp/app/supplier/rpc/internal/svc"
	"erp/app/supplier/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSupplierLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSupplierLogic {
	return &UpdateSupplierLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSupplierLogic) UpdateSupplier(in *pb.UpdateSupplierReq) (*pb.UpdateSupplierResp, error) {
	err := l.svcCtx.SupplierModel.XUpdate(l.ctx, &model.Supplier{
		Id:           in.Id,
		Name:         in.Name,
		Contact:      sql.NullString{String: in.Contact, Valid: true},
		Phone:        sql.NullString{String: in.Phone, Valid: true},
		Address:      sql.NullString{String: in.Address, Valid: true},
		PaymentTerms: sql.NullString{String: in.PaymentTerms, Valid: true},
		Credit:       sql.NullString{String: in.Credit, Valid: true},
		IsActive:     in.IsActive,
		UpdatedBy:    in.UpdatedBy,
	})
	if err != nil {
		return nil, code.UpdateSupplierFail
	}

	return &pb.UpdateSupplierResp{}, nil
}
