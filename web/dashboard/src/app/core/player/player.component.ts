import { CommonModule } from '@angular/common';
import { Component, CUSTOM_ELEMENTS_SCHEMA, effect, ElementRef, input, OnDestroy, ViewChild } from '@angular/core';
import Hls from 'hls.js';
import 'media-chrome';
import 'media-chrome/menu';

@Component({
  selector: 'app-player',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './player.component.html',
  styleUrls: ['./player.component.scss'],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class PlayerComponent implements OnDestroy {
  src = input.required<string>();
  thumbnailSrc = input<string | undefined>();

  @ViewChild('video') videoRef!: ElementRef<HTMLVideoElement>;
  private hls: Hls | null = null;

  constructor() {
    effect(() => {
      const src = this.src();
      if (!src || !this.videoRef) return;

      const video = this.videoRef.nativeElement;

      if (Hls.isSupported()) {
        if (this.hls) {
          this.hls.destroy();
        }
        this.hls = new Hls();
        this.hls.loadSource(src);
        this.hls.attachMedia(video);
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = src;
      }
    });
  }

  ngOnDestroy() {
    if (this.hls) {
      this.hls.destroy();
    }
  }
}
