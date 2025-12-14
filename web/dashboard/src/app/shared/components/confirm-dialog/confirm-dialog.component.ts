import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TuiButton, TuiDialogContext } from '@taiga-ui/core';
import { injectContext } from '@taiga-ui/polymorpheus';

export interface ConfirmDialogData {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
}

@Component({
  selector: 'app-confirm-dialog',
  standalone: true,
  imports: [CommonModule, TuiButton],
  template: `
    <div class="confirm-dialog">
      <div class="dialog-header">
        <h2 class="dialog-title">{{ data.title }}</h2>
      </div>
      
      <div class="dialog-content">
        <p class="dialog-message">{{ data.message }}</p>
      </div>
      
      <div class="dialog-actions">
        <button 
          tuiButton
          appearance="flat"
          size="m"
          class="cancel-btn"
          (click)="onCancel()"
        >
          {{ data.cancelText || 'Cancel' }}
        </button>
        
        <button 
          tuiButton
          appearance="accent"
          size="m"
          class="confirm-btn"
          (click)="onConfirm()"
        >
          {{ data.confirmText || 'Confirm' }}
        </button>
      </div>
    </div>
  `,
  styleUrls: ['./confirm-dialog.component.scss']
})
export class ConfirmDialogComponent {
  private readonly context = injectContext<TuiDialogContext<boolean, ConfirmDialogData>>();
  
  get data(): ConfirmDialogData {
    return this.context.data;
  }
  
  onConfirm(): void {
    this.context.completeWith(true);
  }
  
  onCancel(): void {
    this.context.completeWith(false);
  }
}
