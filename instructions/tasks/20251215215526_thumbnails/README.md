# Task: Thumbnails

## Status
- [x] Defined
- [x] In Progress
- [x] Completed

## Description
As a user, I want to see preview thumbnails when I hover over the video timeline, so that I can easily navigate to specific parts of the video.

## Context
- [Media Chrome Preview Thumbnail](https://www.media-chrome.org/docs/en/components/media-preview-thumbnail)
- We need to generate a VTT file and a set of images (or a sprite sheet) during video processing.
- The VTT file should be exposed via the API as part of the Asset metadata.

## Acceptance Criteria
- [x] `assets` table has a `thumbnail_root` column.
- [x] Worker generates thumbnails (images + VTT) during video processing.
- [x] Thumbnails are uploaded to S3.
- [x] API returns `thumbnail_root` in Asset response. 
