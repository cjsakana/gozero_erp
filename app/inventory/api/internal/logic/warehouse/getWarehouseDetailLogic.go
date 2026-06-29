package warehouse

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/user/rpc/user"
	"erp/common/util"

	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWarehouseDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWarehouseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWarehouseDetailLogic {
	return &GetWarehouseDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWarehouseDetailLogic) GetWarehouseDetail(req *types.GetWarehouseDetailReq) (resp *types.GetWarehouseDetailResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	warehouse := types.WarehouseDetail{
		Id:            util.Int64ToString(ret.Warehouse.Id),
		No:            ret.Warehouse.No,
		Name:          ret.Warehouse.Name,
		Location:      ret.Warehouse.Location,
		ManagerId:     util.Int64ToString(ret.Warehouse.ManagerId),
		ManagerNo:     "",
		ManagerName:   "",
		Contact:       "",
		Capacity:      ret.Warehouse.Capacity,
		UsedCapacity:  ret.Warehouse.UsedCapacity,
		IsActive:      ret.Warehouse.IsActive,
		CreatedAt:     ret.Warehouse.CreatedAt,
		CreatedBy:     util.Int64ToString(ret.Warehouse.CreatedBy),
		CreatedByNo:   "",
		CreatedByName: "",
		UpdatedAt:     ret.Warehouse.UpdatedAt,
		UpdatedBy:     util.Int64ToString(ret.Warehouse.UpdatedBy),
		UpdatedByNo:   "",
		UpdatedByName: "",
	}

	// 添加人员信息 enrichment
	// 1. 获取管理员信息
	if ret.Warehouse.ManagerId > 0 {
		managerDetail, err := l.svcCtx.UserRPC.GetUserByEmployeeId(l.ctx, &user.GetUserByEmployeeIdReq{
			EmployeeId: ret.Warehouse.ManagerId,
		})
		if err != nil {
			logx.Errorf("查询管理员信息失败: employeeId=%d, err=%v", ret.Warehouse.ManagerId, err)
		} else {
			warehouse.ManagerNo = managerDetail.User.EmployeeNo
			warehouse.ManagerName = managerDetail.User.RealName
			warehouse.Contact = managerDetail.User.Phone
		}
	}

	// 2. 获取创建人信息
	if ret.Warehouse.CreatedBy > 0 {
		creatorDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
			Id: ret.Warehouse.CreatedBy,
		})
		if err != nil {
			logx.Errorf("查询创建人信息失败: employeeId=%d, err=%v", ret.Warehouse.CreatedBy, err)
		} else {
			warehouse.CreatedByNo = creatorDetail.EmployeeNonSensitiveDetail.EmployeeNo
			warehouse.CreatedByName = creatorDetail.EmployeeNonSensitiveDetail.Name
		}
	}

	// 3. 获取更新人信息
	if ret.Warehouse.UpdatedBy > 0 {
		updaterDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
			Id: ret.Warehouse.UpdatedBy,
		})
		if err != nil {
			logx.Errorf("查询更新人信息失败: employeeId=%d, err=%v", ret.Warehouse.UpdatedBy, err)
		} else {
			warehouse.UpdatedByNo = updaterDetail.EmployeeNonSensitiveDetail.EmployeeNo
			warehouse.UpdatedByName = updaterDetail.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetWarehouseDetailResp{
		Warehouse: warehouse,
	}
	return
}
