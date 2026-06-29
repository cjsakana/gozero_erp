package resigned

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetResignedApplicationByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetResignedApplicationByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResignedApplicationByIdLogic {
	return &GetResignedApplicationByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResignedApplicationByIdLogic) GetResignedApplicationById(req *types.GetResignedApplicationByIdRequest) (resp *types.GetResignedApplicationByIdResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	applicationById, err := l.svcCtx.HrRPC.ResignedApplicationZrpcClient.GetResignedApplicationById(l.ctx, &pb.GetResignedApplicationByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// 查询员工信息
	employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: applicationById.ResignedApplication.EmployeeId,
	})
	if err != nil {
		return nil, err
	}

	// 查询审批人信息（如果有）
	var approverName, approverNo string
	if applicationById.ResignedApplication.ApproverId > 0 {
		approverDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
			Id: applicationById.ResignedApplication.ApproverId,
		})
		if err != nil {
			logx.Errorf("查询审批人信息失败: approverId=%d, err=%v", applicationById.ResignedApplication.ApproverId, err)
		} else {
			approverNo = approverDetail.EmployeeNonSensitiveDetail.EmployeeNo
			approverName = approverDetail.EmployeeNonSensitiveDetail.Name
		}
	}

	resp = &types.GetResignedApplicationByIdResponse{
		Application: types.ResignedApplication{
			Id:            util.Int64ToString(applicationById.ResignedApplication.Id),
			EmployeeId:    util.Int64ToString(applicationById.ResignedApplication.EmployeeId),
			EmployeeNo:    employeeDetail.EmployeeNonSensitiveDetail.EmployeeNo,
			EmployeeName:  employeeDetail.EmployeeNonSensitiveDetail.Name,
			Reason:        applicationById.ResignedApplication.Reason,
			LeaveDate:     applicationById.ResignedApplication.LeaveDate,
			Evidence:      applicationById.ResignedApplication.Evidence,
			Status:        applicationById.ResignedApplication.Status,
			ApproverId:    util.Int64ToString(applicationById.ResignedApplication.ApproverId),
			ApproverNo:    approverNo,
			ApproverName:  approverName,
			ApproveTime:   applicationById.ResignedApplication.ApproveTime,
			ApproveRemark: applicationById.ResignedApplication.ApproveRemark,
			CreatedAt:     applicationById.ResignedApplication.CreatedAt,
		},
	}

	return
}
