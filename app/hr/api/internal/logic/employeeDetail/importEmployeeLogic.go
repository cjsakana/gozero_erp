package employeeDetail

import (
	"context"
	"erp/app/hr/rpc/pb"
	pb2 "erp/app/user/rpc/pb"
	"erp/common/excel"
	"erp/common/util"
	"erp/common/xtypes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportEmployeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportEmployeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportEmployeeLogic {
	return &ImportEmployeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const maxFileSize = 5 * 1024 * 1024 // 5MB
func (l *ImportEmployeeLogic) ImportEmployee(req *http.Request) (resp *types.ImportEmployeeResponse, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	_ = req.ParseMultipartForm(maxFileSize)
	file, handler, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 检查文件类型
	if !strings.HasSuffix(handler.Filename, ".xlsx") {
		return nil, errors.New("只支持.xlsx格式的Excel文件")
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 将上传内容拷贝到临时文件
	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, fmt.Errorf("保存临时文件失败: %v", err)
	}

	// 解析Excel文件
	type employee struct {
		Name           string  `excel:"姓名"`
		Phone          string  `excel:"手机号"`
		IdCard         string  `excel:"身份证"`
		Account        string  `excel:"银行卡"`
		Gender         string  `excel:"性别"`
		DepartmentName string  `excel:"部门"`
		PositionName   string  `excel:"岗位"`
		Salary         float64 `excel:"基本工资"`
		HireDate       int64   `excel:"入职日期（年/月/日）"`
	}
	employees, err := excel.ParseExcelToStruct(tempFile.Name(), "Sheet1", 1, employee{})
	if err != nil {
		return nil, fmt.Errorf("解析Excel失败: %v", err)
	}

	deptRet, err := l.svcCtx.HrRPC.DepartmentZrpcClient.SearchDepartment(l.ctx, &pb.SearchDepartmentReq{
		Limit: -1,
	})
	if err != nil {
		return nil, err
	}
	positionRet, err := l.svcCtx.HrRPC.PositionZrpcClient.SearchPosition(l.ctx, &pb.SearchPositionReq{
		Limit: -1,
	})
	if err != nil {
		return nil, err
	}

	// 基于名称构建部门与岗位的映射，便于通过名称找到对应的ID
	deptNameToID := make(map[string]int64)
	for _, d := range deptRet.Department {
		deptNameToID[d.Name] = d.Id
	}
	positionNameToID := make(map[string]int64)
	for _, p := range positionRet.Position {
		positionNameToID[p.Name] = p.Id
	}

	dateStr := time.Now().Format("20060102")

	var pbEs []*pb.AddEmployeeDetailItem
	var pbUs []*pb2.BulkInsertUser
	for _, e := range employees {
		v := e.(employee)
		employeeNo := fmt.Sprintf("E%s%s", dateStr, util.RandomNumeric(4))
		password := v.Account[len(v.Account)-8:]
		gener := 0
		if v.Gender == "男" {
			gener = 1
		} else if v.Gender == "女" {
			gener = 2
		}
		// 通过名称映射得到部门与岗位ID
		deptID := deptNameToID[v.DepartmentName]
		posID := positionNameToID[v.PositionName]

		employeeId := util.GenerateSnowflake()

		pbEs = append(pbEs, &pb.AddEmployeeDetailItem{
			Id:           employeeId,
			EmployeeNo:   employeeNo,
			Name:         v.Name,
			IdCard:       v.IdCard,
			Account:      v.Account,
			Gender:       int64(gener),
			DepartmentId: deptID,
			PositionId:   posID,
			Salary:       v.Salary,
			HireDate:     v.HireDate,
		})
		pbUs = append(pbUs, &pb2.BulkInsertUser{
			EmployeeId: employeeId,
			EmployeeNo: employeeNo,
			Username:   v.Name,
			RealName:   v.Name,
			Password:   password,
			Phone:      v.Phone,
			Email:      "",
			CreatedBy:  createdBy,
		})
	}

	ret, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.BulkAddEmployeeDetail(l.ctx, &pb.BulkAddEmployeeDetailReq{
		EmployeeDetails: pbEs,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.UserRPC.BulkInsertUser(l.ctx, &pb2.BulkInsertUserReq{
		Users: pbUs,
	})
	if err != nil {
		return nil, err
	}

	items := []*types.ImportEmployeeItem{}
	for _, v := range ret.Items {
		items = append(items, &types.ImportEmployeeItem{
			Index: v.Index + 1,
			Error: v.Error,
		})
	}
	resp = &types.ImportEmployeeResponse{
		SuccessCount: ret.SuccessCount,
		FailCount:    ret.ErrorCount,
		Items:        items,
	}
	return
}
