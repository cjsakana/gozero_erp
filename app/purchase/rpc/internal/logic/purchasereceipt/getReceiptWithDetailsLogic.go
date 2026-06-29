package purchasereceiptlogic

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

type GetReceiptWithDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReceiptWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiptWithDetailsLogic {
	return &GetReceiptWithDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取入库单及明细
func (l *GetReceiptWithDetailsLogic) GetReceiptWithDetails(in *pb.GetReceiptWithDetailsReq) (*pb.GetReceiptWithDetailsResp, error) {
	receipt, err := l.svcCtx.PurchaseReceiptModel.FindOne(l.ctx, in.ReceiptId)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.PurchaseReceiptNotFound
		}
		return nil, code.PurchaseReceiptNotFound
	}

	// 优先读ids缓存
	key := fmt.Sprintf(types.CacheErpPurchasePurchaseReceiptDetailIdsByReceiptId, in.ReceiptId)
	idsStr, err := l.svcCtx.BizRedis.LrangeCtx(l.ctx, key, 0, -1)
	useModel := false
	if err != nil || len(idsStr) == 0 {
		useModel = true
	}

	var pbDetails []*pb.PurchaseReceiptDetail
	if useModel {
		pbDetails, err = l.ModelGetDetails(receipt.Id)
		if err != nil {
			return nil, err
		}
	} else {
		pbDetails, err = l.CacheGetDetails(idsStr)
		if err != nil {
			return nil, err
		}
	}

	pbReceipt := &pb.PurchaseReceipt{
		Id:            receipt.Id,
		ReceiptNo:     receipt.ReceiptNo,
		OrderId:       receipt.OrderId.Int64,
		WarehouseId:   receipt.WarehouseId,
		ReceiptDate:   receipt.ReceiptDate.Unix(),
		TotalQuantity: receipt.TotalQuantity,
		TotalAmount:   receipt.TotalAmount,
		Status:        receipt.Status,
		CreatedAt:     receipt.CreatedAt.Unix(),
		CreatedBy:     receipt.CreatedBy,
	}
	if receipt.OrderId.Valid {
		pbReceipt.OrderId = receipt.OrderId.Int64
	}

	return &pb.GetReceiptWithDetailsResp{
		Receipt: pbReceipt,
		Details: pbDetails,
	}, nil
}

func (l *GetReceiptWithDetailsLogic) ModelGetDetails(receiptId int64) ([]*pb.PurchaseReceiptDetail, error) {
	details, err := l.svcCtx.PurchaseReceiptDetailModel.ListByReceiptId(l.ctx, receiptId)
	if err != nil {
		return nil, err
	}
	var pbDetails []*pb.PurchaseReceiptDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &pb.PurchaseReceiptDetail{
			Id:           d.Id,
			ReceiptId:    d.ReceiptId,
			ProductId:    d.ProductId,
			ProductName:  d.ProductName.String,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			BatchId: func() int64 {
				if d.BatchId.Valid {
					return d.BatchId.Int64
				}
				return 0
			}(),
		})
		// 缓存明细
		key := fmt.Sprintf(types.CacheErpPurchasePurchaseReceiptDetailIdPrefix, d.Id)
		bytes, _ := json.Marshal(pbDetails)
		_ = l.svcCtx.BizRedis.SetexCtx(l.ctx, key, string(bytes), 24*60*60*3) // 3天

		// 缓存 ids
		key = fmt.Sprintf(types.CacheErpPurchasePurchaseReceiptDetailIdsByReceiptId, receiptId)
		_, _ = l.svcCtx.BizRedis.LpushCtx(l.ctx, key, strconv.FormatInt(d.Id, 10))
	}
	return pbDetails, nil
}

func (l *GetReceiptWithDetailsLogic) CacheGetDetails(idsStr []string) ([]*pb.PurchaseReceiptDetail, error) {
	generate := func(source chan<- int64) {
		for _, idS := range idsStr {
			id, _ := strconv.ParseInt(idS, 10, 64)
			source <- id
		}
	}

	mapper := func(id int64, writer mr.Writer[*pb.PurchaseReceiptDetail], cancel func(error)) {
		d, err := l.svcCtx.PurchaseReceiptDetailModel.FindOne(l.ctx, id)
		if err != nil {
			return
		}
		writer.Write(&pb.PurchaseReceiptDetail{
			Id:           d.Id,
			ReceiptId:    d.ReceiptId,
			ProductId:    d.ProductId,
			ProductName:  d.ProductName.String,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			BatchId: func() int64 {
				if d.BatchId.Valid {
					return d.BatchId.Int64
				}
				return 0
			}(),
		})
	}

	reducer := func(pipe <-chan *pb.PurchaseReceiptDetail, writer mr.Writer[[]*pb.PurchaseReceiptDetail], cancel func(error)) {
		result := []*pb.PurchaseReceiptDetail{}
		for p := range pipe {
			result = append(result, p)
		}
		writer.Write(result)
	}
	details, err := mr.MapReduce[int64, *pb.PurchaseReceiptDetail, []*pb.PurchaseReceiptDetail](generate, mapper, reducer)
	if err != nil {
		return nil, err
	}
	return details, nil
}
