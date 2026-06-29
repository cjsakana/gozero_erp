package supplier

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"
	"erp/app/supplier/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSupplierListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSupplierListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupplierListLogic {
	return &GetSupplierListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSupplierListLogic) GetSupplierList(req *types.GetSupplierListReq) (resp *types.GetSupplierListResp, err error) {

	ret, err := l.svcCtx.SupplierRPC.SearchSupplier(l.ctx, &pb.SearchSupplierReq{
		Page:     req.Page,
		Limit:    req.Limit,
		Code:     req.Code,
		Name:     req.Name,
		Contact:  req.Contact,
		Address:  req.Address,
		Credit:   req.Credit,
		IsActive: req.IsActive,
		Uscc:     req.Uscc,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetSupplierListResp{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	for _, supplier := range ret.Supplier {
		if _, ok := employeeMap[supplier.CreatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: supplier.CreatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", supplier.CreatedBy, err)
			}
			employeeMap[supplier.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		if _, ok := employeeMap[supplier.UpdatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: supplier.UpdatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", supplier.UpdatedBy, err)
			}
			employeeMap[supplier.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		resp.Suppliers = append(resp.Suppliers, &types.Supplier{
			Id:            util.Int64ToString(supplier.Id),
			Code:          supplier.Code,
			Uscc:          supplier.Uscc,
			Name:          supplier.Name,
			Contact:       supplier.Contact,
			Phone:         supplier.Phone,
			Address:       supplier.Address,
			PaymentTerms:  supplier.PaymentTerms,
			Credit:        supplier.Credit,
			IsActive:      supplier.IsActive,
			CreatedAt:     supplier.CreatedAt,
			CreatedBy:     util.Int64ToString(supplier.CreatedBy),
			CreatedByNo:   employeeMap[supplier.CreatedBy].EmployeeNo,
			CreatedByName: employeeMap[supplier.CreatedBy].Name,
			UpdatedAt:     supplier.UpdatedAt,
			UpdatedBy:     util.Int64ToString(supplier.UpdatedBy),
			UpdatedByNo:   employeeMap[supplier.UpdatedBy].EmployeeNo,
			UpdatedByName: employeeMap[supplier.UpdatedBy].Name,
		})
	}

	return
}
