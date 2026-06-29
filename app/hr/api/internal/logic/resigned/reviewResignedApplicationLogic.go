package resigned

import (
	"context"
	pb3 "erp/app/auth/rpc/pb"
	"erp/app/hr/rpc/pb"
	pb2 "erp/app/user/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"fmt"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewResignedApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReviewResignedApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewResignedApplicationLogic {
	return &ReviewResignedApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReviewResignedApplicationLogic) ReviewResignedApplication(req *types.ReviewResignedApplicationRequest) (resp *types.EmptyResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	// 先更新离职申请的审批状态
	_, err = l.svcCtx.HrRPC.ResignedApplicationZrpcClient.UpdateResignedApplication(l.ctx, &pb.UpdateResignedApplicationReq{
		Id:            id,
		Status:        req.Status,
		ApproveTime:   time.Now().Unix(),
		ApproveRemark: req.ApproveRemark,
	})
	if err != nil {
		return nil, err
	}

	// 若未通过审批，直接返回
	if req.Status != 2 {
		return &types.EmptyResponse{}, nil
	}

	// 获取离职申请详情
	ra, err := l.svcCtx.HrRPC.ResignedApplicationZrpcClient.GetResignedApplicationById(l.ctx, &pb.GetResignedApplicationByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	employeeId := ra.ResignedApplication.EmployeeId
	leaveDate := ra.ResignedApplication.LeaveDate

	// 查询用户信息以获取 employeeNo
	u, err := l.svcCtx.UserRPC.GetUserByEmployeeId(l.ctx, &pb2.GetUserByEmployeeIdReq{EmployeeId: employeeId})
	if err != nil {
		return nil, err
	}
	userId := u.User.Id

	// 构建两个步骤的请求体
	hrUpdate := &pb.UpdateEmployeeDetailReq{
		Id:        employeeId,
		LeaveDate: leaveDate,
	}
	userUpdate := &pb2.UpdateUserReq{
		Id:       userId,
		Resigned: true,
	}

	// 目标地址
	hrTarget, err := l.svcCtx.Config.HrRPC.BuildTarget()
	if err != nil {
		return nil, err
	}
	userTarget, err := l.svcCtx.Config.UserRPC.BuildTarget()
	if err != nil {
		return nil, err
	}
	authTarget, err := l.svcCtx.Config.AuthRPC.BuildTarget()
	if err != nil {
		return nil, err
	}

	// 方法全名 URL（仅正向动作）
	hrAction := fmt.Sprintf("%s%s", hrTarget, pb.EmployeeDetail_UpdateEmployeeDetail_FullMethodName)
	userAction := fmt.Sprintf("%s%s", userTarget, pb2.User_UpdateUser_FullMethodName)
	authAction := fmt.Sprintf("%s%s", authTarget, pb3.UserRole_DelUserRoleByUserId_FullMethodName)

	// 使用 DTM 消息型事务（Msg），具备自动重试
	gid := dtmgrpc.MustGenGid(l.svcCtx.Config.DtmServer)
	msg := dtmgrpc.NewMsgGrpc(l.svcCtx.Config.DtmServer, gid).
		Add(hrAction, hrUpdate).
		Add(userAction, userUpdate).
		Add(authAction, &pb3.DelUserRoleByUserIdReq{UserId: userId})

	if err := msg.Submit(); err != nil {

		return nil, err
	}

	l.svcCtx.BizRedis.Del(fmt.Sprintf(xtypes.CacheJWTVersionKey, userId))

	return &types.EmptyResponse{}, nil
}
