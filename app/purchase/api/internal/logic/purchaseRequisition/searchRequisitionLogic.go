package purchaseRequisition

import (
	"context"
	pb2 "erp/app/hr/rpc/pb"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchRequisitionLogic {
	return &SearchRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchRequisitionLogic) SearchRequisition(req *types.SearchRequisitionReq) (resp *types.SearchRequisitionResp, err error) {
	var departmentId, applicantId, approverId int64
	if req.DepartmentId != "" {
		departmentId, err = util.StringToInt64(req.DepartmentId)
		if err != nil {
			return nil, err
		}
	}
	if req.ApplicantId != "" {
		applicantId, err = util.StringToInt64(req.ApplicantId)
		if err != nil {
			return nil, err
		}
	}
	if req.ApproverId != "" {
		approverId, err = util.StringToInt64(req.ApproverId)
		if err != nil {
			return nil, err
		}
	}

	ret, err := l.svcCtx.PurchaseRPC.SearchRequisition(l.ctx, &pb.SearchRequisitionReq{
		Page:          req.Page,
		Limit:         req.Limit,
		RequisitionNo: req.RequisitionNo,
		DepartmentId:  departmentId,
		ApplicantId:   applicantId,
		ApproverId:    approverId,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.PurchaseRequisitionWithDetails, 0, int(ret.Total))
	for _, rd := range ret.RequisitionsWithDetails {
		//details := make([]*types.PurchaseRequisitionDetail, 0, len(rd.Details))
		//for _, d := range rd.Details {
		//	details = append(details, &types.PurchaseRequisitionDetail{
		//		Id:            util.Int64ToString(d.Id),
		//		RequisitionId: util.Int64ToString(d.RequisitionId),
		//		ProductId:     util.Int64ToString(d.ProductId),
		//		ProductNo:     "",
		//		ProductName:   d.ProductName,
		//		CategoryType:  d.CategoryType,
		//		Quantity:      d.Quantity,
		//		UnitPrice:     d.UnitPrice,
		//		Amount:        d.Amount,
		//		Remark:        d.Remark,
		//	})
		//}

		// 获取申请人信息
		var applicantNo, applicantName string
		if rd.Requisitions.ApplicantId > 0 {
			empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
				Id: rd.Requisitions.ApplicantId,
			})
			if err == nil && empResp.EmployeeNonSensitiveDetail != nil {
				applicantNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
				applicantName = empResp.EmployeeNonSensitiveDetail.Name
			}
		}
		// 获取审批人信息
		var approverNo, approverName string
		if rd.Requisitions.ApproverId > 0 {
			empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
				Id: rd.Requisitions.ApproverId,
			})
			if err == nil && empResp.EmployeeNonSensitiveDetail != nil {
				approverNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
				approverName = empResp.EmployeeNonSensitiveDetail.Name
			}
		}

		departmentByIdResp, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &pb2.GetDepartmentByIdReq{
			Id: rd.Requisitions.DepartmentId,
		})
		if err != nil {
			return nil, err
		}

		tRd := &types.PurchaseRequisitionWithDetails{
			Requisition: types.PurchaseRequisition{
				Id:             util.Int64ToString(rd.Requisitions.Id),
				RequisitionNo:  rd.Requisitions.RequisitionNo,
				DepartmentId:   util.Int64ToString(rd.Requisitions.DepartmentId),
				DepartmentName: departmentByIdResp.Department.Name,
				ApplicantId:    util.Int64ToString(rd.Requisitions.ApplicantId),
				ApplicantNo:    applicantNo,
				ApplicantName:  applicantName,
				RequestDate:    rd.Requisitions.RequestDate,
				TotalAmount:    rd.Requisitions.TotalAmount,
				Status:         rd.Requisitions.Status,
				ApproverId:     util.Int64ToString(rd.Requisitions.ApproverId),
				ApproverNo:     approverNo,
				ApproverName:   approverName,
				ApproveTime:    rd.Requisitions.ApproveTime,
				ApproveRemark:  rd.Requisitions.ApproveRemark,
				CreatedAt:      rd.Requisitions.CreatedAt,
				UpdatedAt:      rd.Requisitions.UpdatedAt,
			},
			Total: rd.Total,
			//Details: details,
			Details: nil,
		}
		list = append(list, tRd)
	}

	resp = &types.SearchRequisitionResp{
		RequisitionsWithDetails: list,
		Total:                   ret.Total,
	}
	return
}
