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
export class ChannelService {

  constructor(private http: Http, private data:DataService) { }

  getChannels():void {
    let headers: Headers = new Headers();
    headers.append('Authorization', this.data.user.token);
    headers.append('Content-Type', 'application/json');

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let url = environment.apiUrl + 'channels?peer=peer1';

    this.http.get(url, opts)
      .subscribe((res: Response) => {
        let json = res.json();
        console.log(json);
        this.data.channel.channels = json.channels;
        this.data.channel.lastResult = JSON.stringify(json);
        if (this.data.channel.channels.length > 0) {
          this.data.channel.currentChannel = Object(this.data.channel.channels[0]).channel_id;
          this.data.channel.readyToConnect = true;
        }
        else {
          this.data.channel.readyToConnect = false;
          this.data.channel.readyToCreate = true;
        }
      });
  }

  createChannel():void {
    let headers: Headers = new Headers();
    headers.append('Authorization', this.data.user.token);
    headers.append('Content-Type', 'application/json');

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let url = environment.apiUrl + 'channels';

    this.http.post(url, JSON.stringify({channelName: "mychannel", channelConfigPath: "../artifacts/channel/mychannel.tx"}), opts)
      .subscribe((res: Response) => {
        let json = res.json();
        console.log(json);
        if (json.success) {
          console.log("Channel created. Please join the channel");
        }
        this.data.channel.lastResult = JSON.stringify(json);
      });
  }

  joinChannel():void {
    let headers: Headers = new Headers();
    headers.append('Authorization', this.data.user.token);
    headers.append('Content-Type', 'application/json');

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let url = environment.apiUrl + 'channels/mychannel/peers';

    this.http.post(url, JSON.stringify({peers: ["localhost:7051","localhost:7056"]}), opts)
      .subscribe((res: Response) => {
        let json = res.json();
        console.log(json);
        if (json.success) {
          this.data.channel.isConnected = true;
          console.log("Successfully connected to channel");
        }
        this.data.channel.lastResult = JSON.stringify(json);
      });
  }
}
