import { DatePipe, NgFor, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatCheckbox } from '@angular/material/checkbox';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule, MatIconRegistry } from '@angular/material/icon';
import { DomSanitizer } from '@angular/platform-browser';
import { FeatureDialogComponent } from './feature-dialog/feature-dialog.component';
import { FeatureToggle, FeatureToggleService } from './services/feature.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports:
    [
      NgFor,
      NgIf,
      DatePipe,
      FormsModule,

      MatCardModule,
      MatButtonModule,
      MatIconModule,
      MatCheckbox,
    ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})

export class AppComponent {
  title = 'featureToggles';
  features: FeatureToggle[] = [];
  showArchived = false;

  constructor(
    // Fetching data
    private featureToggleService: FeatureToggleService,
    // Icons
    private iconRegistry: MatIconRegistry,
    private sanitizer: DomSanitizer,
    //
    public dialog: MatDialog
  ) { }

  updateList() { }

  openDialog(selectedFeature?: FeatureToggle) { // update only if saved
    this.featureToggleService.getCustomers().subscribe(customers => {
      this.dialog.open(FeatureDialogComponent, {
        width: '400px',
        data: {
          customers: customers,
          selectedFeature: selectedFeature
        }
      }).afterClosed().subscribe((feature?: FeatureToggle) => {
        // If dialog closed and feature is undefined, that means it was canceled
        if (!feature) {
          return;
        }

        if (selectedFeature) {
          // If feature was selected, then update it
          const idx = this.features.findIndex(feature => feature.ID === selectedFeature.ID);
          this.features[idx] = feature;
        } else {
          // If feature is new, then add it
          feature.ID = this.features.length + 1;
          this.features.push(feature);
        }
      });
    });
  }

  archive(feature: FeatureToggle) {
    feature.isArchived = !feature.isArchived;
    this.featureToggleService.updateFeature(feature).subscribe();
  }

  ngOnInit() {
    this.featureToggleService.getFeatures().subscribe((features) => this.features = features);
    this.iconRegistry.addSvgIcon(
      'archive',
      this.sanitizer.bypassSecurityTrustResourceUrl('assets/archive.svg')
    );
    this.iconRegistry.addSvgIcon(
      'edit',
      this.sanitizer.bypassSecurityTrustResourceUrl('assets/edit.svg')
    )
  }
}
