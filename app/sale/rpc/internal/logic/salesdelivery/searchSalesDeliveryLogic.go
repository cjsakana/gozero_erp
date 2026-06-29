package salesdeliverylogic

import (
	"context"
	"encoding/json"
	"erp/app/sale/rpc/internal/code"
	"erp/app/sale/rpc/internal/types"
	"time"

	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/pb"

	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type SearchSalesDeliveryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchSalesDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSalesDeliveryLogic {
	return &SearchSalesDeliveryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchSalesDeliveryLogic) SearchSalesDelivery(in *pb.SearchSalesDeliveryReq) (*pb.SearchSalesDeliveryResp, error) {
	deliveries, total, err := l.svcCtx.SalesDeliveryModel.Search(l.ctx, &types.SearchDeliveryParams{
		SearchComm: types.SearchComm{
			Page:  in.Page,
			Limit: in.Limit,
		},
		DeliveryNo:   in.DeliveryNo,
		OrderId:      in.OrderId,
		WarehouseId:  in.WarehouseId,
		DeliveryDate: time.Unix(in.DeliveryDate, 0),
		Status:       in.Status,
	})
	if err != nil {
		return nil, code.GetDeliveryFail
	}

	out := &pb.SearchSalesDeliveryResp{
		Total:               total,
		DeliveryWithDetails: nil,
	}

	pbDeliveries := []*pb.DeliveryWithDetails{}
	for _, delivery := range deliveries {
		// 缓存 ids key
		key := fmt.Sprintf(types.CacheErpSaleSalesDeliveryDetailIdsByDeliveryId, delivery.Id)
		idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)

		useModel := false
		if err != nil || len(idsStr) == 0 {
			useModel = true
		}

		var pbDetails []*pb.SalesDeliveryDetail
		if useModel {
			pbDetails, err = l.ModelGetDetails(delivery.Id)
			if err != nil {

				return nil, code.GetDeliveryFail
			}
		} else {
			pbDetails, err = l.CacheGetDetails(idsStr)
			if err != nil {
				return nil, code.GetDeliveryFail
			}
		}

		salesOrder, err := l.svcCtx.SalesOrderModel.FindOne(l.ctx, delivery.OrderId.Int64)
		if err != nil {
			return nil, err
		}
		pbDeliveries = append(pbDeliveries, &pb.DeliveryWithDetails{
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
			Total:               int64(len(pbDetails)),
			SalesDeliveryDetail: pbDetails,
		})
	}

	out.DeliveryWithDetails = pbDeliveries

	return out, nil
}

// ModelGetDetails 复用 getSalesDeliveryByIdLogic.go 逻辑
func (l *SearchSalesDeliveryLogic) ModelGetDetails(deliveryId int64) ([]*pb.SalesDeliveryDetail, error) {
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

// CacheGetDetails 复用 getSalesDeliveryByIdLogic.go 逻辑
func (l *SearchSalesDeliveryLogic) CacheGetDetails(idsStr []string) ([]*pb.SalesDeliveryDetail, error) {
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
