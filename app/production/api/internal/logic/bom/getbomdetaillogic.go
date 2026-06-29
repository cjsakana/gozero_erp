package bom

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/app/product/rpc/client/product"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBomDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取BOM详情
func NewGetBomDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBomDetailLogic {
	return &GetBomDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBomDetailLogic) GetBomDetail(req *types.IdReq) (resp *types.BomInfo, err error) {

	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	bomInfo, err := l.svcCtx.ProductionRPC.GetBom(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	productMap := make(map[int64]*product.Product)
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	// 转换明细数据
	items := make([]types.BomItemInfo, 0, len(bomInfo.Items))
	for _, item := range bomInfo.Items {
		if _, ok := productMap[item.MaterialId]; !ok {
			ret4, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
				Id: item.MaterialId,
			})
			if err != nil {
				return nil, err
			}
			productMap[item.MaterialId] = ret4.Product
		}

		items = append(items, types.BomItemInfo{
			Id:           util.Int64ToString(item.Id),
			BomId:        util.Int64ToString(item.BomId),
			MaterialId:   util.Int64ToString(item.MaterialId),
			MaterialNo:   productMap[item.MaterialId].ProductNo,
			MaterialName: item.MaterialName,
			Quantity:     item.Quantity,
			Unit:         item.Unit,
			ScrapRate:    item.ScrapRate,
			Remark:       item.Remark,
		})
	}

	if _, ok := productMap[bomInfo.ProductId]; !ok {
		ret4, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
			Id: bomInfo.ProductId,
		})
		if err != nil {
			return nil, err
		}
		productMap[bomInfo.ProductId] = ret4.Product
	}

	if _, ok := employeeMap[bomInfo.CreatedBy]; !ok {
		employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: bomInfo.CreatedBy,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", bomInfo.CreatedBy, err)
		}
		employeeMap[bomInfo.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
	}

	resp = &types.BomInfo{
		Id:            util.Int64ToString(bomInfo.Id),
		BomNo:         bomInfo.BomNo,
		ProductId:     util.Int64ToString(bomInfo.ProductId),
		ProductNo:     productMap[bomInfo.ProductId].ProductNo,
		ProductName:   bomInfo.ProductName,
		Version:       bomInfo.Version,
		UnitCost:      bomInfo.UnitCost,
		IsActive:      bomInfo.IsActive,
		Remark:        bomInfo.Remark,
		CreatedAt:     bomInfo.CreatedAt,
		CreatedById:   util.Int64ToString(bomInfo.CreatedBy),
		CreatedByNo:   employeeMap[bomInfo.CreatedBy].EmployeeNo,
		CreatedByName: employeeMap[bomInfo.CreatedBy].Name,

		Items: items,
	}

	if bomInfo.UpdatedBy > 0 {
		if _, ok := employeeMap[bomInfo.UpdatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: bomInfo.UpdatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", bomInfo.UpdatedBy, err)
			}
			employeeMap[bomInfo.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		resp.UpdatedAt = bomInfo.UpdatedAt
		resp.UpdatedById = util.Int64ToString(bomInfo.UpdatedBy)
		resp.UpdatedByNo = employeeMap[bomInfo.CreatedBy].EmployeeNo
		resp.UpdatedByName = employeeMap[bomInfo.CreatedBy].Name
	}
	return
}
