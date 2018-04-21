// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package rsort

import "github.com/Rukenshia/graw/reddit"

type messagesThingImpl struct {
	e *reddit.Message
}

func (g messagesThingImpl) Name() string { return g.e.Name }

func (g messagesThingImpl) Birth() uint64 { return g.e.CreatedUTC }

func messagesAsThings(gs []*reddit.Message) []redditThing {
	things := make([]redditThing, len(gs))
	for i, g := range gs {
		things[i] = &messagesThingImpl{g}
	}
	return things
}
