import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { TuiButton } from '@taiga-ui/core';
import { AuthService } from '../auth/auth.service';
import { UploadUiService } from '../services/upload-ui.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [TuiButton, RouterLink],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss'],
})
export class NavbarComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly uploadUiService = inject(UploadUiService);

  logout(): void {
    this.authService.logout().subscribe({
      next: () => {
        this.router.navigate(['/login']);
      },
      error: (err) => {
        console.error('Logout failed', err);
        // Force navigation even if logout fails
        this.router.navigate(['/login']);
      }
    });
  }

  openUpload(): void {
    this.uploadUiService.open();
  }
}
