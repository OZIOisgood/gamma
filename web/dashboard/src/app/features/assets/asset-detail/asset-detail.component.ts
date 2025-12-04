import { CommonModule, NgForOf, NgSwitch, NgSwitchCase, NgSwitchDefault } from '@angular/common';
import { Component, computed, ElementRef, inject, OnDestroy, OnInit, signal, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TuiTable } from '@taiga-ui/addon-table';
import { TuiLoader } from '@taiga-ui/core';
import { TuiBadge, TuiStatus } from '@taiga-ui/kit';
import Hls from 'hls.js';
import { of, Subscription } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Asset, AssetsService } from '../../../core/assets/assets.service';
import { NavbarComponent } from '../../../core/navbar/navbar.component';

@Component({
  selector: 'app-asset-detail',
  standalone: true,
  imports: [
    CommonModule, 
    TuiLoader, 
    TuiBadge, 
    TuiStatus, 
    NavbarComponent, 
    TuiTable,
    NgForOf,
    NgSwitch,
    NgSwitchCase,
    NgSwitchDefault
  ],
  template: `
    <app-navbar></app-navbar>
    <div class="container">
      <tui-loader [overlay]="true" [showLoader]="loading()">
        <div *ngIf="asset() as asset" class="content">
          <header>
            <h1>Asset #{{ getShortId(asset.ID) }}</h1>
          </header>

          <div class="details-grid">
            <div class="player-card" *ngIf="isReady(asset)">
              <div class="video-container">
                <video #videoPlayer controls class="video-js"></video>
              </div>
            </div>
            
            <div class="player-card" *ngIf="!isReady(asset)">
               <p>Video is processing or failed. Player unavailable.</p>
            </div>

            <div class="info-card">
              <table tuiTable [columns]="columns" class="info-table">
                <tbody tuiTbody [data]="infoData()">
                  <tr *ngFor="let item of infoData()" tuiTr>
                    <td *tuiCell="'label'" tuiTd class="label-cell">{{ item.label }}</td>
                    <td *tuiCell="'value'" tuiTd>
                      <ng-container [ngSwitch]="item.type">
                        <tui-badge *ngSwitchCase="'status'" [appearance]="getStatusAppearance(asset.Status)" tuiStatus>
                          {{ asset.Status }}
                        </tui-badge>
                        <span *ngSwitchCase="'date'">{{ item.value | date:'medium' }}</span>
                        <span *ngSwitchDefault class="break-all">{{ item.value }}</span>
                      </ng-container>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        
        <div *ngIf="error()" class="error-message">
          {{ error() }}
        </div>
      </tui-loader>
    </div>
  `,
  styles: [`
    .container {
      padding: 2rem;
      max-width: 1200px;
      margin: 0 auto;
    }
    
    header {
      margin-bottom: 1rem;
      border-bottom: 1px solid var(--tui-base-03);
      padding-bottom: 0.5rem;
    }
    
    .details-grid {
      display: flex;
      flex-direction: column;
      gap: 0;
    }
    
    .info-card, .player-card {
      background: var(--tui-base-01);
      border-radius: 0;
      overflow: hidden;
    }

    .info-table {
      width: 100%;
      border-collapse: collapse;
    }

    .info-table tr:first-child td {
      border-top: 1px solid var(--tui-base-03);
    }

    .label-cell {
      font-weight: bold;
      width: 150px;
      color: var(--tui-text-02);
    }

    .break-all {
      word-break: break-all;
    }
    
    .video-container {
      width: 100%;
      aspect-ratio: 16/9;
      background: black;
    }
    
    video {
      width: 100%;
      height: 100%;
    }
    
    .error-message {
      color: var(--tui-error-fill);
      text-align: center;
      margin-top: 2rem;
    }
  `]
})
export class AssetDetailComponent implements OnInit, OnDestroy {
  @ViewChild('videoPlayer') videoElement?: ElementRef<HTMLVideoElement>;
  
  private route = inject(ActivatedRoute);
  private assetsService = inject(AssetsService);
  
  asset = signal<Asset | null>(null);
  loading = signal(true);
  error = signal<string | null>(null);
  
  readonly columns = ['label', 'value'];
  infoData = computed(() => {
    const asset = this.asset();
    if (!asset) return [];
    return [
      { label: 'ID', value: asset.ID },
      { label: 'Status', value: asset.Status, type: 'status' as const },
      { label: 'Upload ID', value: asset.UploadID },
      { label: 'Created At', value: asset.CreatedAt, type: 'date' as const },
      { label: 'Updated At', value: asset.UpdatedAt, type: 'date' as const },
      { label: 'HLS Root', value: asset.HlsRoot }
    ];
  });

  private hls: Hls | null = null;
  private sub?: Subscription;

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      this.error.set('No asset ID provided');
      this.loading.set(false);
      return;
    }

    this.sub = this.assetsService.getAsset(id).pipe(
      tap(asset => {
        this.asset.set(asset);
        if (this.isReady(asset)) {
          this.initPlayer(id);
        } else {
          this.loading.set(false);
        }
      }),
      catchError(err => {
        console.error('Error loading asset', err);
        this.error.set('Failed to load asset details');
        this.loading.set(false);
        return of(null);
      })
    ).subscribe();
  }

  ngOnDestroy() {
    if (this.hls) {
      this.hls.destroy();
    }
    this.sub?.unsubscribe();
  }

  getShortId(id: string): string {
    return id.substring(0, 8);
  }

  isReady(asset: Asset): boolean {
    return asset.Status.toLowerCase() === 'ready';
  }

  getStatusAppearance(status: string): string {
    switch (status.toLowerCase()) {
      case 'ready': return 'positive';
      case 'failed': return 'negative';
      case 'processing': return 'warning';
      default: return 'neutral';
    }
  }

  private initPlayer(id: string) {
    this.assetsService.getAssetPlaylist(id).subscribe({
      next: (response) => {
        this.loading.set(false);
        // Give Angular a tick to render the video element
        setTimeout(() => {
          if (this.videoElement) {
            this.setupHls(response.url);
          }
        });
      },
      error: (err) => {
        console.error('Failed to get playlist', err);
        this.loading.set(false);
      }
    });
  }

  private setupHls(url: string) {
    const video = this.videoElement!.nativeElement;

    if (Hls.isSupported()) {
      this.hls = new Hls();
      this.hls.loadSource(url);
      this.hls.attachMedia(video);
      this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
        // video.play().catch(e => console.log('Auto-play prevented', e));
      });
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = url;
      video.addEventListener('loadedmetadata', () => {
        // video.play().catch(e => console.log('Auto-play prevented', e));
      });
    }
  }
}
