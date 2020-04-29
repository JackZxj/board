import {HTTP_INTERCEPTORS, HttpErrorResponse, HttpInterceptor, HttpRequest, HttpResponse} from '@angular/common/http';
import {HttpHandler} from '@angular/common/http/src/backend';
import {HttpEvent} from '@angular/common/http/src/response';
import {Injectable} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';
import {Observable, of, throwError, TimeoutError} from 'rxjs';
import {catchError, timeout} from 'rxjs/operators';
import {MessageService} from './message/message.service';
import {GlobalAlertType} from './message/message.types';

@Injectable()
export class HttpClientInterceptor implements HttpInterceptor {

  constructor(private messageService: MessageService,
              private translateService: TranslateService) {

  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let authReq: HttpRequest<any> = req.clone({
      headers: req.headers
    });
    authReq = authReq.clone({
      params: authReq.params.set('Timestamp', Date.now().toString())
    });
    return next.handle(authReq)
      .pipe(
        timeout(120 * 1000),
        catchError((err: HttpErrorResponse | TimeoutError) => {
          if (err instanceof HttpErrorResponse) {
            if (err.status >= 200 && err.status < 300) {
              const res = new HttpResponse({
                body: null,
                headers: err.headers,
                status: err.status,
                statusText: err.statusText,
                url: err.url
              });
              return of(res);
            } else if (err.status === 502) {
              this.messageService.showGlobalMessage('ERROR.HTTP_502', {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status === 504) {
              this.messageService.showGlobalMessage('ERROR.HTTP_504', {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status === 500) {
              this.messageService.showGlobalMessage('ERROR.HTTP_500', {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status === 400) {
              this.messageService.showGlobalMessage(`ERROR.HTTP_400`, {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status === 401) {
              this.messageService.showGlobalMessage(`ERROR.HTTP_401`, {
                globalAlertType: GlobalAlertType.gatLogin,
                alertType: 'warning'
              });
            } else if (err.status === 403) {
              this.messageService.showAlert(`ERROR.HTTP_403`, {alertType: 'danger'});
            } else if (err.status === 404) {
              this.messageService.showAlert(`ERROR.HTTP_404`, {alertType: 'danger'});
            } else if (err.status === 412) {
              this.messageService.showAlert(`ERROR.HTTP_412`, {alertType: 'warning'});
            } else if (err.status === 422) {
              this.translateService.get(`ERROR.HTTP_422`).subscribe((msg: string) => {
                const alertMsg = `${msg},${err.error}`;
                this.messageService.showAlert(alertMsg, {alertType: 'danger'});
              });
            }
          } else {
            this.messageService.showGlobalMessage(`ERROR.HTTP_TIME_OUT`, {
              globalAlertType: GlobalAlertType.gatShowDetail,
              errorObject: err,
              endMessage: req.url
            });
          }
          return throwError(err);
        }));
  }
}

export const HttpInterceptorService = {
  provide: HTTP_INTERCEPTORS,
  useClass: HttpClientInterceptor,
  deps: [MessageService, TranslateService],
  multi: true
};
