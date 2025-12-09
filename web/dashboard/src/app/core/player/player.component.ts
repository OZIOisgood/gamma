import { Component, CUSTOM_ELEMENTS_SCHEMA, input } from '@angular/core';
import 'hls-video-element';
import 'media-chrome';
import 'media-chrome/menu';

@Component({
  selector: 'app-player',
  standalone: true,
  templateUrl: './player.component.html',
  styleUrls: ['./player.component.scss'],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class PlayerComponent {
  src = input.required<string>();
}
