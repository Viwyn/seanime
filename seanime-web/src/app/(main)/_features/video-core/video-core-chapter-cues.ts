import { MKVParser_ChapterInfo } from "@/api/generated/types"
import { VideoCoreChapterCue } from "@/app/(main)/_features/video-core/video-core"

export const vc_createChapterCues = (chapters: Array<MKVParser_ChapterInfo> | undefined, duration: number): VideoCoreChapterCue[] => {
    if (!chapters || chapters.length === 0 || duration === 0) {
        return []
    }

    // Chapters sharing a start (e.g. a zero-length "Part 1" before "Recap") keep the last one, like mpv does
    const sorted = [...chapters]
        .sort((a, b) => a.start - b.start)
        .filter((chapter, index, arr) => arr[index + 1]?.start !== chapter.start)

    return vc_fillChapterCues(sorted.map((chapter, index) => {
        // A chapter ends where the next one starts; MKV chapters often have no end (0) or one that overlaps the next
        const nextStart = sorted[index + 1]?.start ?? duration
        return {
            startTime: chapter.start,
            endTime: chapter.end > chapter.start ? Math.min(chapter.end, nextStart) : nextStart,
            text: chapter.text || ``,
        }
    })).filter(c => c.startTime !== undefined && c.endTime !== undefined && c.startTime <= duration && c.endTime <= duration)
}

export const vc_fillChapterCues = (chapters: VideoCoreChapterCue[]): VideoCoreChapterCue[] => {
    if (!chapters || chapters.length === 0) {
        return []
    }

    const EPS = 1e-6
    const sorted = [...chapters].sort((a, b) => (a.startTime ?? 0) - (b.startTime ?? 0))
    const out: VideoCoreChapterCue[] = []

    for (let i = 0; i < sorted.length; i++) {
        const cur = sorted[i]
        const next = sorted[i + 1]
        const prevEnd = out[out.length - 1]?.endTime ?? 0
        // clip overlaps and drop duplicates (e.g. MKV files with multiple chapter editions)
        // so the chapter widths never add up to more than the duration
        const start = Math.max(cur.startTime ?? 0, prevEnd)
        const end = typeof cur.endTime === "number" ? cur.endTime : (next ? (next.startTime ?? start) : start)
        if (end <= start + EPS) continue

        // gap between previous output end (or 0) and this start
        if (start > prevEnd + EPS) {
            out.push({ startTime: prevEnd, endTime: start, text: "" })
        }

        out.push({ startTime: start, endTime: end, text: cur.text ?? "" })
    }

    return out
}
