import { CommonModule } from '@angular/common';
import { Component, input, output } from '@angular/core';
import { TuiButton } from '@taiga-ui/core';

@Component({
  selector: 'app-drawer',
  standalone: true,
  imports: [CommonModule, TuiButton],
  template: `
    <div class="backdrop" *ngIf="isOpen()" (click)="onBackdropClick($event)">
      <div class="drawer">
        <div class="drawer-header">
          <div class="header-title">{{ title() }}</div>
          <button 
            tuiIconButton 
            appearance="flat" 
            iconStart="@tui.material.sharp.close" 
            size="s" 
            (click)="close.emit()"
            class="close-btn"
          ></button>
        </div>
        
        <div class="drawer-content">
          <ng-content></ng-content>
        </div>
      </div>
    </div>
  `,
  styleUrls: ['./drawer.component.scss']
})
export class DrawerComponent {
  isOpen = input.required<boolean>();
  title = input.required<string>();
  close = output<void>();

  onBackdropClick(event: MouseEvent) {
    if ((event.target as HTMLElement).classList.contains('backdrop')) {
      this.close.emit();
    }
  }
}
