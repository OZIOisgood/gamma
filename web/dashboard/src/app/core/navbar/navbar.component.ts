import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NavigationEnd, Router, RouterLink } from '@angular/router';
import { TuiButton, TuiDataList, TuiDropdown, TuiOption } from '@taiga-ui/core';
import { filter } from 'rxjs/operators';
import { AuthService } from '../auth/auth.service';
import { Realm, RealmService } from '../services/realm.service';
import { UploadUiService } from '../services/upload-ui.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, FormsModule, TuiButton, RouterLink, TuiDropdown, TuiDataList, TuiOption],
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss'],
})
export class NavbarComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly uploadUiService = inject(UploadUiService);
  private readonly realmService = inject(RealmService);

  realms: Realm[] = [];
  currentRealmName: string = 'default';
  open = false;

  ngOnInit() {
    this.loadRealms();
    
    this.realmService.realmsUpdated$.subscribe(() => {
      this.loadRealms();
    });

    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe(() => {
      this.updateCurrentRealm();
    });
  }

  private loadRealms() {
    this.realmService.list().subscribe(realms => {
      this.realms = realms;
      this.updateCurrentRealm();
      
      // If current realm is no longer valid, redirect to first available realm
      const currentRealmExists = realms.some(r => r.name === this.currentRealmName);
      if (!currentRealmExists && realms.length > 0) {
        const firstRealm = realms[0].name;
        const url = this.router.url;
        const parts = url.split('/');
        if (parts.length > 1) {
          parts[1] = firstRealm;
          this.router.navigateByUrl(parts.join('/'));
        }
      }
    });
  }

  private updateCurrentRealm() {
    const url = this.router.url;
    const parts = url.split('/');
    if (parts.length > 1) {
      const realmFromUrl = parts[1];
      // Only update if the realm from URL exists in our realms list, or if we don't have realms loaded yet
      if (this.realms.length === 0 || this.realms.some(r => r.name === realmFromUrl)) {
        this.currentRealmName = realmFromUrl;
      } else if (this.realms.length > 0) {
        // If current realm doesn't exist, use the first available realm
        this.currentRealmName = this.realms[0].name;
      }
    }
  }

  getValidRealmName(): string {
    // Return a valid realm name for navigation
    if (this.realms.length > 0 && this.realms.some(r => r.name === this.currentRealmName)) {
      return this.currentRealmName;
    }
    return this.realms.length > 0 ? this.realms[0].name : 'default';
  }

  onRealmChange(realmName: string) {
    if (!realmName) return;
    const url = this.router.url;
    const parts = url.split('/');
    if (parts.length > 1) {
      parts[1] = realmName;
      this.router.navigateByUrl(parts.join('/'));
    }
    this.open = false;
  }

  openUpload() {
    this.uploadUiService.toggle();
  }

  logout() {
    this.authService.logout().subscribe({
      next: () => {
        this.router.navigate(['/login']);
      },
      error: (err) => {
        console.error('Logout failed', err);
        this.router.navigate(['/login']);
      }
    });
  }
}
