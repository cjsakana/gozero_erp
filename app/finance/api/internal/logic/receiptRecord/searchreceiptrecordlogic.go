package receiptRecord

import (
	"context"
	customerpb "erp/app/customer/rpc/pb"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	hrpb "erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchReceiptRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchReceiptRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchReceiptRecordLogic {
	return &SearchReceiptRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type receiptEmployeeInfo struct {
	EmployeeNo string
	Name       string
}

func (l *SearchReceiptRecordLogic) SearchReceiptRecord(req *types.SearchReceiptRecordReq) (resp *types.SearchReceiptRecordResp, err error) {
	customerId, err := util.StringToInt64(req.CustomerId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.SearchReceiptRecord(l.ctx, &pb.SearchReceiptRecordReq{
		Page:          req.Page,
		Limit:         req.Limit,
		ReceiptNo:     req.ReceiptNo,
		CustomerId:    customerId,
		ReceiptType:   req.ReceiptType,
		ReceiptMethod: req.ReceiptMethod,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 批量获取客户名称（去重）
	customerMap := make(map[int64]string)
	for _, rr := range ret.ReceiptRecord {
		if rr.CustomerId > 0 {
			if _, ok := customerMap[rr.CustomerId]; !ok {
				customerResp, err := l.svcCtx.CustomerRPC.GetCustomerById(l.ctx, &customerpb.GetCustomerByIdReq{Id: rr.CustomerId})
				if err != nil {
					logx.Errorf("查询客户信息失败: customerId=%d, err=%v", rr.CustomerId, err)
					customerMap[rr.CustomerId] = ""
				} else if customerResp.Customer != nil {
					customerMap[rr.CustomerId] = customerResp.Customer.Name
				}
			}
		}
	}

	// 批量获取操作人（员工）信息（去重）
	employeeMap := make(map[int64]*receiptEmployeeInfo)
	for _, rr := range ret.ReceiptRecord {
		if rr.OperatorId > 0 {
			if _, ok := employeeMap[rr.OperatorId]; !ok {
				empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{Id: rr.OperatorId})
				if err != nil {
					logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", rr.OperatorId, err)
					employeeMap[rr.OperatorId] = &receiptEmployeeInfo{}
				} else if empResp.EmployeeNonSensitiveDetail != nil {
					employeeMap[rr.OperatorId] = &receiptEmployeeInfo{
						EmployeeNo: empResp.EmployeeNonSensitiveDetail.EmployeeNo,
						Name:       empResp.EmployeeNonSensitiveDetail.Name,
					}
				}
			}
		}
	}

	list := make([]*types.ReceiptRecord, 0, len(ret.ReceiptRecord))
	for _, rr := range ret.ReceiptRecord {
		var customerName, operatorNo, operatorName string
		if name, ok := customerMap[rr.CustomerId]; ok {
			customerName = name
		}
		if info, ok := employeeMap[rr.OperatorId]; ok {
			operatorNo = info.EmployeeNo
			operatorName = info.Name
		}
		item := &types.ReceiptRecord{
			Id:            util.Int64ToString(rr.Id),
			ReceiptNo:     rr.ReceiptNo,
			CustomerId:    util.Int64ToString(rr.CustomerId),
			CustomerName:  customerName,
			ReceiptType:   rr.ReceiptType,
			Amount:        rr.Amount,
			ReceiptDate:   rr.ReceiptDate,
			ReceiptMethod: rr.ReceiptMethod,
			OrderId:       util.Int64ToString(rr.OrderId),
			Status:        rr.Status,
			VerifyStatus:  rr.VerifyStatus,
			OperatorId:    util.Int64ToString(rr.OperatorId),
			OperatorNo:    operatorNo,
			OperatorName:  operatorName,
			CreatedAt:     rr.CreatedAt,
		}

		list = append(list, item)
	}

	resp = &types.SearchReceiptRecordResp{
		ReceiptRecord: list,
		Total:         ret.Total,
	}
	return
}
