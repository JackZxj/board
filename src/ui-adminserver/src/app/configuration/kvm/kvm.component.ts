import { Component, OnInit, Input, OnChanges, SimpleChanges } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { Kvm } from '../cfg.models';

@Component({
  selector: 'app-kvm',
  templateUrl: './kvm.component.html',
  styleUrls: [
    './kvm.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class KvmComponent extends CfgComponent implements OnInit {
  cfgNum = 3;
  tipsList: Array<boolean>;
  @Input() kvm: Kvm;

  constructor() {
    super();
    this.tipsList = initTipsArray(this.cfgNum, false);
  }

  ngOnInit() {
  }

  onEdit(num: number) {
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.tipsList[num] = true;
  }
}