package salesorderlogic

import (
	"context"
	"encoding/json"
	"erp/app/sale/rpc/internal/code"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/mr"

	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/internal/types"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSalesOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSalesOrderLogic {
	return &SearchSalesOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchSalesOrderLogic) SearchSalesOrder(in *pb.SearchSalesOrderReq) (*pb.SearchSalesOrderResp, error) {
	params := &types.SearchOrderParams{
		SearchComm: types.SearchComm{Page: in.Page, Limit: in.Limit},
		OrderNo:    in.OrderNo,
		CustomerId: in.CustomerId,
		Status:     in.Status,
		SalesmanId: in.SalesmanId,
		StartOrderDate: func() time.Time {
			if in.StartOrderDate == 0 {
				return time.Time{}
			}
			return time.Unix(in.StartOrderDate, 0)
		}(),
		EndOrderDate: func() time.Time {
			if in.EndOrderDate == 0 {
				return time.Time{}
			}
			return time.Unix(in.EndOrderDate, 0)
		}(),
	}
	orders, total, err := l.svcCtx.SalesOrderModel.Search(l.ctx, params)
	if err != nil {
		return nil, code.GetSalesOrderFail
	}

	resp := &pb.SearchSalesOrderResp{Total: total}
	for _, o := range orders {
		// 优先读明细缓存id
		key := fmt.Sprintf(types.CacheErpSaleSalesOrderDetailIdsByOrderId, o.Id)
		idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)

		useModel := false
		// 缓存找不到就得从model找
		if err != nil || len(idsStr) == 0 {
			useModel = true
		}
		var pbDetails []*pb.SalesOrderDetail
		if useModel {
			pbDetails, err = l.ModelGetDetails(o.Id)
			if err != nil {
				return nil, code.GetSalesOrderFail
			}
		} else {
			pbDetails, err = l.CacheGetDetails(idsStr)
			if err != nil {
				return nil, code.GetSalesOrderFail
			}
		}
		pbOrder := &pb.SalesOrder{
			Id:           o.Id,
			OrderNo:      o.OrderNo,
			CustomerId:   o.CustomerId,
			OrderDate:    o.OrderDate.Unix(),
			PromisedDate: o.PromisedDate.Time.Unix(),
			TotalAmount:  o.TotalAmount,
			Status:       o.Status,
			SalesmanId:   o.SalesmanId,
			ContractUrl:  o.ContractUrl.String,
			CreatedAt:    o.CreatedAt.Unix(),
		}
		resp.OrderWithDetails = append(resp.OrderWithDetails, &pb.OrderWithDetails{
			SalesOrder:       pbOrder,
			Total:            int64(len(pbDetails)),
			SalesOrderDetail: pbDetails,
		})
	}
	return resp, nil
}

func (l *SearchSalesOrderLogic) ModelGetDetails(orderId int64) ([]*pb.SalesOrderDetail, error) {
	details, err := l.svcCtx.SalesOrderDetailModel.ListByOrderId(l.ctx, orderId)
	if err != nil {
		return nil, code.GetSalesOrderFail
	}
	var pbDetails []*pb.SalesOrderDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &pb.SalesOrderDetail{
			Id:           d.Id,
			OrderId:      d.OrderId,
			ProductId:    d.ProductId,
			ProductName:  d.ProductName.String,
			Unit:         d.Unit,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			DeliveredQty: d.DeliveredQty,
			Remark:       d.Remark.String,
		})
		// 缓存明细
		key := fmt.Sprintf(types.CacheErpSaleSalesOrderDetailIdPrefix, d.Id)
		bytes, _ := json.Marshal(pbDetails)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, string(bytes), 24*60*60*3) // 3天

		// 缓存 ids
		key = fmt.Sprintf(types.CacheErpSaleSalesOrderDetailIdsByOrderId, orderId)
		_, _ = l.svcCtx.BizRedis.LpushCtx(l.ctx, key, strconv.FormatInt(d.Id, 10))
	}
	return pbDetails, nil
}

func (l *SearchSalesOrderLogic) CacheGetDetails(idsStr []string) ([]*pb.SalesOrderDetail, error) {
	generate := func(source chan<- int64) {
		for _, idS := range idsStr {
			id, _ := strconv.ParseInt(idS, 10, 64)
			source <- id
		}
	}

	mapper := func(id int64, writer mr.Writer[*pb.SalesOrderDetail], cancel func(error)) {
		d, err := l.svcCtx.SalesOrderDetailModel.FindOne(l.ctx, id)
		if err != nil {
			return
		}
		writer.Write(&pb.SalesOrderDetail{
			Id:           d.Id,
			OrderId:      d.OrderId,
			ProductId:    d.ProductId,
			ProductName:  d.ProductName.String,
			Unit:         d.Unit,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			DeliveredQty: d.DeliveredQty,
			Remark:       d.Remark.String,
		})
	}

	reducer := func(pipe <-chan *pb.SalesOrderDetail, writer mr.Writer[[]*pb.SalesOrderDetail], cancel func(error)) {
		result := []*pb.SalesOrderDetail{}
		for p := range pipe {
			result = append(result, p)
		}
		writer.Write(result)
	}
	details, err := mr.MapReduce[int64, *pb.SalesOrderDetail, []*pb.SalesOrderDetail](generate, mapper, reducer)
	if err != nil {
		return nil, err
	}
	return details, nil
}
