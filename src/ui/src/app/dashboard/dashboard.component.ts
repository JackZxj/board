import { OnInit, AfterViewInit, Component, OnDestroy } from '@angular/core';
import { Assist } from "./dashboard-assist"
import { scaleOption } from "app/dashboard/time-range-scale.component/time-range-scale.component";
import { DashboardService, ServiceListModel, LinesData, LineDataModel } from "app/dashboard/dashboard.service";
import { TranslateService } from "@ngx-translate/core";
import { Subscription } from "rxjs/Subscription";

const MAX_COUNT_PER_PAGE: number = 200;
const MAX_COUNT_PER_DRAG: number = 100;
const REFRESH_SEED_SERVICE: number = 10;
@Component({
  selector: 'dashboard',
  templateUrl: 'dashboard.component.html',
  styleUrls: ['dashboard.component.css']
})

export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy {
  _intervalRead: any;
  _onLangChangeSubscription: Subscription;
  scaleOptions: Array<scaleOption> = [
    {"id": 1, "description": "DASHBOARD.MIN", "value": "second", valueOfSecond: 5},
    {"id": 2, "description": "DASHBOARD.HR", "value": "minute", valueOfSecond: 60},
    {"id": 3, "description": "DASHBOARD.DAY", "value": "hour", valueOfSecond: 60 * 60},
    {"id": 4, "description": "DASHBOARD.MTH", "value": "day", valueOfSecond: 60 * 60 * 24}];

  memoryPercent: string = '70%';
  cpuPercent: string = '40%';
  usageVolume: string = '3T';
  totalVolume: string = '10T';
  nodeBtnValue: string;
  nodeOptions: object = {};

  storageBtnValue: string;
  storageOptions: object = {};

  _serviceQuery: {model: ServiceListModel, scale: scaleOption, time_count: number, timestamp_base: number};
  _serviceOptionsBuffer: {lastZoomStart: number, lastZoomEnd: number};
  _serviceIntervalSeed: number = REFRESH_SEED_SERVICE;
  _serviceInRefresh: boolean = false;
  _serviceDropInfo: {isInDrop: boolean, isDropNext: boolean};
  serviceCurPod: number = 0;
  serviceCurContainer: number = 0;
  serviceDropdownText: string;
  serviceList: Array<ServiceListModel> = Array<ServiceListModel>();
  serviceOptions: object = {};
  serviceAlready: boolean = false;
  serviceNoData: boolean = false;
  serviceNoDataErrMsg: string = "";
  serviceData: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];

  constructor(private service: DashboardService,
              private translateService: TranslateService) {
  }

  private static concatLineData(source: LinesData, res: LinesData, isAppend: boolean): LinesData {
    let result: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
    if (source[0].length == 1 && source[0][0][1] == 0) {
      source[0] = source[0].slice(0, 0);
    }
    if (source[1].length == 1 && source[1][0][1] == 0) {
      source[1] = source[1].slice(0, 0);
    }
    if (isAppend) {
      let bufArr1: LineDataModel[] = source[0].slice(0, res[0].length);
      let bufArr2: LineDataModel[] = source[1].slice(0, res[1].length);
      result[0] = res[0].concat(bufArr1);
      result[1] = res[1].concat(bufArr2);
    } else {
      let bufArr1: LineDataModel[] = source[0].slice(res[0].length, source[0].length);
      let bufArr2: LineDataModel[] = source[1].slice(res[1].length, source[1].length);
      result[0] = bufArr1.concat(res[0]);
      result[1] = bufArr2.concat(res[1]);
    }
    return result;
  }

  private static getTimeStamp(source: LinesData, isMinTimeStamp: boolean): number {
    let timeArr: LineDataModel[] = source[0];
    if (timeArr.length == 0) {
      return Math.round(new Date().getTime() / 1000);
    }
    let date: Date = isMinTimeStamp ? timeArr[0][0] : timeArr[timeArr.length - 1][0];
    return Math.round(date.getTime() / 1000);
  }

  private static getValue(source: LinesData, lineNum: 0 | 1, isMin: boolean): number {
    let valueArr: LineDataModel[] = source[lineNum];
    if (valueArr.length == 0) {
      return 0;
    }
    return isMin ? valueArr[0][1] : valueArr[valueArr.length - 1][1]
  }

  private static haveLessMinTimeValue(timeBase, timeStepAsSecond, count: number, res: LinesData): boolean {
    let curTime = timeBase - timeStepAsSecond * 90;
    let result: boolean = false;
    res[0].forEach((item: [Date, number]) => {
      if (Math.round(item[0].getTime() / 1000) < curTime) {
        result = true;
      }
    });
    return result;
  }

  private static calculateZoom(source: LinesData, res: LinesData): {start: number, end: number} {
    let result = {start: 100, end: 80};
    if (res[0].length > 0) {
      result.start = Math.max((res[0].length / source[0].length) * 100 + 10, 20);
      result.start = result.start > 100 ? 100 : result.start;
      result.end = result.start - 20;
    }
    return result;
  }

  private resetServiceState() {
    this._serviceIntervalSeed = REFRESH_SEED_SERVICE;
    this._serviceQuery.time_count = MAX_COUNT_PER_PAGE;
    this._serviceQuery.timestamp_base = Math.round(new Date().getTime() / 1000);
    this._serviceDropInfo.isInDrop = false;
    this._serviceDropInfo.isDropNext = false;
    this.serviceAlready = false;
    this.serviceNoData = false;
    this._serviceOptionsBuffer.lastZoomEnd = 80;
    this._serviceOptionsBuffer.lastZoomStart = 100;
    this.serviceData = [[[new Date(), 1]], [[new Date(), 1]]];
  }

  serviceScaleChange(data: scaleOption) {
    if (this.serviceAlready) {
      this._serviceQuery.scale = data;
      this.resetServiceState();
      this.refreshServiceData();
    }
  }

  serviceDropDownChange(service: ServiceListModel) {
    if (this.serviceAlready) {
      this._serviceQuery.model = service;
      this.serviceDropdownText = service.service_name;
      this.resetServiceState();
      this.refreshServiceData();
    }
  }

  setServiceOption() {
    this.translateService.get(["DASHBOARD.CONTAINERS", "DASHBOARD.PODS"])
      .subscribe(res => {
        let podsTranslate: string = res["DASHBOARD.PODS"];
        let containersTranslate: string = res["DASHBOARD.CONTAINERS"];
        this.serviceOptions = Assist.getServiceOptions();
        this.serviceOptions["tooltip"] = Assist.getTooltip(podsTranslate, containersTranslate);
        this.serviceOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
        this.serviceOptions["series"][0]["name"] = podsTranslate;
        this.serviceOptions["series"][1]["name"] = containersTranslate;
        this.serviceOptions["legend"] = {data: [podsTranslate, containersTranslate], x: "left"};
      });
  }

  serviceChartDataZoom(event: object) {
    this._serviceOptionsBuffer.lastZoomEnd = event["end"];
    this._serviceOptionsBuffer.lastZoomStart = event["start"];
    if (event["start"] == 0 && !this._serviceInRefresh) {//get backup data
      this.serviceAlready = false;
      this.serviceNoData = false;
      this._serviceDropInfo.isInDrop = true;
      this._serviceDropInfo.isDropNext = true;
      this._serviceQuery.timestamp_base = DashboardComponent.getTimeStamp(this.serviceData, true);
      this._serviceQuery.time_count = MAX_COUNT_PER_DRAG;
      this.refreshServiceData();
    }
    else if (event["end"] == 100 && this._serviceDropInfo.isInDrop && !this._serviceInRefresh) {//get forward data
      this.serviceAlready = false;
      this.serviceNoData = false;
      this._serviceDropInfo.isInDrop = true;
      this._serviceDropInfo.isDropNext = false;
      this._serviceQuery.timestamp_base = DashboardComponent.getTimeStamp(this.serviceData, false) +
        this._serviceQuery.scale.valueOfSecond * MAX_COUNT_PER_DRAG;
      this._serviceQuery.time_count = MAX_COUNT_PER_DRAG;
      this.refreshServiceData();
    }
  }

  refreshServiceData() {
    let query = {
      time_count: this._serviceQuery.time_count,
      time_unit: this._serviceQuery.scale.value,
      service_name: this._serviceQuery.model.service_name == "total" ? "" : this._serviceQuery.model.service_name,
      timestamp_base: this._serviceQuery.timestamp_base
    };
    this._serviceInRefresh = true;
    this.service.getServiceData(query)
      .then(res => {
        // if (this._serviceDropInfo.isInDrop && !this._serviceDropInfo.isDropNext &&
        //   DashboardComponent.haveLessMinTimeValue(this._serviceQuery.timestamp_base,
        //     this._serviceQuery.scale.valueOfSecond, this._serviceQuery.time_count, res)) {
        //   this.resetServiceState();
        //   this.refreshServiceData();
        // } else {
          if (this._serviceDropInfo.isInDrop) {
            this.serviceOptions["dataZoom"][0]["start"] = DashboardComponent.calculateZoom(this.serviceData, res).start;
            this.serviceOptions["dataZoom"][0]["end"] = DashboardComponent.calculateZoom(this.serviceData, res).end;
          }
          else {
            this.serviceOptions["dataZoom"][0]["start"] = this._serviceOptionsBuffer.lastZoomStart;
            this.serviceOptions["dataZoom"][0]["end"] = this._serviceOptionsBuffer.lastZoomEnd;
          }
          this.serviceData = DashboardComponent.concatLineData(this.serviceData, res, this._serviceDropInfo.isDropNext);
          this.serviceCurPod = DashboardComponent.getValue(this.serviceData, 0, true);
          this.serviceCurContainer = DashboardComponent.getValue(this.serviceData, 1, true);

          this.serviceNoData = false;
          this.serviceAlready = true;
          this._serviceInRefresh = false;
          this._serviceIntervalSeed = REFRESH_SEED_SERVICE;
        // }
      })
      .catch((reason: string) => {
        this.serviceData = [[[new Date(), 1]], [[new Date(), 1]]];
        this.serviceCurPod = 0;
        this.serviceCurContainer = 0;
        this.serviceNoData = true;
        this.serviceAlready = true;
        this._serviceInRefresh = false;
        this._serviceIntervalSeed = REFRESH_SEED_SERVICE;
        this.serviceNoDataErrMsg = reason;
      });
  }

  ngOnInit() {
    this._onLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.setServiceOption();
    });
  }

  ngOnDestroy() {
    clearInterval(this._intervalRead);
    if (this._onLangChangeSubscription) {
      this._onLangChangeSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this.setServiceOption();
    this._serviceOptionsBuffer = Object.create({
      lastZoomStart: 100,
      lastZoomEnd: 80
    });
    this._serviceQuery = Object.create({
      time_count: MAX_COUNT_PER_PAGE,
      model: {service_name: ""},//default total
      timestamp_base: Math.round(new Date().getTime() / 1000),
      scale: this.scaleOptions[0]
    });
    this._serviceDropInfo = {isInDrop: false, isDropNext: false};
    this._intervalRead = setInterval(() => {
      if (this._serviceIntervalSeed > 0 && this.serviceAlready) {
        this._serviceIntervalSeed--;
        if (this._serviceIntervalSeed == 0 && !this._serviceDropInfo.isInDrop) {
          this.refreshServiceData();
        }
      }
    }, 1000);
    this.refreshServiceData();
    this.service.getServiceList()
      .then(res => {
        this.serviceList = res;
        this.serviceDropdownText = this.serviceList[0].service_name;
        this.nodeBtnValue = this.serviceList[0].service_name;
        this.storageBtnValue = this.serviceList[0].service_name;
      })
      .catch((reason: string) => {
        this.serviceDropdownText = reason;
      });

    let serviceSimulateData = DashboardService.getBySimulateData(0, 1);
    this.nodeOptions = Assist.getBaseOptions();
    this.nodeOptions["tooltip"] = Assist.getTooltip("CPU", "Memory");
    this.nodeOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
    this.nodeOptions["series"][0]["data"] = serviceSimulateData[0];
    this.nodeOptions["series"][1]["data"] = serviceSimulateData[1];

    this.storageOptions = Assist.getBaseOptions();
    this.storageOptions["tooltip"] = Assist.getTooltip("", "Total");
    this.storageOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
    this.storageOptions["series"][0]["data"] = serviceSimulateData[0];
    this.storageOptions["series"][1]["data"] = serviceSimulateData[1];
  }
}