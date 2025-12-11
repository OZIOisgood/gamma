import { CommonModule } from '@angular/common';
import { HttpEventType } from '@angular/common/http';
import { Component, ElementRef, HostListener, ViewChild, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TuiButton, TuiIcon } from '@taiga-ui/core';
import { AssetsService } from '../../../core/assets/assets.service';
import { UploadUiService } from '../../../core/services/upload-ui.service';
import { UploadService } from '../upload.service';

@Component({
  selector: 'app-upload-drawer',
  standalone: true,
  imports: [CommonModule, TuiButton, TuiIcon],
  templateUrl: './upload-drawer.component.html',
  styleUrls: ['./upload-drawer.component.less']
})
export class UploadDrawerComponent {
  private uploadUiService = inject(UploadUiService);
  private uploadService = inject(UploadService);
  private assetsService = inject(AssetsService);
  private route = inject(ActivatedRoute);

  isOpen$ = this.uploadUiService.isOpen$;
  isDragging = false;
  uploading = false;
  progress = 0;

  @ViewChild('fileInput') fileInput!: ElementRef<HTMLInputElement>;

  close() {
    if (!this.uploading) {
      this.uploadUiService.close();
    }
  }

  onBackdropClick(event: MouseEvent) {
    if ((event.target as HTMLElement).classList.contains('backdrop')) {
      this.close();
    }
  }

  triggerFileInput() {
    this.fileInput.nativeElement.click();
  }

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.handleFile(input.files[0]);
    }
  }

  @HostListener('dragover', ['$event'])
  onDragOver(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = true;
  }

  @HostListener('dragleave', ['$event'])
  onDragLeave(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
  }

  @HostListener('drop', ['$event'])
  onDrop(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
    
    if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
      const file = event.dataTransfer.files[0];
      if (file.type.startsWith('video/')) {
        this.handleFile(file);
      } else {
        alert('Only video files are allowed');
      }
    }
  }

  handleFile(file: File) {
    // Since this component might be outside the router outlet or not have direct access to route params
    // we might need to get the realm from the URL manually or inject a service that holds the current realm.
    // For now, let's try to get it from the window location or assume 'default' if not found.
    // A better approach would be a RealmService that holds the current state.
    
    const pathParts = window.location.pathname.split('/');
    const realm = pathParts[1] || 'default';

    this.uploading = true;
    this.progress = 0;

    this.uploadService.uploadVideo(realm, file).subscribe({
      next: (event: any) => {
        if (event.type === HttpEventType.UploadProgress) {
          this.progress = Math.round(100 * event.loaded / event.total);
        } else if (event.type === HttpEventType.Response) {
          this.uploading = false;
          this.close();
          this.assetsService.refresh(); 
        }
      },
      error: (err) => {
        console.error('Upload failed', err);
        this.uploading = false;
        alert('Upload failed');
      }
    });
  }
}
