package requisition

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRequisitionDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取领料单详情
func NewGetRequisitionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRequisitionDetailLogic {
	return &GetRequisitionDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRequisitionDetailLogic) GetRequisitionDetail(req *types.IdReq) (resp *types.RequisitionInfo, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	reqInfo, err := l.svcCtx.ProductionRPC.GetRequisition(l.ctx, &production.IdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.RequisitionItemInfo, 0, len(reqInfo.Items))
	for _, item := range reqInfo.Items {
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

	whRet, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
		Id: reqInfo.WarehouseId,
	})
	if err != nil {
		return nil, err
	}
	woInfo, err := l.svcCtx.ProductionRPC.GetWorkOrder(l.ctx, &production.IdReq{
		Id: reqInfo.WorkOrderId,
	})
	if err != nil {
		return nil, err
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	resp = &types.RequisitionInfo{
		Id:              util.Int64ToString(reqInfo.Id),
		RequisitionNo:   reqInfo.RequisitionNo,
		WorkOrderId:     util.Int64ToString(reqInfo.WorkOrderId),
		WorkOrderNo:     reqInfo.WorkOrderNo,
		WarehouseId:     util.Int64ToString(reqInfo.WarehouseId),
		WarehouseNo:     whRet.Warehouse.No,
		WarehouseName:   whRet.Warehouse.Name,
		ProductId:       util.Int64ToString(woInfo.ProductId),
		ProductName:     woInfo.ProductName,
		RequisitionDate: reqInfo.RequisitionDate,
		Status:          reqInfo.Status,
		ApprovedById:    util.Int64ToString(reqInfo.ApprovedBy),
		ApprovedByNo:    "",
		ApprovedByName:  "",
		ApprovedAt:      reqInfo.ApprovedAt,
		CreatedAt:       reqInfo.CreatedAt,
		CreatedById:     util.Int64ToString(reqInfo.CreatedBy),
		CreatedByNo:     "",
		CreatedByName:   "",
		UpdatedAt:       reqInfo.UpdatedAt,
		UpdatedById:     util.Int64ToString(reqInfo.UpdatedBy),
		UpdatedByNo:     "",
		UpdatedByName:   "",
		Items:           items,
	}

	if _, ok := employeeMap[reqInfo.CreatedBy]; !ok {
		employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: reqInfo.CreatedBy,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", reqInfo.CreatedBy, err)
		}
		employeeMap[reqInfo.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
	}
	resp.CreatedByNo = employeeMap[reqInfo.CreatedBy].EmployeeNo
	resp.CreatedByName = employeeMap[reqInfo.CreatedBy].Name

	if reqInfo.UpdatedBy > 0 {
		if _, ok := employeeMap[reqInfo.UpdatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: reqInfo.UpdatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", reqInfo.CreatedBy, err)
			}
			employeeMap[reqInfo.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}
		resp.UpdatedByNo = employeeMap[reqInfo.UpdatedBy].EmployeeNo
		resp.UpdatedByName = employeeMap[reqInfo.UpdatedBy].Name
	}

	if _, ok := employeeMap[reqInfo.ApprovedBy]; !ok {
		employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: reqInfo.ApprovedBy,
		})
		if err != nil {
			logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", reqInfo.ApprovedBy, err)
		}
		employeeMap[reqInfo.ApprovedBy] = employeeDetail.EmployeeNonSensitiveDetail
	}
	resp.ApprovedByNo = employeeMap[reqInfo.ApprovedBy].EmployeeNo
	resp.ApprovedByName = employeeMap[reqInfo.ApprovedBy].Name
	return
}
