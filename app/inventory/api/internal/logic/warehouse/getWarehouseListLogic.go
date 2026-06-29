package warehouse

import (
	"context"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/pb"
	"erp/app/user/rpc/user"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWarehouseListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWarehouseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWarehouseListLogic {
	return &GetWarehouseListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type UserInfo struct {
	EmployeeId int64
	EmployeeNo string
	RealName   string
	Contact    string
}

func (l *GetWarehouseListLogic) GetWarehouseList(req *types.GetWarehouseListReq) (resp *types.GetWarehouseListResp, err error) {
	ret, err := l.svcCtx.InventoryRPC.SearchWarehouse(l.ctx, &pb.SearchWarehouseReq{
		Page:     req.Page,
		Limit:    req.Limit,
		No:       req.Keyword,
		Name:     req.Keyword,
		Location: req.Keyword,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetWarehouseListResp{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*UserInfo)

	for _, one := range ret.Warehouse {
		// 获取管理员信息
		if one.ManagerId > 0 {
			if _, ok := employeeMap[one.ManagerId]; !ok {
				userDetail, err := l.svcCtx.UserRPC.GetUserByEmployeeId(l.ctx, &user.GetUserByEmployeeIdReq{
					EmployeeId: one.ManagerId,
				})
				if err != nil {
					logx.Errorf("查询管理员信息失败: employeeId=%d, err=%v", one.ManagerId, err)
				} else {
					employeeMap[one.ManagerId] = &UserInfo{
						EmployeeId: userDetail.User.EmployeeId,
						EmployeeNo: userDetail.User.EmployeeNo,
						RealName:   userDetail.User.RealName,
						Contact:    userDetail.User.Phone,
					}
				}
			}
		}
		// 获取创建人信息
		if one.CreatedBy > 0 {
			if _, ok := employeeMap[one.CreatedBy]; !ok {
				userDetail, err := l.svcCtx.UserRPC.GetUserByEmployeeId(l.ctx, &user.GetUserByEmployeeIdReq{
					EmployeeId: one.CreatedBy,
				})
				if err != nil {
					logx.Errorf("查询管理员信息失败: employeeId=%d, err=%v", one.ManagerId, err)
				} else {
					employeeMap[one.CreatedBy] = &UserInfo{
						EmployeeId: userDetail.User.EmployeeId,
						EmployeeNo: userDetail.User.EmployeeNo,
						RealName:   userDetail.User.RealName,
						Contact:    userDetail.User.Phone,
					}
				}
			}
		}
		// 获取更新人信息
		if one.UpdatedBy > 0 {
			if _, ok := employeeMap[one.UpdatedBy]; !ok {
				userDetail, err := l.svcCtx.UserRPC.GetUserByEmployeeId(l.ctx, &user.GetUserByEmployeeIdReq{
					EmployeeId: one.UpdatedBy,
				})
				if err != nil {
					logx.Errorf("查询管理员信息失败: employeeId=%d, err=%v", one.ManagerId, err)
				} else {
					employeeMap[one.UpdatedBy] = &UserInfo{
						EmployeeId: userDetail.User.EmployeeId,
						EmployeeNo: userDetail.User.EmployeeNo,
						RealName:   userDetail.User.RealName,
						Contact:    userDetail.User.Phone,
					}
				}
			}
		}

		warehouse := types.WarehouseDetail{
			Id:           util.Int64ToString(one.Id),
			No:           one.No,
			Name:         one.Name,
			Location:     one.Location,
			ManagerId:    util.Int64ToString(one.ManagerId),
			Capacity:     one.Capacity,
			UsedCapacity: one.UsedCapacity,
			IsActive:     one.IsActive,
			CreatedAt:    one.CreatedAt,
			CreatedBy:    util.Int64ToString(one.CreatedBy),
			UpdatedAt:    one.UpdatedAt,
			UpdatedBy:    util.Int64ToString(one.UpdatedBy),
		}

		// 填充人员信息
		if emp, ok := employeeMap[one.ManagerId]; ok && emp != nil {
			warehouse.ManagerNo = emp.EmployeeNo
			warehouse.ManagerName = emp.RealName
			warehouse.Contact = emp.Contact
		}
		if emp, ok := employeeMap[one.CreatedBy]; ok && emp != nil {
			warehouse.CreatedByNo = emp.EmployeeNo
			warehouse.CreatedByName = emp.RealName
		}
		if emp, ok := employeeMap[one.UpdatedBy]; ok && emp != nil {
			warehouse.UpdatedByNo = emp.EmployeeNo
			warehouse.UpdatedByName = emp.RealName
		}

		resp.Warehouses = append(resp.Warehouses, warehouse)
	}
	return
}
