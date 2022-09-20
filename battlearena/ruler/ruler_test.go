package ruler

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/controller"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
)

func TestRuler(t *testing.T) {
	ruler := NewRuler()
	ctrl := controller.New()
	ctrl2 := controller.New()

	res := make(chan message.Message)
	ruler.GetMQ().SendAndForget(message.Create(nil, rulermethods.AddController{
		Controller: ctrl,
	}, nil))
	ruler.GetMQ().Send(message.Create(nil, rulermethods.AddController{
		Controller: ctrl2,
	}, nil), res)

	<-res
	t.Fail()
}
