import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  // Assuming the API is proxied or we use absolute URL. 
  // For local dev with docker-compose, it might be http://localhost:8080
  private apiUrl = 'http://localhost:8080/auth';

  constructor(private http: HttpClient) {}

  get isAuthenticated(): boolean {
    return !!localStorage.getItem('isLoggedIn');
  }

  login(credentials: any): Observable<any> {
    return new Observable(observer => {
      this.http.post(`${this.apiUrl}/login`, credentials, {
        withCredentials: true,
        responseType: 'text'
      }).subscribe({
        next: (res) => {
          localStorage.setItem('isLoggedIn', 'true');
          observer.next(res);
          observer.complete();
        },
        error: (err) => {
          observer.error(err);
        }
      });
    });
  }

  logout(): Observable<any> {
    return new Observable(observer => {
      this.http.post(`${this.apiUrl}/logout`, {}, {
        withCredentials: true,
        responseType: 'text'
      }).subscribe({
        next: (res) => {
          localStorage.removeItem('isLoggedIn');
          observer.next(res);
          observer.complete();
        },
        error: (err) => {
          // Even if logout fails on server, we clear local state
          localStorage.removeItem('isLoggedIn');
          observer.next(err);
          observer.complete();
        }
      });
    });
  }
}
