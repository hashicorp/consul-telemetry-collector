package confhelper

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

const (
	ballastID = "memory_ballast"
)

// Ballast creates a helper to configure the ballast extension.
// It will allocate a block of the 10% of available memory to attempt to reduce GC.
// Documentation about ballast can be found in this blog post: [Go memory ballast: How I learnt to stop worrying and
// love the heap](https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap/)
func Ballast(c *confresolver.Config) confresolver.ComponentConfig {
	const (
		sizeInPercentage = "size_in_percentage"
	)
	// TODO: Replace ballast with a call to runtime/debug.SetMemoryLimit
	ballast := c.NewExtensions(component.NewID(ballastID))
	ballast.Set(sizeInPercentage, 10)
	return ballast
}
