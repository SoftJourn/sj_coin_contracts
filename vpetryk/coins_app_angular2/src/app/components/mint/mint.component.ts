import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../services/data.service";
import {MintService} from "../../services/mint.service";

@Component({
  selector: 'app-mint',
  templateUrl: './mint.component.html',
  styleUrls: ['./mint.component.css']
})
export class MintComponent implements OnInit {

  public amount: number = 10;

  mintForm: FormGroup;

  constructor(private mintService: MintService, fb: FormBuilder, public data:DataService) {
    this.mintForm = fb.group({
      'amount':  ['', Validators.required]
    });

  }

  ngOnInit() {
  }

  mint(formData:string) {
    this.mintService.mint(formData)
  }

}
