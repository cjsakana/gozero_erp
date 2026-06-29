package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportDepartmentLogic {
	return &ImportDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const maxFileSize = 5 * 1024 * 1024 // 5MB

func (l *ImportDepartmentLogic) ImportDepartment(req *http.Request) (resp *types.ImportDepartmentResponse, err error) {
	//_ = req.ParseMultipartForm(maxFileSize)
	//file, handler, err := req.FormFile("file")
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()
	//
	//// 检查文件类型
	//if !strings.HasSuffix(handler.Filename, ".xlsx") {
	//	return nil, errors.New("只支持.xlsx格式的Excel文件")
	//}
	//
	//// 创建临时文件
	//tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	//if err != nil {
	//	return nil, fmt.Errorf("创建临时文件失败: %v", err)
	//}
	//defer os.Remove(tempFile.Name())
	//defer tempFile.Close()
	//
	//// 将上传内容拷贝到临时文件
	//if _, err := io.Copy(tempFile, file); err != nil {
	//	return nil, fmt.Errorf("保存临时文件失败: %v", err)
	//}
	//
	//options := &excel.ParseOptions{
	//	FieldMappings: []excel.FieldMapping{
	//		{"部门名", "Name"},
	//		{"部门编码", "Code"},
	//		{"负责人工号", "ManagerNo"},
	//		{"负责人姓名", "ManagerName"},
	//	},
	//}
	//// 解析Excel文件
	//departments, err := excel.ParseExcelToStruct(tempFile.Name(), "Sheet1", 1, pb.Department{}, options)
	//if err != nil {
	//	return nil, fmt.Errorf("解析Excel失败: %v", err)
	//}
	//
	//var data []*pb.AddDepartmentItem
	//for _, department := range departments {
	//	v := department.(*pb.Department)
	//	data = append(data, &pb.AddDepartmentItem{
	//		Name:      v.Name,
	//		Code:      v.Code,
	//		ManagerNo: v.ManagerNo,
	//	})
	//}
	//ret, err := l.svcCtx.HrRPC.DepartmentZrpcClient.BulkAddDepartment(l.ctx, &pb.BulkAddDepartmentReq{
	//	Departments: data,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//items := []*types.ImportDepartmentItem{}
	//for _, v := range ret.Items {
	//	items = append(items, &types.ImportDepartmentItem{
	//		Index: v.Index,
	//		Error: v.Error,
	//	})
	//}
	//
	//resp = &types.ImportDepartmentResponse{
	//	SuccessCount: ret.SuccessCount,
	//	FailCount:    ret.ErrorCount,
	//	Items:        items,
	//}
	return
}
