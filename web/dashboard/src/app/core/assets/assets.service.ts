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

@Injectable({
  providedIn: 'root',
})
export class AssetsService {
  private readonly apiUrl = 'http://localhost:8080/uploads';
  private refreshSubject = new BehaviorSubject<void>(undefined);

  constructor(private http: HttpClient) {}

  getUploads(): Observable<Upload[]> {
    return this.refreshSubject.pipe(
      switchMap(() => this.http.get<Upload[]>(this.apiUrl, {
        withCredentials: true,
      }))
    );
  }

  refresh() {
    this.refreshSubject.next();
  }
}
