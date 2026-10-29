import { MKVParser_Metadata } from "@/api/generated/types"
import { describe, expect, it } from "vitest"
import { shouldTranscodeStream } from "./stream-transcode"

function metadata(video: string, audio: string): MKVParser_Metadata {
    return {
        videoTracks: [{ codecID: video }],
        audioTracks: [{ codecID: audio, default: true }],
    } as unknown as MKVParser_Metadata
}

const canPlayNothing = () => ""
const canPlayEverything = () => "probably"

describe("shouldTranscodeStream", () => {
    it("respects off and always", () => {
        const md = metadata("V_MPEGH/ISO/HEVC", "A_AC3")
        expect(shouldTranscodeStream("off", md, canPlayNothing)).toBe(false)
        expect(shouldTranscodeStream("always", metadata("V_MPEG4/ISO/AVC", "A_AAC"), canPlayEverything)).toBe(true)
    })

    it("does not transcode browser-safe codecs in auto", () => {
        expect(shouldTranscodeStream("auto", metadata("V_MPEG4/ISO/AVC", "A_AAC"), canPlayNothing)).toBe(false)
        expect(shouldTranscodeStream("auto", metadata("V_VP9", "A_OPUS"), canPlayNothing)).toBe(false)
    })

    it("transcodes HEVC or AC3 only when the browser can't play them", () => {
        expect(shouldTranscodeStream("auto", metadata("V_MPEGH/ISO/HEVC", "A_AAC"), canPlayNothing)).toBe(true)
        expect(shouldTranscodeStream("auto", metadata("V_MPEGH/ISO/HEVC", "A_AAC"), canPlayEverything)).toBe(false)
        expect(shouldTranscodeStream("auto", metadata("V_MPEG4/ISO/AVC", "A_AC3"), canPlayNothing)).toBe(true)
    })

    it("always transcodes codecs browsers never decode", () => {
        expect(shouldTranscodeStream("auto", metadata("V_MPEG4/ISO/AVC", "A_DTS"), canPlayEverything)).toBe(true)
    })

    it("does not transcode without MKV metadata in auto", () => {
        expect(shouldTranscodeStream("auto", undefined, canPlayNothing)).toBe(false)
    })
})
