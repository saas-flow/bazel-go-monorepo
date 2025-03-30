package pprof

import (
	"context"
	"net/http"

	_ "net/http/pprof"

	"go.uber.org/fx"
)

var Module = fx.Module("pprof", fx.Invoke(Register))

func Register(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) (err error) {
			go http.ListenAndServe(":8080", nil)
			return
		},
	})
}
