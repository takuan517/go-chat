package main

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するチャンネルです。
	forward chan []byte
	//joinはチャットルームに参加しようとしているクライアントのためのチャンネルです。
	join chan *client
	//leaveはチャットルームから退出しようとしているクライアントのためのチャンネルです。
	leave chan *client
	//clientsに在室している全てのクライアントが保持されます。
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <=r.join:
			//参加
			r.clients[client] = true
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			//全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}