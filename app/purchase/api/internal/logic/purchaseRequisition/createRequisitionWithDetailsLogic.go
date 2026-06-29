package purchaseRequisition

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRequisitionWithDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRequisitionWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRequisitionWithDetailsLogic {
	return &CreateRequisitionWithDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRequisitionWithDetailsLogic) CreateRequisitionWithDetails(req *types.CreateRequisitionWithDetailsReq) (resp *types.CreateRequisitionWithDetailsResp, err error) {
	details := make([]*pb.RequisitionDetailInput, 0, len(req.Details))
	for _, d := range req.Details {
		productId, err := util.StringToInt64(d.ProductId)
		if err != nil {
			return nil, err
		}
		fmt.Println("1fas351f3a1d3as", d)
		details = append(details, &pb.RequisitionDetailInput{
			ProductId:    productId,
			ProductName:  d.ProductName,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Remark:       d.Remark,
		})
	}

	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	applicantId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}

	no := util.GenerateNo("PR")
	ret, err := l.svcCtx.PurchaseRPC.CreateRequisitionWithDetails(l.ctx, &pb.CreateRequisitionWithDetailsReq{
		RequisitionNo: no,
		DepartmentId:  departmentId,
		ApplicantId:   applicantId,
		ApproverId:    approverId,
		RequestDate:   time.Now().Unix(),
		TotalAmount:   req.TotalAmount,
		Status:        req.Status,
		Details:       details,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateRequisitionWithDetailsResp{
		RequisitionId: util.Int64ToString(ret.RequisitionId),
	}
	return
}
