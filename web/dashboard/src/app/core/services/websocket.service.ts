import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

export interface WsMessage {
  type: string;
  data: any;
}

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private socket: WebSocket | null = null;
  private messagesSubject = new Subject<WsMessage>();
  public messages$ = this.messagesSubject.asObservable();

  constructor() {
    this.connect();
  }

  private connect() {
    this.socket = new WebSocket('ws://localhost:8080/ws');

    this.socket.onopen = () => {
      console.log('WebSocket connected');
    };

    this.socket.onmessage = (event) => {
      try {
        const message: WsMessage = JSON.parse(event.data);
        this.messagesSubject.next(message);
      } catch (e) {
        console.error('Failed to parse WebSocket message', e);
      }
    };

    this.socket.onclose = () => {
      console.log('WebSocket disconnected, retrying in 5s...');
      setTimeout(() => this.connect(), 5000);
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error', error);
    };
  }
}
