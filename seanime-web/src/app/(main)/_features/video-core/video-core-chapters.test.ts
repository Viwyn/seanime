import { describe, expect, it } from "vitest"
import { vc_createChapterCues } from "./video-core-chapter-cues"

function totalLength(cues: { startTime: number, endTime: number }[]) {
    return cues.reduce((sum, c) => sum + (c.endTime - c.startTime), 0)
}

describe("chapter cues", () => {
    it("drops duplicated chapter editions", () => {
        const edition = [
            { uid: 1, start: 0, end: 90, text: "Opening" },
            { uid: 2, start: 90, end: 1300, text: "Part A" },
            { uid: 3, start: 1300, end: 1440, text: "Ending" },
        ]
        const cues = vc_createChapterCues([...edition, ...edition], 1440)

        expect(cues.map(c => c.text)).toEqual(["Opening", "Part A", "Ending"])
        expect(totalLength(cues)).toBe(1440)
    })

    it("handles chapters without end times sharing a start (ToonsHub NF release)", () => {
        const cues = vc_createChapterCues([
            { uid: 1, start: 0, end: 0, text: "Part 1" },
            { uid: 2, start: 0, end: 0, text: "Recap" },
            { uid: 3, start: 19.978, end: 0, text: "Part 2" },
            { uid: 4, start: 1345, end: 0, text: "Credits" },
        ], 1542.124)

        expect(cues).toEqual([
            { startTime: 0, endTime: 19.978, text: "Recap" },
            { startTime: 19.978, endTime: 1345, text: "Part 2" },
            { startTime: 1345, endTime: 1542.124, text: "Credits" },
        ])
    })

    it("clips overlapping chapters and fills gaps", () => {
        const cues = vc_createChapterCues([
            { uid: 1, start: 10, end: 100, text: "Opening" },
            { uid: 2, start: 80, end: 1000, text: "Part A" },
        ], 1000)

        expect(cues).toEqual([
            { startTime: 0, endTime: 10, text: "" },
            { startTime: 10, endTime: 80, text: "Opening" },
            { startTime: 80, endTime: 1000, text: "Part A" },
        ])
    })
})
