import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

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

  constructor(private http: HttpClient) {}

  getUploads(): Observable<Upload[]> {
    return this.http.get<Upload[]>(this.apiUrl, {
      withCredentials: true,
    });
  }
}
