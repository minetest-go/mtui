package mail

func (p *PlayerEntry) FindMessage(id string) *Message {
	list := append(append(p.Inbox, p.Outbox...), p.Drafts...)
	for _, m := range list {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func (p *PlayerEntry) RemoveMessage(id string) {
	p.Inbox = filter(p.Inbox, func(m *Message) bool {
		return m.ID != id
	})
	p.Outbox = filter(p.Outbox, func(m *Message) bool {
		return m.ID != id
	})
	p.Drafts = filter(p.Drafts, func(m *Message) bool {
		return m.ID != id
	})
}
