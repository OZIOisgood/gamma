# Resolution

## Summary
Implemented dynamic resolution generation in the worker. The worker now probes the input video using `ffprobe` and only generates HLS variants that are equal to or lower than the input resolution.

## Technical Details
1.  **`ProbeData` & `Quality` structs**: Added structures to handle `ffprobe` output and define available quality presets.
2.  **`getVideoDimensions`**: Added a helper method to run `ffprobe` and extract video width/height.
3.  **Dynamic `ffmpeg` command**:
    -   Modified `processVideo` to call `getVideoDimensions`.
    -   Filters `availableQualities` (1080p, 720p, etc.) against the input height.
    -   Dynamically constructs the `filter_complex`, `-map` arguments, and `-var_stream_map` for `ffmpeg` based on the filtered qualities.
    -   Ensures at least one quality is generated (fallback to lowest if input is extremely small).

## Verification
- [x] Build passed (`go build ./cmd/worker`)
- [ ] Verified in UI/API (Requires running the full stack and uploading a video)
