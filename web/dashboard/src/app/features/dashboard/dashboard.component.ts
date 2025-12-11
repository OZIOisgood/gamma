import { AsyncPipe, DatePipe, NgForOf, NgIf } from '@angular/common';
import { HttpEventType } from '@angular/common/http';
import { Component, ElementRef, inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { TuiTable } from '@taiga-ui/addon-table';
import { TuiButton, TuiIcon, TuiLoader } from '@taiga-ui/core';
import { TuiBadge, TuiStatus } from '@taiga-ui/kit';
import { BehaviorSubject, EMPTY, map, Observable, Subscription, switchMap } from 'rxjs';
import { AssetsService, Upload } from '../../core/assets/assets.service';
import { UploadUiService } from '../../core/services/upload-ui.service';
import { WebsocketService } from '../../core/services/websocket.service';
import { PageHeaderComponent } from '../../shared/components/page-header/page-header.component';
import { UploadService } from '../upload/upload.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    TuiTable, 
    NgForOf, 
    AsyncPipe, 
    DatePipe, 
    TuiLoader, 
    NgIf, 
    TuiBadge, 
    TuiStatus,
    TuiIcon,
    TuiButton,
    PageHeaderComponent
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  private readonly assetsService = inject(AssetsService);
  private readonly websocketService = inject(WebsocketService);
  private readonly uploadService = inject(UploadService);
  private readonly uploadUiService = inject(UploadUiService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  
  private readonly refresh$ = new BehaviorSubject<void>(undefined);
  private wsSubscription?: Subscription;

  readonly columns = ['ID', 'Title', 'Status', 'CreatedAt'];
  
  private getParentParamMap(): Observable<ParamMap> {
    return this.route.parent?.paramMap || EMPTY;
  }
  
  readonly data$ = this.getParentParamMap().pipe(
    switchMap((params: ParamMap) => {
      const realm = params.get('realm') || 'default';
      return this.refresh$.pipe(
        switchMap(() => this.assetsService.getUploads(realm).pipe(
          map(data => data || [])
        ))
      );
    })
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
    const realm = this.route.parent?.snapshot.paramMap.get('realm') || 'default';
    this.router.navigate(['/', realm, 'assets', item.ID]);
  }

  openUploadDrawer() {
    this.uploadUiService.open();
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
    const realm = this.route.parent?.snapshot.paramMap.get('realm') || 'default';
    this.uploading = true;
    this.progress = 0;

    this.uploadService.uploadVideo(realm, file).subscribe({
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
