package adapters

import (
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// Ensure that all item types implements components.ListItem interface
var _ components.ListItem = (*appliedDirectiveItem)(nil)
var _ components.ListItem = (*argumentItem)(nil)
var _ components.ListItem = (*directiveDefItem)(nil)
var _ components.ListItem = (*fieldItem)(nil)
var _ components.ListItem = (*typeDefItem)(nil)
var _ components.ListItem = (*usageItem)(nil)
