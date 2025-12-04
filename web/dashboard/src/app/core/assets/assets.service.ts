import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, switchMap } from 'rxjs';

export interface Upload {
  ID: string;
  Title: string;
  S3Key: string;
  Status: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface Asset {
  ID: string;
  UploadID: string;
  HlsRoot: string;
  Status: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface PlaylistResponse {
  url: string;
}

@Injectable({
  providedIn: 'root',
})
export class AssetsService {
  private readonly apiUrl = 'http://localhost:8080/uploads';
  private readonly assetsUrl = 'http://localhost:8080/assets';
  private refreshSubject = new BehaviorSubject<void>(undefined);

  constructor(private http: HttpClient) {}

  getUploads(): Observable<Upload[]> {
    return this.refreshSubject.pipe(
      switchMap(() => this.http.get<Upload[]>(this.apiUrl, {
        withCredentials: true,
      }))
    );
  }

  getAsset(id: string): Observable<Asset> {
    return this.http.get<Asset>(`${this.assetsUrl}/${id}`, {
      withCredentials: true,
    });
  }

  getAssetPlaylist(id: string): Observable<PlaylistResponse> {
    return this.http.get<PlaylistResponse>(`${this.assetsUrl}/${id}/playlist`, {
      withCredentials: true,
    });
  }

  refresh() {
    this.refreshSubject.next();
  }
}
