package purchaserequisitionlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"erp/app/purchase/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchRequisitionLogic {
	return &SearchRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询采购申请
func (l *SearchRequisitionLogic) SearchRequisition(in *pb.SearchRequisitionReq) (*pb.SearchRequisitionResp, error) {
	params := &types.SearchRequisitionParams{
		SearchComm:    types.SearchComm{Page: in.Page, Limit: in.Limit},
		RequisitionNo: in.RequisitionNo,
		DepartmentId:  in.DepartmentId,
		ApplicantId:   in.ApplicantId,
		ApproverId:    in.ApproverId,
		Status:        in.Status,
	}
	requisitions, total, err := l.svcCtx.PurchaseRequisitionModel.Search(l.ctx, params)
	if err != nil {

		return nil, code.GetRequisitionFail

	}

	var pbRequisitionD []*pb.PurchaseRequisitionWithDetails
	for _, req := range requisitions {
		pbRequisitionWithDetails := &pb.PurchaseRequisitionWithDetails{
			Requisitions: &pb.PurchaseRequisition{
				Id:            req.Id,
				RequisitionNo: req.RequisitionNo,
				DepartmentId:  req.DepartmentId,
				ApplicantId:   req.ApplicantId,
				RequestDate:   req.RequestDate.Unix(),
				TotalAmount:   req.TotalAmount,
				Status:        req.Status,
				ApproverId:    req.ApproverId.Int64,
				ApproveTime: func() int64 {
					if req.ApproveTime.Valid {
						return req.ApproveTime.Time.Unix()
					}
					return 0
				}(),
				ApproveRemark: req.ApproveRemark.String,
				CreatedAt:     req.CreatedAt.Unix(),
				UpdatedAt:     req.UpdatedAt.Unix(),
			},
		}
		details, err := l.svcCtx.PurchaseRequisitionDetailModel.ListByRequisitionId(l.ctx, req.Id)
		if err != nil {

			return nil, code.GetRequisitionFail

		}
		pbRequisitionWithDetails.Total = int64(len(details))
		for _, detail := range details {
			pbRequisitionWithDetails.Details = append(pbRequisitionWithDetails.Details, &pb.PurchaseRequisitionDetail{
				Id:            detail.Id,
				RequisitionId: detail.RequisitionId,
				ProductId:     detail.ProductId.Int64,
				ProductName:   detail.ProductName.String,
				CategoryType:  detail.CategoryType,
				Quantity:      detail.Quantity,
				UnitPrice:     detail.UnitPrice.Float64,
				Amount:        detail.Amount.Float64,
				Remark:        detail.Remark.String,
			})
		}

		pbRequisitionD = append(pbRequisitionD, pbRequisitionWithDetails)
	}

	return &pb.SearchRequisitionResp{
		RequisitionsWithDetails: pbRequisitionD,
		Total:                   total,
	}, nil
}
