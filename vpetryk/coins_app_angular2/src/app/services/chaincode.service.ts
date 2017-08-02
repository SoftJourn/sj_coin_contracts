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
export class ChaincodeService {

  constructor(private http: Http, private data:DataService) { }

  deployChaincode(formData:string){
    let headers: Headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Authorization', this.data.user.token);

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let dataObject = Object(formData);
    dataObject.peers = ["localhost:7051","localhost:7056"]

    let url = environment.apiUrl + 'chaincodes';

    this.http.post(url, JSON.stringify(dataObject), opts)
      .subscribe((res: Response) => {
        console.log(res);
        this.data.chaincode.lastResult = res.text();
        this.data.chaincode.chaincodeName = dataObject.chaincodeName;
        this.data.chaincode.chaincodePath = dataObject.chaincodePath;
        this.data.chaincode.chaincodeVersion = dataObject.chaincodeVersion;
      });
  }

  initializeChaincode() {
    let headers: Headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Authorization', this.data.user.token);

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let dataObject = {
      peers: ["localhost:7051"],
      chaincodeName: this.data.chaincode.chaincodeName,
      chaincodeVersion: this.data.chaincode.chaincodeVersion,
      functionName: "init",
      args: ["6bc374ddc7bec9d9ef29b7025db02e99411314358a282f9a55a6c9925d09e679", "100"],
      channel: this.data.channel.currentChannel
    }

    let url = environment.apiUrl + 'channels/' + this.data.channel.currentChannel + '/chaincodes';

    this.data.chaincode.lastResult = "Initializing...";

    this.http.post(url, JSON.stringify(dataObject), opts)
      .subscribe((res: Response) => {
        console.log(res);
        this.data.chaincode.lastResult = res.text();
      });
  }

}
