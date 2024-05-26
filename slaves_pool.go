package forjislave

type SlavesPool struct {
	handler SlaveHandlerFunc
	slaves  map[string]*Slave
}

func NewSlavesPool() *SlavesPool {
	return &SlavesPool{
		slaves: make(map[string]*Slave),
	}
}

func (sp *SlavesPool) SetHandler(handler SlaveHandlerFunc) {
	sp.handler = handler
	for _, slave := range sp.slaves {
		slave.handler = handler
	}
}

func (sp *SlavesPool) SetActions(actionUrls []string) {
	// Remove slaves that are not in the new list
	for actionUrl, slave := range sp.slaves {
		found := false
		for _, newActionUrl := range actionUrls {
			if actionUrl == newActionUrl {
				found = true
				break
			}
		}
		if !found {
			slave.Stop()
			delete(sp.slaves, actionUrl)
		}
	}

	// Add new slaves
	for _, actionUrl := range actionUrls {
		if _, ok := sp.slaves[actionUrl]; !ok {
			slave := NewSlave()
			slave.Start(actionUrl, sp.handler)
			sp.slaves[actionUrl] = slave
		}
	}
}
