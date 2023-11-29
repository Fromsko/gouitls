package design

import (
	"sync"
)

// Observer 定义观察者接口
type Observer interface {
	Update(interface{})
}

// Subject 定义主题接口
type Subject interface {
	RegisterObserver(Observer)
	RemoveObserver(Observer)
	NotifyObservers()
}

// ConcreteSubject 具体主题实现
type ConcreteSubject struct {
	observers sync.Map
}

// RegisterObserver 注册观察者
func (s *ConcreteSubject) RegisterObserver(observer Observer) {
	s.observers.Store(observer, struct{}{})
}

// RemoveObserver 移除观察者
func (s *ConcreteSubject) RemoveObserver(observer Observer) {
	s.observers.Delete(observer)
}

// NotifyObservers 通知观察者
func (s *ConcreteSubject) NotifyObservers() {
	s.observers.Range(func(key, value interface{}) bool {
		observer := key.(Observer)
		observer.Update(s)
		return true
	})
}

/*
// ConcreteObserver 具体观察者实现
type ConcreteObserver struct {
	name string
}

// Update 更新观察者状态
func (o *ConcreteObserver) Update(subject interface{}) {
	fmt.Printf("%s received update from subject\n", o.name)
}

func main() {
	// 创建具体主题
	subject := &ConcreteSubject{}

	// 创建观察者池
	pool, _ := ants.NewPool(10)

	// 向主题注册观察者
	for i := 1; i <= 5; i++ {
		observer := &ConcreteObserver{name: fmt.Sprintf("Observer %d", i)}
		subject.RegisterObserver(observer)

		// 使用观察者池并发执行观察者更新操作
		pool.Submit(func() {
			observer.Update(subject)
		})
	}

	// 通知所有观察者
	subject.NotifyObservers()

	// 关闭观察者池
	pool.Release()
}
*/
