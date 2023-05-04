package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправленя останнього разу у Receiver

	mq messageQueue
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.mq = messageQueue{}
	go l.processEvents()

	// TODO: ініціалізувати чергу подій.
	// TODO: запустити рутину обробки повідомлень у черзі подій.
}

func (l *Loop) processEvents() {
	for {
		op := l.mq.pull()
		l.Post(op)
	}
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
    if op == nil {
        return
    }

    l.mq.push(op)

    update := op.Do(l.next)
    if update {
        l.Receiver.Update(l.next)
        l.next, l.prev = l.prev, l.next
    }
}

// StopAndWait сигналізує
func (l *Loop) StopAndWait() {

}

// TODO: реалізувати власну чергу повідомлень.
type messageQueue struct {
	queue []Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.queue = append(mq.queue, op)
}

func (mq *messageQueue) pull() Operation {
	if len(mq.queue) == 0 {
		return nil
	}

	op := mq.queue[0]
	mq.queue = mq.queue[1:]

	return op
}

