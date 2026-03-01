package matreshka_api_impl

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service/v1/subscription"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (s *Impl) SubscribeOnChanges(stream api.MatreshkaApi_SubscribeOnChangesServer) error {
	sub := subscription.NewSubscriber()

	defer s.subService.StopSubscription(sub)
	defer sub.NotifyStopped()

	subscriberEvents := consumeSubscriberStream(stream)

	errorCount := 0

	for {
		select {
		case req, ok := <-subscriberEvents:
			if !ok {
				return nil
			}

			s.subService.Subscribe(sub, req.SubscribeConfigNames...)
			s.subService.Unsubscribe(sub, req.UnsubscribeConfigNames...)

		case updates := <-sub.GetUpdateChan():
			patch := &api.SubscribeOnChanges_Response{
				ConfigName: updates.ConfigName.Name(),
				Timestamp:  uint32(time.Now().UTC().Unix()),
				Patches:    toPatches(updates),
			}

			err := stream.Send(patch)
			if err != nil {
				logrus.Errorf("error sending update to subscriber %s", err)
				errorCount++
				if errorCount >= 3 {
					return rerrors.Wrap(err)
				}
			}
		}
	}
}

func consumeSubscriberStream(stream api.MatreshkaApi_SubscribeOnChangesServer) <-chan *api.SubscribeOnChanges_Request {
	subscriberEvents := make(chan *api.SubscribeOnChanges_Request, 1)

	errorCount := 0

	go func() {
		defer close(subscriberEvents)

		for {
			req, err := stream.Recv()
			if err != nil {
				if !rerrors.Is(err, io.EOF) {
					logrus.Error(rerrors.Wrap(err, "error receiving names from subscription"))
					errorCount++
					if errorCount > 3 {
						return
					}

					continue
				}

				return
			}
			subscriberEvents <- req
		}
	}()

	return subscriberEvents

}

func toPatches(req domain.PatchConfigRequest) []*api.PatchConfig_Patch {
	out := make([]*api.PatchConfig_Patch, 0, len(req.Upsert)+len(req.Delete)+len(req.RenameTo))

	for _, up := range req.Upsert {
		out = append(out, &api.PatchConfig_Patch{
			FieldName: up.FieldName,
			Patch: &api.PatchConfig_Patch_UpdateValue{
				UpdateValue: up.FieldValue,
			},
		})
	}

	for _, fieldName := range req.Delete {
		out = append(out, &api.PatchConfig_Patch{
			FieldName: fieldName,
			Patch: &api.PatchConfig_Patch_Delete{
				Delete: true,
			},
		})
	}

	for _, rp := range req.RenameTo {
		out = append(out, &api.PatchConfig_Patch{
			FieldName: rp.OldName,
			Patch: &api.PatchConfig_Patch_Rename{
				Rename: rp.NewName,
			},
		})
	}

	return out
}
