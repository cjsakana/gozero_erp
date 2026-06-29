package supplier

import (
	"context"
	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"
	"erp/app/supplier/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSupplierLogic {
	return &CreateSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSupplierLogic) CreateSupplier(req *types.CreateSupplierReq) (resp *types.CreateSupplierResp, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	code := util.GenerateNo("SUP")
	ret, err := l.svcCtx.SupplierRPC.AddSupplier(l.ctx, &pb.AddSupplierReq{
		Code:         code,
		Uscc:         req.Uscc,
		Name:         req.Name,
		Contact:      req.Contact,
		Phone:        req.Phone,
		Address:      req.Address,
		PaymentTerms: req.PaymentTerms,
		Credit:       req.Credit,
		IsActive:     req.IsActive,
		CreatedBy:    employeeId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.CreateSupplierResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
