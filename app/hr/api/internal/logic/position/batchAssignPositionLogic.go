package position

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchAssignPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const maxFileSize = 5 * 1024 * 1024 // 5MB

func NewBatchAssignPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchAssignPositionLogic {
	return &BatchAssignPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchAssignPositionLogic) BatchAssignPosition(req *types.BatchAssignPositionRequest) (resp *types.BatchAssignPositionResponse, err error) {
	positionId, err := util.StringToInt64(req.PositionId)
	if err != nil {
		return nil, err
	}
	
	var successCount, failCount int64
	items := []*types.BatchAssignPositionErrItem{}
	for _, employeeIdStr := range req.EmployeeIds {
		employeeId, err := util.StringToInt64(employeeIdStr)
		if err != nil {
			failCount++
			items = append(items, &types.BatchAssignPositionErrItem{
				EmployeeId: employeeIdStr,
				Error:      err.Error(),
			})
			continue
		}
		
		_, err = l.svcCtx.HrRPC.EmployeeDetailZrpcClient.UpdateEmployeeDetail(l.ctx, &pb.UpdateEmployeeDetailReq{
			Id:         employeeId,
			PositionId: positionId,
		})
		if err != nil {
			failCount++
			items = append(items, &types.BatchAssignPositionErrItem{
				EmployeeId: employeeIdStr,
				Error:      err.Error(),
			})
		} else {
			successCount++
		}
	}
	resp = &types.BatchAssignPositionResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		Items:        items,
	}

	return
}
