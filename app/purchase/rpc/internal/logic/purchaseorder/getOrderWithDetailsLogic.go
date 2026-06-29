package purchaseorderlogic

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

type GetOrderWithDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderWithDetailsLogic {
	return &GetOrderWithDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取采购订单及明细
func (l *GetOrderWithDetailsLogic) GetOrderWithDetails(in *pb.GetOrderWithDetailsReq) (*pb.GetOrderWithDetailsResp, error) {
	order, err := l.svcCtx.PurchaseOrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PurchaseOrderNotFound
		}
		return nil, code.PurchaseOrderNotFound
	}

	// 优先读ids缓存
	key := fmt.Sprintf(types.CacheErpPurchasePurchaseOrderDetailIdsByOrderId, in.OrderId)
	idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)
	useModel := false
	if err != nil || len(idsStr) == 0 {
		useModel = true
	}

	var pbDetails []*pb.PurchaseOrderDetail
	if useModel {
		pbDetails, err = l.ModelGetDetails(order.Id)
		if err != nil {
			return nil, err
		}
	} else {
		pbDetails, err = l.CacheGetDetails(idsStr)
		if err != nil {
			return nil, err
		}
	}

	pbOrder := &pb.PurchaseOrder{
		Id:           order.Id,
		OrderNo:      order.OrderNo,
		SupplierId:   order.SupplierId,
		OrderDate:    order.OrderDate.Unix(),
		ExpectedDate: order.ExpectedDate.Time.Unix(),
		TotalAmount:  order.TotalAmount,
		Status:       order.Status,
		PurchaserId:  order.PurchaserId,
		ContractUrl:  order.ContractUrl.String,
		CreatedAt:    order.CreatedAt.Unix(),
		UpdatedAt:    order.UpdatedAt.Unix(),
	}
	if order.ExpectedDate.Valid {
		pbOrder.ExpectedDate = order.ExpectedDate.Time.Unix()
	}

	return &pb.GetOrderWithDetailsResp{
		Order:   pbOrder,
		Details: pbDetails,
	}, nil
}

func (l *GetOrderWithDetailsLogic) ModelGetDetails(orderId int64) ([]*pb.PurchaseOrderDetail, error) {
	details, err := l.svcCtx.PurchaseOrderDetailModel.ListByOrderId(l.ctx, orderId)
	if err != nil {
		return nil, err
	}
	var pbDetails []*pb.PurchaseOrderDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &pb.PurchaseOrderDetail{
			Id:           d.Id,
			OrderId:      d.OrderId,
			ProductId:    d.ProductId.Int64,
			ProductName:  d.ProductName.String,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			ReceivedQty:  d.ReceivedQty,
			Remark:       d.Remark.String,
		})
		// 缓存明细
		key := fmt.Sprintf(types.CacheErpPurchasePurchaseOrderDetailIdPrefix, d.Id)
		bytes, _ := json.Marshal(pbDetails)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, string(bytes), 24*60*60*3) // 3天

		// 缓存 ids
		key = fmt.Sprintf(types.CacheErpPurchasePurchaseOrderDetailIdsByOrderId, orderId)
		_, _ = l.svcCtx.BizRedis.LpushCtx(l.ctx, key, strconv.FormatInt(d.Id, 10))
	}
	return pbDetails, nil
}

func (l *GetOrderWithDetailsLogic) CacheGetDetails(idsStr []string) ([]*pb.PurchaseOrderDetail, error) {
	generate := func(source chan<- int64) {
		for _, idS := range idsStr {
			id, _ := strconv.ParseInt(idS, 10, 64)
			source <- id
		}
	}

	mapper := func(id int64, writer mr.Writer[*pb.PurchaseOrderDetail], cancel func(error)) {
		d, err := l.svcCtx.PurchaseOrderDetailModel.FindOne(l.ctx, id)
		if err != nil {
			return
		}
		writer.Write(&pb.PurchaseOrderDetail{
			Id:           d.Id,
			OrderId:      d.OrderId,
			ProductId:    d.ProductId.Int64,
			ProductName:  d.ProductName.String,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			ReceivedQty:  d.ReceivedQty,
			Remark:       d.Remark.String,
		})
	}

	reducer := func(pipe <-chan *pb.PurchaseOrderDetail, writer mr.Writer[[]*pb.PurchaseOrderDetail], cancel func(error)) {
		result := []*pb.PurchaseOrderDetail{}
		for p := range pipe {
			result = append(result, p)
		}
		writer.Write(result)
	}
	details, err := mr.MapReduce[int64, *pb.PurchaseOrderDetail, []*pb.PurchaseOrderDetail](generate, mapper, reducer)
	if err != nil {
		return nil, err
	}
	return details, nil
}
