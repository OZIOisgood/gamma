import { AsyncPipe, DatePipe, NgForOf, NgIf } from '@angular/common';
import { HttpEventType } from '@angular/common/http';
import { Component, ElementRef, inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { TuiTable } from '@taiga-ui/addon-table';
import { TuiButton, TuiIcon, TuiLoader } from '@taiga-ui/core';
import { TuiBadge, TuiStatus } from '@taiga-ui/kit';
import { BehaviorSubject, map, Subscription, switchMap } from 'rxjs';
import { AssetsService, Upload } from '../../core/assets/assets.service';
import { NavbarComponent } from '../../core/navbar/navbar.component';
import { WebsocketService } from '../../core/services/websocket.service';
import { UploadDrawerComponent } from '../upload/upload-drawer/upload-drawer.component';
import { UploadService } from '../upload/upload.service';

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
    UploadDrawerComponent,
    TuiIcon,
    TuiButton
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  private readonly assetsService = inject(AssetsService);
  private readonly websocketService = inject(WebsocketService);
  private readonly uploadService = inject(UploadService);
  private readonly router = inject(Router);
  
  private readonly refresh$ = new BehaviorSubject<void>(undefined);
  private wsSubscription?: Subscription;

  readonly columns = ['ID', 'Title', 'Status', 'CreatedAt'];
  readonly data$ = this.refresh$.pipe(
    switchMap(() => this.assetsService.getUploads().pipe(
      map(data => data || [])
    ))
  );

  isDragging = false;
  uploading = false;
  progress = 0;

  @ViewChild('fileInput') fileInput!: ElementRef<HTMLInputElement>;

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

  getShortId(id: string): string {
    return id.substring(0, 8);
  }

  getStatusAppearance(status: string): string {
    switch (status) {
      case 'completed':
      case 'ready': return 'success';
      case 'processing': return 'info';
      case 'pending': return 'warning';
      case 'failed': return 'error';
      default: return 'neutral';
    }
  }

  onRowClick(item: Upload) {
    this.router.navigate(['/assets', item.ID]);
  }

  triggerFileInput() {
    this.fileInput.nativeElement.click();
  }

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.handleFile(input.files[0]);
    }
  }

  onDragOver(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = true;
  }

  onDragLeave(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
  }

  onDrop(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
    
    if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
      const file = event.dataTransfer.files[0];
      if (file.type.startsWith('video/')) {
        this.handleFile(file);
      } else {
        alert('Only video files are allowed');
      }
    }
  }

  handleFile(file: File) {
    this.uploading = true;
    this.progress = 0;

    this.uploadService.uploadVideo(file).subscribe({
      next: (event: any) => {
        if (event.type === HttpEventType.UploadProgress) {
          this.progress = Math.round(100 * event.loaded / event.total);
        } else if (event.type === HttpEventType.Response) {
          this.uploading = false;
          this.refresh$.next();
        }
      },
      error: (err) => {
        console.error('Upload failed', err);
        this.uploading = false;
        alert('Upload failed');
      }
    });
  }
}
