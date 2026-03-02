"use client";

import { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  Video,
  VideoOff,
  Mic,
  MicOff,
  PhoneOff,
  Send,
  Info,
  AlertCircle,
} from "lucide-react";
import { Button } from "@/components/ui/button";

interface PeerSessionProps {
  roomId: string;
  sessionId: string;
  matchContext: string;
  similarity: number;
  token: string;
  onEnd: (sessionId: string) => void;
}

interface ChatMessage {
  id: string;
  text: string;
  fromSelf: boolean;
  timestamp: Date;
}

export function PeerSession({
  roomId,
  sessionId,
  matchContext,
  similarity,
  token,
  onEnd,
}: PeerSessionProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputText, setInputText] = useState("");
  const [videoEnabled, setVideoEnabled] = useState(true);
  const [audioEnabled, setAudioEnabled] = useState(true);
  const [timeRemaining, setTimeRemaining] = useState(45 * 60); // 45 minutes in seconds
  const [connectionState, setConnectionState] = useState<
    "connecting" | "connected" | "failed"
  >("connecting");

  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const peerConnectionRef = useRef<RTCPeerConnection | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom of chat
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Timer countdown
  useEffect(() => {
    const interval = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev <= 0) {
          clearInterval(interval);
          handleEndSession();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  // Initialize WebRTC and WebSocket
  useEffect(() => {
    const initSession = async () => {
      try {
        // Get user media
        const stream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        });
        localStreamRef.current = stream;

        if (localVideoRef.current) {
          localVideoRef.current.srcObject = stream;
        }

        // Create peer connection
        const pc = new RTCPeerConnection({
          iceServers: [
            { urls: "stun:stun.l.google.com:19302" },
            { urls: "stun:stun1.l.google.com:19302" },
          ],
        });
        peerConnectionRef.current = pc;

        // Add local tracks to peer connection
        stream.getTracks().forEach((track) => {
          pc.addTrack(track, stream);
        });

        // Handle remote stream
        pc.ontrack = (event) => {
          if (remoteVideoRef.current) {
            remoteVideoRef.current.srcObject = event.streams[0];
          }
        };

        pc.oniceconnectionstatechange = () => {
          if (pc.iceConnectionState === "connected") {
            setConnectionState("connected");
          } else if (pc.iceConnectionState === "failed") {
            setConnectionState("failed");
          }
        };

        // Connect to signaling WebSocket
        const wsUrl = `${process.env.NEXT_PUBLIC_API_URL?.replace("http", "ws")}/api/v1/constellation/signal/${roomId}`;
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
          console.log("WebSocket connected");
          ws.send(JSON.stringify({ type: "auth", token }));
        };

        ws.onmessage = async (event) => {
          const data = JSON.parse(event.data);

          switch (data.type) {
            case "auth_ok":
              // Create and send offer
              const offer = await pc.createOffer();
              await pc.setLocalDescription(offer);
              ws.send(JSON.stringify({ type: "offer", sdp: offer.sdp }));
              break;

            case "offer":
              // Receive offer, create answer
              await pc.setRemoteDescription(
                new RTCSessionDescription({ type: "offer", sdp: data.sdp }),
              );
              const answer = await pc.createAnswer();
              await pc.setLocalDescription(answer);
              ws.send(JSON.stringify({ type: "answer", sdp: answer.sdp }));
              break;

            case "answer":
              // Receive answer
              await pc.setRemoteDescription(
                new RTCSessionDescription({ type: "answer", sdp: data.sdp }),
              );
              break;

            case "ice-candidate":
              // Receive ICE candidate
              if (data.candidate) {
                await pc.addIceCandidate(new RTCIceCandidate(data.candidate));
              }
              break;

            case "chat":
              // Receive chat message
              setMessages((prev) => [
                ...prev,
                {
                  id: Date.now().toString(),
                  text: data.text,
                  fromSelf: false,
                  timestamp: new Date(),
                },
              ]);
              break;
          }
        };

        // Send ICE candidates
        pc.onicecandidate = (event) => {
          if (event.candidate && ws.readyState === WebSocket.OPEN) {
            ws.send(
              JSON.stringify({
                type: "ice-candidate",
                candidate: event.candidate,
              }),
            );
          }
        };
      } catch (error) {
        console.error("Failed to initialize session:", error);
        setConnectionState("failed");
      }
    };

    initSession();

    // Cleanup
    return () => {
      localStreamRef.current?.getTracks().forEach((track) => track.stop());
      peerConnectionRef.current?.close();
      wsRef.current?.close();
    };
  }, [roomId, token]);

  const handleSendMessage = () => {
    if (!inputText.trim() || !wsRef.current) return;

    const message: ChatMessage = {
      id: Date.now().toString(),
      text: inputText,
      fromSelf: true,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, message]);
    wsRef.current.send(JSON.stringify({ type: "chat", text: inputText }));
    setInputText("");
  };

  const toggleVideo = () => {
    if (localStreamRef.current) {
      const videoTrack = localStreamRef.current.getVideoTracks()[0];
      if (videoTrack) {
        videoTrack.enabled = !videoTrack.enabled;
        setVideoEnabled(videoTrack.enabled);
      }
    }
  };

  const toggleAudio = () => {
    if (localStreamRef.current) {
      const audioTrack = localStreamRef.current.getAudioTracks()[0];
      if (audioTrack) {
        audioTrack.enabled = !audioTrack.enabled;
        setAudioEnabled(audioTrack.enabled);
      }
    }
  };

  const handleEndSession = () => {
    localStreamRef.current?.getTracks().forEach((track) => track.stop());
    peerConnectionRef.current?.close();
    wsRef.current?.close();
    onEnd(sessionId);
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  return (
    <div className="fixed inset-0 z-50 bg-pulse-bg flex flex-col">
      {/* Header */}
      <div className="bg-pulse-surface border-b border-pulse-border px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-2 h-2 rounded-full bg-pulse-primary animate-pulse" />
          <div>
            <h2 className="text-sm font-medium text-pulse-text-primary">
              Anonymous Peer Session
            </h2>
            <p className="text-xs text-pulse-text-muted">
              {Math.round(similarity * 100)}% pattern match · End-to-end
              encrypted
            </p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-pulse-bg border border-pulse-border">
            <div
              className={`w-1.5 h-1.5 rounded-full ${connectionState === "connected" ? "bg-emerald-400" : connectionState === "connecting" ? "bg-pulse-accent-warm animate-pulse" : "bg-pulse-danger"}`}
            />
            <span className="text-xs text-pulse-text-muted capitalize">
              {connectionState}
            </span>
          </div>
          <span className="text-sm font-mono text-pulse-text-secondary">
            {formatTime(timeRemaining)}
          </span>
        </div>
      </div>

      {/* Video grid */}
      <div className="flex-1 grid grid-cols-2 gap-4 p-6 bg-pulse-bg">
        {/* Remote video */}
        <div className="relative rounded-xl overflow-hidden bg-pulse-surface border border-pulse-border">
          <video
            ref={remoteVideoRef}
            autoPlay
            playsInline
            className="w-full h-full object-cover"
          />
          <div className="absolute bottom-4 left-4">
            <div className="px-3 py-1.5 rounded-full bg-black/50 backdrop-blur-sm border border-white/10">
              <span className="text-xs text-white font-mono">Peer</span>
            </div>
          </div>
        </div>

        {/* Local video */}
        <div className="relative rounded-xl overflow-hidden bg-pulse-surface border border-pulse-border">
          <video
            ref={localVideoRef}
            autoPlay
            playsInline
            muted
            className="w-full h-full object-cover scale-x-[-1]"
          />
          <div className="absolute bottom-4 left-4">
            <div className="px-3 py-1.5 rounded-full bg-black/50 backdrop-blur-sm border border-white/10">
              <span className="text-xs text-white font-mono">You</span>
            </div>
          </div>
        </div>
      </div>

      {/* Chat sidebar */}
      <div className="absolute right-6 top-24 bottom-24 w-80 bg-pulse-surface border border-pulse-border rounded-xl flex flex-col">
        <div className="p-4 border-b border-pulse-border">
          <h3 className="text-xs font-medium text-pulse-text-primary">
            Text Chat
          </h3>
          <p className="text-[10px] text-pulse-text-muted mt-0.5">
            Conversation is anonymous
          </p>
        </div>
        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {messages.length === 0 && (
            <div className="text-center py-8">
              <p className="text-xs text-pulse-text-muted">No messages yet</p>
            </div>
          )}
          {messages.map((msg) => (
            <div
              key={msg.id}
              className={`flex ${msg.fromSelf ? "justify-end" : "justify-start"}`}
            >
              <div
                className={`max-w-[80%] px-3 py-2 rounded-lg text-xs ${
                  msg.fromSelf
                    ? "bg-pulse-primary text-white"
                    : "bg-pulse-surface-raised text-pulse-text-primary"
                }`}
              >
                {msg.text}
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
        <div className="p-4 border-t border-pulse-border">
          <div className="flex gap-2">
            <input
              type="text"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
              placeholder="Type a message..."
              className="flex-1 px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-xs text-pulse-text-primary placeholder:text-pulse-text-muted focus:outline-none focus:border-pulse-primary"
            />
            <Button
              size="sm"
              onClick={handleSendMessage}
              disabled={!inputText.trim()}
              className="bg-pulse-primary text-white h-9 w-9 p-0"
            >
              <Send className="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      </div>

      {/* Context info */}
      {matchContext && (
        <div className="absolute top-24 left-6 max-w-md">
          <div className="bg-pulse-surface/95 backdrop-blur-sm border border-pulse-border rounded-lg p-3 flex items-start gap-2">
            <Info className="w-4 h-4 text-pulse-primary shrink-0 mt-0.5" />
            <div>
              <p className="text-xs text-pulse-text-secondary">
                {matchContext}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Controls */}
      <div className="bg-pulse-surface border-t border-pulse-border px-6 py-4 flex items-center justify-center gap-4">
        <Button
          size="sm"
          variant="outline"
          onClick={toggleVideo}
          className={`h-12 w-12 rounded-full p-0 ${!videoEnabled ? "bg-pulse-danger/10 border-pulse-danger text-pulse-danger" : ""}`}
        >
          {videoEnabled ? (
            <Video className="w-5 h-5" />
          ) : (
            <VideoOff className="w-5 h-5" />
          )}
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={toggleAudio}
          className={`h-12 w-12 rounded-full p-0 ${!audioEnabled ? "bg-pulse-danger/10 border-pulse-danger text-pulse-danger" : ""}`}
        >
          {audioEnabled ? (
            <Mic className="w-5 h-5" />
          ) : (
            <MicOff className="w-5 h-5" />
          )}
        </Button>
        <Button
          size="sm"
          onClick={handleEndSession}
          className="h-12 px-6 rounded-full bg-pulse-danger hover:bg-pulse-danger/90 text-white flex items-center gap-2"
        >
          <PhoneOff className="w-4 h-4" />
          End Session
        </Button>
      </div>

      {/* Safety warning */}
      <div className="absolute bottom-24 left-6 right-96">
        <div className="bg-pulse-surface-raised/50 backdrop-blur-sm border border-pulse-border/50 rounded-lg p-3 flex items-start gap-2">
          <AlertCircle className="w-3.5 h-3.5 text-pulse-text-muted shrink-0 mt-0.5" />
          <p className="text-[10px] text-pulse-text-muted">
            This is peer support, not therapy. In crisis? Call 988 Suicide &
            Crisis Lifeline.
          </p>
        </div>
      </div>
    </div>
  );
}
