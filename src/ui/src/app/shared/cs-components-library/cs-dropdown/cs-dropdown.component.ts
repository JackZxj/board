/**
 * Dropdown Component
 * v1.0
 * Created by liyanq on 9/4/17.
 */

import { Component, Input, Output, EventEmitter } from "@angular/core"

export const ONLY_FOR_CLICK = "OnlyClick";
const DROP_DOWN_SHOW_COUNT = 20;
export type EnableSelectCallBack = (item: Object) => boolean;
@Component({
  selector: "cs-dropdown",
  templateUrl: "./cs-dropdown.component.html",
  styleUrls: ["./cs-dropdown.component.css"]
})
export class CsDropdownComponent {
  isShowDefaultText: boolean = true;
  dropdownText: string;
  dropdownShowTimes: number = 1;
  _dropdownSearchText: string = "";
  set dropdownSearchText(value:string){
    this._dropdownSearchText = value;
  }
  get dropdownSearchText():string{
    return this._dropdownSearchText;
  }

  @Input() dropdownCanSelect: EnableSelectCallBack;
  @Input() dropdownWidth: number = 100;
  @Input() dropdownDefaultText;
  @Input() dropdownDisabled: boolean = false;
  @Input() dropdownList: Array<any>;
  @Input() dropdownListTextKey;
  @Input() dropdownTitleFontSize: number = 14;
  @Output("onChange") dropdownChange: EventEmitter<any>;
  @Output("onOnlyClickItem") dropdownClick: EventEmitter<any>;

  constructor() {
    this.dropdownChange = new EventEmitter<any>();
    this.dropdownClick = new EventEmitter<any>();
  }

  getItemClass(item: any) {
    return {
      'special': (typeof item == "object") && item['isSpecial'],
      'active': this.dropdownText == this.getItemDescription(item) ||
      this.dropdownDefaultText == this.getItemDescription(item)
    }
  }

  get dropdownShowItems(): Array<any> {
    if (this.dropdownSearchText == "") {
      return this.dropdownList ? this.dropdownList.filter(
        (value, index) => index < this.dropdownShowTimes * DROP_DOWN_SHOW_COUNT) : null;

    } else {
      let result = this.dropdownList.filter(value => {
        let text = this.getItemDescription(value);
        return text.indexOf(this.dropdownSearchText) > -1;
      });
      return result.filter((value, index) => index < this.dropdownShowTimes * DROP_DOWN_SHOW_COUNT)
    }
  }

  getItemDescription(item: any): string {
    if (typeof item == "object") {
      return item[this.dropdownListTextKey].toString();
    }
    return item.toString();
  }

  changeSelect(item: any) {
    if (item[ONLY_FOR_CLICK]) {
      this.dropdownClick.emit(item);
    } else {
      if (this.dropdownCanSelect && !this.dropdownCanSelect(item)) {
        return;
      }
      this.isShowDefaultText = false;
      if (this.dropdownText != this.getItemDescription(item)) {
        this.dropdownText = this.getItemDescription(item);
        this.dropdownChange.emit(item);
      }
    }
  }

  incShowTimes(): void {
    this.dropdownShowTimes += 1;
  }
}