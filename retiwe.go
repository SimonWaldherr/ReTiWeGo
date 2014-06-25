package retiwe

type Connections struct {
	Clients      map[chan string]bool
	AddClient    chan chan string
	RemoveClient chan chan string
	Messages     chan string
}

func (hub *Connections) Init() {
	go func() {
		for {
			select {
			case s := <-hub.AddClient:
				hub.Clients[s] = true
			case s := <-hub.RemoveClient:
				delete(hub.Clients, s)
			case msg := <-hub.Messages:
				for s, _ := range hub.Clients {
					s <- msg
				}
			}
		}
	}()
}
