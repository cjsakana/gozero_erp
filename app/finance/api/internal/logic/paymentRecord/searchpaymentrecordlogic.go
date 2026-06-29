package paymentRecord

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	hrpb "erp/app/hr/rpc/pb"
	supplierpb "erp/app/supplier/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchPaymentRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchPaymentRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPaymentRecordLogic {
	return &SearchPaymentRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type paymentEmployeeInfo struct {
	EmployeeNo string
	Name       string
}

func (l *SearchPaymentRecordLogic) SearchPaymentRecord(req *types.SearchPaymentRecordReq) (resp *types.SearchPaymentRecordResp, err error) {
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.SearchPaymentRecord(l.ctx, &pb.SearchPaymentRecordReq{
		Page:          req.Page,
		Limit:         req.Limit,
		PaymentNo:     req.PaymentNo,
		SupplierId:    supplierId,
		PaymentType:   req.PaymentType,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 批量获取供应商名称（去重）
	supplierMap := make(map[int64]string)
	for _, pr := range ret.PaymentRecord {
		if pr.SupplierId > 0 {
			if _, ok := supplierMap[pr.SupplierId]; !ok {
				supplierResp, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &supplierpb.GetSupplierByIdReq{Id: pr.SupplierId})
				if err != nil {
					logx.Errorf("查询供应商信息失败: supplierId=%d, err=%v", pr.SupplierId, err)
					supplierMap[pr.SupplierId] = ""
				} else if supplierResp.Supplier != nil {
					supplierMap[pr.SupplierId] = supplierResp.Supplier.Name
				}
			}
		}
	}

	// 批量获取操作人（员工）信息（去重）
	employeeMap := make(map[int64]*paymentEmployeeInfo)
	for _, pr := range ret.PaymentRecord {
		if pr.OperatorId > 0 {
			if _, ok := employeeMap[pr.OperatorId]; !ok {
				empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{Id: pr.OperatorId})
				if err != nil {
					logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", pr.OperatorId, err)
					employeeMap[pr.OperatorId] = &paymentEmployeeInfo{}
				} else if empResp.EmployeeNonSensitiveDetail != nil {
					employeeMap[pr.OperatorId] = &paymentEmployeeInfo{
						EmployeeNo: empResp.EmployeeNonSensitiveDetail.EmployeeNo,
						Name:       empResp.EmployeeNonSensitiveDetail.Name,
					}
				}
			}
		}
	}

	list := make([]*types.PaymentRecord, 0, len(ret.PaymentRecord))
	for _, pr := range ret.PaymentRecord {
		var supplierName, operatorNo, operatorName string
		if name, ok := supplierMap[pr.SupplierId]; ok {
			supplierName = name
		}
		if info, ok := employeeMap[pr.OperatorId]; ok {
			operatorNo = info.EmployeeNo
			operatorName = info.Name
		}
		item := &types.PaymentRecord{
			Id:            util.Int64ToString(pr.Id),
			PaymentNo:     pr.PaymentNo,
			SupplierId:    util.Int64ToString(pr.SupplierId),
			SupplierName:  supplierName,
			PaymentType:   pr.PaymentType,
			Amount:        pr.Amount,
			PaymentDate:   pr.PaymentDate,
			PaymentMethod: pr.PaymentMethod,
			OrderId:       util.Int64ToString(pr.OrderId),
			Status:        pr.Status,
			VerifyStatus:  pr.VerifyStatus,
			OperatorId:    util.Int64ToString(pr.OperatorId),
			OperatorNo:    operatorNo,
			OperatorName:  operatorName,
			CreatedAt:     pr.CreatedAt,
		}

		list = append(list, item)
	}

	resp = &types.SearchPaymentRecordResp{
		PaymentRecord: list,
		Total:         ret.Total,
	}
	return
}
