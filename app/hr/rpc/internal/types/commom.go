package types

type BulkInsertResult struct {
	Index   int  // data 索引
	Success bool // 是否成功
	Err     error
}

type SearchCom struct {
	Page  int64
	Limit int64
}
