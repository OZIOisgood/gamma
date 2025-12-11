import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, Subject, tap } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface Realm {
  id: string;
  name: string;
  created_at: string;
}

@Injectable({
  providedIn: 'root'
})
export class RealmService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;
  
  private readonly _realmsUpdated = new Subject<void>();
  readonly realmsUpdated$ = this._realmsUpdated.asObservable();

  list(): Observable<Realm[]> {
    return this.http.get<Realm[]>(`${this.apiUrl}/realms`);
  }

  create(name: string): Observable<Realm> {
    return this.http.post<Realm>(`${this.apiUrl}/realms`, { name }).pipe(
      tap(() => this._realmsUpdated.next())
    );
  }

  get(id: string): Observable<Realm> {
    return this.http.get<Realm>(`${this.apiUrl}/realms/${id}`);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/realms/${id}`).pipe(
      tap(() => this._realmsUpdated.next())
    );
  }
}
