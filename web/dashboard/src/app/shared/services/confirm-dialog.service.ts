import { inject, Injectable } from '@angular/core';
import { TuiDialogService } from '@taiga-ui/core';
import { PolymorpheusComponent } from '@taiga-ui/polymorpheus';
import { Observable } from 'rxjs';
import { ConfirmDialogComponent } from '../components/confirm-dialog/confirm-dialog.component';

export interface ConfirmDialogData {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
}

@Injectable({
  providedIn: 'root'
})
export class ConfirmDialogService {
  private readonly dialogs = inject(TuiDialogService);

  confirm(data: ConfirmDialogData): Observable<boolean> {
    return this.dialogs.open<boolean>(
      new PolymorpheusComponent(ConfirmDialogComponent),
      {
        data,
        dismissible: true,
        closeable: true,
        size: 's',
        label: data.title,
        appearance: 'confirm-dialog'
      }
    );
  }
}
