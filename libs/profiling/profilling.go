package profilling

import (
	"github.com/grafana/pyroscope-go"
	"github.com/saas-flow/shared-libs/config"
)

func NewPyroscope() (*pyroscope.Profiler, error) {
	return pyroscope.Start(pyroscope.Config{
		ApplicationName: config.GetString("SERVICE_NAME"),
		ServerAddress:   config.GetString("PYROSCOPE_SERVER_ADDRESS"),
		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
		Tags: map[string]string{},
	})
}
