import { AsyncPipe, DatePipe, NgForOf, NgIf } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { TuiTable } from '@taiga-ui/addon-table';
import { TuiLoader } from '@taiga-ui/core';
import { TuiBadge, TuiStatus } from '@taiga-ui/kit';
import { BehaviorSubject, Subscription, switchMap } from 'rxjs';
import { AssetsService, Upload } from '../../core/assets/assets.service';
import { NavbarComponent } from '../../core/navbar/navbar.component';
import { WebsocketService } from '../../core/services/websocket.service';
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
export class DashboardComponent implements OnInit, OnDestroy {
  private readonly assetsService = inject(AssetsService);
  private readonly websocketService = inject(WebsocketService);
  private readonly router = inject(Router);
  
  private readonly refresh$ = new BehaviorSubject<void>(undefined);
  private wsSubscription?: Subscription;

  readonly columns = ['ID', 'Title', 'Status', 'CreatedAt'];
  readonly data$ = this.refresh$.pipe(
    switchMap(() => this.assetsService.getUploads())
  );

  ngOnInit() {
    this.wsSubscription = this.websocketService.messages$.subscribe(msg => {
      if (msg.type === 'asset_processed') {
        this.refresh$.next();
      }
    });
  }

  ngOnDestroy() {
    this.wsSubscription?.unsubscribe();
  }

  onRowClick(item: Upload) {
    this.router.navigate(['/assets', item.ID]);
  }

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
