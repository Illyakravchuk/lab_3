package painter

import (
    "image"
    "sync"

    "golang.org/x/exp/shiny/screen"
)

type Receiver interface {
    Update(screen.Texture)
}

type Loop struct {
    Receiver Receiver
    next     screen.Texture
    prev     screen.Texture
    Mq       MessageQueue
    done     chan struct{}
    stopped  bool
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
    l.next, _ = s.NewTexture(size)
    l.prev, _ = s.NewTexture(size)

    l.done = make(chan struct{})

    go func() {
        for !l.stopped || !l.Mq.isEmpty() {
            op := l.Mq.Pull()
            update := op.Do(l.next)
            if update {
                l.Receiver.Update(l.next)
                l.next, l.prev = l.prev, l.next
            }
        }
        close(l.done)
    }()
}

func (l *Loop) Post(op Operation) {
    l.Mq.Push(op)
}

func (l *Loop) StopAndWait() {
    l.Post(OperationFunc(func(t screen.Texture) {
        l.stopped = true
    }))
    l.stopped = true
    <-l.done
}

type MessageQueue struct {
    Operations []Operation
    mu         sync.Mutex
    blocked    chan struct{}
}

func (Mq *MessageQueue) Push(op Operation) {
    Mq.mu.Lock()
    defer Mq.mu.Unlock()

    Mq.Operations = append(Mq.Operations, op)

    if Mq.blocked != nil {
        close(Mq.blocked)
        Mq.blocked = nil
    }
}

func (Mq *MessageQueue) Pull() Operation {
    Mq.mu.Lock()
    defer Mq.mu.Unlock()

    for len(Mq.Operations) == 0 {
        Mq.blocked = make(chan struct{})
        Mq.mu.Unlock()
        <-Mq.blocked
        Mq.mu.Lock()
    }

    op := Mq.Operations[0]
    Mq.Operations[0] = nil
    Mq.Operations = Mq.Operations[1:]
    return op
}

func (Mq *MessageQueue) isEmpty() bool {
    Mq.mu.Lock()
    defer Mq.mu.Unlock()

    return len(Mq.Operations) == 0
}
