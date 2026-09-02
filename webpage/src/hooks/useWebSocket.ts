import { useEffect, useRef, useState } from 'react';

interface UseWebSocketProps {
  containerName: string;
  onMessage: (data: string) => void;
  onError?: (error: Event) => void;
}

export function useWebSocket({ containerName, onMessage, onError }: UseWebSocketProps) {
  const ws = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    if (!containerName) return;

    const token = localStorage.getItem('token');
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/ws/terminal/${containerName}?token=${token}`;

    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      setIsConnected(true);
    };

    ws.current.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === 'stdout' || msg.type === 'stderr') {
          const decoded = atob(msg.data);
          onMessage(decoded);
        } else if (msg.type === 'error') {
          onMessage(`\r\n[Error] ${msg.message}\r\n`);
        } else if (msg.type === 'exit') {
          onMessage(`\r\n[Process exited with code ${msg.code}]\r\n`);
          setIsConnected(false);
        }
      } catch (e) {
        console.error('Failed to parse WebSocket message:', e);
      }
    };

    ws.current.onerror = (error) => {
      setIsConnected(false);
      onError?.(error);
    };

    ws.current.onclose = () => {
      setIsConnected(false);
    };

    return () => {
      ws.current?.close();
    };
  }, [containerName, onMessage, onError]);

  const send = (data: string) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      const encoded = btoa(data);
      ws.current.send(JSON.stringify({ type: 'stdin', data: encoded }));
    }
  };

  const resize = (cols: number, rows: number) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: 'resize', cols, rows }));
    }
  };

  return { send, resize, isConnected };
}
