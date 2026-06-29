package salesorderlogic

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

type GetSalesOrderByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSalesOrderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesOrderByIdLogic {
	return &GetSalesOrderByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSalesOrderByIdLogic) GetSalesOrderById(in *pb.GetSalesOrderByIdReq) (*pb.GetSalesOrderByIdResp, error) {
	order, err := l.svcCtx.SalesOrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.SalesOrderNotFound
		}
		return nil, code.GetSalesOrderFail
	}
	// 优先读ids缓存
	key := fmt.Sprintf(types.CacheErpSaleSalesOrderDetailIdsByOrderId, in.Id)
	idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)
	useModel := true
	if err != nil || len(idsStr) == 0 {
		useModel = false
	}

	var pbDetails []*pb.SalesOrderDetail
	if useModel {
		pbDetails, err = l.ModelGetDetails(order.Id)
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
		Id:           order.Id,
		OrderNo:      order.OrderNo,
		CustomerId:   order.CustomerId,
		OrderDate:    order.OrderDate.Unix(),
		PromisedDate: order.PromisedDate.Time.Unix(),
		TotalAmount:  order.TotalAmount,
		Status:       order.Status,
		SalesmanId:   order.SalesmanId,
		ContractUrl:  order.ContractUrl.String,
		CreatedAt:    order.CreatedAt.Unix(),
	}

	return &pb.GetSalesOrderByIdResp{OrderWithDetails: &pb.OrderWithDetails{
		SalesOrder:       pbOrder,
		Total:            int64(len(pbDetails)),
		SalesOrderDetail: pbDetails,
	}}, nil
}

func (l *GetSalesOrderByIdLogic) ModelGetDetails(orderId int64) ([]*pb.SalesOrderDetail, error) {
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

func (l *GetSalesOrderByIdLogic) CacheGetDetails(idsStr []string) ([]*pb.SalesOrderDetail, error) {
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
