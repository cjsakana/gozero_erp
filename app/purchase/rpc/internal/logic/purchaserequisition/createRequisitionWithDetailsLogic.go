package purchaserequisitionlogic

import (
	"context"
	"erp/common/util"
	"fmt"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"erp/app/purchase/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRequisitionWithDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRequisitionWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRequisitionWithDetailsLogic {
	return &CreateRequisitionWithDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建采购申请及明细（事务）
func (l *CreateRequisitionWithDetailsLogic) CreateRequisitionWithDetails(in *pb.CreateRequisitionWithDetailsReq) (*pb.CreateRequisitionWithDetailsResp, error) {
	// 生成主表雪花ID
	requisitionId := util.GenerateSnowflake()

	param := &types.CreateRequisitionWithDetailsParam{
		RequisitionNo: in.RequisitionNo,
		DepartmentId:  in.DepartmentId,
		ApplicantId:   in.ApplicantId,
		ApproverId:    in.ApproverId,
		RequestDate:   in.RequestDate,
		TotalAmount:   in.TotalAmount,
		Status:        in.Status,
	}
	// 为每个明细生成雪花ID
	for _, d := range in.Details {
		fmt.Println("1111", d)
		amount := d.Quantity * d.UnitPrice
		param.Details = append(param.Details, types.RequisitionDetailParam{
			Id:           util.GenerateSnowflake(),
			ProductId:    d.ProductId,
			ProductName:  d.ProductName,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       amount,
			Remark:       d.Remark,
		})
	}

	err := l.svcCtx.PurchaseRequisitionModel.CreateWithDetails(l.ctx, requisitionId, param)
	if err != nil {

		return nil, code.CreateRequisitionFail

	}
	return &pb.CreateRequisitionWithDetailsResp{RequisitionId: requisitionId}, nil
}
