import { AsyncPipe, DatePipe, NgForOf, NgIf } from '@angular/common';
import { Component, inject } from '@angular/core';
import { TuiTable } from '@taiga-ui/addon-table';
import { TuiLoader } from '@taiga-ui/core';
import { TuiBadge, TuiStatus } from '@taiga-ui/kit';
import { AssetsService } from '../../core/assets/assets.service';
import { NavbarComponent } from '../../core/navbar/navbar.component';
import { UploadDrawerComponent } from '../upload/upload-drawer/upload-drawer.component';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    NavbarComponent, 
    TuiTable, 
    NgForOf, 
    AsyncPipe, 
    DatePipe, 
    TuiLoader, 
    NgIf, 
    TuiBadge, 
    TuiStatus,
    UploadDrawerComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent {
  private readonly assetsService = inject(AssetsService);
  
  readonly columns = ['ID', 'Title', 'Status', 'CreatedAt'];
  readonly data$ = this.assetsService.getUploads();

  getShortId(id: string): string {
    return id.split('-')[0];
  }

  getStatusAppearance(status: string): string {
    switch (status.toLowerCase()) {
      case 'ready':
      case 'uploaded':
        return 'positive';
      case 'failed':
        return 'negative';
      case 'processing':
      case 'pending':
        return 'warning';
      default:
        return 'neutral';
    }
  }
}
