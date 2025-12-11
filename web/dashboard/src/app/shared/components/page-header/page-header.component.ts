import { CommonModule } from '@angular/common';
import { Component, input, output } from '@angular/core';

@Component({
  selector: 'app-page-header',
  standalone: true,
  imports: [CommonModule],
  template: `
    <header class="page-header">
      <h1 class="page-title">{{ title() }}</h1>
      <button 
        *ngIf="buttonText() && showButton()"
        class="btn-primary header-btn"
        [disabled]="buttonDisabled()"
        (click)="buttonClick.emit()"
      >
        {{ buttonText() }}
      </button>
    </header>
  `,
  styles: [`
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 2rem;
    }
    .page-title {
      font-size: 3rem;
      font-weight: 900;
      text-transform: uppercase;
      letter-spacing: -0.05em;
    }
    .header-btn {
      width: auto;
      padding: 0 2rem;
    }
  `]
})
export class PageHeaderComponent {
  title = input.required<string>();
  buttonText = input<string>('');
  buttonDisabled = input<boolean>(false);
  showButton = input<boolean>(true);
  buttonClick = output<void>();
}
