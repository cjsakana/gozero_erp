package code

import "erp/common/xcode"

var (
	// 统一错误码范围：140101 - 140110 为 image API 侧（如需单独暴露 API 错误）
	ImageNotFound   = xcode.New(140101, "图片不存在")
	AddImageFail    = xcode.New(140102, "添加图片失败")
	UpdateImageFail = xcode.New(140103, "更新图片失败")
	DeleteImageFail = xcode.New(140104, "删除图片失败")
	GetImageFail    = xcode.New(140105, "获取图片失败")
	SearchImageFail = xcode.New(140106, "搜索图片失败")
	UploadImageFail = xcode.New(140107, "上传图片失败")
)
