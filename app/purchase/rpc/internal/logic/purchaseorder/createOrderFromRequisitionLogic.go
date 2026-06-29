package purchaseorderlogic

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"erp/app/purchase/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderFromRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderFromRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderFromRequisitionLogic {
	return &CreateOrderFromRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 从采购申请创建采购订单（通过 requisition_id 获取申请明细数据，保证一致性）
func (l *CreateOrderFromRequisitionLogic) CreateOrderFromRequisition(in *pb.CreateOrderFromRequisitionReq) (*pb.CreateOrderFromRequisitionResp, error) {
	// 获取申请明细
	id := util.GenerateSnowflake()

	reqDetails, err := l.svcCtx.PurchaseRequisitionDetailModel.ListByRequisitionId(l.ctx, in.RequisitionId)
	if err != nil {

		return nil, code.CreateOrderFail

	}

	param := &types.CreateOrderFromRequisitionParam{
		RequisitionId: in.RequisitionId,
		OrderNo:       in.OrderNo,
		SupplierId:    in.SupplierId,
		OrderDate:     in.OrderDate,
		ExpectedDate:  in.ExpectedDate,
		PurchaserId:   in.PurchaserId,
	}

	// 如果传入了明细，则使用传入的明细（支持覆盖）
	if len(in.Details) > 0 {
		for _, d := range in.Details {
			param.Details = append(param.Details, types.OrderDetailParam{
				Id:           util.GenerateSnowflake(),
				ProductId:    d.ProductId,
				ProductName:  d.ProductName,
				CategoryType: d.CategoryType,
				Quantity:     d.Quantity,
				UnitPrice:    d.UnitPrice,
				Amount:       d.Amount,
				Remark:       d.Remark,
			})
		}
	} else {
		// 使用申请的明细数据
		for _, d := range reqDetails {
			param.Details = append(param.Details, types.OrderDetailParam{
				Id:           util.GenerateSnowflake(),
				ProductId:    d.ProductId.Int64,
				ProductName:  d.ProductName.String,
				CategoryType: d.CategoryType,
				Quantity:     d.Quantity,
				UnitPrice:    d.UnitPrice.Float64,
				Amount:       d.Amount.Float64,
				Remark:       d.Remark.String,
			})
		}
	}

	err = l.svcCtx.PurchaseOrderModel.CreateFromRequisition(l.ctx, id, param)
	if err != nil {

		return nil, code.CreateOrderFail

	}
	return &pb.CreateOrderFromRequisitionResp{OrderId: id}, nil
}
