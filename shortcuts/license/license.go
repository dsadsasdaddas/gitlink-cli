package license

import (
	"net/url"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

// Shortcuts returns license management shortcuts.
//
// These shortcuts provide access to the GitLink license registry,
// which lists all available open-source licenses supported by the platform.
func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: "List available licenses",
			Flags: []common.Flag{
				{Name: "name", Short: "n", Usage: "Filter licenses by name"},
			},
			Run: runList,
		},
	}
}

func runList(ctx *common.RuntimeContext) error {
	q := url.Values{}
	if name := ctx.Arg("name"); name != "" {
		q.Set("name", name)
	}
	env, err := ctx.CallAPIWithQuery("GET", "/licenses", q)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}
