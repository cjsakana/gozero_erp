package logic

import (
	"context"

	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/svc"
	"erp/app/supplier/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSupplierByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSupplierByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupplierByIdLogic {
	return &GetSupplierByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc DelSupplier(DelSupplierReq) returns (DelSupplierResp);
func (l *GetSupplierByIdLogic) GetSupplierById(in *pb.GetSupplierByIdReq) (*pb.GetSupplierByIdResp, error) {
	one, err := l.svcCtx.SupplierModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, code.SupplierNotFound
	}

	return &pb.GetSupplierByIdResp{
		Supplier: &pb.Supplier{
			Id:           one.Id,
			Code:         one.Code.String,
			Uscc:         one.Uscc.String,
			Name:         one.Name,
			Contact:      one.Contact.String,
			Phone:        one.Phone.String,
			Address:      one.Address.String,
			PaymentTerms: one.PaymentTerms.String,
			Credit:       one.Credit.String,
			IsActive:     one.IsActive,
			CreatedAt:    one.CreatedAt.Unix(),
			CreatedBy:    one.CreatedBy,
			UpdatedAt:    one.UpdatedAt.Unix(),
			UpdatedBy:    one.UpdatedBy,
		},
	}, nil
}
