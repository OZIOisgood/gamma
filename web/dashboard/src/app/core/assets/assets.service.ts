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
  private readonly baseUrl = 'http://localhost:8080';
  private refreshSubject = new BehaviorSubject<void>(undefined);

  constructor(private http: HttpClient) {}

  getUploads(realm: string): Observable<Upload[]> {
    return this.refreshSubject.pipe(
      switchMap(() => this.http.get<Upload[]>(`${this.baseUrl}/${realm}/uploads`, {
        withCredentials: true,
      }))
    );
  }

  getAsset(realm: string, id: string): Observable<Asset> {
    return this.http.get<Asset>(`${this.baseUrl}/${realm}/assets/${id}`, {
      withCredentials: true,
    });
  }

  getPlaylist(realm: string, id: string): Observable<PlaylistResponse> {
    return this.http.get<PlaylistResponse>(`${this.baseUrl}/${realm}/assets/${id}/playlist`, {
      withCredentials: true,
    });
  }

  deleteAsset(realm: string, id: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${realm}/assets/${id}`, {
      withCredentials: true,
    });
  }

  refresh() {
    this.refreshSubject.next();
  }
}
