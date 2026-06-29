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

type GetBomListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取BOM列表
func NewGetBomListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBomListLogic {
	return &GetBomListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBomListLogic) GetBomList(req *types.BomListReq) (resp *types.BomListResp, err error) {

	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}

	listResp, err := l.svcCtx.ProductionRPC.GetBomList(l.ctx, &production.BomListReq{
		Page:      req.Page,
		PageSize:  req.PageSize,
		ProductId: productId,
		IsActive:  req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	productMap := make(map[int64]*product.Product)
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	// 转换列表数据
	list := make([]types.BomInfo, 0, len(listResp.List))
	for _, bom := range listResp.List {
		items := make([]types.BomItemInfo, 0, len(bom.Items))
		for _, item := range bom.Items {
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

		if _, ok := productMap[bom.ProductId]; !ok {
			ret4, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
				Id: bom.ProductId,
			})
			if err != nil {
				return nil, err
			}
			productMap[bom.ProductId] = ret4.Product
		}

		if _, ok := employeeMap[bom.CreatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: bom.CreatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", bom.CreatedBy, err)
			}
			employeeMap[bom.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}

		bomInfo := types.BomInfo{
			Id:            util.Int64ToString(bom.Id),
			BomNo:         bom.BomNo,
			ProductId:     util.Int64ToString(bom.ProductId),
			ProductNo:     productMap[bom.ProductId].ProductNo,
			ProductName:   bom.ProductName,
			Version:       bom.Version,
			UnitCost:      bom.UnitCost,
			IsActive:      bom.IsActive,
			Remark:        bom.Remark,
			CreatedAt:     bom.CreatedAt,
			CreatedById:   util.Int64ToString(bom.CreatedBy),
			CreatedByNo:   employeeMap[bom.CreatedBy].EmployeeNo,
			CreatedByName: employeeMap[bom.CreatedBy].Name,

			Items: items,
		}
		if bom.UpdatedBy > 0 {
			if _, ok := employeeMap[bom.UpdatedBy]; !ok {
				employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
					Id: bom.UpdatedBy,
				})
				if err != nil {
					logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", bom.UpdatedBy, err)
				}
				employeeMap[bom.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
			}
			bomInfo.UpdatedAt = bom.UpdatedAt
			bomInfo.UpdatedById = util.Int64ToString(bom.UpdatedBy)
			bomInfo.UpdatedByNo = employeeMap[bom.CreatedBy].EmployeeNo
			bomInfo.UpdatedByName = employeeMap[bom.CreatedBy].Name
		}
		list = append(list, bomInfo)
	}

	resp = &types.BomListResp{
		Total: listResp.Total,
		List:  list,
	}
	return
}
