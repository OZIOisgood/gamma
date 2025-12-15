# Task: Cut max resolution

## Status
- [x] Defined
- [ ] In Progress
- [ ] Completed

## Description
Currently, the worker generates a fixed set of HLS variants (1080p, 720p, 480p, 360p, 240p, 144p) regardless of the input video's resolution. This leads to upscaling (e.g., 480p input -> 1080p output), which wastes storage/processing and provides poor quality.

The goal is to automatically detect the input video's resolution and only generate variants that are equal to or lower than the input resolution.

## Context
- `internal/worker/handler.go`: Contains the `processVideo` function with the hardcoded `ffmpeg` command.
- `Dockerfile`: Confirms `ffmpeg` (and `ffprobe`) is installed.

## Acceptance Criteria
- [ ] Detect input video resolution using `ffprobe`.
- [ ] Filter the target resolutions list (1080p, 720p, 480p, 360p, 240p, 144p) to only include those <= input height.
- [ ] Dynamically construct the `ffmpeg` command arguments based on the filtered list.
- [ ] Ensure at least one variant is generated (even if input is very low res).
