package supplier

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	pb2 "erp/app/hr/rpc/pb"
	"erp/app/supplier/rpc/pb"
	"erp/common/util"

	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSupplierDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSupplierDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupplierDetailLogic {
	return &GetSupplierDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSupplierDetailLogic) GetSupplierDetail(req *types.GetSupplierDetailReq) (resp *types.GetSupplierDetailResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &pb.GetSupplierByIdReq{Id: id})
	if err != nil {
		return nil, err
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	if _, ok := employeeMap[ret.Supplier.CreatedBy]; !ok {
		employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
			Id: ret.Supplier.CreatedBy,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", ret.Supplier.CreatedBy, err)
		}
		employeeMap[ret.Supplier.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
	}
	if _, ok := employeeMap[ret.Supplier.UpdatedBy]; !ok {
		employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
			Id: ret.Supplier.UpdatedBy,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", ret.Supplier.UpdatedBy, err)
		}
		employeeMap[ret.Supplier.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
	}

	resp = &types.GetSupplierDetailResp{
		Supplier: types.Supplier{
			Id:            util.Int64ToString(ret.Supplier.Id),
			Code:          ret.Supplier.Code,
			Uscc:          ret.Supplier.Uscc,
			Name:          ret.Supplier.Name,
			Contact:       ret.Supplier.Contact,
			Phone:         ret.Supplier.Phone,
			Address:       ret.Supplier.Address,
			PaymentTerms:  ret.Supplier.PaymentTerms,
			Credit:        ret.Supplier.Credit,
			IsActive:      ret.Supplier.IsActive,
			CreatedAt:     ret.Supplier.CreatedAt,
			CreatedBy:     util.Int64ToString(ret.Supplier.CreatedBy),
			CreatedByNo:   employeeMap[ret.Supplier.CreatedBy].EmployeeNo,
			CreatedByName: employeeMap[ret.Supplier.CreatedBy].Name,
			UpdatedAt:     ret.Supplier.UpdatedAt,
			UpdatedBy:     util.Int64ToString(ret.Supplier.UpdatedBy),
			UpdatedByNo:   employeeMap[ret.Supplier.UpdatedBy].EmployeeNo,
			UpdatedByName: employeeMap[ret.Supplier.UpdatedBy].Name,
		},
	}

	return
}
