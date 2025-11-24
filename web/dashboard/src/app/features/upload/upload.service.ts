import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, switchMap } from 'rxjs';

export interface CreateUploadResponse {
  id: string;
  upload_url: string;
  key: string;
}

@Injectable({
  providedIn: 'root'
})
export class UploadService {
  private apiUrl = 'http://localhost:8080'; // Adjust if needed

  constructor(private http: HttpClient) {}

  uploadVideo(file: File): Observable<any> {
    // 1. Create upload session
    return this.http.post<CreateUploadResponse>(`${this.apiUrl}/uploads`, {
      filename: file.name
    }, { withCredentials: true }).pipe(
      // 2. Upload file to presigned URL
      switchMap(response => {
        return this.http.put(response.upload_url, file, {
          headers: {
            'Content-Type': file.type
          },
          reportProgress: true,
          observe: 'events'
        });
      })
    );
  }
}
