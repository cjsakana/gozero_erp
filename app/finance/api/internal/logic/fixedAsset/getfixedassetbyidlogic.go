package fixedAsset

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

type GetFixedAssetByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFixedAssetByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFixedAssetByIdLogic {
	return &GetFixedAssetByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFixedAssetByIdLogic) GetFixedAssetById(req *types.GetFixedAssetByIdReq) (resp *types.GetFixedAssetByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.FinanceRPC.GetFixedAssetById(l.ctx, &pb.GetFixedAssetByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 获取供应商名称
	var supplierName string
	if ret.FixedAsset.SupplierId > 0 {
		supplierResp, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &supplierpb.GetSupplierByIdReq{
			Id: ret.FixedAsset.SupplierId,
		})
		if err != nil {
			logx.Errorf("查询供应商信息失败: supplierId=%d, err=%v", ret.FixedAsset.SupplierId, err)
		} else if supplierResp.Supplier != nil {
			supplierName = supplierResp.Supplier.Name
		}
	}

	// 获取使用人（员工）信息
	var userNo, userName string
	if ret.FixedAsset.UserId > 0 {
		empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{
			Id: ret.FixedAsset.UserId,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", ret.FixedAsset.UserId, err)
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			userNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
			userName = empResp.EmployeeNonSensitiveDetail.Name
		}
	}

	// 获取使用部门名称
	var departmentName string
	if ret.FixedAsset.DepartmentId > 0 {
		hrResp, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &hrpb.GetDepartmentByIdReq{
			Id: ret.FixedAsset.DepartmentId,
		})
		if err != nil {
			logx.Errorf("查询部门失败: departmentId=%d, err=%v", ret.FixedAsset.DepartmentId, err)
		} else if hrResp.Department != nil {
			departmentName = hrResp.Department.Name
		}
	}
	resp = &types.GetFixedAssetByIdResp{
		FixedAsset: types.FixedAsset{
			Id:                 util.Int64ToString(ret.FixedAsset.Id),
			AssetNo:            ret.FixedAsset.AssetNo,
			AssetName:          ret.FixedAsset.AssetName,
			Category:           ret.FixedAsset.Category,
			PurchaseDate:       ret.FixedAsset.PurchaseDate,
			PurchasePrice:      ret.FixedAsset.PurchasePrice,
			SupplierId:         util.Int64ToString(ret.FixedAsset.SupplierId),
			SupplierName:       supplierName,
			DepartmentId:       util.Int64ToString(ret.FixedAsset.DepartmentId),
			DepartmentName:     departmentName,
			UserId:             util.Int64ToString(ret.FixedAsset.UserId),
			UserName:           userName,
			UserNo:             userNo,
			Status:             ret.FixedAsset.Status,
			Location:           ret.FixedAsset.Location,
			DepreciationMethod: ret.FixedAsset.DepreciationMethod,
			UsefulLife:         ret.FixedAsset.UsefulLife,
			CreatedAt:          ret.FixedAsset.CreatedAt,
		},
	}
	return
}
