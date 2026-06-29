package requisition

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRequisitionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取领料单列表
func NewGetRequisitionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequisitionListLogic {
	return &GetRequisitionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRequisitionListLogic) GetRequisitionList(req *types.RequisitionListReq) (resp *types.RequisitionListResp, err error) {
	workOrderId, err := util.StringToInt64(req.WorkOrderId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	listResp, err := l.svcCtx.ProductionRPC.GetRequisitionList(l.ctx, &production.RequisitionListReq{
		Page:        req.Page,
		PageSize:    req.PageSize,
		WorkOrderId: workOrderId,
		WarehouseId: warehouseId,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	warehouseMap := make(map[int64]*inventory.WarehouseDetail)
	workOrderMap := make(map[int64]*production.WorkOrderInfo)
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	list := make([]types.RequisitionInfo, 0, len(listResp.List))
	for _, reqItem := range listResp.List {
		items := make([]types.RequisitionItemInfo, 0, len(reqItem.Items))
		for _, item := range reqItem.Items {
			items = append(items, types.RequisitionItemInfo{
				Id:             util.Int64ToString(item.Id),
				RequisitionId:  util.Int64ToString(item.RequisitionId),
				MaterialId:     util.Int64ToString(item.MaterialId),
				MaterialName:   item.MaterialName,
				PlanQuantity:   item.PlanQuantity,
				ActualQuantity: item.ActualQuantity,
				Unit:           item.Unit,
				BatchNo:        item.BatchNo,
				Remark:         item.Remark,
			})
		}

		if _, ok := warehouseMap[reqItem.WarehouseId]; !ok {
			whRet, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
				Id: reqItem.WarehouseId,
			})
			if err != nil {
				return nil, err
			}
			warehouseMap[reqItem.WarehouseId] = whRet.Warehouse
		}

		if _, ok := workOrderMap[reqItem.WorkOrderId]; !ok {
			woInfo, err := l.svcCtx.ProductionRPC.GetWorkOrder(l.ctx, &production.IdReq{
				Id: reqItem.WorkOrderId,
			})
			if err != nil {
				return nil, err
			}
			workOrderMap[reqItem.WorkOrderId] = woInfo
		}

		requisitionInfo := types.RequisitionInfo{
			Id:              util.Int64ToString(reqItem.Id),
			RequisitionNo:   reqItem.RequisitionNo,
			WorkOrderId:     util.Int64ToString(reqItem.WorkOrderId),
			WorkOrderNo:     reqItem.WorkOrderNo,
			WarehouseId:     util.Int64ToString(reqItem.WarehouseId),
			WarehouseNo:     warehouseMap[reqItem.WarehouseId].No,
			WarehouseName:   warehouseMap[reqItem.WarehouseId].Name,
			ProductId:       util.Int64ToString(workOrderMap[reqItem.WorkOrderId].ProductId),
			ProductName:     workOrderMap[reqItem.WorkOrderId].ProductName,
			RequisitionDate: reqItem.RequisitionDate,
			Status:          reqItem.Status,
			ApprovedById:    util.Int64ToString(reqItem.ApprovedBy),
			ApprovedByNo:    "",
			ApprovedByName:  "",
			ApprovedAt:      reqItem.ApprovedAt,
			CreatedAt:       reqItem.CreatedAt,
			CreatedById:     util.Int64ToString(reqItem.CreatedBy),
			CreatedByNo:     "",
			CreatedByName:   "",
			UpdatedAt:       reqItem.UpdatedBy,
			UpdatedById:     util.Int64ToString(reqItem.UpdatedBy),
			UpdatedByNo:     "",
			UpdatedByName:   "",
			Items:           items,
		}

		if _, ok := employeeMap[reqItem.CreatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: reqItem.CreatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", reqItem.CreatedBy, err)
			}
			employeeMap[reqItem.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		requisitionInfo.CreatedByNo = employeeMap[reqItem.CreatedBy].EmployeeNo
		requisitionInfo.CreatedByName = employeeMap[reqItem.CreatedBy].Name

		if reqItem.UpdatedBy > 0 {
			if _, ok := employeeMap[reqItem.UpdatedBy]; !ok {
				employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
					Id: reqItem.UpdatedBy,
				})
				if err != nil {
					logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", reqItem.CreatedBy, err)
				}
				employeeMap[reqItem.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
			}
			requisitionInfo.UpdatedByNo = employeeMap[reqItem.UpdatedBy].EmployeeNo
			requisitionInfo.UpdatedByName = employeeMap[reqItem.UpdatedBy].Name
		}

		if _, ok := employeeMap[reqItem.ApprovedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: reqItem.ApprovedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", reqItem.ApprovedBy, err)
			}
			employeeMap[reqItem.ApprovedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		requisitionInfo.ApprovedByNo = employeeMap[reqItem.ApprovedBy].EmployeeNo
		requisitionInfo.ApprovedByName = employeeMap[reqItem.ApprovedBy].Name

		list = append(list, requisitionInfo)
	}

	resp = &types.RequisitionListResp{
		Total: listResp.Total,
		List:  list,
	}
	return
}
