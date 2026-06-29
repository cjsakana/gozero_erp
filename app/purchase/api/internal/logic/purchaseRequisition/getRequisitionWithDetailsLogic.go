package purchaseRequisition

import (
	"context"
	"erp/app/hr/rpc/client/department"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/product/rpc/client/product"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRequisitionWithDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRequisitionWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequisitionWithDetailsLogic {
	return &GetRequisitionWithDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRequisitionWithDetailsLogic) GetRequisitionWithDetails(req *types.GetRequisitionWithDetailsReq) (resp *types.GetRequisitionWithDetailsResp, err error) {
	requisitionId, err := util.StringToInt64(req.RequisitionId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.PurchaseRPC.GetRequisitionWithDetails(l.ctx, &pb.GetRequisitionWithDetailsReq{
		RequisitionId: requisitionId,
	})
	if err != nil {
		return nil, err
	}

	// 组装响应
	requisition := types.PurchaseRequisition{
		Id:             util.Int64ToString(ret.Requisition.Id),
		RequisitionNo:  ret.Requisition.RequisitionNo,
		DepartmentId:   util.Int64ToString(ret.Requisition.DepartmentId),
		DepartmentName: "",
		ApplicantId:    "",
		ApplicantNo:    "",
		ApplicantName:  "",
		RequestDate:    ret.Requisition.RequestDate,
		TotalAmount:    ret.Requisition.TotalAmount,
		Status:         ret.Requisition.Status,
		ApproverId:     "",
		ApproverNo:     "",
		ApproverName:   "",
		ApproveTime:    ret.Requisition.ApproveTime,
		ApproveRemark:  ret.Requisition.ApproveRemark,
		CreatedAt:      ret.Requisition.CreatedAt,
		UpdatedAt:      ret.Requisition.UpdatedAt,
	}

	// 填充申请人信息
	empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
		Id: ret.Requisition.ApplicantId,
	})
	if err == nil && empResp.EmployeeNonSensitiveDetail != nil {
		requisition.ApplicantId = util.Int64ToString(ret.Requisition.ApplicantId)
		requisition.ApplicantNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
		requisition.ApplicantName = empResp.EmployeeNonSensitiveDetail.Name
	}

	// 填充审批人信息
	empResp2, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
		Id: ret.Requisition.ApproverId,
	})
	if err == nil && empResp2.EmployeeNonSensitiveDetail != nil {
		requisition.ApproverId = util.Int64ToString(ret.Requisition.ApproverId)
		requisition.ApproverNo = empResp2.EmployeeNonSensitiveDetail.EmployeeNo
		requisition.ApproverName = empResp2.EmployeeNonSensitiveDetail.Name
	}

	departmentByIdResp, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &department.GetDepartmentByIdReq{
		Id: ret.Requisition.DepartmentId,
	})
	if err == nil && departmentByIdResp.Department != nil {
		requisition.DepartmentName = departmentByIdResp.Department.Name
	}

	resp = &types.GetRequisitionWithDetailsResp{
		Requisition: requisition,
		Details: func() []*types.PurchaseRequisitionDetail {
			list := make([]*types.PurchaseRequisitionDetail, 0, len(ret.Details))
			productMap := make(map[int64]*product.Product)

			for _, d := range ret.Details {

				detail := &types.PurchaseRequisitionDetail{
					Id:            util.Int64ToString(d.Id),
					RequisitionId: util.Int64ToString(d.RequisitionId),

					ProductName:  d.ProductName,
					CategoryType: d.CategoryType,
					Quantity:     d.Quantity,
					UnitPrice:    d.UnitPrice,
					Amount:       d.Amount,
					Remark:       d.Remark,
				}
				if d.ProductId > 0 {
					if _, ok := productMap[d.ProductId]; !ok {
						prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
							Id: d.ProductId,
						})
						if err == nil {
							productMap[d.ProductId] = prod.Product
						}
					}
					detail.ProductId = util.Int64ToString(d.ProductId)
					detail.ProductNo = productMap[d.ProductId].ProductNo
				}
				list = append(list, detail)
			}
			return list
		}(),
	}
	return
}
