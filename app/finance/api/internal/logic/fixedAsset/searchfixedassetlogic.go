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

type SearchFixedAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchFixedAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFixedAssetLogic {
	return &SearchFixedAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type employeeInfo struct {
	EmployeeNo string
	Name       string
}

func (l *SearchFixedAssetLogic) SearchFixedAsset(req *types.SearchFixedAssetReq) (resp *types.SearchFixedAssetResp, err error) {
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.FinanceRPC.SearchFixedAsset(l.ctx, &pb.SearchFixedAssetReq{
		Page:         req.Page,
		Limit:        req.Limit,
		AssetNo:      req.AssetNo,
		AssetName:    req.AssetName,
		Category:     req.Category,
		SupplierId:   supplierId,
		DepartmentId: departmentId,
		Status:       req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 批量获取供应商名称（去重）
	supplierMap := make(map[int64]string)
	for _, fa := range ret.FixedAsset {
		if fa.SupplierId > 0 {
			if _, ok := supplierMap[fa.SupplierId]; !ok {
				supplierResp, err := l.svcCtx.SupplierRPC.GetSupplierById(l.ctx, &supplierpb.GetSupplierByIdReq{Id: fa.SupplierId})
				if err != nil {
					logx.Errorf("查询供应商信息失败: supplierId=%d, err=%v", fa.SupplierId, err)
					supplierMap[fa.SupplierId] = ""
				} else if supplierResp.Supplier != nil {
					supplierMap[fa.SupplierId] = supplierResp.Supplier.Name
				}
			}
		}
	}

	// 批量获取使用人（员工）信息（去重）
	employeeMap := make(map[int64]*employeeInfo)
	for _, fa := range ret.FixedAsset {
		if fa.UserId > 0 {
			if _, ok := employeeMap[fa.UserId]; !ok {
				empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{Id: fa.UserId})
				if err != nil {
					logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", fa.UserId, err)
					employeeMap[fa.UserId] = &employeeInfo{}
				} else if empResp.EmployeeNonSensitiveDetail != nil {
					employeeMap[fa.UserId] = &employeeInfo{
						EmployeeNo: empResp.EmployeeNonSensitiveDetail.EmployeeNo,
						Name:       empResp.EmployeeNonSensitiveDetail.Name,
					}
				}
			}
		}
	}

	// 批量获取使用部门信息（去重）
	departmentMap := make(map[int64]string)
	for _, fa := range ret.FixedAsset {
		if fa.DepartmentId > 0 {
			if _, ok := departmentMap[fa.DepartmentId]; !ok {
				hrResp, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &hrpb.GetDepartmentByIdReq{
					Id: fa.DepartmentId,
				})
				if err != nil {
					logx.Errorf("查询部门失败: departmentId=%d, err=%v", fa.DepartmentId, err)
				} else if hrResp.Department != nil {
					departmentMap[fa.DepartmentId] = hrResp.Department.Name
				}
			}
		}
	}

	list := make([]*types.FixedAsset, 0, len(ret.FixedAsset))
	for _, fa := range ret.FixedAsset {
		var supplierName, userNo, userName, departmentName string
		if info, ok := supplierMap[fa.SupplierId]; ok {
			supplierName = info
		}
		if info, ok := employeeMap[fa.UserId]; ok {
			userNo = info.EmployeeNo
			userName = info.Name
		}
		if info, ok := departmentMap[fa.DepartmentId]; ok {
			departmentName = info
		}
		list = append(list, &types.FixedAsset{
			Id:                 util.Int64ToString(fa.Id),
			AssetNo:            fa.AssetNo,
			AssetName:          fa.AssetName,
			Category:           fa.Category,
			PurchaseDate:       fa.PurchaseDate,
			PurchasePrice:      fa.PurchasePrice,
			SupplierId:         util.Int64ToString(fa.SupplierId),
			SupplierName:       supplierName,
			DepartmentId:       util.Int64ToString(fa.DepartmentId),
			DepartmentName:     departmentName,
			UserId:             util.Int64ToString(fa.UserId),
			UserName:           userName,
			UserNo:             userNo,
			Status:             fa.Status,
			Location:           fa.Location,
			DepreciationMethod: fa.DepreciationMethod,
			UsefulLife:         fa.UsefulLife,
			CreatedAt:          fa.CreatedAt,
		})
	}

	resp = &types.SearchFixedAssetResp{
		FixedAsset: list,
		Total:      ret.Total,
	}
	return
}
