import { CommonModule } from '@angular/common';
import { HttpEventType } from '@angular/common/http';
import { Component, ElementRef, HostListener, ViewChild, inject } from '@angular/core';
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
    this.uploading = true;
    this.progress = 0;

    this.uploadService.uploadVideo(file).subscribe({
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
