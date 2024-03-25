import { AsyncPipe, NgFor } from '@angular/common';
import { Component, Inject } from '@angular/core';
import { FormBuilder, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { provideNativeDateAdapter } from '@angular/material/core';
import { MatDatepickerModule } from '@angular/material/datepicker';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose, MatDialogContent,
  MatDialogRef, MatDialogTitle
} from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { catchError, of } from 'rxjs';
import { Customer, FeatureToggle, FeatureToggleService } from '../services/toogles.service';

@Component({
  selector: 'app-feature-dialog',
  standalone: true,
  imports: [
    FormsModule,
    MatButtonModule,
    MatDialogTitle,
    MatDialogContent,
    MatDialogActions,
    MatDialogClose,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatCheckboxModule,
    MatIconModule,
    MatDatepickerModule,
    AsyncPipe,
    NgFor,
  ],
  providers: [provideNativeDateAdapter()],
  templateUrl: './feature-dialog.component.html',
  styleUrl: './feature-dialog.component.css'
})
export class FeatureDialogComponent {
  featureForm = this._formBuilder.group({
    displayName: [this.data.selectedFeature?.displayName],
    technicalName: [this.data.selectedFeature?.technicalName, Validators.required],
    isInverted: [this.data.selectedFeature?.isInverted || false, Validators.required],
    expiresOn: [this.data.selectedFeature?.expiresOn],
    description: [this.data.selectedFeature?.description],
    customers: [this.data.selectedFeature?.customers.map(customer => { return customer.ID }) || []],
  })

  constructor(
    private _formBuilder: FormBuilder,
    private featureService: FeatureToggleService,
    private dialogRef: MatDialogRef<FeatureDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: {
      customers: Customer[],
      selectedFeature: FeatureToggle
    },
  ) { }


  onSubmit() {
    if (!this.featureForm.valid) {
      console.error("Form is invalid");
      return;
    }

    const featureToggle: FeatureToggle = {
      displayName: this.featureForm.value.displayName!,
      technicalName: this.featureForm.value.technicalName!,
      description: this.featureForm.value.description!,
      expiresOn: this.featureForm.value.expiresOn!,
      isInverted: this.featureForm.value.isInverted!,
      customers: this.featureForm.value.customers!.map(id => ({ ID: id })) || [],
      isArchived: false,
    };

    let featureAction = this.featureService.createFeature.bind(this.featureService)
    // If we have a selected feature, we are updating it
    if (this.data.selectedFeature) {
      featureToggle.ID = this.data.selectedFeature.ID;
      featureAction = this.featureService.updateFeature.bind(this.featureService)
    }

    featureAction(featureToggle).pipe(
      catchError(error => {
        alert(error);
        return of(null); // return a new observable
      })
    ).subscribe(() => {
      this.dialogRef.close(featureToggle);
    });
  }
}
