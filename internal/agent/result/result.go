// Пакет result предоставляет структуру для подсчета
// отправленных метрик и ошибок
package result

import "sync/atomic"

// Структура Result хранит в себе сведения об
// отправленных метрик и ошибок
type Result struct {
	all, errors int64
}

// Функция New возвращает экземпляр Result
func New() *Result {
	return &Result{
		all:    0,
		errors: 0,
	}
}

// Метод All возвращает кол-во отправленных запросов
func (r *Result) All() int64 {
	return atomic.LoadInt64(&r.all)
}

// Метод Errors возвращает кол-во ошибок
func (r *Result) Errors() int64 {
	return atomic.LoadInt64(&r.errors)
}

// Метод IncAll инкрементирует кол-во отправленных метрик
func (r *Result) IncAll() {
	atomic.AddInt64(&r.all, 1)
}

// Метод IncErrors инкрементирует кол-во ошибок
func (r *Result) IncErrors() {
	atomic.AddInt64(&r.errors, 1)
}

// Метод Reset сбрасывает все счетчики
func (r *Result) Reset() {
	r.all = 0
	r.errors = 0
}
