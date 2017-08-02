import {Component} from '@angular/core';
import {UserService} from "../../services/user.service";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../services/data.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {
  title = 'SJ Coins fabric application';

  public username: string = "Jim";
  public orgName: string = "org1";
  userForm: FormGroup;

  constructor(private userService: UserService, fb: FormBuilder, public data:DataService) {
    this.userForm = fb.group({
      'username':  ['', Validators.required],
      'orgName':  ['', Validators.required]
    });
  }

  enrollUser(formData: string): void {
    let object = Object(formData);
    console.log('you submitted value: ', formData);
    console.log('object.username: ', object.username);
    console.log('object.orgName: ', object.orgName);

    this.userService.getUser(object.username, object.orgName)
  }
}
