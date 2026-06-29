package code

import "erp/common/xcode"

var (
    ImageNotFound   = xcode.New(140001, "图片不存在")
    AddImageFail    = xcode.New(140002, "添加图片失败")
    UpdateImageFail = xcode.New(140003, "更新图片失败")
    DeleteImageFail = xcode.New(140004, "删除图片失败")
    GetImageFail    = xcode.New(140005, "获取图片失败")
    SearchImageFail = xcode.New(140006, "搜索图片失败")
    UploadImageFail = xcode.New(140007, "上传图片失败")
)
