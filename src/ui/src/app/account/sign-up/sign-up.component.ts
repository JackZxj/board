import { Component, OnDestroy, OnInit } from '@angular/core';
import { Router} from '@angular/router';
import { Message } from '../../shared/message-service/message';
import { MessageService } from '../../shared/message-service/message.service';
import { SignUp } from './sign-up';
import { Account } from '../account';
import { AccountService } from '../account.service';
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { HttpErrorResponse } from "@angular/common/http";
import { Subscription } from "rxjs/Subscription";
import {AppInitService} from "../../app.init.service";
import { CsComponentBase } from "../../shared/cs-components-library/cs-component-base";

@Component({
   templateUrl: './sign-up.component.html',
   styleUrls: [ './sign-up.component.css' ]
})
export class SignUpComponent extends CsComponentBase implements OnDestroy,OnInit {
  isSignUpWIP:boolean = false;
  signUpModel: SignUp = new SignUp();
  _subscription: Subscription;

  constructor(private accountService: AccountService,
              private messageService: MessageService,
              private appInitService: AppInitService,
              private router: Router) {
    super();
    this._subscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.SIGN_UP_SUCCESSFUL) {
        this.router.navigate(['/sign-in']);
      }
    });
  }

  ngOnInit(): void {
    if(this.appInitService.systemInfo["auth_mode"] != 'db_auth') {
      this.router.navigate(['/sign-in']);
    }
  }

  ngOnDestroy() {
    this._subscription.unsubscribe();
  }
  
  signUp(): void {
    if (this.verifyInputValid()){
      this.isSignUpWIP = true;
      let account: Account = {
        username: this.signUpModel.username,
        email: this.signUpModel.email,
        password: this.signUpModel.password,
        realname: this.signUpModel.realname,
        comment: this.signUpModel.comment
      };
      this.accountService
        .signUp(account)
        .then(() => {
          this.isSignUpWIP = false;
          let confirmationMessage = new Message();
          confirmationMessage.title = "ACCOUNT.SIGN_UP";
          confirmationMessage.buttons = BUTTON_STYLE.ONLY_CONFIRM;
          confirmationMessage.target = MESSAGE_TARGET.SIGN_UP_SUCCESSFUL;
          confirmationMessage.message = 'ACCOUNT.SUCCESS_TO_SIGN_UP';
          this.messageService.announceMessage(confirmationMessage);
        })
        .catch((err: HttpErrorResponse)=>{
          this.isSignUpWIP = false;
          let confirmationMessage = new Message();
          confirmationMessage.title = "ACCOUNT.ERROR";
          confirmationMessage.buttons = BUTTON_STYLE.ONLY_CONFIRM;
          confirmationMessage.target = MESSAGE_TARGET.SIGN_UP_ERROR;
          if(err && err.status === 409) {
            confirmationMessage.message = 'ACCOUNT.USERNAME_ALREADY_EXISTS';
          } else {
            confirmationMessage.message = "ACCOUNT.FAILED_TO_SIGN_UP";
          }
          this.messageService.announceMessage(confirmationMessage);
        });
    }
  }

  goBack(): void {
    this.router.navigate(['/sign-in']);
  }
}