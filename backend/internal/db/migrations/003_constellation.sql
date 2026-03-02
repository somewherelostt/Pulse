-- Layer 5: Constellation peer matching and WebRTC session coordination
-- Requires pgvector extension (available on Supabase by default)

CREATE EXTENSION IF NOT EXISTS vector;

-- -----------------------------------------------------------------------
-- peer_pool: users who have opted in to peer support matching
-- -----------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS public.peer_pool (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    fingerprint     vector(6),                          -- 6-dim behavioral fingerprint
    full_embedding  vector(1536),                       -- optional rich embedding (nullable)
    is_available    BOOLEAN     NOT NULL DEFAULT true,
    is_recovering   BOOLEAN     NOT NULL DEFAULT false,
    mood_recovered  BOOLEAN     NOT NULL DEFAULT false,  -- true when mood trend improved
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_active     TIMESTAMPTZ NOT NULL DEFAULT now(),
    opt_in_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id)
);

-- Index for availability filter (most queries filter on is_available = true)
CREATE INDEX IF NOT EXISTS idx_peer_pool_available
    ON public.peer_pool (is_available)
    WHERE is_available = true;

-- Vector similarity index on the 6-dim fingerprint
-- IVFFlat with lists=1 works even with very few rows
CREATE INDEX IF NOT EXISTS idx_peer_pool_fingerprint_cos
    ON public.peer_pool USING ivfflat (fingerprint vector_cosine_ops)
    WITH (lists = 1);

-- -----------------------------------------------------------------------
-- constellation_sessions: matched peer support sessions
-- -----------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS public.constellation_sessions (
    id                   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    seeker_id            UUID        NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    supporter_id         UUID        NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    room_id              UUID        UNIQUE,             -- NULL until session/start is called
    matched_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at           TIMESTAMPTZ,
    ended_at             TIMESTAMPTZ,
    seeker_rating        SMALLINT    CHECK (seeker_rating    BETWEEN 1 AND 5),
    supporter_rating     SMALLINT    CHECK (supporter_rating BETWEEN 1 AND 5),
    seeker_would_again   BOOLEAN,
    supporter_would_again BOOLEAN,
    context_hint         TEXT,
    similarity           DOUBLE PRECISION
);

CREATE INDEX IF NOT EXISTS idx_const_sessions_seeker
    ON public.constellation_sessions (seeker_id);
CREATE INDEX IF NOT EXISTS idx_const_sessions_supporter
    ON public.constellation_sessions (supporter_id);
CREATE INDEX IF NOT EXISTS idx_const_sessions_room
    ON public.constellation_sessions (room_id)
    WHERE room_id IS NOT NULL;

-- -----------------------------------------------------------------------
-- constellation_session_log: content-free audit trail
-- Only event_type is recorded, never message content.
-- -----------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS public.constellation_session_log (
    id          BIGSERIAL   PRIMARY KEY,
    session_id  UUID        NOT NULL REFERENCES public.constellation_sessions(id) ON DELETE CASCADE,
    event_type  TEXT        NOT NULL, -- 'joined','offer','answer','ice','heartbeat','end','disconnected','timeout'
    logged_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_const_session_log_session
    ON public.constellation_session_log (session_id);
