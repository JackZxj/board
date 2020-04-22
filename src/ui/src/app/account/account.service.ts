import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Account } from './account';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { ReqSignIn } from './sign-in/sign-in.types';
export const BASE_URL = '/api/v1';

@Injectable()
export class AccountService {

  constructor(private http: HttpClient) {
  }

  signIn(signInUser: ReqSignIn): Observable<{ token: string }> {
    return this.http.post<{ token: string }>(BASE_URL + '/sign-in',
      {user_name: signInUser.username, user_password: signInUser.password},
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        params: {captcha_id: signInUser.captchaId, captcha: signInUser.captcha}
      }
    );
  }

  signUp(account: Account): Observable<any> {
    return this.http.post(
      BASE_URL + '/sign-up',
      {
        user_name: account.username,
        user_email: account.email,
        user_password: account.password,
        user_realname: account.realname,
        user_comment: account.comment
      },
      {headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE)}
    );
  }

  postEmail(credential: string): Observable<any> {
    return this.http.post(BASE_URL + `/forgot-password?credential=${credential}`, null,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response'
      });
  }

  resetPassword(password, resetUuid): Observable<any> {
    const httpParams = new HttpParams().append('password', password).append('reset_uuid', resetUuid);
    return this.http.post(BASE_URL + '/reset-password', null, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response',
      params: httpParams
    });
  }

  getCaptcha(): Observable<{ captcha_id: string }> {
    return this.http.get<{ captcha_id: string }>(BASE_URL + '/captcha');
  }

  getVerifyPicture(captchaId: string): Observable<any> {
    return this.http.get(`${BASE_URL}/captcha/${captchaId}.png`);
  }
}
