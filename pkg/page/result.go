package page

import (
	"math"
)

// 定义一个泛型的分页结果结构体
type PagedResult[T any] struct {
	Data       []T   // 泛型类型的切片，用于存储实际的数据
	Page       int   // 当前页码
	PageSize   int   // 每页大小
	TotalCount int64 // 数据总条数
	TotalPages int   // 总页数
}

func (p *PagedResult[T]) Offset() int {
	return Offset(p.Page)
}

// 计算总页数
func (p *PagedResult[T]) CalculateTotalPages() {
	if p.PageSize == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = int(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
	}
}

// 获取指定页的数据
func (p *PagedResult[T]) GetDataForPage(page int) []T {
	if page < 1 || page > p.TotalPages {
		return nil
	}

	startIndex := (page - 1) * p.PageSize
	endIndex := startIndex + p.PageSize
	if endIndex > int(p.TotalCount) {
		endIndex = int(p.TotalCount)
	}

	return p.Data[startIndex:endIndex]
}

//
//// 扫描并填充分页结果
//func (p *PagedResult[T]) ScanByPage(offset int, limit int) (err error) {
//	// 调用 ScanByPage 方法获取数据和总数
//	var count int64
//	data, err := IAssistantDo.ScanByPage((*T)(nil), offset, limit)
//	if err != nil {
//		return err
//	}
//
//	p.TotalCount = count
//	p.PageSize = limit
//	p.Page = offset/limit + 1
//	p.CalculateTotalPages()
//
//	// 将数据转换为切片
//	p.Data = data.([]T)
//
//	return nil
//}
