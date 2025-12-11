import { CommonModule } from '@angular/common';
import { Component, inject, output, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { TuiTextfield } from '@taiga-ui/core';
import { RealmService } from '../../../core/services/realm.service';

@Component({
  selector: 'app-create-realm',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TuiTextfield
  ],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()" class="create-realm-form">
      <div class="form-group">
        <label class="input-label">REALM NAME *</label>
        <tui-textfield class="theme-input" [tuiTextfieldSize]="'l'">
          <input tuiTextfield formControlName="name" placeholder="Enter realm name" />
        </tui-textfield>
      </div>

      <button 
        type="submit" 
        class="btn-primary"
        [disabled]="form.invalid || loading()"
      >
        {{ loading() ? 'CREATING...' : 'CREATE REALM' }}
      </button>

      <div *ngIf="error()" class="error-message">
        {{ error() }}
      </div>
    </form>
  `,
  styles: [`
    .create-realm-form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
    .error-message {
      color: red;
      font-size: 0.875rem;
      margin-top: 0.5rem;
    }
  `]
})
export class CreateRealmComponent {
  private readonly realmService = inject(RealmService);
  
  created = output<void>();
  
  form = new FormGroup({
    name: new FormControl('', [Validators.required])
  });

  loading = signal(false);
  error = signal('');

  onSubmit() {
    if (this.form.valid) {
      this.loading.set(true);
      this.error.set('');
      
      this.realmService.create(this.form.value.name!).subscribe({
        next: () => {
          this.loading.set(false);
          this.form.reset();
          this.created.emit();
        },
        error: (err) => {
          this.loading.set(false);
          this.error.set('Failed to create realm');
          console.error(err);
        }
      });
    }
  }
}
