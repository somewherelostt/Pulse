"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import {
  Video,
  VideoOff,
  Mic,
  MicOff,
  PhoneOff,
  Send,
  Info,
  AlertCircle,
  Copy,
  Check,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { createClient } from "@/lib/supabase/client";
import type { RealtimeChannel } from "@supabase/supabase-js";

interface PeerSessionProps {
  roomId: string;
  sessionId: string;
  matchContext: string;
  similarity: number;
  token: string;
  isCreator?: boolean;
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
  isCreator = false,
  onEnd,
}: PeerSessionProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputText, setInputText] = useState("");
  const [videoEnabled, setVideoEnabled] = useState(true);
  const [audioEnabled, setAudioEnabled] = useState(true);
  const [timeRemaining, setTimeRemaining] = useState(45 * 60);
  const [connectionState, setConnectionState] = useState<
    "waiting" | "connecting" | "connected" | "failed"
  >("waiting");
  const [chatEnabled, setChatEnabled] = useState(false);
  const [copied, setCopied] = useState(false);

  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const peerConnectionRef = useRef<RTCPeerConnection | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);
  const channelRef = useRef<RealtimeChannel | null>(null);
  const myPeerIdRef = useRef(
    `peer-${Math.random().toString(36).slice(2, 10)}`,
  );
  const makingOfferRef = useRef(false);
  const isEndingRef = useRef(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const shareUrl =
    typeof window !== "undefined"
      ? `${window.location.origin}/dashboard/constellation?room=${roomId}`
      : "";

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // 45-min countdown
  useEffect(() => {
    const interval = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev <= 0) {
          clearInterval(interval);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  const handleEndSession = useCallback(() => {
    if (isEndingRef.current) return;
    isEndingRef.current = true;

    try {
      channelRef.current?.send({
        type: "broadcast",
        event: "end",
        payload: {},
      });
    } catch {}

    localStreamRef.current?.getTracks().forEach((t) => t.stop());
    peerConnectionRef.current?.close();

    if (channelRef.current) {
      const supabase = createClient();
      supabase.removeChannel(channelRef.current);
      channelRef.current = null;
    }

    onEnd(sessionId);
  }, [sessionId, onEnd]);

  // Main session init: WebRTC + Supabase Realtime signaling
  useEffect(() => {
    const supabase = createClient();
    // Explicitly set JWT so Realtime WebSocket is authenticated
    if (token) {
      supabase.realtime.setAuth(token);
    }

    // --- WebRTC peer connection ---
    const pc = new RTCPeerConnection({
      iceServers: [
        { urls: "stun:stun.l.google.com:19302" },
        { urls: "stun:stun1.l.google.com:19302" },
      ],
    });
    peerConnectionRef.current = pc;

    // Get camera/mic (gracefully fails — text chat still works)
    navigator.mediaDevices
      .getUserMedia({ video: true, audio: true })
      .then((stream) => {
        localStreamRef.current = stream;
        if (localVideoRef.current) localVideoRef.current.srcObject = stream;
        stream.getTracks().forEach((track) => pc.addTrack(track, stream));
      })
      .catch((err) => {
        console.warn("Media access denied — text chat only:", err);
      });

    pc.ontrack = (event) => {
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = event.streams[0];
      }
    };

    pc.oniceconnectionstatechange = () => {
      if (
        pc.iceConnectionState === "connected" ||
        pc.iceConnectionState === "completed"
      ) {
        setConnectionState("connected");
      } else if (pc.iceConnectionState === "failed") {
        // Ice failed — but text chat may still work via Supabase
        pc.restartIce();
      }
    };

    // --- Supabase Realtime channel ---
    const channel = supabase.channel(`constellation:${roomId}`, {
      config: {
        broadcast: { self: false, ack: false },
        presence: { key: myPeerIdRef.current },
      },
    });
    channelRef.current = channel;

    // Forward local ICE candidates through the channel
    pc.onicecandidate = (event) => {
      if (event.candidate) {
        channel.send({
          type: "broadcast",
          event: "ice",
          payload: { candidate: event.candidate.toJSON() },
        });
      }
    };

    // --- Signaling message handlers ---
    channel.on("broadcast", { event: "offer" }, async ({ payload }) => {
      if (!payload?.sdp) return;
      try {
        await pc.setRemoteDescription(
          new RTCSessionDescription({ type: "offer", sdp: payload.sdp }),
        );
        const answer = await pc.createAnswer();
        await pc.setLocalDescription(answer);
        channel.send({
          type: "broadcast",
          event: "answer",
          payload: { sdp: answer.sdp },
        });
      } catch (err) {
        console.error("Error handling offer:", err);
      }
    });

    channel.on("broadcast", { event: "answer" }, async ({ payload }) => {
      if (!payload?.sdp) return;
      try {
        await pc.setRemoteDescription(
          new RTCSessionDescription({ type: "answer", sdp: payload.sdp }),
        );
      } catch (err) {
        console.error("Error handling answer:", err);
      }
    });

    channel.on("broadcast", { event: "ice" }, async ({ payload }) => {
      if (!payload?.candidate) return;
      try {
        await pc.addIceCandidate(new RTCIceCandidate(payload.candidate));
      } catch {}
    });

    // --- Chat handler (goes through Supabase, not WebRTC) ---
    channel.on("broadcast", { event: "chat" }, ({ payload }) => {
      if (!payload?.text) return;
      setMessages((prev) => [
        ...prev,
        {
          id: `${Date.now()}-peer`,
          text: payload.text,
          fromSelf: false,
          timestamp: new Date(),
        },
      ]);
    });

    // Peer ended the session
    channel.on("broadcast", { event: "end" }, () => {
      if (!isEndingRef.current) handleEndSession();
    });

    // --- Presence: detect when 2nd peer joins ---
    channel.on("presence", { event: "sync" }, async () => {
      const state = channel.presenceState();
      const peerIds = Object.keys(state);
      const count = peerIds.length;

      if (count >= 2) {
        setChatEnabled(true);
        setConnectionState((prev) =>
          prev === "waiting" ? "connecting" : prev,
        );

        // Deterministic: alphabetically-last peerId creates the offer
        const sorted = [...peerIds].sort();
        const shouldOffer =
          myPeerIdRef.current === sorted[sorted.length - 1];

        if (
          shouldOffer &&
          !makingOfferRef.current &&
          pc.signalingState === "stable"
        ) {
          makingOfferRef.current = true;
          try {
            const offer = await pc.createOffer();
            await pc.setLocalDescription(offer);
            channel.send({
              type: "broadcast",
              event: "offer",
              payload: { sdp: offer.sdp },
            });
          } catch (err) {
            console.error("Error creating offer:", err);
            makingOfferRef.current = false;
          }
        }
      }
    });

    channel.subscribe(async (status) => {
      if (status === "SUBSCRIBED") {
        await channel.track({
          peerId: myPeerIdRef.current,
          joined_at: Date.now(),
        });
      } else if (status === "CHANNEL_ERROR" || status === "TIMED_OUT") {
        setConnectionState("failed");
      }
    });

    return () => {
      localStreamRef.current?.getTracks().forEach((t) => t.stop());
      pc.close();
      supabase.removeChannel(channel);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [roomId]);

  const handleSendMessage = () => {
    if (!inputText.trim() || !channelRef.current || !chatEnabled) return;
    const text = inputText.trim();

    setMessages((prev) => [
      ...prev,
      {
        id: `${Date.now()}-self`,
        text,
        fromSelf: true,
        timestamp: new Date(),
      },
    ]);

    channelRef.current.send({
      type: "broadcast",
      event: "chat",
      payload: { text },
    });

    setInputText("");
  };

  const copyShareLink = () => {
    navigator.clipboard.writeText(shareUrl).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2500);
    });
  };

  const toggleVideo = () => {
    const track = localStreamRef.current?.getVideoTracks()[0];
    if (track) {
      track.enabled = !track.enabled;
      setVideoEnabled(track.enabled);
    }
  };

  const toggleAudio = () => {
    const track = localStreamRef.current?.getAudioTracks()[0];
    if (track) {
      track.enabled = !track.enabled;
      setAudioEnabled(track.enabled);
    }
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  const stateColor = {
    waiting: "bg-pulse-accent-warm animate-pulse",
    connecting: "bg-pulse-accent-warm animate-pulse",
    connected: "bg-emerald-400",
    failed: "bg-pulse-danger",
  }[connectionState];

  const stateLabel = {
    waiting: "Waiting for peer",
    connecting: "Connecting...",
    connected: "Connected",
    failed: "Failed",
  }[connectionState];

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

        <div className="flex items-center gap-3">
          {/* Share link button — always visible for creator */}
          {isCreator && (
            <button
              onClick={copyShareLink}
              className="flex items-center gap-1.5 text-xs text-pulse-primary border border-pulse-primary/30 px-3 py-1.5 rounded-lg hover:bg-pulse-primary/5 transition-all"
            >
              {copied ? (
                <Check className="w-3.5 h-3.5" />
              ) : (
                <Copy className="w-3.5 h-3.5" />
              )}
              {copied ? "Copied!" : "Copy invite link"}
            </button>
          )}

          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-pulse-bg border border-pulse-border">
            <div className={`w-1.5 h-1.5 rounded-full ${stateColor}`} />
            <span className="text-xs text-pulse-text-muted">{stateLabel}</span>
          </div>

          <span className="text-sm font-mono text-pulse-text-secondary">
            {formatTime(timeRemaining)}
          </span>
        </div>
      </div>

      {/* Waiting overlay — shown until a second peer joins */}
      {connectionState === "waiting" && (
        <div className="absolute inset-x-0 top-[73px] bottom-[80px] z-10 bg-pulse-bg/95 backdrop-blur-sm flex items-center justify-center">
          <div className="text-center max-w-sm px-6">
            <div className="w-16 h-16 rounded-full bg-pulse-primary/10 border border-pulse-primary/20 flex items-center justify-center mx-auto mb-4">
              <div className="w-6 h-6 border-2 border-pulse-primary/30 border-t-pulse-primary rounded-full animate-spin" />
            </div>
            <h3 className="text-base font-medium text-pulse-text-primary mb-2">
              {isCreator ? "Waiting for peer" : "Joining session…"}
            </h3>
            {isCreator ? (
              <>
                <p className="text-sm text-pulse-text-muted mb-5">
                  Share the invite link so your peer can join this anonymous
                  session.
                </p>
                <button
                  onClick={copyShareLink}
                  className="flex items-center gap-2 mx-auto text-sm font-medium text-pulse-primary bg-pulse-primary/10 border border-pulse-primary/20 px-5 py-2.5 rounded-lg hover:bg-pulse-primary/20 transition-all"
                >
                  {copied ? (
                    <Check className="w-4 h-4" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                  {copied ? "Link copied!" : "Copy invite link"}
                </button>
                <p className="text-[10px] text-pulse-text-muted mt-3 font-mono break-all">
                  {shareUrl}
                </p>
                <button
                  onClick={handleEndSession}
                  className="mt-5 text-xs text-pulse-text-muted hover:text-pulse-danger transition-colors"
                >
                  Cancel →
                </button>
              </>
            ) : (
              <p className="text-sm text-pulse-text-muted">
                Connecting to session room…
              </p>
            )}
          </div>
        </div>
      )}

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
          {connectionState !== "connected" && (
            <div className="absolute inset-0 flex items-center justify-center bg-pulse-surface/60">
              <span className="text-xs text-pulse-text-muted">
                Peer video pending…
              </span>
            </div>
          )}
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
      <div className="absolute right-6 top-[88px] bottom-[88px] w-80 bg-pulse-surface border border-pulse-border rounded-xl flex flex-col">
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
              placeholder={
                chatEnabled ? "Type a message…" : "Waiting for peer…"
              }
              disabled={!chatEnabled}
              className="flex-1 px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-xs text-pulse-text-primary placeholder:text-pulse-text-muted focus:outline-none focus:border-pulse-primary disabled:opacity-40"
            />
            <Button
              size="sm"
              onClick={handleSendMessage}
              disabled={!inputText.trim() || !chatEnabled}
              className="bg-pulse-primary text-white h-9 w-9 p-0"
            >
              <Send className="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      </div>

      {/* Match context hint */}
      {matchContext && connectionState !== "waiting" && (
        <div className="absolute top-[88px] left-6 max-w-sm">
          <div className="bg-pulse-surface/95 backdrop-blur-sm border border-pulse-border rounded-lg p-3 flex items-start gap-2">
            <Info className="w-4 h-4 text-pulse-primary shrink-0 mt-0.5" />
            <p className="text-xs text-pulse-text-secondary">{matchContext}</p>
          </div>
        </div>
      )}

      {/* Controls bar */}
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
      <div className="absolute bottom-[88px] left-6 right-[340px]">
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
