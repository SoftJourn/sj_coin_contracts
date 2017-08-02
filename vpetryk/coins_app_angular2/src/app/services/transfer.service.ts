import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import {
  Response,
  RequestOptions,
  Headers, Http
} from '@angular/http';
import {DataService} from "./data.service";

@Injectable()
export class TransferService {

  constructor(private http: Http, private data:DataService) { }

  transfer(formData:string){
    let headers: Headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Authorization', this.data.user.token);

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let dataObject = Object(formData);
    let argsObject = {
      fcn: "transfer",
      peers:  ["localhost:7051","localhost:7056"],
      args: [dataObject.transferTo, dataObject.amount.toString()]
    };

    let url = environment.apiUrl + 'channels/' + this.data.channel.currentChannel + '/chaincodes/' + this.data.chaincode.chaincodeName;

    this.http.post(url, JSON.stringify(argsObject), opts)
      .subscribe((res: Response) => {
        console.log(res);
        this.data.chaincode.lastResult = res.text();
      });

  }

}
