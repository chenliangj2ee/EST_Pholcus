package ignorecase_test

import (
	"testing"
	"time"

	tp "github.com/henrylee2cn/teleport"
	"github.com/henrylee2cn/teleport/plugin/ignorecase"
)

type Home struct {
	tp.CallCtx
}

func (h *Home) Test(arg *map[string]string) (map[string]interface{}, *tp.Rerror) {
	h.Session().Push("/push/test", map[string]string{
		"your_id": string(h.PeekMeta("peer_id")),
	})
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}

func TestIngoreCase(t *testing.T) {
	// Server
	srv := tp.NewPeer(tp.PeerConfig{ListenPort: 9090}, ignorecase.NewIgnoreCase())
	srv.RouteCall(new(Home))
	go srv.ListenAndServe()
	time.Sleep(1e9)

	// Client
	cli := tp.NewPeer(tp.PeerConfig{}, ignorecase.NewIgnoreCase())
	cli.RoutePush(new(Push))
	sess, err := cli.Dial(":9090")
	if err != nil {
		if err != nil {
			t.Error(err)
		}
	}
	var result interface{}
	rerr := sess.Call("/home/TesT",
		map[string]string{
			"author": "henrylee2cn",
		},
		&result,
		tp.WithAddMeta("peer_id", "110"),
	).Rerror()
	if rerr != nil {
		t.Error(rerr)
	}
	t.Logf("result:%v", result)
	time.Sleep(3e9)
}

type Push struct {
	tp.PushCtx
}

func (p *Push) Test(arg *map[string]string) *tp.Rerror {
	tp.Infof("receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}
