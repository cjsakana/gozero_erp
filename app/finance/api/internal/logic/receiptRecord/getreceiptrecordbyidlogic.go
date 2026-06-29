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

type GetReceiptRecordByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReceiptRecordByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiptRecordByIdLogic {
	return &GetReceiptRecordByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReceiptRecordByIdLogic) GetReceiptRecordById(req *types.GetReceiptRecordByIdReq) (resp *types.GetReceiptRecordByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.GetReceiptRecordById(l.ctx, &pb.GetReceiptRecordByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 获取客户名称
	var customerName string
	if ret.ReceiptRecord.CustomerId > 0 {
		customerResp, err := l.svcCtx.CustomerRPC.GetCustomerById(l.ctx, &customerpb.GetCustomerByIdReq{
			Id: ret.ReceiptRecord.CustomerId,
		})
		if err != nil {
			logx.Errorf("查询客户信息失败: customerId=%d, err=%v", ret.ReceiptRecord.CustomerId, err)
		} else if customerResp.Customer != nil {
			customerName = customerResp.Customer.Name
		}
	}

	// 获取操作人（员工）信息
	var operatorNo, operatorName string
	if ret.ReceiptRecord.OperatorId > 0 {
		empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{
			Id: ret.ReceiptRecord.OperatorId,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", ret.ReceiptRecord.OperatorId, err)
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			operatorNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
			operatorName = empResp.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetReceiptRecordByIdResp{
		ReceiptRecord: types.ReceiptRecord{
			Id:            util.Int64ToString(ret.ReceiptRecord.Id),
			ReceiptNo:     ret.ReceiptRecord.ReceiptNo,
			CustomerId:    util.Int64ToString(ret.ReceiptRecord.CustomerId),
			CustomerName:  customerName,
			ReceiptType:   ret.ReceiptRecord.ReceiptType,
			Amount:        ret.ReceiptRecord.Amount,
			ReceiptDate:   ret.ReceiptRecord.ReceiptDate,
			ReceiptMethod: ret.ReceiptRecord.ReceiptMethod,
			OrderId:       util.Int64ToString(ret.ReceiptRecord.OrderId),
			Status:        ret.ReceiptRecord.Status,
			VerifyStatus:  ret.ReceiptRecord.VerifyStatus,
			OperatorId:    util.Int64ToString(ret.ReceiptRecord.OperatorId),
			OperatorNo:    operatorNo,
			OperatorName:  operatorName,
			CreatedAt:     ret.ReceiptRecord.CreatedAt,
		},
	}
	return
}
