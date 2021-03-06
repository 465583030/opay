package handles

import (
	"github.com/henrylee2cn/opay"
)

/*
 * 提现
 */
type Withdraw struct {
	Background
}

// 编译期检查接口实现
var _ Handler = (*Withdraw)(nil)

// 执行入口
func (w *Withdraw) ServeOpay(ctx *opay.Context) error {

	if ctx.HasStakeholder() {
		return opay.ErrExtraStakeholder
	}
	if ctx.GreaterOrEqual(ctx.Request.Initiator.GetAmount(), 0) {
		return opay.ErrIncorrectAmount
	}
	return w.Call(w, ctx)
}

// 新建订单，并标记为等待处理状态，
// 先从账户扣除提现金额。
func (w *Withdraw) Pend() error {
	// 操作账户
	err := w.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	// 创建订单
	return w.Background.Context.Pend()
}

// 处理账户并标记订单为成功状态
func (w *Withdraw) Succeed() error {
	return w.Background.Context.Succeed()
}

// 标记订单为撤销状态
func (w *Withdraw) Cancel() error {
	// 回滚账户
	err := w.Background.Context.RollbackBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return w.Background.Context.Cancel()
}

// 标记订单为失败状态
func (w *Withdraw) Fail() error {
	// 回滚账户
	err := w.Background.Context.RollbackBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return w.Background.Context.Fail()
}
