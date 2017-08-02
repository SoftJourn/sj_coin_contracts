import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../services/data.service";
import {ChaincodeService} from "../../services/chaincode.service";

@Component({
  selector: 'app-contract',
  templateUrl: './chaincode.component.html',
  styleUrls: ['./chaincode.component.css']
})
export class ChaincodeComponent implements OnInit {

  public chaincodeName: string = "coin";
  public chaincodePath: string = "github.com/coins";
  public chaincodeVersion: string = "v0";

  chaincodeForm: FormGroup;

  constructor(private chaincodeService: ChaincodeService, fb: FormBuilder, public data:DataService) {
    this.chaincodeForm = fb.group({
      'chaincodeName':  ['', Validators.required],
      'chaincodePath':  ['', Validators.required],
      'chaincodeVersion':  ['', Validators.required]
    });
  }

  ngOnInit() {
  }

  deployChaincode(formData:string):void {
    this.chaincodeService.deployChaincode(formData);
  }

  initializeChaincode():void {
    this.chaincodeService.initializeChaincode();
  }


}
