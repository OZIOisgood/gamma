import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { TuiTable } from '@taiga-ui/addon-table';
import { Realm, RealmService } from '../../core/services/realm.service';
import { DrawerComponent } from '../../shared/components/drawer/drawer.component';
import { PageHeaderComponent } from '../../shared/components/page-header/page-header.component';
import { CreateRealmComponent } from './create-realm/create-realm.component';

@Component({
  selector: 'app-realms',
  standalone: true,
  imports: [
    CommonModule, 
    TuiTable,
    PageHeaderComponent,
    DrawerComponent,
    CreateRealmComponent
  ],
  templateUrl: './realms.component.html',
  styleUrls: ['./realms.component.scss']
})
export class RealmsComponent implements OnInit {
  private readonly realmService = inject(RealmService);
  private readonly router = inject(Router);

  realms = signal<Realm[]>([]);
  isDrawerOpen = signal(false);
  columns = ['name', 'created_at'];

  ngOnInit() {
    this.loadRealms();
  }

  loadRealms() {
    this.realmService.list().subscribe(realms => {
      this.realms.set(realms);
    });
  }

  onRowClick(realm: Realm) {
    const currentUrl = this.router.url;
    this.router.navigate([currentUrl, realm.id]);
  }

  openDrawer() {
    this.isDrawerOpen.set(true);
  }

  closeDrawer() {
    this.isDrawerOpen.set(false);
  }

  onRealmCreated() {
    this.closeDrawer();
    this.loadRealms();
  }
}
