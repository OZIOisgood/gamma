import { CommonModule, NgForOf, NgSwitch, NgSwitchCase, NgSwitchDefault } from '@angular/common';
import { Component, computed, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TuiTable } from '@taiga-ui/addon-table';
import { TuiButton, TuiLoader } from '@taiga-ui/core';
import { TuiBadge, TuiStatus } from '@taiga-ui/kit';
import { of, Subscription } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Asset, AssetsService } from '../../../core/assets/assets.service';
import { PlayerComponent } from '../../../core/player/player.component';

@Component({
  selector: 'app-asset-detail',
  standalone: true,
  imports: [
    CommonModule, 
    TuiLoader, 
    TuiButton,
    TuiBadge, 
    TuiStatus, 
    TuiTable,
    NgForOf,
    NgSwitch,
    NgSwitchCase,
    NgSwitchDefault,
    PlayerComponent
  ],
  templateUrl: './asset-detail.component.html',
  styleUrls: ['./asset-detail.component.scss']
})
export class AssetDetailComponent implements OnInit, OnDestroy {
  
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private assetsService = inject(AssetsService);
  
  asset = signal<Asset | null>(null);
  playlistUrl = signal<string>('');
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

  private sub?: Subscription;

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    const realm = this.route.parent?.snapshot.paramMap.get('realm') || 'default';
    if (!id) {
      this.error.set('No asset ID provided');
      this.loading.set(false);
      return;
    }

    this.sub = this.assetsService.getAsset(realm, id).pipe(
      tap(asset => {
        this.asset.set(asset);
        if (this.isReady(asset)) {
          this.initPlayer(realm, id);
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

  private initPlayer(realm: string, id: string) {
    this.assetsService.getPlaylist(realm, id).subscribe({
      next: (response) => {
        this.loading.set(false);
        this.playlistUrl.set(response.url);
      },
      error: (err) => {
        console.error('Failed to get playlist', err);
        this.loading.set(false);
      }
    });
  }

  deleteAsset() {
    const asset = this.asset();
    const realm = this.route.parent?.snapshot.paramMap.get('realm') || 'default';
    if (!asset) return;

    if (confirm('Are you sure you want to delete this asset?')) {
      this.assetsService.deleteAsset(realm, asset.ID).subscribe({
        next: () => {
          this.router.navigate(['/', realm, 'dashboard']);
        },
        error: (err) => {
          console.error('Failed to delete asset', err);
        }
      });
    }
  }
}
