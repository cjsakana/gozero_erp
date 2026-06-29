package purchaseOrder

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderFromRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderFromRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderFromRequisitionLogic {
	return &CreateOrderFromRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderFromRequisitionLogic) CreateOrderFromRequisition(req *types.CreateOrderFromRequisitionReq) (resp *types.CreateOrderFromRequisitionResp, err error) {
	// 从采购申请创建订单：如果未提供明细，系统会自动从申请中获取明细数据
	// 如果提供了明细，则使用提供的明细（用于调整价格、数量等）
	var details []*pb.OrderDetailInput
	if req.Details != nil && len(req.Details) > 0 {
		details = make([]*pb.OrderDetailInput, 0, len(req.Details))
		for _, d := range req.Details {
			productId, err := util.StringToInt64(d.ProductId)
			if err != nil {
				return nil, err
			}
			details = append(details, &pb.OrderDetailInput{
				ProductId:    productId,
				ProductName:  d.ProductName,
				CategoryType: d.CategoryType,
				Quantity:     d.Quantity,
				UnitPrice:    d.UnitPrice,
				Amount:       d.Amount,
				Remark:       d.Remark,
			})
		}
	}
	requisitionId, err := util.StringToInt64(req.RequisitionId)
	if err != nil {
		return nil, err
	}
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	purchaserId, err := util.StringToInt64(req.PurchaserId)
	if err != nil {
		return nil, err
	}
	// 如果 details 为 nil 或空，RPC 层会从 requisitionId 对应的申请中自动获取明细
	no := util.GenerateNo("PO")
	ret, err := l.svcCtx.PurchaseRPC.CreateOrderFromRequisition(l.ctx, &pb.CreateOrderFromRequisitionReq{
		RequisitionId: requisitionId,
		OrderNo:       no,
		SupplierId:    supplierId,
		OrderDate:     req.OrderDate,
		ExpectedDate:  req.ExpectedDate,
		PurchaserId:   purchaserId,
		Details:       details, // nil 或空数组时，RPC 会从申请中获取明细
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateOrderFromRequisitionResp{
		OrderId: util.Int64ToString(ret.OrderId),
	}
	return
}
