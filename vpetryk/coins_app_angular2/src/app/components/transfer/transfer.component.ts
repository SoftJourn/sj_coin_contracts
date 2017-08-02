import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TransferService} from "../../services/transfer.service";
import {DataService} from "../../services/data.service";

@Component({
  selector: 'app-transfer',
  templateUrl: './transfer.component.html',
  styleUrls: ['./transfer.component.css']
})
export class TransferComponent implements OnInit {

  public amount: number = 10;
  public transferTo: string = "4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5";

  transferForm: FormGroup;

  constructor(private transferService: TransferService, fb: FormBuilder, public data:DataService) {
    this.transferForm = fb.group({
      'amount':  ['', Validators.required],
      'transferTo':  ['', Validators.required]
    });
  }

  ngOnInit() {
  }

  transfer(formData:string) {
    this.transferService.transfer(formData);
  }

}
