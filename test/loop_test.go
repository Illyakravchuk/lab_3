package test

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/Illyakravchuk/lab_3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

type receiverMock struct {
	mock.Mock
}

func (rm *receiverMock) Update(t screen.Texture) {
	rm.Called(t)
}

type screenMock struct {
	mock.Mock
}

func (sm *screenMock) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (sm *screenMock) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func (sm *screenMock) NewTexture(size image.Point) (screen.Texture, error) {
	args := sm.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

type textureMock struct {
	mock.Mock
}

func (tm *textureMock) Release() {
	tm.Called()
}

func (tm *textureMock) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	tm.Called(dp, src, sr)
}

func (tm *textureMock) Bounds() image.Rectangle {
	args := tm.Called()
	return args.Get(0).(image.Rectangle)
}

func (tm *textureMock) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	tm.Called(dr, src, op)
}

func (tm *textureMock) Size() image.Point {
	args := tm.Called()
	return args.Get(0).(image.Point)
}

type operationMock struct {
	mock.Mock
}

func (om *operationMock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

type operationQueueMock struct{}

func (m *operationQueueMock) Do(t screen.Texture) (ready bool) {
	return false
} 

func TestLoop_Post(t *testing.T) {
	screenMock := new(screenMock)
	textureMock := new(textureMock)
	receiverMock := new(receiverMock)
	tx := image.Pt(800, 800)
	l := painter.Loop{
		Receiver: receiverMock,
	}

	screenMock.On("NewTexture", tx).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()

	l.Start(screenMock)

	op1Mock := new(operationMock)
	op2Mock := new(operationMock)
	op3Mock := new(operationMock)

	textureMock.On("Bounds").Return(image.Rectangle{})
	op1Mock.On("Do", textureMock).Return(false)
	op2Mock.On("Do", textureMock).Return(true)
	op3Mock.On("Do", textureMock).Return(true)

	assert.Empty(t, l.Mq.Operations)
	l.Post(op1Mock)
	l.Post(op2Mock)
	l.Post(op3Mock)
	time.Sleep(1 * time.Second)
	assert.Empty(t, l.Mq.Operations)

	op1Mock.AssertCalled(t, "Do", textureMock)
	op2Mock.AssertCalled(t, "Do", textureMock)
	op3Mock.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(800, 800))
}

func TestMessageQueue_Push(t *testing.T) {
	mq := &painter.MessageQueue{}

	op1Mock := &operationQueueMock{}
	mq.Push(op1Mock)
	if len(mq.Operations) != 1 {
		t.Errorf("Expected 1 operation in the queue, but got %d", len(mq.Operations))
	}

	op2Mock := &operationQueueMock{}
	mq.Push(op2Mock)
	if len(mq.Operations) != 2 {
		t.Errorf("Expected 2 operations in the queue, but got %d", len(mq.Operations))
	}

	op3Mock := &operationQueueMock{}
	mq.Push(op3Mock)
	if len(mq.Operations) != 3 {
		t.Errorf("Expected 3 operations in the queue, but got %d", len(mq.Operations))
	}
}
