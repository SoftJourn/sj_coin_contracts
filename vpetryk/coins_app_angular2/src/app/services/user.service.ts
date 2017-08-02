import { Injectable } from '@angular/core';
import {
  Http,
  Response,
  RequestOptions,
  Headers
} from '@angular/http';
import {DataService} from "./data.service";
import { environment } from '../../environments/environment';

@Injectable()
export class UserService {

  constructor(private http: Http, private data:DataService) { }

  public getUser(username: string, orgName: string):void {

    let headers: Headers = new Headers();
    headers.append('Content-Type', 'application/json');

    let opts: RequestOptions = new RequestOptions();
    opts.headers = headers;

    let url = environment.apiUrl + 'users';

    this.http.post(url, JSON.stringify({username: username, orgName: orgName}), opts)
      .subscribe((res: Response) => {
        let json = res.json();
        if (json.success) {
          console.log(json);
          this.data.user.token = "Bearer " + json.token;
          this.data.user.username = username;
          this.data.user.orgName = orgName;
          console.log(this.data.user.username);
          console.log(this.data.user.orgName);
          this.data.user.lastResult = JSON.stringify(json);
        }
      });
  }
}
