import { Component, OnInit } from '@angular/core';
import {BalanceService} from "../../services/balance.service";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../services/data.service";

@Component({
  selector: 'app-balance',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.css']
})
export class BalanceComponent implements OnInit {

  public balanceOf: string = "4e923c618bac62daeab4651c8e82d9c26e674f5cb9faf9eb0ef120a8ba00cba5";

  balanceForm: FormGroup;

  constructor(private balanceService: BalanceService, fb: FormBuilder, public data:DataService) {
    this.balanceForm = fb.group({
      'balanceOf':  ['', Validators.required]
    });
  }

  ngOnInit() {
  }

  getBalance(formData:string) {
    this.balanceService.getBalance(formData);
  }

  getTransactionInfo() {
    this.balanceService.getTransactionInfo();
  }
}
