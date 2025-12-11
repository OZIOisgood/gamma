import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { UploadDrawerComponent } from '../../../features/upload/upload-drawer/upload-drawer.component';
import { NavbarComponent } from '../../navbar/navbar.component';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [RouterOutlet, NavbarComponent, UploadDrawerComponent],
  template: `
    <app-navbar></app-navbar>
    <app-upload-drawer></app-upload-drawer>
    <router-outlet></router-outlet>
  `,
  styles: []
})
export class MainLayoutComponent {}
