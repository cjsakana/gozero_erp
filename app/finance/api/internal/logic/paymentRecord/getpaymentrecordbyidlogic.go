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

type GetPaymentRecordByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentRecordByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentRecordByIdLogic {
	return &GetPaymentRecordByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentRecordByIdLogic) GetPaymentRecordById(req *types.GetPaymentRecordByIdReq) (resp *types.GetPaymentRecordByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.GetPaymentRecordById(l.ctx, &pb.GetPaymentRecordByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 获取供应商名称
	var supplierName string
	if ret.PaymentRecord.SupplierId > 0 {
		supplierResp, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &supplierpb.GetSupplierByIdReq{
			Id: ret.PaymentRecord.SupplierId,
		})
		if err != nil {
			logx.Errorf("查询供应商信息失败: supplierId=%d, err=%v", ret.PaymentRecord.SupplierId, err)
		} else if supplierResp.Supplier != nil {
			supplierName = supplierResp.Supplier.Name
		}
	}

	// 获取操作人（员工）信息
	var operatorNo, operatorName string
	if ret.PaymentRecord.OperatorId > 0 {
		empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{
			Id: ret.PaymentRecord.OperatorId,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", ret.PaymentRecord.OperatorId, err)
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			operatorNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
			operatorName = empResp.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetPaymentRecordByIdResp{
		PaymentRecord: types.PaymentRecord{
			Id:            util.Int64ToString(ret.PaymentRecord.Id),
			PaymentNo:     ret.PaymentRecord.PaymentNo,
			SupplierId:    util.Int64ToString(ret.PaymentRecord.SupplierId),
			SupplierName:  supplierName,
			PaymentType:   ret.PaymentRecord.PaymentType,
			Amount:        ret.PaymentRecord.Amount,
			PaymentDate:   ret.PaymentRecord.PaymentDate,
			PaymentMethod: ret.PaymentRecord.PaymentMethod,
			OrderId:       util.Int64ToString(ret.PaymentRecord.OrderId),
			Status:        ret.PaymentRecord.Status,
			VerifyStatus:  ret.PaymentRecord.VerifyStatus,
			OperatorId:    util.Int64ToString(ret.PaymentRecord.OperatorId),
			OperatorNo:    operatorNo,
			OperatorName:  operatorName,
			CreatedAt:     ret.PaymentRecord.CreatedAt,
		},
	}
	return
}
