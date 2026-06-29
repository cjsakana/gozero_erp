package purchaserequisitionlogic

import (
	"context"
	"encoding/json"
	"erp/app/purchase/rpc/internal/types"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/mr"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/pb"

	"erp/app/purchase/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetRequisitionWithDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRequisitionWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequisitionWithDetailsLogic {
	return &GetRequisitionWithDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取采购申请及明细
func (l *GetRequisitionWithDetailsLogic) GetRequisitionWithDetails(in *pb.GetRequisitionWithDetailsReq) (*pb.GetRequisitionWithDetailsResp, error) {
	requisition, err := l.svcCtx.PurchaseRequisitionModel.FindOne(l.ctx, in.RequisitionId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PurchaseRequisitionNotFound
		}
		return nil, code.PurchaseRequisitionNotFound
	}

	// 优先读ids缓存
	key := fmt.Sprintf(types.CacheErpPurchasePurchaseRequisitionDetailIdsByRequisitionId, in.RequisitionId)
	idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)
	useModel := false
	if err != nil || len(idsStr) == 0 {
		useModel = true
	}

	var pbDetails []*pb.PurchaseRequisitionDetail
	if useModel {
		pbDetails, err = l.ModelGetDetails(requisition.Id)
		if err != nil {
			return nil, err
		}
	} else {
		pbDetails, err = l.CacheGetDetails(idsStr)
		if err != nil {
			return nil, err
		}
	}

	pbRequisition := &pb.PurchaseRequisition{
		Id:            requisition.Id,
		RequisitionNo: requisition.RequisitionNo,
		DepartmentId:  requisition.DepartmentId,
		ApplicantId:   requisition.ApplicantId,
		RequestDate:   requisition.RequestDate.Unix(),
		TotalAmount:   requisition.TotalAmount,
		Status:        requisition.Status,
		ApproverId:    requisition.ApproverId.Int64,
		ApproveTime:   requisition.ApproveTime.Time.Unix(),
		ApproveRemark: requisition.ApproveRemark.String,
		CreatedAt:     requisition.CreatedAt.Unix(),
		UpdatedAt:     requisition.UpdatedAt.Unix(),
	}
	if requisition.ApproveTime.Valid {
		pbRequisition.ApproveTime = requisition.ApproveTime.Time.Unix()
	}

	return &pb.GetRequisitionWithDetailsResp{
		Requisition: pbRequisition,
		Details:     pbDetails,
	}, nil
}

func (l *GetRequisitionWithDetailsLogic) ModelGetDetails(requisitionId int64) ([]*pb.PurchaseRequisitionDetail, error) {
	details, err := l.svcCtx.PurchaseRequisitionDetailModel.ListByRequisitionId(l.ctx, requisitionId)
	if err != nil {
		return nil, err
	}
	var pbDetails []*pb.PurchaseRequisitionDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &pb.PurchaseRequisitionDetail{
			Id:            d.Id,
			RequisitionId: d.RequisitionId,
			ProductId:     d.ProductId.Int64,
			ProductName:   d.ProductName.String,
			CategoryType:  d.CategoryType,
			Quantity:      d.Quantity,
			UnitPrice:     d.UnitPrice.Float64,
			Amount:        d.Amount.Float64,
			Remark:        d.Remark.String,
		})
		// 缓存明细
		key := fmt.Sprintf(types.CacheErpPurchasePurchaseRequisitionDetailIdPrefix, d.Id)
		bytes, _ := json.Marshal(pbDetails)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, string(bytes), 24*60*60*3) // 3天

		// 缓存 ids
		key = fmt.Sprintf(types.CacheErpPurchasePurchaseRequisitionDetailIdsByRequisitionId, requisitionId)
		_, _ = l.svcCtx.BizRedis.LpushCtx(l.ctx, key, strconv.FormatInt(d.Id, 10))
	}
	return pbDetails, nil
}

func (l *GetRequisitionWithDetailsLogic) CacheGetDetails(idsStr []string) ([]*pb.PurchaseRequisitionDetail, error) {
	generate := func(source chan<- int64) {
		for _, idS := range idsStr {
			id, _ := strconv.ParseInt(idS, 10, 64)
			source <- id
		}
	}

	mapper := func(id int64, writer mr.Writer[*pb.PurchaseRequisitionDetail], cancel func(error)) {
		d, err := l.svcCtx.PurchaseRequisitionDetailModel.FindOne(l.ctx, id)
		if err != nil {
			return
		}
		writer.Write(&pb.PurchaseRequisitionDetail{
			Id:            d.Id,
			RequisitionId: d.RequisitionId,
			ProductId:     d.ProductId.Int64,
			ProductName:   d.ProductName.String,
			CategoryType:  d.CategoryType,
			Quantity:      d.Quantity,
			UnitPrice:     d.UnitPrice.Float64,
			Amount:        d.Amount.Float64,
			Remark:        d.Remark.String,
		})
	}

	reducer := func(pipe <-chan *pb.PurchaseRequisitionDetail, writer mr.Writer[[]*pb.PurchaseRequisitionDetail], cancel func(error)) {
		result := []*pb.PurchaseRequisitionDetail{}
		for p := range pipe {
			result = append(result, p)
		}
		writer.Write(result)
	}
	details, err := mr.MapReduce[int64, *pb.PurchaseRequisitionDetail, []*pb.PurchaseRequisitionDetail](generate, mapper, reducer)
	if err != nil {
		return nil, err
	}
	return details, nil
}
