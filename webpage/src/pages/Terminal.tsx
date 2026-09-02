import { useEffect, useRef, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';
import { useWebSocket } from '@/hooks/useWebSocket';
import api, { Container } from '@/lib/api';

export default function Terminal() {
  const [searchParams] = useSearchParams();
  const [containers, setContainers] = useState<Container[]>([]);
  const [selectedContainer, setSelectedContainer] = useState<string>('');
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<XTerm | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);

  useEffect(() => {
    loadContainers();
  }, []);

  useEffect(() => {
    const containerParam = searchParams.get('container');
    if (containerParam && containers.length > 0) {
      setSelectedContainer(containerParam);
    } else if (containers.length > 0 && !selectedContainer) {
      const running = containers.find(c => c.status === 'running');
      if (running) setSelectedContainer(running.name);
    }
  }, [searchParams, containers]);

  const loadContainers = async () => {
    try {
      const data: any = await api.get('/containers');
      if (data.success) {
        const runningContainers = (data.data || []).filter((c: Container) => c.status === 'running');
        setContainers(runningContainers);
      }
    } catch (error) {
      console.error('Failed to load containers:', error);
    }
  };

  const { send, resize, isConnected } = useWebSocket({
    containerName: selectedContainer,
    onMessage: (data) => {
      xtermRef.current?.write(data);
    },
    onError: (error) => {
      console.error('WebSocket error:', error);
      xtermRef.current?.write('\r\n\x1b[31m[连接错误]\x1b[0m\r\n');
    },
  });

  useEffect(() => {
    if (!terminalRef.current || !selectedContainer) return;

    // Initialize xterm
    const term = new XTerm({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'JetBrains Mono, monospace',
      theme: {
        background: '#0A0A0F',
        foreground: '#F8FAFC',
        cursor: '#6366F1',
        selectionBackground: '#6366F1',
        black: '#0A0A0F',
        red: '#EF4444',
        green: '#10B981',
        yellow: '#F59E0B',
        blue: '#3B82F6',
        magenta: '#A855F7',
        cyan: '#06B6D4',
        white: '#F8FAFC',
        brightBlack: '#64748B',
        brightRed: '#F87171',
        brightGreen: '#34D399',
        brightYellow: '#FBBF24',
        brightBlue: '#60A5FA',
        brightMagenta: '#C084FC',
        brightCyan: '#22D3EE',
        brightWhite: '#FFFFFF',
      },
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    xtermRef.current = term;
    fitAddonRef.current = fitAddon;

    // Handle terminal input
    term.onData((data) => {
      send(data);
    });

    // Handle resize
    term.onResize(({ cols, rows }) => {
      resize(cols, rows);
    });

    // Fit on window resize
    const handleResize = () => fitAddon.fit();
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      term.dispose();
    };
  }, [selectedContainer, send, resize]);

  useEffect(() => {
    if (isConnected) {
      xtermRef.current?.write('\r\n\x1b[32m[已连接]\x1b[0m\r\n');
    }
  }, [isConnected]);

  if (containers.length === 0) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="text-text-muted">暂无运行中的容器</div>
          <a href="/containers" className="text-accent-indigo hover:underline">
            前往容器管理
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col">
      <div className="flex items-center gap-4 mb-4">
        <h1 className="text-2xl font-bold text-text-primary">Web 终端</h1>
        <select
          value={selectedContainer}
          onChange={(e) => setSelectedContainer(e.target.value)}
          className="px-4 py-2 bg-surface-raised border border-surface-border rounded-lg text-text-primary focus:border-accent-indigo focus:ring-2 focus:ring-accent-indigo/20 outline-none transition-all"
        >
          {containers.map((c) => (
            <option key={c.name} value={c.name}>
              {c.name}
            </option>
          ))}
        </select>
        <div className="flex items-center gap-2 text-sm">
          <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-emerald-500 animate-pulse' : 'bg-red-500'}`} />
          <span className="text-text-secondary">{isConnected ? '已连接' : '未连接'}</span>
        </div>
      </div>

      <div ref={terminalRef} className="flex-1 glass-card p-4 overflow-hidden" />
    </div>
  );
}
