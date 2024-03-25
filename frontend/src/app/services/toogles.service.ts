import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, of } from 'rxjs';

export interface FeatureToggle {
  ID?: number // for backend
  displayName?: string;
  technicalName: string;
  expiresOn?: Date;
  description?: string;
  customers: Customer[];
  isInverted: boolean;
  isArchived: boolean;
}

export interface Customer {
  ID: number
}


@Injectable({
  providedIn: 'root'
})
export class FeatureToggleService {
  private url = 'http://localhost:8081/api/v1';

  constructor(
    private client: HttpClient
  ) { }

  private handleError<T>(result?: T) {
    return (error: HttpErrorResponse): Observable<T> => {
      alert(error.error.error); // ...
      return of(result as T);
    }
  }

  getCustomers(): Observable<Customer[]> {
    return this.client.get<Customer[]>(this.url + '/customers')
      .pipe(
        catchError(this.handleError<Customer[]>([]))
      )
  }

  getFeatures(): Observable<FeatureToggle[]> {
    return this.client.get<FeatureToggle[]>(this.url + '/features')
      .pipe(
        catchError(this.handleError<FeatureToggle[]>([]))
      );
  }

  createFeature(feature: FeatureToggle): Observable<FeatureToggle> {
    return this.client.post<FeatureToggle>(this.url + '/features', feature)
      .pipe(
        catchError(this.handleError<FeatureToggle>())
      );
  }

  updateFeature(feature: FeatureToggle): Observable<FeatureToggle> {
    return this.client.put<FeatureToggle>(this.url + '/features/' + feature.ID, feature)
      .pipe(
        catchError(this.handleError<FeatureToggle>())
      );
  }
}
