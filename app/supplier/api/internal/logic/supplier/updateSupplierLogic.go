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

type UpdateSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSupplierLogic {
	return &UpdateSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSupplierLogic) UpdateSupplier(req *types.UpdateSupplierReq) (resp *types.UpdateSupplierResp, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.SupplierRPC.UpdateSupplier(l.ctx, &pb.UpdateSupplierReq{
		Id:           id,
		Name:         req.Name,
		Contact:      req.Contact,
		Phone:        req.Phone,
		Address:      req.Address,
		PaymentTerms: req.PaymentTerms,
		Credit:       req.Credit,
		IsActive:     req.IsActive,
		UpdatedBy:    employeeId,
	})
	if err != nil {
		return nil, err
	}

	return
}
