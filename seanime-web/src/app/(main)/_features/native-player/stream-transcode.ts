import { MKVParser_Metadata } from "@/api/generated/types"
import { StreamTranscodeMode } from "@/app/(main)/_atoms/playback.atoms"

type CanPlayType = (mimeType: string) => string

// MKV codec IDs every modern browser decodes in <video>
const SAFE_VIDEO_CODECS = ["V_MPEG4/ISO/AVC", "V_VP8", "V_VP9", "V_AV1"]
const SAFE_AUDIO_CODECS = ["A_AAC", "A_OPUS", "A_VORBIS", "A_FLAC", "A_MPEG/L3"]

// Codecs some browsers decode, checked against canPlayType
const PROBED_CODECS: Record<string, string> = {
    "V_MPEGH/ISO/HEVC": "video/mp4; codecs=\"hvc1.1.6.L120.90\"",
    "A_AC3": "audio/mp4; codecs=\"ac-3\"",
    "A_EAC3": "audio/mp4; codecs=\"ec-3\"",
}

function canPlayCodec(codecID: string, safe: string[], canPlayType: CanPlayType) {
    if (safe.some(c => codecID.startsWith(c))) return true
    const probe = PROBED_CODECS[codecID]
    return !!probe && canPlayType(probe) !== ""
}

/**
 * Decides whether a torrent stream should be transcoded on the server before playing in this browser.
 * In "auto" mode, only MKV streams whose video or default audio codec this browser can't decode are transcoded.
 */
export function shouldTranscodeStream(mode: StreamTranscodeMode, metadata: MKVParser_Metadata | undefined, canPlayType: CanPlayType) {
    if (mode === "always") return true
    if (mode === "off" || !metadata) return false

    const video = metadata.videoTracks?.[0]
    if (video && !canPlayCodec(video.codecID, SAFE_VIDEO_CODECS, canPlayType)) return true

    const audio = metadata.audioTracks?.find(t => t.default) ?? metadata.audioTracks?.[0]
    if (audio && !canPlayCodec(audio.codecID, SAFE_AUDIO_CODECS, canPlayType)) return true

    return false
}
