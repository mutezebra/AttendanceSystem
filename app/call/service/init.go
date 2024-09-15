package service

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (svc *CallService) InitCallSystem() {
	for i := 0; i < 1; i++ {
		go svc.initCallSystem()
	}
	svc.readSvcCallEvent()
}

func (svc *CallService) initCallSystem() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	ctx := context.Background()
	p, _ := ants.NewPool(10)
	defer p.Release()

	for {
		select {
		case <-sig:
			close(sig)
			return
		case <-time.After(500 * time.Millisecond):
			svc.mu.Lock()
			events := svc.getTheTimedEvents()
			svc.mu.Unlock()

			for _, v := range events {
				svc.mu.Lock()
				undoUID, err := svc.GetUndoUids(ctx, v.ClassID)
				pack.LogError(err)
				doneUID, err := svc.GetDoneUids(ctx, v.ClassID)
				pack.LogError(err)
				svc.mu.Unlock()
				pack.LogError(svc.cache.DelCallEventSet(ctx, v.ClassID))
				pack.LogError(p.Submit(svc.UpdateEventDone(ctx, v.EventID, v.ClassID, doneUID, undoUID)))
			}
		}
	}
}

func (svc *CallService) readSvcCallEvent() {
	ctx := context.Background()

	datas, _ := svc.cache.ReadSvcCallEvent(ctx)

	for _, data := range datas {
		e := JsonToSvcCallEvent(data)
		svc.unDoneEvents[e.ClassID] = e
	}
}
