import { CommonModule } from '@angular/common';
import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TuiButton } from '@taiga-ui/core';
import { Realm, RealmService } from '../../../core/services/realm.service';
import { ConfirmDialogService } from '../../../shared/services/confirm-dialog.service';

@Component({
  selector: 'app-realm-detail',
  standalone: true,
  imports: [CommonModule, TuiButton],
  templateUrl: './realm-detail.component.html',
  styleUrls: ['./realm-detail.component.scss']
})
export class RealmDetailComponent implements OnInit {
  private readonly realmService = inject(RealmService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);

  realm = signal<Realm | null>(null);
  canDelete = signal(false);

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.realmService.get(id).subscribe(realm => {
        this.realm.set(realm);
      });
    }

    this.realmService.list().subscribe(realms => {
      this.canDelete.set(realms.length > 1);
    });
  }

  deleteRealm() {
    const currentRealm = this.realm();
    if (!currentRealm || !this.canDelete()) return;

    this.confirmDialog.confirm({
      title: 'Delete Realm',
      message: `Are you sure you want to delete realm "${currentRealm.name}"? This action cannot be undone.`,
      confirmText: 'Delete',
      cancelText: 'Cancel'
    }).subscribe(confirmed => {
      if (confirmed) {
        this.realmService.delete(currentRealm.id).subscribe(() => {
          this.realmService.list().subscribe(realms => {
            const firstRealm = realms.length > 0 ? realms[0].name : 'default';
            this.router.navigate([`/${firstRealm}/realms`]);
          });
        });
      }
    });
  }
}
