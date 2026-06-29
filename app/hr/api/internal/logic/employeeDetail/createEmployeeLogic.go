package employeeDetail

import (
	"context"
	pb3 "erp/app/auth/rpc/pb"
	"erp/app/hr/rpc/pb"
	pb2 "erp/app/user/rpc/pb"
	"erp/app/user/rpc/user"
	"erp/common/util"
	"erp/common/xtypes"
	"fmt"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmployeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEmployeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmployeeLogic {
	return &CreateEmployeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEmployeeLogic) CreateEmployee(req *types.CreateEmployeeRequest) (resp *types.CreateEmployeeResponse, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	positionId, err := util.StringToInt64(req.PositionId)
	if err != nil {
		return nil, err
	}

	dateStr := time.Now().Format("20060102")
	employeeNo := fmt.Sprintf("E%s%s", dateStr, util.RandomNumeric(4))
	password := req.IdCard[len(req.IdCard)-8:]

	// 目标地址（go-zero driver-gozero 支持 service 发现）
	hrTarget, err := l.svcCtx.Config.HrRPC.BuildTarget()
	if err != nil {
		return nil, err
	}
	userTarget, err := l.svcCtx.Config.UserRPC.BuildTarget()
	if err != nil {
		return nil, err
	}

	employeeId := util.GenerateSnowflake()

	// 构建请求体
	hrAdd := &pb.AddEmployeeDetailReq{
		EmployeeDetail: &pb.AddEmployeeDetailItem{
			Id:           employeeId,
			EmployeeNo:   employeeNo,
			Name:         req.Name,
			IdCard:       req.IdCard,
			Account:      req.Account,
			Gender:       req.Gender,
			DepartmentId: departmentId,
			PositionId:   positionId,
			Salary:       req.Salary,
			HireDate:     req.HireDate,
		},
	}
	userAdd := &pb2.AddUserReq{
		EmployeeId: employeeId,
		EmployeeNo: employeeNo,
		Username:   req.Name,
		RealName:   req.Name,
		Password:   password,
		Phone:      req.Phone,
		Email:      "",
		CreatedBy:  createdBy,
	}

	// 业务 gRPC URL（目标 + 方法全名）
	hrAddUrl := fmt.Sprintf("%s%s", hrTarget, pb.EmployeeDetail_AddEmployeeDetail_FullMethodName)
	userAddUrl := fmt.Sprintf("%s%s", userTarget, pb2.User_AddUser_FullMethodName)

	// 创建 Msg 并提交（使用 DTM 消息型事务，具备自动重试）
	gid := dtmgrpc.MustGenGid(l.svcCtx.Config.DtmServer)
	msg := dtmgrpc.NewMsgGrpc(l.svcCtx.Config.DtmServer, gid).
		Add(hrAddUrl, hrAdd).
		Add(userAddUrl, userAdd)
	if err := msg.Submit(); err != nil {
		return nil, err
	}

	// 在事务提交后，启动一个带 context 控制的 goroutine
	go func() {
		// 1. 创建一个带超时的 context，比如最多等 30 秒
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel() // 确保资源释放

		var userRet *user.GetUserByEmployeeIdResp
		ticker := time.NewTicker(2 * time.Second) // 每 2 秒查一次
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// 超时或被取消
				l.Logger.Errorf("wait GetUserByEmployeeId timeout or cancelled, employeeId=%v, err=%v", employeeId, ctx.Err())
				return
			case <-ticker.C:
				userRet, err = l.svcCtx.UserRPC.GetUserByEmployeeId(ctx,
					&pb2.GetUserByEmployeeIdReq{EmployeeId: employeeId})
				if err == nil && userRet != nil && userRet.User != nil {
					// 查到了，跳出循环
					goto done
				}
				// 如果是 NotFound 错误，继续等；其他错误可以记录日志
				if err != nil {
					l.Logger.Errorf("GetUserByEmployeeId failed, employeeId=%v, err=%v", employeeId, err)
				}
			}
		}

	done:
		// 2. 查到用户后，再调用 AddUserRole
		_, err = l.svcCtx.AuthRPC.AddUserRole(ctx, &pb3.AddUserRoleReq{
			UserId: userRet.User.Id,
			RoleId: 2,
		})
		if err != nil {
			l.Logger.Errorf("AddUserRole failed, userId=%v, err=%v", userRet.User.Id, err)
		}
	}()

	resp = &types.CreateEmployeeResponse{
		EmployeeId: util.Int64ToString(employeeId),
	}
	return
}
