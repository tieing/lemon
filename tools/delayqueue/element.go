// slot中存储的元素

package delayqueue

// Element 队列元素结构体
type Element struct {
	// 生存周期
	cycleNum int

	// 存储数据
	data interface{}
}

// NewElement 创建新的队列元素
func NewElement(cycleNum int, data interface{}) *Element {
	return &Element{
		cycleNum: cycleNum,
		data:     data,
	}
}

// subCycleNum 生命周期减1
func (e *Element) subCycleNum() {
	if e.cycleNum < 1 {
		return
	}
	e.cycleNum--
}
