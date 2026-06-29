package salesdeliverylogic

import (
	"context"
	"encoding/json"
	"erp/app/sale/rpc/internal/code"
	"erp/app/sale/rpc/internal/types"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/mr"

	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetSalesDeliveryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSalesDeliveryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesDeliveryByIdLogic {
	return &GetSalesDeliveryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	// 7 * 24 * 60 * +0
	cacheSec = 604800
)

func (l *GetSalesDeliveryByIdLogic) GetSalesDeliveryById(in *pb.GetSalesDeliveryByIdReq) (*pb.GetSalesDeliveryByIdResp, error) {
	delivery, err := l.svcCtx.SalesDeliveryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.SalesDeliveryNotFound
		}
		return nil, code.GetDeliveryFail
	}

	// 缓存 delivery
	key := fmt.Sprintf(types.CacheErpSaleSalesDeliveryIdPrefix, delivery.Id)
	bytes, _ := json.Marshal(delivery)
	_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, string(bytes), cacheSec)

	salesOrder, err := l.svcCtx.SalesOrderModel.FindOne(l.ctx, delivery.OrderId.Int64)
	if err != nil {
		return nil, err
	}

	// 返回结果组装
	out := &pb.GetSalesDeliveryByIdResp{
		DeliveryWithDetails: &pb.DeliveryWithDetails{
			SalesDelivery: &pb.SalesDelivery{
				Id:            delivery.Id,
				DeliveryNo:    delivery.DeliveryNo,
				OrderId:       delivery.OrderId.Int64,
				WarehouseId:   delivery.WarehouseId,
				DeliveryDate:  delivery.DeliveryDate.Unix(),
				TotalQuantity: delivery.TotalQuantity,
				TotalAmount:   delivery.TotalAmount,
				Status:        delivery.Status,
				CreatedBy:     delivery.CreatedBy,
				CreatedAt:     delivery.CreatedAt.Unix(),
				OrderNo:       salesOrder.OrderNo,
			},
			Total:               0,
			SalesDeliveryDetail: nil,
		},
	}

	// 标记是走Redis还是走model
	useModel := true

	// 从缓存中读取 details 的 id
	key = fmt.Sprintf(types.CacheErpSaleSalesDeliveryDetailIdsByDeliveryId, delivery.Id)
	idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)

	if err != nil {
		useModel = false
	}
	if len(idsStr) == 0 {
		useModel = false
	}

	var details []*pb.SalesDeliveryDetail
	if useModel {
		details, err = l.ModelGetDetails(delivery.Id)
		if err != nil {
			return nil, code.GetDeliveryFail
		}
	} else {
		details, err = l.CacheGetDetails(idsStr)
		if err != nil {
			return nil, code.GetDeliveryFail
		}
	}
	out.DeliveryWithDetails.Total = int64(len(out.DeliveryWithDetails.SalesDeliveryDetail))
	out.DeliveryWithDetails.SalesDeliveryDetail = details

	return out, nil
}

func (l *GetSalesDeliveryByIdLogic) ModelGetDetails(deliveryId int64) ([]*pb.SalesDeliveryDetail, error) {
	details, err := l.svcCtx.SalesDeliveryDetailModel.FindByDeliveryId(l.ctx, deliveryId)
	if err != nil {
		return nil, code.GetDeliveryFail
	}

	var pbDetails []*pb.SalesDeliveryDetail

	for _, d := range details {
		pbDetails = append(pbDetails, &pb.SalesDeliveryDetail{
			Id:          d.Id,
			DeliveryId:  d.DeliveryId,
			ProductId:   d.ProductId,
			ProductName: d.ProductName.String,
			Unit:        d.Unit,
			Quantity:    d.Quantity,
			UnitPrice:   d.UnitPrice,
			Amount:      d.Amount,
			BatchId:     d.BatchId.Int64,
		})

		// 缓存 detail
		key := fmt.Sprintf(types.CacheErpSaleSalesDeliveryDetailIdPrefix, d.Id)
		bytes, _ := json.Marshal(pbDetails)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, string(bytes), cacheSec)

		// 缓存 ids
		key = fmt.Sprintf(types.CacheErpSaleSalesDeliveryDetailIdsByDeliveryId, deliveryId)
		_, _ = l.svcCtx.BizRedis.LpushCtx(l.ctx, key, strconv.FormatInt(d.Id, 10))
	}

	return pbDetails, nil
}

func (l *GetSalesDeliveryByIdLogic) CacheGetDetails(idsStr []string) ([]*pb.SalesDeliveryDetail, error) {
	generate := func(source chan<- int64) {
		for _, idS := range idsStr {
			id, _ := strconv.ParseInt(idS, 10, 64)
			source <- id
		}
	}

	mapper := func(id int64, writer mr.Writer[*pb.SalesDeliveryDetail], cancel func(error)) {
		d, err := l.svcCtx.SalesDeliveryDetailModel.FindOne(l.ctx, id)
		if err != nil {
			return
		}
		writer.Write(&pb.SalesDeliveryDetail{
			Id:          d.Id,
			DeliveryId:  d.DeliveryId,
			ProductId:   d.ProductId,
			ProductName: d.ProductName.String,
			Unit:        d.Unit,
			Quantity:    d.Quantity,
			UnitPrice:   d.UnitPrice,
			Amount:      d.Amount,
			BatchId:     d.BatchId.Int64,
		})
	}

	reducer := func(pipe <-chan *pb.SalesDeliveryDetail, writer mr.Writer[[]*pb.SalesDeliveryDetail], cancel func(error)) {
		result := []*pb.SalesDeliveryDetail{}
		for p := range pipe {
			result = append(result, p)
		}
		writer.Write(result)
	}
	details, err := mr.MapReduce[int64, *pb.SalesDeliveryDetail, []*pb.SalesDeliveryDetail](generate, mapper, reducer)
	if err != nil {
		return nil, err
	}
	return details, nil
}
