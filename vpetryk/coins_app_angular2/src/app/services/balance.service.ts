import { Injectable } from '@angular/core';
import {DataService} from "./data.service";
import {
  Http,
  Response,
  RequestOptions,
  Headers
} from '@angular/http';
import { environment } from '../../environments/environment';

@Injectable()
export class BalanceService {

  constructor(private http: Http, private data:DataService) { }

  getBalance(formData:string) {
    let headers: Headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Authorization', this.data.user.token);

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let dataObject = Object(formData);
    let argsObject = {
      fcn: "balanceOf",
      peers:  ["localhost:7051","localhost:7056"],
      args: [dataObject.balanceOf]
    };

    let url = environment.apiUrl + 'channels/' + this.data.channel.currentChannel + '/chaincodes/' + this.data.chaincode.chaincodeName;

    this.http.post(url, JSON.stringify(argsObject), opts)
      .subscribe((res: Response) => {
        console.log(res);
        this.data.chaincode.lastResult = res.text();
      });
  }

  getTransactionInfo() {
    let headers: Headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Authorization', this.data.user.token);

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let url = environment.apiUrl + 'channels/' + this.data.channel.currentChannel + '/transactions/' +
      this.data.chaincode.lastResult + "?peer=peer1";

    this.http.get(url, opts)
      .subscribe((res: Response) => {
        console.log(res);
        this.data.chaincode.lastResult = res.text();
      });
  }

}
